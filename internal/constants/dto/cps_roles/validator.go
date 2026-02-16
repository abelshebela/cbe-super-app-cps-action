package cpsroles

import (
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (r CreateCPSRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.RoleCode,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.Description,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
	)
}

func (r UpdateCPSRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.RoleCode,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
		validation.Field(&r.Description,
			validation.Required,
			validation.By(utils.NoSpecialChars),
			validation.By(utils.TrimWhiteSpace),
		),
	)
}

func (r ToggleServiceAccessRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AccessListKeys,
			validation.Required,
			validation.Length(1, 0),
		),
	)
}
