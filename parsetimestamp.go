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
	"strings"
	"time"
)

// Function to parse the timestamp and extract hour, minute, second
func parseTimestamp(timestamp string) (hour, minute, second string, err error) {
	// List of possible time formats
	timeFormats := []string{
		"02/Jan/2006:15:04:05 -0700", // Format like 23/Oct/2024:12:19:08 +0200
		"2006-01-02 15:04:05 UTC",    // Format like 2024-09-10 18:13:19 UTC
	}

	var parsedTime time.Time

	// Try parsing with the different formats
	for _, format := range timeFormats {
		parsedTime, err = time.Parse(format, timestamp)
		if err == nil {
			// If we successfully parsed the time, break out of the loop
			break
		}
	}

	if err != nil {
		return "", "", "", err
	}

	// Extract hour, minute, second from the parsed time
	hour = parsedTime.Format("15")   // Extract hour in 24-hour format
	minute = parsedTime.Format("04") // Extract minute
	second = parsedTime.Format("05") // Extract second

	return hour, minute, second, nil
}

func formatTimestamp(timestamp string) (string, error) {
	var parsedTime time.Time
	var err error

	// Try multiple possible formats (e.g., Apache, Nginx, custom)
	formatOne := "2006-01-02 15:04:05 UTC"
	formatTwo := "02/Jan/2006:15:04:05 -0700"

	if strings.Contains(timestamp, "UTC") {
		parsedTime, err = time.Parse(formatOne, timestamp)
	} else {
		parsedTime, err = time.Parse(formatTwo, timestamp)
	}

	if err != nil {
		return "", err
	}

	// Format time in the required format for Matomo "YYYY-MM-DD HH:MM:SS"
	return parsedTime.Format("2006-01-02 15:04:05"), nil
}
