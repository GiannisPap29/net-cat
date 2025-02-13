package utils

import (
	"fmt"
	"strings"
)

// ValidateName checks if the provided name is valid
func ValidateName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' ||
			char == '-') {
			return false
		}
	}
	return true
}

// NormalizeName converts a name to lowercase for consistent comparison
func NormalizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// FormatMessage formats a chat message with timestamp
func FormatMessage(timestamp, name, message string) string {
	return fmt.Sprintf("[%s][%s]: %s\n", timestamp, name, message)
}

// to klassiko Atoi apo ti pisina ! ! ! ! !
func Atoi(s string) int {
	sign := 1
	var result int
	if len(s) > 0 {
		if s[0] == '-' {
			sign = -1
		} else if s[0] == '+' {
			sign = 1
		} else if s[0] >= '0' && s[0] <= '9' {
			result = result*10 + int(s[0]-'0')
		}
	} else {
		return 0
	}
	if len(s) > 1 {
		for i := 1; i < len(s); i++ {
			if s[i] >= '0' && s[i] <= '9' {
				result = result*10 + int(s[i]-'0')
			} else {
				return 0
			}
		}
	}
	result *= sign
	return result
}
