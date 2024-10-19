package main

import (
	"regexp"
)

// Struct of log data.
type LogData struct {
	IP        string
	Timestamp string
	Method    string
	URL       string
	Protocol  string
	Status    string
	Size      string
	Referrer  string
	UserAgent string
}

// Parse log line for Nginx or Apache - for now these are the same.
func parseLog(line, logFormat string) *LogData {
	var logPattern *regexp.Regexp
	if logFormat == "nginx" {
		logPattern = regexp.MustCompile(`(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<url>\S+) (?P<protocol>\S+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^"]*)" "(?P<user_agent>[^"]*)"`)
	} else if logFormat == "apache" {
		logPattern = regexp.MustCompile(`(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<url>\S+) (?P<protocol>\S+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^"]*)" "(?P<user_agent>[^"]*)"`)
	} else {
		logger.Warn("Unknown log format")
		return nil
	}

	match := logPattern.FindStringSubmatch(line)
	if match != nil {
		return &LogData{
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
	}
	return nil
}
