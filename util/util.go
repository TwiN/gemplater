package util

import (
	"fmt"
	"strings"
)

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

func ExtractVariablesFromString(s, wrapper string) (variableNames []string, err error) {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	// Instead of doing it all at once, we'll do it line by line to reduce the odds of picking up a multiline variable
	for _, line := range lines {
		variablesInLine := GetAllBetween(line, wrapper, wrapper)
		for _, variable := range variablesInLine {
			if strings.Contains(variable, " ") {
				continue
			}
			variableNames = append(variableNames, variable)
		}
	}
	return
}

func GenerateCuteHeader(s string) string {
	bread := strings.Repeat("#", len(s)+4)
	return fmt.Sprintf("%s\n# %s #\n%s\n", bread, s, bread)
}
