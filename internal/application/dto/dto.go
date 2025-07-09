package dto

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateDepartmentRequest struct {
	Department  string   `json:"department"`
	PortalCards []string `json:"portal_cards"`
}

func (i CreateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Department, validation.Required),
		validation.Field(&i.PortalCards, validation.Required, validation.Each(validation.Required)),
	)
}
