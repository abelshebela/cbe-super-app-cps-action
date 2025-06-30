package validator

import (
	"strconv"
	"strings"
)

func PinStrengthValidator(input int) []string {
	var errs []string

	pin := strconv.Itoa(input)

	if strings.TrimSpace(pin) == "0" {
		errs = append(errs, "Pin is required")
	} else if len(pin) < 6 || len(pin) > 6 {
		errs = append(errs, "Pin 2 must be only 6 characters")
	}

	return errs
}
