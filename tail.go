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
	if config.Log.LogFormat == "nginx" {
		logFilePath = config.Log.LogPath
	} else if config.Log.LogFormat == "apache" {
		logFilePath = config.Log.LogPath
	} else {
		logger.Fatal("Invalid log format in config")
	}

	t, err := tail.TailFile(logFilePath, tail.Config{Follow: true})
	if err != nil {
		logger.Fatal("Failed to open log file:", err)
	}

	for line := range t.Lines {
		logData := parseLog(line.Text, config.Log.LogFormat)
		if logData != nil {
			sendToMatomo(logData, config)
		}
	}
}
