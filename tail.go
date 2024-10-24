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
	"github.com/tenebris-tech/tail"
)

// Tail the log file based on configuration and send to Matomo
func tailLogFile(config *Config) {
	var logFilePath string
	if config.Log.LogFormat == "nginx" || config.Log.LogFormat == "apache" || config.Log.LogFormat == "csv" {
		logFilePath = config.Log.LogPath
	} else {
		logger.Fatal("Invalid log format in config")
	}

	// Open the log file for tailing
	t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
	if err != nil {
		logger.Fatal("Failed to open log file:", err)
	}

	// Process each line from the log file
	for line := range t.Lines {
		// Parse the log line
		logData := parseLog(line.Text, config.Log.LogFormat)
		if logData == nil {
			logger.Warnf("Failed to parse log line: %s", line.Text)
			continue
		}

		// Check if the request URL contains an ignored media file extension (without query params)
		if isIgnored(logData.URL) {
			logger.Debugf("Skipping media file request: %s", logData.URL)
			continue
		}

		// Send parsed log to Matomo
		sendToMatomo(logData, config)
	}
}
