package budget

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var hexColorRegex = regexp.MustCompile(`^#?([a-fA-F\d]{2}){3}$`)

func (b BudgetCreateColor) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Color,
			validation.Required.Error("color is required"),
			validation.Length(7, 7).Error("color must be 7 characters long"),
			validation.Match(hexColorRegex).Error("invalid color format, please enter a valid hex color"),
		),
	)
}

func (r UpdateColorRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Color,
			validation.Required.Error("color is required"),
			validation.Length(7, 7).Error("color must be 7 characters"),
			validation.Match(hexColorRegex).Error("invalid hex color format, please enter a valid hex color"),
		),
	)
}

