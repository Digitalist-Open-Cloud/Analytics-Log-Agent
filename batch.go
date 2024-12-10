package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const batchSize = 200

var logBuffer []url.Values
var bufferMutex sync.Mutex

func sendBatch(config *Config) {
	// Check if there's anything to send
	if len(logBuffer) == 0 {
		return
	}

	batchRequests := make([]string, len(logBuffer))
	for i, log := range logBuffer {
		// Use url.Values.Encode() to format the query string correctly
		encodedLog := log.Encode() // Automatically handles URL encoding and concatenates with '&'

		// Add the properly formatted request to the batchRequests array
		batchRequests[i] = "?" + encodedLog // Prepend '?' for Matomo-compatible query strings
	}

	// Create the final payload as a map
	payload := map[string]interface{}{
		"requests":   batchRequests,
		"token_auth": config.Matomo.TokenAuth,
	}

	// Marshal the payload into JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Errorf("Error marshalling JSON payload: %v", err)
		return
	}

	payloadString := string(jsonData)

	// Decode the JSON payload string to make it more readable and replace \\u0026 with &
	payloadString = strings.ReplaceAll(payloadString, "\\u0026", "&")
	payloadString = strings.ReplaceAll(payloadString, "\\u003D", "=")

	// Generate the curl command
	curlCommand := fmt.Sprintf(
		"curl -i -X POST -d '%s' %s",
		payloadString,            // the JSON payload
		config.Matomo.TrackerURL, // the Matomo endpoint
	)

	logger.Infof("Curl command to send this request: %s", curlCommand) // Log the curl command

	// Log the final JSON payload for debugging
	logger.Infof("Sending batch request with %d logs: %s", len(logBuffer), string(jsonData))

	// Send the JSON payload to Matomo
	targetURL := config.Matomo.TrackerURL
	req, err := http.NewRequest("POST", targetURL+"matomo.php", bytes.NewBuffer(jsonData))

	if err != nil {
		logger.Errorf("Error creating HTTP request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Log-Agent/1.0")

	// Execute the HTTP request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Error sending batch to Matomo: %v", err)
		return
	}
	defer resp.Body.Close()

	logger.Infof("Batch sent: %d logs, Status: %s", len(logBuffer), resp.Status)

	// Clear the log buffer after sending
	logBuffer = nil
}

func addLogToBatch(log url.Values, config *Config) {
	logger.Infof("Log added to batch")

	// Locking the buffer for safe access in concurrent environments
	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	logBuffer = append(logBuffer, log)

	// Check if the batch size is reached
	if len(logBuffer) >= batchSize {
		logger.Infof("Batch length %d", len(logBuffer))
		sendBatch(config)
	}
}

func flushBatch(config *Config) {
	sendBatch(config) // Send any remaining logs
}
