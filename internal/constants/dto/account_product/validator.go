package account_product_dto

import (
	"errors"
	"regexp"
	"strings"

	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var (
	apSafeStringRe = regexp.MustCompile(`^[a-zA-Z0-9 _\-\.]+$`)
	apCodeRe       = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	currencyRe     = regexp.MustCompile(`^[A-Z]{3}$`)
)

func apNoSpecialChars(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !apSafeStringRe.MatchString(s) {
		return validation.NewError("validation_special_chars", "contains invalid characters")
	}
	return nil
}

func apCodeFormat(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !apCodeRe.MatchString(s) {
		return errors.New("must contain only alphanumeric characters, hyphens, or underscores")
	}
	return nil
}

func apValidCurrency(value interface{}) error {
	s, _ := value.(string)
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	if !currencyRe.MatchString(s) {
		return validation.NewError("validation_currency", "must be a 3-letter uppercase currency code")
	}
	return nil
}

func apValidIcon(value interface{}) error {
	fh, ok := value.(interface{ GetHeader() interface{} })
	_ = fh
	_ = ok
	return nil
}

func validateIconFile(req *CreateAPRequest) error {
	if req.Icon == nil {
		return validation.NewError("validation_icon_required", "icon image is required")
	}
	if !utils.IsValidImage(req.Icon) {
		return validation.NewError("validation_icon_invalid", "icon must be a valid image (jpeg/png/gif/webp)")
	}
	return nil
}

func (r CreateAPRequest) Validate() error {
	if err := validateIconFile(&r); err != nil {
		return err
	}
	return validation.ValidateStruct(&r,
		validation.Field(&r.CBSProductCode,
			validation.Required,
			validation.Length(1, 64),
			validation.By(apCodeFormat),
		),
		validation.Field(&r.ProductName,
			validation.Required,
			validation.Length(1, 128),
			validation.By(apNoSpecialChars),
		),
		validation.Field(&r.ProductTagLine,
			validation.Required,
			validation.Length(1, 128),
			validation.By(apNoSpecialChars),
		),
		validation.Field(&r.AccountCategoryID,
			validation.Required,
		),
		validation.Field(&r.AccountCurrency,
			validation.Required,
			validation.By(apValidCurrency),
		),
		validation.Field(&r.MinimumOpeningBalance,
			validation.Required,
			validation.Min(float64(0)),
		),
		validation.Field(&r.MinimumMaintenanceFee,
			validation.Min(float64(0)),
		),
		validation.Field(&r.FaqURL,
			validation.Length(0, 512),
		),
		validation.Field(&r.ProductFeatures,
			validation.Length(0, 1024),
		),
	)
}

func (r UpdateAPRequest) Validate() error {
	if r.Icon != nil && !utils.IsValidImage(r.Icon) {
		return validation.NewError("validation_icon_invalid", "icon must be a valid image (jpeg/png/gif/webp)")
	}
	return validation.ValidateStruct(&r,
		validation.Field(&r.CBSProductCode,
			validation.Length(0, 64),
			validation.When(r.CBSProductCode != "", validation.By(apCodeFormat)),
		),
		validation.Field(&r.ProductName,
			validation.Length(0, 128),
			validation.When(r.ProductName != "", validation.By(apNoSpecialChars)),
		),
		validation.Field(&r.ProductTagLine,
			validation.Length(0, 128),
			validation.When(r.ProductTagLine != "", validation.By(apNoSpecialChars)),
		),
		validation.Field(&r.AccountCurrency,
			validation.When(r.AccountCurrency != "", validation.By(apValidCurrency)),
		),
		validation.Field(&r.MinimumOpeningBalance,
			validation.Min(float64(0)),
		),
		validation.Field(&r.MinimumMaintenanceFee,
			validation.Min(float64(0)),
		),
		validation.Field(&r.FaqURL,
			validation.Length(0, 512),
		),
		validation.Field(&r.ProductFeatures,
			validation.Length(0, 1024),
		),
	)
}
