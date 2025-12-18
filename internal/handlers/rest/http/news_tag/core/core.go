package core

import (
	"errors"
	"strings"
	"unicode"
)

func hasSpecialChar(s string) bool {
	for _, c := range s {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			return true
		}
	}
	return false
}

func ValidateString(s []string) error {
	for _, tag := range s {
		tag = strings.TrimSpace(tag)
		if tag == "" {
			return errors.New("string is empty")
		}
		if hasSpecialChar(tag) {
			return errors.New("string contains special characters")
		}
	}
	return nil
}
