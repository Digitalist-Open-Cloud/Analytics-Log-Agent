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
	"encoding/csv"
	"io"
	"regexp"
	"strings"
)

// Struct of log data.
type LogData struct {
	IP        string
	Timestamp string
	Host      string
	Method    string
	URL       string
	Protocol  string
	Status    string
	Size      string
	Referrer  string
	UserAgent string
	Hour      string
	Minute    string
	Second    string
}

// Parse log line for Nginx, Apache, or CSV.
func parseLog(line, logFormat string) *LogData {
	var logPattern *regexp.Regexp
	var logData *LogData

	// Nginx and Apache log formats (same regex for both)
	if logFormat == "nginx" || logFormat == "apache" {
		logPattern = regexp.MustCompile(`(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<url>\S+) (?P<protocol>\S+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^"]*)" "(?P<user_agent>[^"]*)"`)
	} else if logFormat == "csv" {
		// CSV log format: timestamp, req_method, final_host, req_uri, resp_status, client_ip, req_referer, req_user_agent

		// Example: 2024-09-10 18:13:19 UTC,GET,example.org,/news/2024/09/09/south-sudan-parliament-approves-transitional-justice-laws,200,2.86.73.63,https://my.url/,Mozilla/5.0...
		reader := csv.NewReader(strings.NewReader(line))
		reader.TrimLeadingSpace = true // In case there are extra spaces

		// Read one record (log line)
		record, err := reader.Read()
		if err == io.EOF {
			return nil // End of input
		} else if err != nil {
			logger.Warnf("Failed to parse CSV log line: %s", err)
			return nil
		}
		if len(record) < 8 {
			logger.Warn("Invalid CSV log format")
			return nil
		}

		logData = &LogData{
			Timestamp: record[0],
			Method:    record[1],
			Host:      record[2],
			URL:       record[3],
			Status:    record[4],
			IP:        record[5],
			Referrer:  record[5],
			UserAgent: record[7],
		}
	} else {
		logger.Warn("Unknown log format")

		return nil
	}

	if logFormat == "nginx" || logFormat == "apache" {
		// Match the Nginx/Apache log format
		match := logPattern.FindStringSubmatch(line)
		if match != nil {
			logData = &LogData{
				IP:        match[1],
				Timestamp: match[2],
				Method:    match[3],
				URL:       match[4],
				Protocol:  match[5],
				Status:    match[6],
				Size:      match[7],
				Referrer:  match[8],
				UserAgent: match[9],
			}
		} else {
			return nil
		}
	}

	// Parse the timestamp and extract hour, minute, second
	if logData != nil && logData.Timestamp != "" {
		h, m, s, err := parseTimestamp(logData.Timestamp)
		if err != nil {
			logger.Warnf("Error parsing timestamp: %v", err)
		} else {
			logData.Hour = h
			logData.Minute = m
			logData.Second = s
		}
	}

	return logData
}
