package account_product_dto

import (
	"errors"
	"regexp"
	"strings"

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


func (r CreateAPRequest) Validate() error {
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
			validation.Min(float64(1)),
		),
		validation.Field(&r.InterestRate,
			validation.Required,
			validation.Min(float64(1)),
		),
		validation.Field(&r.MinimumMaintenanceFee,
			validation.Min(float64(1)),
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
	if r.CBSProductCode != "" {
		if len(r.CBSProductCode) > 64 {
			return errors.New("CBS Product Code must be at most 64 characters")
		}
		if err := apCodeFormat(r.CBSProductCode); err != nil {
			return errors.New("CBS Product Code: " + err.Error())
		}
	}
	if r.ProductName != "" {
		if len(r.ProductName) > 128 {
			return errors.New("Product Name must be at most 128 characters")
		}
		if err := apNoSpecialChars(r.ProductName); err != nil {
			return errors.New("Product Name: " + err.Error())
		}
	}
	if r.InterestRate < 0{
		return errors.New("Interest rate can not less that zero")
	}
	if r.ProductTagLine != "" {
		if len(r.ProductTagLine) > 128 {
			return errors.New("Product Tag Line must be at most 128 characters")
		}
		if err := apNoSpecialChars(r.ProductTagLine); err != nil {
			return errors.New("Product Tag Line: " + err.Error())
		}
	}
	if r.AccountCurrency != "" {
		if err := apValidCurrency(r.AccountCurrency); err != nil {
			return errors.New("Account Currency: " + err.Error())
		}
	}
	if r.MinimumOpeningBalance != 0 && r.MinimumOpeningBalance < 1 {
		return errors.New("Minimum Opening Balance must be at least 1")
	}
	if r.MinimumMaintenanceFee != 0 && r.MinimumMaintenanceFee < 1 {
		return errors.New("Minimum Maintenance Fee must be at least 1")
	}
	if len(r.FaqURL) > 512 {
		return errors.New("FAQ URL must be at most 512 characters")
	}
	if len(r.ProductFeatures) > 1024 {
		return errors.New("Product Features must be at most 1024 characters")
	}
	return nil
}
