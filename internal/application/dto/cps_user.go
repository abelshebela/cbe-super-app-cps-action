package dto

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved"`
	Reason  *string `json:"reason,omitempty"`
}

type CreateUserRequest struct {
	UserName           string   `json:"user_name"`
	FullName           string   `json:"full_name"`
	PhoneNumber        string   `json:"phone_number"`
	UserRole           string   `json:"user_role"`
	Department         string   `json:"department"`
	PermissionCategory []string `json:"permission_category"`
	PermissionGroups   []string `json:"permission_groups"`
}

type UpdateUserRequest struct {
	UserName           *string   `json:"user_name"`
	FullName           *string   `json:"full_name"`
	PhoneNumber        *string   `json:"phone_number"`
	UserRole           *string   `json:"user_role"`
	Department         *string   `json:"department"`
	PermissionCategory *[]string `json:"permission_category"`
	PermissionGroups   *[]string `json:"permission_groups"`
}

func (r ApproveUserActionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Approve,
			validation.Required.Error("approve is required"),
			validation.In(true, false).Error("approve must be true or false"),
		),
		validation.Field(&r.Reason,
			validation.When(!r.Approve, validation.Required.Error("reason must not be empty")),
		),
	)
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserName, validation.Required.Error("user_name is required")),
		validation.Field(&r.FullName, validation.Required.Error("full_name is required")),
		validation.Field(&r.PhoneNumber, validation.Required.Error("phone_number is required")),
		validation.Field(&r.UserRole, validation.Required.Error("user_role is required")),
		validation.Field(&r.Department, validation.Required.Error("department is required")),
		validation.Field(&r.PermissionCategory, validation.Required.Error("permission_category is required")),
		validation.Field(&r.PermissionGroups, validation.Required.Error("permission_groups is required")),
	)
}

func (r UpdateUserRequest) Validate() error {
	if r.UserName == nil && r.FullName == nil && r.PhoneNumber == nil &&
		r.UserRole == nil && r.Department == nil &&
		r.PermissionCategory == nil && r.PermissionGroups == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	return validation.ValidateStruct(&r,
		validation.Field(&r.UserName,
			validation.When(r.UserName != nil, validation.Length(1, 100).Error("user_name cannot be empty"))),
		validation.Field(&r.FullName,
			validation.When(r.FullName != nil, validation.Length(1, 100).Error("full_name cannot be empty")),
		),
		validation.Field(&r.PhoneNumber,
			validation.When(r.PhoneNumber != nil, validation.Length(1, 20).Error("phone_number cannot be empty")),
		),
		validation.Field(&r.UserRole,
			validation.When(r.UserRole != nil, validation.Length(1, 50).Error("user_role cannot be empty")),
		),
		validation.Field(&r.Department,
			validation.When(r.Department != nil, validation.Length(1, 50).Error("department cannot be empty")),
		),
		validation.Field(&r.PermissionCategory),
		validation.Field(&r.PermissionGroups),
	)
}
