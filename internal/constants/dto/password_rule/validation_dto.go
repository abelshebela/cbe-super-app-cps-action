package passwordrule

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var allowedChars = "a-zA-Z0-9\\s._-"

func noSpecialChars(value any) error {
	str, ok := value.(string)
	if !ok {
		return validation.NewError("validation", "invalid type")
	}
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}

	re := regexp.MustCompile("^[" + allowedChars + "]+$")
	if !re.MatchString(str) {
		return validation.NewError("validation", "contains invalid characters")
	}
	return nil
}

func (r PasswordRuleUpdate) Validate() error {
	// Validate Name
	if err := noSpecialChars(r.Rule.Name); err != nil {
		return errors.New("NO special character allowed for name")
	}

	// Validate MinLength
	if r.Rule.MinLength < 1 {
		return errors.New("minimum length must be at least 1")
	}

	// Validate MaxLength
	if r.Rule.MaxLength < 1 {
		return errors.New("maximum length must be at least 1")
	}
	if r.Rule.MaxLength > 20 {
		return errors.New("maximum length must be at most 20")
	}
	if r.Rule.MinLength >= r.Rule.MaxLength {
		return errors.New("maximum length must be greater than minimum length")
	}

	return nil
}
