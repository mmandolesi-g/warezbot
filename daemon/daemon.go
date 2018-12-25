package daemon

import (
	"fmt"
	"warezbot/emby"
	"warezbot/radarr"
	"warezbot/slack"
	"warezbot/warez"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

type WarezDaemon struct {
	*HTTPSDaemon
	apiKey string
	logger log.Logger
}

func NewWarezDaemon(logger log.Logger, requestLogPath string, config string) (*WarezDaemon, error) {
	cfg, err := loadConfig(config)
	if err != nil {
		return nil, err
	}

	logger = SetLoggerLevel(logger, cfg.LogLevel)

	embyClient, err := emby.NewClient(cfg.Emby.Path, cfg.Emby.Token, cfg.Emby.AdminID)
	if err != nil {
		return nil, err
	}
	radarrClient, err := radarr.NewClient(cfg.Radarr.Path, cfg.Radarr.APIKey)
	if err != nil {
		return nil, err
	}
	slackClient, err := slack.NewClient(cfg.Slack.BotToken, cfg.Slack.ChannelID, cfg.Slack.BotID)
	if err != nil {
		return nil, err
	}

	svc, err := warez.NewService(embyClient, radarrClient, slackClient)
	if err != nil {
		return nil, err
	}

	d := &WarezDaemon{
		logger: logger,
	}

	d.HTTPSDaemon, err = NewHTTPDaemon(HTTPSConfig{
		Handler: d.setupHTTP(svc),
		Logger:  logger,
		LogPath: requestLogPath,
		TLSCfg:  cfg.TLSConfig,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize HTTPS Daemon: %v", err)
	}

	return d, nil
}

func SetLoggerLevel(logger log.Logger, levelName string) log.Logger {
	switch levelName {
	case "debug":
		return level.NewFilter(logger, level.AllowDebug())
	case "info":
		return level.NewFilter(logger, level.AllowInfo())
	case "error":
		return level.NewFilter(logger, level.AllowError())
	default:
		return level.NewFilter(logger, level.AllowWarn())
	}
}
