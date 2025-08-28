package passwordrule

import (
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
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.By(noSpecialChars),
		),
		validation.Field(&r.MinLength,
			validation.Min(1),
		),
		validation.Field(&r.MaxLength,
			validation.Min(1),
			validation.By(func(value interface{}) error {
				if r.MinLength > r.MaxLength {
					return validation.NewError("validation_max_length", "max_length must be greater than or equal to min_length")
				}
				return nil
			}),
		),
	)
}
