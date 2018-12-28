package daemon

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/mmandolesi-g/warezbot/warez"
)

const (
	slackProcessPath = "/verify"
	slackInteractive = "/slack/actions"

	DefaultHTTPIdleTimeout       = 30 * time.Second // The timeout before unused open connections are close
	DefaultHTTPReadHeaderTimeout = 5 * time.Second  // The max time to read the request header
	DefaultHTTPWriteTimeout      = 15 * time.Second // The max time to read and respond to the request, including, db/cache lookup
)

type TLSConfig struct {
	TLSCA   string `json:"tlsca"`
	TLSCert string `json:"tlscert"`
	TLSKey  string `json:"tlskey"`
}

type HTTPSConfig struct {
	DisableCompression bool
	DisableLogging     bool
	Handler            http.Handler
	IdleTimeout        time.Duration
	ReadHeaderTimeout  time.Duration
	WriteTimeout       time.Duration
	Logger             log.Logger
	LogPath            string
	TLSCfg             TLSConfig
}

type HTTPSDaemon struct {
	caPool            *x509.CertPool
	cert              tls.Certificate
	handler           http.Handler
	idleTimeout       time.Duration
	logger            log.Logger
	quit              chan bool
	readHeaderTimeout time.Duration
	writeTimeout      time.Duration
}

func (dm *WarezDaemon) setupHTTP(svc warez.Service) http.Handler {
	router := mux.NewRouter()

	var slackEventEndpoint endpoint.Endpoint
	{
		slackEventEndpoint = slackProcessEndpoint(svc.ProcessSlackEvents)
	}
	var slackEventHandler http.Handler
	{
		slackEventHandler = httptransport.NewServer(
			slackEventEndpoint,
			dm.decodeSlackEvent,
			dm.encodeWarezResponse)
	}
	router.Methods("POST").Path(slackProcessPath).Handler(slackEventHandler)

	var slackActionEndpoint endpoint.Endpoint
	{
		slackActionEndpoint = slackProcessActionEndpoint(svc.ProcessSlackActions)
	}
	var slackActionHandler http.Handler
	{
		slackActionHandler = httptransport.NewServer(
			slackActionEndpoint,
			dm.decodeSlackAction,
			dm.encodeWarezNilResponse)
	}
	router.Methods("POST").Path(slackInteractive).Handler(slackActionHandler)

	return router
}

func (dm *WarezDaemon) decodeSlackAction(ctx context.Context, r *http.Request) (interface{}, error) {
	var s warez.SlackAction

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	b, err := url.QueryUnescape(string(body))
	if err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("error unescaping payload: %v", err)
	}

	b = strings.TrimLeft(b, "payload=")
	if err := json.Unmarshal([]byte(b), &s); err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("error unmarshaling action request: %v", err)
	}

	return s, nil
}

func (dm *WarezDaemon) decodeSlackEvent(ctx context.Context, r *http.Request) (interface{}, error) {
	var s warez.SlackEvent

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("error reading request body: %v", err)
	}

	if err := json.Unmarshal(body, &s); err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("error unmarshaling request: %v", err)
	}

	return s, nil
}

func (dm *WarezDaemon) encodeWarezNilResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp, ok := response.(warez.Response)
	if !ok {
		return errors.New("endpoint response error")
	}
	if resp.StatusCode != http.StatusAccepted {
		w.WriteHeader(resp.StatusCode)
		return errors.New("something went wrong")
	}
	w.WriteHeader(resp.StatusCode)
	return nil
}

func (dm *WarezDaemon) encodeWarezResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp, ok := response.(warez.Response)
	if !ok {
		return errors.New("endpoint response error")
	}

	return json.NewEncoder(w).Encode(resp)
}

func slackProcessActionEndpoint(searchFunc warez.SlackActionFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(warez.SlackAction)
		if !ok {
			return nil, fmt.Errorf("unknown request data")
		}

		return searchFunc(ctx, req)
	}
}

func slackProcessEndpoint(searchFunc warez.SlackEventFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(warez.SlackEvent)
		if !ok {
			return nil, fmt.Errorf("unknown request data")
		}

		return searchFunc(ctx, req)
	}
}

// http related functions below

func NewHTTPDaemon(cfg HTTPSConfig) (*HTTPSDaemon, error) {
	handler := cfg.Handler

	cert, caPool, err := loadTLSCertificates(cfg.TLSCfg.TLSCert, cfg.TLSCfg.TLSKey, cfg.TLSCfg.TLSCA)
	if err != nil {
		return nil, err
	}

	d := &HTTPSDaemon{
		caPool:            caPool,
		cert:              cert,
		handler:           handler,
		logger:            log.With(cfg.Logger, "source", "HTTPSDaemon"),
		quit:              make(chan bool),
		idleTimeout:       cfg.IdleTimeout,
		readHeaderTimeout: cfg.ReadHeaderTimeout,
		writeTimeout:      cfg.WriteTimeout,
	}

	if d.idleTimeout == 0 {
		d.idleTimeout = DefaultHTTPIdleTimeout
	}
	if d.readHeaderTimeout == 0 {
		d.readHeaderTimeout = DefaultHTTPReadHeaderTimeout
	}
	if d.writeTimeout == 0 {
		d.writeTimeout = DefaultHTTPWriteTimeout
	}

	return d, nil
}

// Run does this and that
func (d *HTTPSDaemon) Run(httpListenAddr string) error {
	errorChan := make(chan error)

	go func() {
		level.Info(d.logger).Log("event", "Starting HTTPS server", "httpListen", httpListenAddr)
		srv := &http.Server{
			Addr:              httpListenAddr,
			Handler:           d.handler,
			IdleTimeout:       d.idleTimeout,
			ReadHeaderTimeout: d.readHeaderTimeout,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{d.cert},
			},
			WriteTimeout: d.writeTimeout,
		}
		go func() {
			if err := srv.ListenAndServeTLS("", ""); err != nil {
				level.Error(d.logger).Log("event", "failed to start HTTPS server", "error", err)
				close(d.quit)
			}
		}()

		d.handleSignals()

		<-d.quit
		level.Info(d.logger).Log("event", "Terminating HTTPS Server")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		errorChan <- srv.Shutdown(ctx)
	}()
	return <-errorChan
}

func (d *HTTPSDaemon) handleSignals() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		s := <-c
		level.Info(d.logger).Log("signal", s, "event", "received os signal, stopping execution")
		close(d.quit)
	}()
}

func loadTLSCertificates(cert, key, ca string) (tls.Certificate, *x509.CertPool, error) {
	tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return tls.Certificate{}, nil, fmt.Errorf("failed to load TLS certificates and key: %v", err)
	}

	pool := x509.NewCertPool()
	if ok := pool.AppendCertsFromPEM([]byte(ca)); !ok {
		return tls.Certificate{}, nil, errors.New("failed to load TLS CA: %v")
	}

	return tlsCert, pool, nil
}
