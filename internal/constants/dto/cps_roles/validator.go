package cpsroles

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (r *CreateCPSRoleRequest) Normalize() {
	r.Name = strings.TrimSpace(r.Name)
	r.RoleCode = strings.TrimSpace(r.RoleCode)
	r.Description = strings.TrimSpace(r.Description)
}
func (r *UpdateCPSRoleRequest) Normalize() {
	if r.Name != nil {
		*r.Name = strings.TrimSpace(*r.Name)
	}
	if r.RoleCode != nil {
		*r.RoleCode = strings.TrimSpace(*r.RoleCode)
	}
	if r.Description != nil {
		*r.Description = strings.TrimSpace(*r.Description)
	}
}

func (r CreateCPSRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.RoleCode,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.Description,
			validation.Required,
			validation.By(utils.NoSpecialChars),
		),
	)
}

func (r UpdateCPSRoleRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.RoleCode,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&r.Description,
			validation.By(utils.NoSpecialChars),
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
