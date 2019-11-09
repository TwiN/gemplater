package util

import "strings"

// Get all substrings between two strings
// This variation does not strip the suffix and the prefix from the substring
func GetAllBetween(haystack, prefix, suffix string) (needles []string) {
	for {
		if len(haystack) < len(prefix)+len(suffix) {
			break
		}
		start := strings.Index(haystack, prefix) + len(prefix)
		if start-len(prefix) == -1 {
			break
		}
		end := strings.Index(haystack[start:], suffix) + start
		if end-start == -1 || start >= end {
			break
		}
		needles = append(needles, haystack[start-len(prefix):end+len(suffix)])
		if len(haystack) <= end {
			break
		}
		haystack = haystack[end+len(suffix):]
	}
	return needles
}
