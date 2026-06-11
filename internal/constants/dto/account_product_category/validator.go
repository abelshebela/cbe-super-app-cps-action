package account_product_category_dto

import (
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
	return validation.ValidateStruct(&r,
		validation.Field(&r.ProductLine,
			validation.Required,
			validation.By(validProductLine),
		),
		validation.Field(&r.CBSCategoryCode,
			validation.Required,
			validation.Length(1, 64),
			validation.By(codeFormat),
		),
		validation.Field(&r.CategoryName,
			validation.Required,
			validation.Length(1, 128),
			validation.By(noSpecialChars),
		),
		validation.Field(&r.Description,
			validation.Length(0, 512),
			validation.By(noSpecialChars),
		),
	)
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
