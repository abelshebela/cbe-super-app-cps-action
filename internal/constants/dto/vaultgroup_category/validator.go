package vaultgroupcategory

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r *CreateVaultGroupCategoryRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	var err error
	if r.Name, err = sanitizeString(r.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if r.CategoryType, err = sanitizeString(r.CategoryType); err != nil {
		return fmt.Errorf("category_type: %w", err)
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required.Error("name is required"),
			validation.Length(3, 100).Error("name must be between 3 and 100 characters"),
			validation.By(noSpecialChars),
		),
		validation.Field(&r.CategoryType,
			validation.Required.Error("category_type is required"),
			validation.In("GROUP", "PERSONAL").Error("category_type must be GROUP or PERSONAL"),
		),
		validation.Field(&r.CoverImage,
			validation.Required.Error("cover_image is required"),
			validation.By(func(value interface{}) error { return validateCoverImage(value) }),
		),
	)
}

func (r *UpdateVaultGroupCategoryRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	if r.Name == "" && r.CoverImage == nil && r.CategoryType == "" {
		return errors.New("at least one field (name, or cover_image) must be provided")
	}

	if r.Name != "" {
		sanitized, err := sanitizeString(r.Name)
		if err != nil {
			return fmt.Errorf("name: %w", err)
		}
		r.Name = sanitized
	}
	if r.CategoryType != "" {
		sanitized, err := sanitizeString(r.CategoryType)
		if err != nil {
			return fmt.Errorf("category_type: %w", err)
		}
		r.CategoryType = sanitized
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.When(r.Name != "",
				validation.Length(3, 100).Error("name must be between 3 and 100 characters"),
				validation.By(noSpecialChars),
			),
		),
		validation.Field(&r.CategoryType,
			validation.When(r.CategoryType != "",
				validation.Required.Error("category_type is required"),
				validation.In("GROUP", "PERSONAL").Error("category_type must be GROUP or PERSONAL"),
			),
		),
		validation.Field(&r.CoverImage,
			validation.When(r.CoverImage != nil,
				validation.By(func(value interface{}) error { return validateCoverImage(value) })),
		),
	)
}

func (r *UpdateVaultGroupCategoryRequest) HasUpdates() bool {
	return r.Name != "" || r.CoverImage != nil || r.CategoryType != ""
}

func sanitizeString(s string) (string, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return "", errors.New("cannot be empty or only whitespace")
	}
	return trimmed, nil
}

func noSpecialChars(value interface{}) error {
	var s string

	switch v := value.(type) {
	case string:
		s = v
	case *string:
		if v == nil {
			return nil
		}
		s = *v
	default:
		return nil
	}

	if s == "" {
		return nil
	}

	// Allow letters, numbers, space, dot, underscore, dash
	re := regexp.MustCompile(`^[a-zA-Z0-9 ._-]+$`)
	if !re.MatchString(s) {
		return errors.New("contains invalid characters (only letters, numbers, spaces, dots, underscores, and dashes are allowed)")
	}
	return nil
}

func validateCoverImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorBankImageMissingOrInvalid
	}

	if !utils.IsValidImage(file) {
		return errors.New(localization.MsgBankImageRequiredOrMissing)
	}

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}
