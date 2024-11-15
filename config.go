/**
 * A log agent for Matomo.
 *
 * Copyright (C) 2024 Digitalist Open Cloud <cloud@digitalist.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Matomo struct {
		URL        string `mapstructure:"url"`
		TrackerURL string `mapstructure:"tracker_url"`
		AgentURL   string
		SiteID     string `mapstructure:"site_id"`
		WebSite    string `mapstructure:"website_url"`
		TokenAuth  string `mapstructure:"token_auth"`
		Plugin     bool   `mapstructure:"plugin"`
		Downloads  bool   `mapstructure:"downloads"`
	}
	Log struct {
		LogFormat    string   `mapstructure:"log_format"`
		LogPath      string   `mapstructure:"log_path"`
		UserAgents   []string `mapstructure:"user_agents"`
		ExcludedURLs []string `mapstructure:"excluded_urls"`
	}
	Agent struct {
		LogLevel string `mapstructure:"log_level"`
		LogFile  string `mapstructure:"log_file"`
	}
	Title struct {
		Collect bool   `mapstructure:"collect_titles"`
		Domain  string `mapstructure:"title_domain"`
		Cache   string `mapstructure:"cache_file"`
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
