package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApproveUserActionRequest struct {
	Approve bool    `json:"approved"`
	Reason  *string `json:"reason,omitempty"`
}

type CreateUserRequest struct {
	UserName           string          `json:"username" bson:"username,omitempty"`
	FullName           string          `json:"full_name" bson:"full_name,omitempty"`
	Department         bson.ObjectID   `json:"department" bson:"department,omitempty"`
	PhoneNumber        string          `json:"phone_number" bson:"phone_number,omitempty"`
	Role               string          `json:"role" bson:"role,omitempty"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category" bson:"permission_category,omitempty"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups" bson:"permission_groups,omitempty"`
}

type UpdateUserRequest struct {
	UserCode           string          `json:"user_code,omitempty" bson:"user_code,omitempty"`
	UserName           string          `json:"username,omitempty" bson:"username,omitempty"`
	FullName           string          `json:"full_name,omitempty" bson:"full_name,omitempty"`
	Department         bson.ObjectID   `json:"department,omitempty" bson:"department,omitempty"`
	PhoneNumber        string          `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Role               string          `json:"role,omitempty" bson:"role,omitempty"`
	Gender             string          `json:"gender,omitempty" bson:"gender,omitempty"`
	Email              string          `json:"email,omitempty" bson:"email,omitempty"`
	PermissionCategory []bson.ObjectID `json:"permission_category,omitempty" bson:"permission_category,omitempty"`
	PermissionGroups   []bson.ObjectID `json:"permission_groups,omitempty" bson:"permission_groups,omitempty"`
}

type CPSUserDTO struct {
	ID                 bson.ObjectID   `json:"id"`
	UserCode           string          `json:"user_code"`
	FullName           string          `json:"full_name"`
	Role               string          `json:"role"`
	Department         bson.ObjectID   `json:"department"`
	Gender             string          `json:"gender"`
	PhoneNumber        string          `json:"phone_number"`
	Email              string          `json:"email"`
	UserName           string          `json:"username"`
	Realm              string          `json:"realm"`
	PermissionCategory []bson.ObjectID `json:"permission_category"`
	PermissionGroup    []bson.ObjectID `json:"permission_group"`

	Enabled      bool       `json:"enabled"`
	DateJoined   *time.Time `json:"date_joined"`
	LastModified *time.Time `json:"last_modified"`

	Country string `json:"country"`
	Region  string `json:"region"`
}

func NewCPSUserDTO(user model.CPSUser) CPSUserDTO {
	return CPSUserDTO{
		ID:                 user.ID,
		UserCode:           user.UserCode,
		FullName:           user.FullName,
		Role:               user.Role,
		Department:         user.Department,
		Gender:             user.Gender,
		PhoneNumber:        user.PhoneNumber,
		Email:              user.Email,
		UserName:           user.UserName,
		Realm:              user.Realm,
		PermissionCategory: user.PermissionCategory,
		PermissionGroup:    user.PermissionGroup,
		Enabled:            user.Enabled,
		DateJoined:         user.DateJoined,
		LastModified:       user.LastModified,
		Country:            user.Country,
		Region:             user.Region,
	}
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
	r.UserName = strings.TrimSpace(r.UserName)
	r.FullName = strings.TrimSpace(r.FullName)
	r.PhoneNumber = strings.TrimSpace(r.PhoneNumber)
	r.Role = strings.TrimSpace(r.Role)
	r.Gender = strings.TrimSpace(r.Gender)
	r.Email = strings.TrimSpace(r.Email)
}

func IsObjectIDRequired(value interface{}) error {
	_, ok := value.(bson.ObjectID)
	if !ok {
		return validation.NewError("validation_is_objectid_required", "must be a valid ObjectID")
	}
	return nil
}

func IsRequired(fieldName string) validation.RuleFunc {
	return func(value interface{}) error {
		switch v := value.(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				return validation.NewError("validation_required", "is required")
			}
		case bson.ObjectID:
			if v.IsZero() {
				return validation.NewError("validation_required", "is required")
			}
		case *bson.ObjectID:
			if v == nil || v.IsZero() {
				return validation.NewError("validation_required", "is required")
			}
		case []bson.ObjectID:
			if len(v) == 0 {
				return validation.NewError("validation_required", "is required")
			}
			for _, id := range v {
				if id.IsZero() {
					return validation.NewError("validation_required", fmt.Sprintf("%s cannot contain empty ObjectIDs", fieldName))
				}
			}
		case []interface{}:
			if len(v) == 0 {
				return validation.NewError("validation_required", "is required")
			}
		default:
			if value == nil {
				return validation.NewError("validation_required", "is required")
			}
		}
		return nil
	}
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
		validation.Field(&r.UserName, validation.Required.Error("username is required"), validation.By(utils.NoSpecialChars)),
		validation.Field(&r.FullName, validation.Required.Error("full_name is required"), validation.By(utils.NoSpecialChars)),
		validation.Field(&r.PhoneNumber, validation.Required.Error("phone_number is required"), validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Role, validation.Required.Error("user role is required"), validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Department, validation.By(IsRequired("department")), validation.By(IsObjectIDRequired)),
		validation.Field(&r.PermissionCategory, validation.By(IsRequired("permission_category")), validation.By(IsObjectIDSliceRequired)),
		validation.Field(&r.PermissionGroups, validation.By(IsRequired("permission_group")), validation.By(IsObjectIDSliceRequired)),
		validation.Field(&r.Gender, validation.Required.Error("gender is required"), validation.By(utils.NoSpecialChars)),
		validation.Field(&r.Email, validation.Required.Error("email is required")),
	)
}

func (r UpdateUserRequest) Validate() error {
	if r.UserName == "" && r.FullName == "" && r.PhoneNumber == "" &&
		r.Role == "" &&
		r.PermissionCategory == nil && r.PermissionGroups == nil {
		return fmt.Errorf("at least one field must be provided for update")
	}

	return validation.ValidateStruct(&r,
		validation.Field(&r.UserName, validation.When(r.UserName != "", validation.Length(1, 100).Error("user_name cannot be empty"))),
		validation.Field(&r.FullName, validation.When(r.FullName != "", validation.Length(1, 100).Error("full_name cannot be empty"))),
		validation.Field(&r.PhoneNumber, validation.When(r.PhoneNumber != "", validation.Length(1, 20).Error("phone_number cannot be empty"))),
		validation.Field(&r.Role, validation.When(r.Role != "", validation.Length(1, 50).Error("user_role cannot be empty"))),
		validation.Field(&r.Department, validation.By(IsObjectIDRequired)),
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
