package common

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
)

func ResponseMaker(res map[string]interface{}, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func ValidateInputNoSpecialChars(input string) error {
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9 ]*$`)

	if !validPattern.MatchString(input) {
		return errors.New("input contains special characters")
	}

	return nil
}
