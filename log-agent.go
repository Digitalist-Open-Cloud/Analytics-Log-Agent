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
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Function to simulate cat command with rate-limiting
func catLogFile(config *Config, requestsPerSec int) error {
	file, err := os.Open(config.Log.LogPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	// Calculate the delay between requests based on the configured rate
	delay := time.Second / time.Duration(requestsPerSec)

	// Create a ticker to control the request rate
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Parse the log line
		logData := parseLog(line, config.Log.LogFormat)
		if logData == nil {
			logger.Warnf("Failed to parse log line: %s", line)
			continue
		}

		// Wait for the next tick to respect the rate limit
		<-ticker.C

		// Send the parsed log to Matomo
		sendToMatomo(logData, config)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	logger.Info("Finished processing log file in catlog mode")
	return nil
}

func main() {
	// Define flags

	configPath := flag.String("config", "/opt/log-agent/config.toml", "Path to the configuration file")
	catLog := flag.Bool("catlog", false, "Simulate cat command for a log file")
	reqPerSec := flag.Int("rps", 1, "Requests per second limit for catlog mode")
	matomoURL := flag.String("matomo-url", "", "Matomo URL")
	tokenAuth := flag.String("token-auth", "", "Matomo token auth")
	siteID := flag.String("site-id", "", "Matomo site ID")
	pluginEnabled := flag.Bool("plugin", false, "If using the Matomo Agent plugin")
	downloadsEnabled := flag.Bool("downloads", true, "Enable download tracking")
	logFormat := flag.String("log-format", "", "Log format (nginx, apache or csv)")
	logPath := flag.String("log-path", "", "Path to the log file")
	userAgents := flag.String("user-agents", "", "Comma-separated list of user agents to track (Overrides config file)")
	agentLogLevel := flag.String("log-level", "", "Log level (debug, info, warn, error) (Overrides config file)")
	agentLogFile := flag.String("log-file", "", "Path to the agent's log file (Overrides config file)")
	isTitleEnabled := flag.Bool("collect-title", false, "Enable collection of page titles based on URL")
	titleDomain := flag.String("title-domain", "", "Override default domain to fetch title from")
	batchMode := flag.Bool("batch", false, "Enable batch mode for sending logs")

	// Parse the flags first
	flag.Parse()

	// Load the config file
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if *tokenAuth != "" {
		config.Matomo.TokenAuth = *tokenAuth
	}
	if *matomoURL != "" {
		config.Matomo.URL = *matomoURL
	}
	if *siteID != "" {
		config.Matomo.SiteID = *siteID
	}
	if *pluginEnabled {
		config.Matomo.Plugin = *pluginEnabled
	}
	if *downloadsEnabled {
		config.Matomo.Downloads = *downloadsEnabled
	}
	if *logFormat != "" {
		config.Log.LogFormat = *logFormat
	}
	if *logPath != "" {
		config.Log.LogPath = *logPath
	}
	if *userAgents != "" {
		config.Log.UserAgents = strings.Split(*userAgents, ",")
	}
	if *agentLogLevel != "" {
		config.Agent.LogLevel = *agentLogLevel
	}
	if *agentLogFile != "" {
		config.Agent.LogFile = *agentLogFile
	}

	if *isTitleEnabled {
		config.Title.Collect = *isTitleEnabled
	}

	if *titleDomain != "" {
		config.Title.Domain = *titleDomain
	}

	if *batchMode {
		config.Batch.Mode = *batchMode
	}

	// Override config with flag values
	//overrideConfigWithFlags(config)

	// Set up logging (call once after flags and config are processed)
	setupLogging(config.Agent.LogLevel, config.Agent.LogFile)

	// Validate Matomo token
	err = validateTokenAuth(config)
	if err != nil {
		logger.Fatal("Invalid Matomo token:", err)
	}

	// Check if catlog mode is enabled
	if *catLog {
		logger.Infof("Starting in catlog mode, sending %d requests per second", *reqPerSec)
		err = catLogFile(config, *reqPerSec)
		if err != nil {
			logger.Fatalf("Error in catlog mode: %v", err)
		}
	} else {
		// Default log tailing mode
		logger.Infof("Start tailing the log")
		tailLogFile(config)
	}
}
