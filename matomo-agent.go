/**
 * An agent for Matomo.
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
	"flag"
	"log"
)

func main() {
	configPath := flag.String("config", "/opt/matomo-agent/config.toml", "Path to the configuration file")
	flag.Parse()
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up logging
	setupLogging(config.Agent.LogLevel, config.Agent.LogFile)
	// Check if we have a valid token for Matomo.
	err = validateTokenAuth(config)
	if err != nil {
		logger.Fatal("Invalid Matomo token:", err)
	}

	// All set, start tailing the log file
	logger.Infof("Start tailing the log")
	tailLogFile(config)
}
