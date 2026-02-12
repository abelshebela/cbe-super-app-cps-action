package roles

import (
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *RequestRolesCreate) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.JobTitle, validation.Required.Error("Job title is required")),
		validation.Field(&r.Role, validation.Required.Error("Role is required")),
	)
}

func (r *RequestRolesUpdate) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.JobTitle,
			validation.When(r.JobTitle != "", validation.By(utils.NoSpecialChars)),
		),
		validation.Field(&r.Role,
			validation.When(r.Role != "", validation.By(utils.NoSpecialChars)),
		),
	)
}
