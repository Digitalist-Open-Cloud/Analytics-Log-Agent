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
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func getTitleCacheFilePath(config *Config) string {
	if config.Title.Cache != "" {
		return config.Title.Cache
	}
	return "/tmp/matomo_agent-url_title_cache.txt"
}

var titleCache = make(map[string]string)
var cacheMutex = sync.Mutex{}

func loadCache(filePath string) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			titleCache[parts[0]] = parts[1]
		}
	}
	return scanner.Err()
}

func saveCache(filePath, url, title string) error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if URL is already cached in memory to avoid duplicate writes
	if _, exists := titleCache[url]; exists {
		return nil
	}

	// Append new URL and title to cache file
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write to the file and add to in-memory map
	_, err = file.WriteString(url + ":" + title + "\n")
	if err == nil {
		titleCache[url] = title
	}
	return err
}

func collectTitle(url, cacheFilePath string) (string, error) {
	// Check in-memory cache first
	cacheMutex.Lock()
	if title, found := titleCache[url]; found {
		cacheMutex.Unlock()
		return title, nil
	}
	cacheMutex.Unlock()

	// Otherwise, retrieve and cache the title
	title, err := fetchTitleFromURL(url)
	if err != nil {
		return "", err
	}

	// Save to cache (both in-memory and file)
	err = saveCache(cacheFilePath, url, title)
	return title, err
}

func fetchTitleFromURL(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read response body and extract title
	body, err := io.ReadAll(resp.Body) // Using io.ReadAll instead of ioutil.ReadAll
	if err != nil {
		return "", err
	}

	startIdx := strings.Index(string(body), "<title>")
	endIdx := strings.Index(string(body), "</title>")
	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		return "", fmt.Errorf("title not found")
	}

	// Clean the extracted title
	title := string(body[startIdx+7 : endIdx])
	title = strings.ReplaceAll(title, "\n", "")
	title = strings.ReplaceAll(title, "\r", "")
	title = strings.TrimSpace(title)

	return title, nil
}
