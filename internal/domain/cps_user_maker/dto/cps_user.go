package dto

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved"`
	Reason  *string `json:"reason,omitempty"`
}

type CreateUserRequest struct {
	UserCode           string   `json:"user_code"`
	UserName           string   `json:"user_name"`
	FullName           string   `json:"full_name"`
	PhoneNumber        string   `json:"phone_number"`
	UserRole           string   `json:"user_role"`
	Department         string   `json:"department"`
	PermissionCategory []string `json:"permission_category"`
	PermissionGroups   []string `json:"permission_groups"`
}

type UpdateUserRequest struct {
	UserName           *string          `json:"user_name"`
	FullName           *string          `json:"full_name"`
	PhoneNumber        *string          `json:"phone_number"`
	UserRole           *string          `json:"user_role"`
	Department         *bson.ObjectID   `json:"department"`
	PermissionCategory *[]bson.ObjectID `json:"permission_category"`
	PermissionGroups   *[]bson.ObjectID `json:"permission_groups"`
}

type ApproveCPSAction struct {
	Approved bool    `json:"approved"`
	Reason   *string `json:"reason"`
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

func (r *CreateUserRequest) Normalize() {
	r.UserCode = strings.TrimSpace(r.UserCode)
	r.UserName = strings.TrimSpace(r.UserName)
	r.FullName = strings.TrimSpace(r.FullName)
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
	r.UserRole = strings.TrimSpace(r.UserRole)
}

func IsObjectIDRequired(value interface{}) error {
	_, ok := value.(bson.ObjectID)
	if !ok {
		return validation.NewError("validation_is_objectid_required", "must be a valid ObjectID")
	}
	return nil
}

// helper to check if slice of ObjectIDs is not empty
func IsObjectIDSliceRequired(value interface{}) error {
	_, ok := value.([]bson.ObjectID)
	if !ok {
		return validation.NewError("validation_is_objectid_slice_required", "must be a non-empty list of ObjectIDs")
	}
	return nil
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserName, validation.Required.Error("user_name is required")),
		validation.Field(&r.FullName, validation.Required.Error("full_name is required")),
		validation.Field(&r.PhoneNumber, validation.Required.Error("phone_number is required")),
		validation.Field(&r.UserRole, validation.Required.Error("user_role is required")),
		validation.Field(&r.Department, validation.By(IsObjectIDRequired)),
		validation.Field(&r.PermissionCategory, validation.By(IsObjectIDSliceRequired)),
		validation.Field(&r.PermissionGroups, validation.By(IsObjectIDSliceRequired)),
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

func (r ApproveCPSAction) Validate() error {
	if r.Approved && r.Reason == nil {
		empty := ""
		r.Reason = &empty
	}

	if !r.Approved {
		if r.Reason == nil || strings.TrimSpace(*r.Reason) == "" {
			return fmt.Errorf("reason is required")
		}
	}
	return nil
}
