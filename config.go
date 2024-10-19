package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Matomo struct {
		URL       string `mapstructure:"url"`
		AgentURL  string
		SiteID    string `mapstructure:"site_id"`
		WebSite   string `mapstructure:"website_url"`
		TokenAuth string `mapstructure:"token_auth"`
		Plugin    bool   `mapstructure:"plugin"`
	}
	Log struct {
		LogFormat  string   `mapstructure:"log_format"`
		LogPath    string   `mapstructure:"log_path"`
		UserAgents []string `mapstructure:"user_agents"`
	}
	Agent struct {
		LogLevel string `mapstructure:"log_level"`
		LogFile  string `mapstructure:"log_file"`
	}
}

func loadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
