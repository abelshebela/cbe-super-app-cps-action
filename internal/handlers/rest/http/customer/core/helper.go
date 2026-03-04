package core

import (
	"errors"
	"regexp"
)

// ValidateString checks if the string is non-empty and does not contain special characters.
func ValidateString(s string) error {
	if s == "" {
		return errors.New("string is empty")
	}
	// Regex to match any non-alphanumeric character
	specialCharRegex := regexp.MustCompile(`[^a-zA-Z0-9 ]`)
	if specialCharRegex.MatchString(s) {
		return errors.New("string contains special characters")
	}
	return nil
}

func ValidateCustomerLookupRequest(number string) error {
	if number == "" {
		return errors.New("number is empty")
	}
	// Regex to ensure the number contains only digits
	numberRegex := regexp.MustCompile(`^\d+$`)
	if !numberRegex.MatchString(number) {
		return errors.New("number must contain only digits")
	}
	return nil
}
