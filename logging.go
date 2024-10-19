package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger = logrus.New()

func setupLogging(logLevel string, logFile string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Output logs to a file
	if logFile != "" {
		logger.SetOutput(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10, // Megabytes
			MaxBackups: 3,
			MaxAge:     28,   // Days
			Compress:   true, // Compress old logs
		})
	} else {
		// Fallback to stderr if no log file is provided
		logger.Warn("No log file provided, using default stderr")
		logger.SetOutput(os.Stderr)
	}
}
