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

func shouldSendURL(url string, excludedURLs []string) bool {
	for _, excluded := range excludedURLs {
		if strings.Contains(url, excluded) {
			logger.Debugf("Skipping URL %s, contains excluded substring: %s", url, excluded)
			return false // Don't send the URL
		}
	}
	return true // Send the URL
}
