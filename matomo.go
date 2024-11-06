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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func validateTokenAuth(config *Config) error {

	data := url.Values{
		"module":     {"API"},
		"method":     {"API.getMatomoVersion"},
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

	// Define a struct to unmarshal the JSON response
	type VersionResponse struct {
		Value string `json:"value"`
	}

	var versionResp VersionResponse
	err = json.Unmarshal(body, &versionResp)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Log the extracted version
	logger.Infof("Auth token ok, Matomo version is %s", versionResp.Value)

	return nil
}

func InitializeAgentURL(config *Config) {
	// Ensure the Matomo URL ends with a '/', if not, add it.
	if !strings.HasSuffix(config.Matomo.URL, "/") {
		config.Matomo.URL += "/"
	}

	config.Matomo.AgentURL = config.Matomo.URL + "index.php?module=API&method=Agent.postLogData"
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.Contains(item, s) {
			return true
		}
	}
	return false
}

// Matomo Tracking API call
func sendToMatomo(logData *LogData, config *Config) {

	if len(logData.Host) > 0 {
		logData.URL = "https://" + logData.Host + logData.URL
	} else {
		logData.URL = config.Matomo.WebSite + logData.URL
	}
	//logData.URL = config.Matomo.WebSite + logData.URL

	var targetURL string
	InitializeAgentURL(config)

	if len(config.Log.UserAgents) > 0 && !contains(config.Log.UserAgents, logData.UserAgent) {
		logger.Debugf("User agent '%s' not tracked. Skipping log.", logData.UserAgent)
		return
	}

	// Check if the request URL contains an ignored media file extension
	if isIgnored(logData.URL) {
		logger.Debugf("Skipping media file request: %s", logData.URL)
		return
	}

	var isDownload bool
	if config.Matomo.Downloads && isDownloadableFile(logData.URL) {
		isDownload = true
		logger.Debugf("Downloadable file detected: %s", logData.URL)
	}

	formattedTime, err := formatTimestamp(logData.Timestamp)
	if err != nil {
		logger.Warnf("Failed to format timestamp: %v", err)
		return
	}

	data := url.Values{
		"idsite":      {config.Matomo.SiteID},
		"rec":         {"1"},
		"send_image":  {"0"},
		"cip":         {logData.IP},
		"ua":          {logData.UserAgent},
		"url":         {logData.URL},
		"urlref":      {logData.Referrer},
		"token_auth":  {config.Matomo.TokenAuth},
		"status_code": {logData.Status},
		"cdt":         {formattedTime},
	}

	if isDownload {
		data.Set("download", logData.URL)
	}

	errorStatuses := map[string]bool{
		"400": true,
		"401": true,
		"402": true,
		"403": true,
		"404": true,
		"405": true,
		"406": true,
		"407": true,
		"408": true,
		"409": true,
		"410": true,
		"411": true,
		"412": true,
		"413": true,
		"414": true,
		"415": true,
		"416": true,
		"417": true,
		"418": true,
		"421": true,
		"425": true,
		"426": true,
		"428": true,
		"429": true,
		"431": true,
		"451": true,
		"500": true,
		"501": true,
		"502": true,
		"503": true,
		"504": true,
		"505": true,
		"506": true,
		"510": true,
		"511": true,
	}

	// Code that is only executed if you have set plugin = true in config.
	if config.Matomo.Plugin {
		if errorStatuses[logData.Status] {
			targetURL = config.Matomo.AgentURL
			resp, err := http.PostForm(targetURL, data)
			if err != nil {
				logger.Error("Error sending data to Matomo:", err)
				return
			} else {
				var Site string
				if len(logData.Host) > 0 {
					Site = logData.Host
				} else {
					Site = config.Matomo.WebSite
				}
				logger.Debugf("Error log sent for host %s site %s: %s, Status: %s", Site, config.Matomo.SiteID, logData.URL, resp.Status)
			}
			defer resp.Body.Close()
		}
	}
	// Ensure the Matomo URL ends with a '/', if not, add it.
	if !strings.HasSuffix(config.Matomo.URL, "/") {
		config.Matomo.URL += "/"
	}
	targetURL = config.Matomo.URL

	// Post to Tracker API.
	resp, err := http.PostForm(targetURL+"matomo.php", data)
	if err != nil {
		logger.Error("Error sending data to Matomo:", err)
		return
	} else {
		var Site string
		if len(logData.Host) > 0 {
			Site = logData.Host
		} else {
			Site = config.Matomo.WebSite
		}
		logger.Debugf("Log sent host %s and site %s: %s, Status: %s", Site, config.Matomo.SiteID, logData.URL, resp.Status)

	}
	defer resp.Body.Close()

}
