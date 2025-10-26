package utils

import (
	"fmt"
	"strings"
)

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// ContainsAny checks if a string contains any of the substrings
func ContainsAny(s string, substrings []string) bool {
	for _, substr := range substrings {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}

// FormatResourceName formats a resource name with namespace
func FormatResourceName(resourceType, name, namespace string) string {
	if namespace != "" && namespace != "default" {
		return fmt.Sprintf("%s/%s (namespace: %s)", resourceType, name, namespace)
	}
	return fmt.Sprintf("%s/%s", resourceType, name)
}

// RedactSecrets redacts sensitive information from strings
func RedactSecrets(input string) string {
	// Simple redaction - replace common secret patterns
	patterns := []string{
		"password", "token", "secret", "key", "apikey",
	}
	
	result := input
	for _, pattern := range patterns {
		if strings.Contains(strings.ToLower(input), pattern) {
			// Redact the value after the pattern
			result = strings.ReplaceAll(result, pattern, "[REDACTED]")
		}
	}
	
	return result
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Pluralize returns the plural form of a word if count is not 1
func Pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}
