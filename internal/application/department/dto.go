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
		validation.Field(&i.Department, validation.Required.Error("Department name is required")),
		validation.Field(&i.PortalCards, validation.Each(validation.Required.Error("Each portal card must be non-empty"))),
		validation.Field(&i.PermissionGroups, validation.Each(validation.Required.Error("Each permission group must be non-empty"))),
	)
}

type UpdateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}

func (i UpdateDepartmentRequest) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Department, validation.When(i.Department != "", validation.Required.Error("Department name must be non-empty when provided"))),
		validation.Field(&i.PortalCards, validation.When(i.PortalCards != nil, validation.Each(validation.Required.Error("Each portal card must be non-empty when provided")))),
		validation.Field(&i.PermissionGroups, validation.When(i.PermissionGroups != nil, validation.Each(validation.Required.Error("Each permission group must be non-empty when provided")))),
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
		validation.Field(&i.MakerID, validation.Required.Error("Maker ID is required")),
		validation.Field(&i.MakerName, validation.Required.Error("Maker name is required")),
		validation.Field(&i.MakerPhoneNumber, validation.Required.Error("Maker phone number is required")),
		validation.Field(&i.Department, validation.Required.Error("Department is required")),
		validation.Field(&i.DepartmentCode, validation.Required.Error("Department code is required")),
		validation.Field(&i.DepartmentName, validation.Required.Error("Department name is required")),
		validation.Field(&i.PortalCards, validation.Each(validation.Required.Error("Each portal card must be non-empty"))),
	)
}

type CPSActionResponse struct {
	ActionCode string `json:"action_code"`
	Message    string `json:"message"`
}
