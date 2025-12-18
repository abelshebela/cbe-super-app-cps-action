package roles

import validation "github.com/go-ozzo/ozzo-validation/v4"

func (r *RequestRolesCreate) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.JobTitle, validation.Required.Error("Job title is required")),
		validation.Field(&r.Role, validation.Required.Error("Role is required")),
	)
}

func (r *RequestRolesUpdate) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.JobTitle, validation.Required.Error("Job title is required")),
		validation.Field(&r.Role, validation.Required.Error("Role is required")),
	)
}
