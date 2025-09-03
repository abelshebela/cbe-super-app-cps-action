package avatar

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (dto AvatarDTO) IsEmpty() bool {
	return dto.Label == "" && dto.Avatar == nil
}

func (a AvatarDTO) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Label,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)).Error("Label must not contain special characters"),
		),
	)
}
