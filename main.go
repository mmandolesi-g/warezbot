package main

import (
	"flag"
	"os"
	"warezbot/daemon"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

var (
	listenAddr     = flag.String("httpListen", ":3000", "Listen address for the http service")
	requestLogPath = flag.String("requestLog", "/var/log/request.log", "Path to the request log file")
	configFile     = flag.String("configFile", "./daemon/config.json", "Path to the config file")
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	api, err := daemon.NewWarezDaemon(logger, *requestLogPath, *configFile)
	if err != nil {
		level.Error(logger).Log("error", err)
		os.Exit(1)
	}

	if err := api.Run(*listenAddr); err != nil {
		level.Error(logger).Log("error", err)
		os.Exit(1)
	}
}
