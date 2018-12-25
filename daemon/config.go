package daemon

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	LogLevel string `json:"loglevel"`
	Slack    struct {
		BotToken  string `json:"bottoken"`
		BotID     string `json:"botid"`
		ChannelID string `json:"channelid"`
	} `json:"slack"`
	Emby struct {
		AdminID string `json:"adminid"`
		Path    string `json:"path"`
		Token   string `json:"token"`
	} `json:"emby"`
	Radarr struct {
		Path   string `json:"path"`
		APIKey string `json:"apikey"`
	} `json:"radarr"`
	TLSConfig TLSConfig `json:"tlsconfig"`
}

func loadConfig(file string) (*config, error) {
	var config config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode configs: %v", err)
	}

	return &config, nil
}
