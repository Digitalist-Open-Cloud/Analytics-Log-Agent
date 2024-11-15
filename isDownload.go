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

// List of downloadable file extensions
var downloadExtensions = []string{
	".7z", ".aac", ".apk", ".arc", ".arj", ".asf", ".asx", ".avi", ".azw3", ".bin",
	".bz2", ".csv", ".deb", ".dmg", ".doc", ".docx", ".epub", ".exe", ".flac", ".flv",
	".gz", ".gzip", ".hqx", ".ibooks", ".jar", ".json", ".md5", ".mov", ".movie",
	".mp2", ".mp3", ".mp4", ".mpg", ".mpeg", ".mobi", ".msi", ".msp", ".odb", ".odf",
	".odg", ".odp", ".ods", ".odt", ".ogg", ".ogv", ".pdf", ".phps", ".ppt", ".pptx",
	".qt", ".qtm", ".ra", ".ram", ".rar", ".rpm", ".rtf", ".sea", ".sig", ".sit",
	".tar", ".tbz", ".tgz", ".torrent", ".txt", ".wav", ".webm", ".wma", ".wmv",
	".wpd", ".xls", ".xlsx", ".xml", ".xsd", ".z", ".zip",
}

// Helper function to check if a URL is for a downloadable file
func isDownloadableFile(url string) bool {
	// Split the URL at the '?' to remove query strings
	urlWithoutQuery := strings.Split(url, "?")[0]

	// Check if the URL contains a downloadable file extension
	for _, ext := range downloadExtensions {
		if strings.HasSuffix(strings.ToLower(urlWithoutQuery), ext) {
			return true
		}
	}
	return false
}
