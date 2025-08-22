package department

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateDepartmentRequest struct {
	Department  string   `json:"department"`
	PortalCards []string `json:"portal_cards"`
}

func (i CreateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(
			&i.Department,
			validation.Required.Error("Department name is required"),
		),
		validation.Field(
			&i.PortalCards,
			validation.Required.Error("Portal cards are required"),
			validation.Each(validation.Required.Error("Each portal card must be non-empty")),
		),
	)
}

type UpdateDepartmentRequest struct {
	Department  string   `json:"department"`
	PortalCards []string `json:"portal_cards"`
}

func (i UpdateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Department, validation.When(i.Department != "", validation.Required.Error("Department name must be non-empty when provided"))),
		validation.Field(&i.PortalCards, validation.When(i.PortalCards != nil, validation.Each(validation.Required.Error("Each portal card must be non-empty when provided")))),
	)
}

type DepartmentUpdateCPSActionRequest struct {
	Department  string   `json:"department"`
	PortalCards []string `json:"portal_cards"`
}

func (req DepartmentUpdateCPSActionRequest) Validate() error {
	var rules []error
	if req.Department != "" {
		if err := validation.Validate(req.Department, validation.Required); err != nil {
			rules = append(rules, err)
		}
	}
	if req.PortalCards != nil {
		if err := validation.Validate(req.PortalCards, validation.Each(validation.Required)); err != nil {
			rules = append(rules, err)
		}
	}
	if len(rules) > 0 {
		return fmt.Errorf("%v", rules)
	}
	return nil
}

type CPSActionResponse struct {
	ActionCode string `json:"action_code"`
	Message    string `json:"message"`
}
