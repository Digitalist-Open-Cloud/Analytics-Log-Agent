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

import "strings"

// List of media file extensions you don't want to track
var ignoredRequests = []string{
	".jpeg", ".jpg", ".woff", ".woff2", ".ttf", ".gif", ".png", ".webp", ".svg", ".ico", ".js", ".css", ".bmp", ".svgz", ".otf", ".eot", ".xml",
}

// Helper function to check if a URL contains an ignored file extension
func isIgnored(url string) bool {
	// Split the URL at the '?' to remove query strings
	urlWithoutQuery := strings.Split(url, "?")[0]

	if strings.Contains(urlWithoutQuery, "robots.txt") || strings.Contains(strings.ToLower(urlWithoutQuery), "autodiscover") {
		return true
	}

	// Check if the URL without query contains an ignored file extension
	for _, ext := range ignoredRequests {
		if strings.HasSuffix(strings.ToLower(urlWithoutQuery), ext) {
			return true
		}
	}
	return false
}
