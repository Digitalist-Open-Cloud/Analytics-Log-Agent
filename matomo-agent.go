package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tenebris-tech/tail"
)

type Config struct {
	Matomo struct {
		URL       string `mapstructure:"url"`
		SiteID    string `mapstructure:"site_id"`
		WebSite   string `mapstructure:"website_url"`
		TokenAuth string `mapstructure:"token_auth"`
	}
	Log struct {
		LogFormat string `mapstructure:"log_format"`
		LogLevel  string `mapstructure:"log_level"`
		LogFile   string `mapstructure:"log_file"`
	}
	Nginx struct {
		LogPath string `mapstructure:"log_path"`
	}
	Apache struct {
		LogPath string `mapstructure:"log_path"`
	}
}

// Global logger instance
var logger = logrus.New()

// Load configuration from /opt/matomo-agent/config.toml
func loadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}

// Set up logging levels
func setupLogging(logLevel string, logFile string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	// Output logs to a file
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			logger.Warn("Failed to log to file, using default stderr")
		}
	}
}

// Matomo Tracking API call
func sendToMatomo(logData *LogData, config *Config) {
	var Url = config.Matomo.WebSite + logData.URL
	data := url.Values{
		"idsite":     {config.Matomo.SiteID},
		"rec":        {"1"},
		"cip":        {logData.IP},
		"ua":         {logData.UserAgent},
		"url":        {Url},
		"token_auth": {config.Matomo.TokenAuth},
	}

	resp, err := http.Get(config.Matomo.URL + "?" + data.Encode())
	if err != nil {
		logger.Error("Error sending data to Matomo:", err)
	} else {
		logger.Infof("Log sent: %s, Status: %s", logData.URL, resp.Status)
	}
}

// Define a struct to hold the parsed log data
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

// Parse log line for Nginx or Apache
func parseLog(line, logFormat string) *LogData {
	var logPattern *regexp.Regexp
	if logFormat == "nginx" {
		logPattern = regexp.MustCompile(`(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<url>\S+) (?P<protocol>\S+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^"]*)" "(?P<user_agent>[^"]*)"`)
	} else if logFormat == "apache" {
		logPattern = regexp.MustCompile(`(?P<ip>\S+) - - \[(?P<time>[^\]]+)\] "(?P<method>\S+) (?P<url>\S+) (?P<protocol>\S+)" (?P<status>\d+) (?P<size>\d+) "(?P<referrer>[^\"]*)" "(?P<user_agent>[^\"]*)`)
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

// Tail the log file based on configuration and send to Matomo
func tailLogFile(config *Config) {
	var logFilePath string
	if config.Log.LogFormat == "nginx" {
		logFilePath = config.Nginx.LogPath
	} else if config.Log.LogFormat == "apache" {
		logFilePath = config.Apache.LogPath
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

func main() {
	configPath := flag.String("config", "/opt/matomo-agent/config.toml", "Path to the configuration file")
	flag.Parse()
	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up logging
	setupLogging(config.Log.LogLevel, config.Log.LogFile)

	// Start tailing the log file
	tailLogFile(config)
}
