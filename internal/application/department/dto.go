package department

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CreateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}

func (i CreateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Department, validation.Required),
		validation.Field(&i.PortalCards, validation.Required, validation.Each(validation.Required)),
		validation.Field(&i.PermissionGroups, validation.Required, validation.Each(validation.Required)),
	)
}

type UpdateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}

func (i UpdateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Department, validation.Required),
		validation.Field(&i.PortalCards, validation.Required, validation.Each(validation.Required)),
		validation.Field(&i.PermissionGroups, validation.Required, validation.Each(validation.Required)),
	)
}

type DepartmentUpdateCPSActionRequest struct {
	MakerID          string   `json:"maker_id"`
	MakerName        string   `json:"maker_name"`
	MakerPhoneNumber string   `json:"maker_phone_number"`
	Department       string   `json:"department"`
	DepartmentCode   string   `json:"department_code"`
	DepartmentName   string   `json:"department_name"`
	PortalCards      []string `json:"portal_cards"`
}

func (i DepartmentUpdateCPSActionRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.MakerID, validation.Required),
		validation.Field(&i.MakerName, validation.Required),
		validation.Field(&i.MakerPhoneNumber, validation.Required),
		validation.Field(&i.Department, validation.Required),
		validation.Field(&i.DepartmentCode, validation.Required),
		validation.Field(&i.DepartmentName, validation.Required),
		validation.Field(&i.PortalCards, validation.Required, validation.Each(validation.Required)),
	)
}

type CPSActionResponse struct {
	ActionCode string `json:"action_code"`
	Message    string `json:"message"`
}
