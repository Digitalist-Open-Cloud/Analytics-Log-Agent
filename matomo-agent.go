package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tenebris-tech/tail"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Matomo struct {
		URL       string `mapstructure:"url"`
		ErrorURL  string
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

func validateTokenAuth(config *Config) error {

	data := url.Values{
		"module":     {"API"},
		"method":     {"API.getPiwikVersion"},
		"format":     {"JSON"},
		"token_auth": {config.Matomo.TokenAuth},
	}
	validationURL := fmt.Sprintf("%sindex.php", config.Matomo.URL)

	resp, err := http.PostForm(validationURL, data)
	if err != nil {
		return fmt.Errorf("error validating token: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token_auth, received status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	logger.Infof("Matomo version: %s", string(body))

	return nil
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
		logger.SetOutput(&lumberjack.Logger{
			Filename:   logFile, // Log file path
			MaxSize:    10,      // Max size in megabytes before log rotation
			MaxBackups: 3,       // Max number of old log files to keep
			MaxAge:     28,      // Max number of days to retain old log files
			Compress:   true,    // Whether to compress old log files
		})

	} else {
		// Fallback to stderr if no log file is provided
		logger.Warn("No log file provided, using default stderr")
		logger.SetOutput(os.Stderr)
	}
}

// InitializeErrorURL constructs the ErrorURL by combining the base URL with the additional query string.
func InitializeErrorURL(config *Config) {
	// Ensure the base URL ends with a '/' to safely concatenate the query string
	if !strings.HasSuffix(config.Matomo.URL, "/") {
		config.Matomo.URL += "/"
	}

	// Set the ErrorURL by appending the specific endpoint to the base URL
	config.Matomo.ErrorURL = config.Matomo.URL + "index.php?module=API&method=Agent.postLogData"
}

// Matomo Tracking API call
func sendToMatomo(logData *LogData, config *Config) {
	var Url = config.Matomo.WebSite + logData.URL
	var targetURL string
	InitializeErrorURL(config)

	data := url.Values{
		"idsite":      {config.Matomo.SiteID},
		"rec":         {"1"},
		"cip":         {logData.IP},
		"ua":          {logData.UserAgent},
		"url":         {Url},
		"urlref":      {logData.Referrer},
		"token_auth":  {config.Matomo.TokenAuth},
		"status_code": {logData.Status},
	}

	// Common HTTP error statuses you want to handle
	errorStatuses := map[string]bool{
		"404": true,
		"403": true,
		"503": true,
		"500": true,
	}

	// If the status is an error, use the error tracking API endpoint; otherwise, use the regular one
	if errorStatuses[logData.Status] {
		targetURL = config.Matomo.ErrorURL // Define this URL for errors in your config
		resp, err := http.PostForm(targetURL, data)
		if err != nil {
			logger.Error("Error sending data to Matomo:", err)
			return
		} else {
			logger.Infof("Error log sent for %s: %s, Status: %s", config.Matomo.SiteID, logData.URL, resp.Status)
		}
		defer resp.Body.Close()
	} else {
		targetURL = config.Matomo.URL
		resp, err := http.PostForm(targetURL+"matomo.php", data)
		if err != nil {
			logger.Error("Error sending data to Matomo:", err)
			return
		} else {
			logger.Infof("Log sent: %s, Status: %s", logData.URL, resp.Status)
		}
		defer resp.Body.Close()
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
	// Check if we have a valid token for Matomo.
	err = validateTokenAuth(config)
	if err != nil {
		logger.Fatal("Invalid Matomo token_auth:", err)
	}

	// Start tailing the log file
	tailLogFile(config)
}
