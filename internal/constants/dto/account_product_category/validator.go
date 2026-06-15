package account_product_category_dto

import (
	"errors"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	safeStringRe  = regexp.MustCompile(`^[a-zA-Z0-9 _\-\.]+$`)
	codeRe        = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	validAPCTypes = map[string]bool{"IFB": true, "CB": true, "BOTH": true}
)

func noSpecialChars(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !safeStringRe.MatchString(s) {
		return validation.NewError("validation_special_chars", "contains invalid characters")
	}
	return nil
}

func validProductLine(value interface{}) error {
	s, _ := value.(string)
	if s == "" {
		return nil
	}
	if !validAPCTypes[strings.ToUpper(s)] {
		return validation.NewError("validation_product_line", "must be one of IFB, CB, BOTH")
	}
	return nil
}

func codeFormat(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !codeRe.MatchString(s) {
		return validation.NewError("validation_code_format", "must contain only alphanumeric characters, hyphens, or underscores")
	}
	return nil
}

func (r CreateAPCRequest) Validate() error {
	if strings.TrimSpace(r.ProductLine) == "" {
		return errors.New("product line is required")
	}
	if err := validProductLine(r.ProductLine); err != nil {
		return errors.New(err.Error())
	}

	if strings.TrimSpace(r.CBSCategoryCode) == "" {
		return errors.New("cbs category code is required")
	}
	if len(r.CBSCategoryCode) > 64 {
		return errors.New("cbs category code must be at most 64 characters")
	}
	if err := codeFormat(r.CBSCategoryCode); err != nil {
		return errors.New(err.Error())
	}

	if strings.TrimSpace(r.CategoryName) == "" {
		return errors.New("category name is required")
	}
	if len(r.CategoryName) > 128 {
		return errors.New("category name must be at most 128 characters")
	}
	if err := noSpecialChars(r.CategoryName); err != nil {
		return errors.New(err.Error())
	}

	if len(r.Description) > 512 {
		return errors.New("description must be at most 512 characters")
	}
	if err := noSpecialChars(r.Description); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

func (r UpdateAPCRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductLine,
			validation.When(r.ProductLine != "", validation.By(validProductLine)),
		),
		validation.Field(&r.CBSCategoryCode,
			validation.Length(0, 64),
			validation.When(r.CBSCategoryCode != "", validation.By(codeFormat)),
		),
		validation.Field(&r.CategoryName,
			validation.Length(0, 128),
			validation.When(r.CategoryName != "", validation.By(noSpecialChars)),
		),
		validation.Field(&r.Description,
			validation.Length(0, 512),
			validation.When(r.Description != "", validation.By(noSpecialChars)),
		),
	)
}
