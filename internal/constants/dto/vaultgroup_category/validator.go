package vaultgroupcategory

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

func (r *CreateVaultGroupCategoryRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}

	var err error
	if r.Name, err = sanitizeString(r.Name); err != nil {
		return err
	}
	if r.Description, err = sanitizeString(r.Description); err != nil {
		return err
	}

	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Required,
			validation.Length(3, 100),
			validation.By(noSpecialChars),
		),
		validation.Field(&r.Description,
			validation.Required,
			validation.Length(3, 255),
			validation.By(noSpecialChars),
		),
	)
}

func (r *UpdateVaultGroupCategoryRequest) Validate() error {
	if r == nil {
		return errors.New("request is required")
	}
	if r.Name == nil && r.Description == nil {
		return errors.New("at least one field (name or description) must be provided")
	}
	return validation.ValidateStruct(r,
		validation.Field(&r.Name,
			validation.Length(3, 100),
			validation.By(noSpecialChars),
		),
		validation.Field(&r.Description,
			validation.Length(3, 255),
			validation.By(noSpecialChars),
		),
	)
}

func sanitizeString(s string) (string, error) {
	trimmed := strings.TrimSpace(s)
	return trimmed, nil
}

// var reISO4217 = regexp.MustCompile(`^[A-Z]{3}$`)

func noSpecialChars(value interface{}) error {
	var s string

	switch v := value.(type) {
	case string:
		s = v
	case *string:
		if v != nil {
			s = *v
		}
	default:
		return nil
	}

	if s == "" {
		return nil
	}

	// allow letters, numbers, space, dot, underscore, dash
	re := regexp.MustCompile(`^[a-zA-Z0-9 ._-]+$`)
	if !re.MatchString(s) {
		// return validation.NewError("validation_no_special_chars", "contains invalid characters")
		return errors.New("contains invalid characters")
	}
	return nil
}
