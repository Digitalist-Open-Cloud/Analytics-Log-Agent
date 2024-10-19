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
