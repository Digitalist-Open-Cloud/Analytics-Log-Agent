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
