package cpsuser

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
)

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

	r.UserName = strings.ToLower(r.UserName)
	r.FullName = strings.ToUpper(r.FullName)
}

func (r *UpdateUserRequest) Normalize() {
	if r.UserName != "" {
		r.UserName = strings.ToLower(r.UserName)
	}
	if r.FullName != "" {
		r.FullName = strings.ToUpper(r.FullName)
	}
}

func IsObjectIDRequired(value interface{}) error {
	if value == nil {
		return validation.NewError("validation_is_objectid_required", "must be a valid ObjectID")
	}
	switch v := value.(type) {
	case bson.ObjectID:
		if v.Hex() == "" {
			return validation.NewError("validation_is_objectid_required", "must be a valid ObjectID")
		}
		return nil
	case *bson.ObjectID:
		if v == nil || v.Hex() == "" {
			return validation.NewError("validation_is_objectid_required", "must be a valid ObjectID")
		}
		return nil
	default:
		return validation.NewError("validation_is_objectid_required", "must be a valid ObjectID")
	}
}

func IsStringSliceRequired(value interface{}) error {
	slice, ok := value.([]string)
	if !ok {
		return validation.NewError("validation_string_slice_required", "must be a non-empty list of strings")
	}
	if len(slice) > 0 {
		for i, s := range slice {
			trimmed := strings.TrimSpace(s)
			if trimmed == "" {
				return validation.NewError("validation_string_slice_required",
					fmt.Sprintf("item at index %d cannot be empty", i))
			}

			if err := utils.NoSpecialChars(trimmed); err != nil {
				return validation.NewError("validation_string_slice_required",
					fmt.Sprintf("item at index %d is invalid: %s", i, err.Error()))
			}
		}
	}

	return nil
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UserName,
			validation.Required.Error("username is required"),
			validation.Length(4, 0).Error("username length must be greater that 3 character"),
			validation.By(utils.NoSpecialChars),
			validation.By(func(value interface{}) error {
				if s, ok := value.(string); ok {
					if strings.Contains(s, " ") {
						return validation.NewError("validation_no_spaces", "username cannot contain spaces")
					}
				}
				return nil
			}),
		),
		validation.Field(&r.FullName,
			validation.Required.Error("full_name is required"),
			validation.By(utils.NoSpecialChars),
			validation.By(utils.ThreeNamesMinLength),
		),
		validation.Field(&r.PhoneNumber,
			validation.Required.Error("phone_number is required"),
			validation.By(func(value interface{}) error {
				if s, ok := value.(string); ok {
					normalized := utils.FormatPhoneNumber(s)
					if normalized == "" {
						return validation.NewError("validation_phone_format", "unsupported phone number format")
					}
				}
				return nil
			}),
		),
		validation.Field(&r.Gender, validation.Required.Error("gender is required")),
		validation.Field(&r.JobTitle, validation.Required.Error("job_title is required")),
		validation.Field(
			&r.Email,
			validation.Required.Error("email is required"),
			validation.Match(
				regexp.MustCompile(`^[A-Za-z0-9._%+-]+@cbe\.com\.et$`),
			).Error("email must be a valid cbe.com.et email"),
		),
	)
}

func (r UpdateUserRequest) Validate() error {
	if r.UserName == "" && r.FullName == "" && r.PhoneNumber == "" {
		return fmt.Errorf("at least one field must be provided for update")
	}

	return validation.ValidateStruct(&r,
		validation.Field(&r.UserName, validation.When(r.UserName != "",
			validation.Length(4, 0).Error("username length must be greater that 3 character"),
			validation.By(utils.NoSpecialChars),
			validation.By(func(value interface{}) error {
				if s, ok := value.(string); ok {
					if strings.Contains(s, " ") {
						return validation.NewError("validation_no_spaces", "username cannot contain spaces")
					}
				}
				return nil
			}),
		)),
		validation.Field(&r.FullName, validation.When(r.FullName != "",
			validation.Length(1, 100).Error("full_name cannot be empty"),
			validation.By(utils.NoSpecialChars),
		)),
		validation.Field(&r.PhoneNumber, validation.When(r.PhoneNumber != "",
			validation.Length(1, 20).Error("phone_number cannot be empty"),
			validation.By(func(value interface{}) error {
				if s, ok := value.(string); ok {
					normalized := utils.FormatPhoneNumber(s)
					if normalized == "" {
						return validation.NewError("validation_phone_format", "unsupported phone number format")
					}
				}
				return nil
			}),
		)),
		validation.Field(&r.Gender, validation.When(r.Gender != "",
			validation.By(utils.NoSpecialChars),
		)),
		validation.Field(&r.JobTitle, validation.When(r.JobTitle != "",
			validation.By(utils.NoSpecialChars),
		)),
		validation.Field(
			&r.Email,
			validation.When(r.Email != "",
				validation.Match(
					regexp.MustCompile(`^[A-Za-z0-9._%+-]+@cbe\.com\.et$`),
				).Error("email must be a valid cbe.com.et email"),
			),
		),
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
