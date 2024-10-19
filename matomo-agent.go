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
