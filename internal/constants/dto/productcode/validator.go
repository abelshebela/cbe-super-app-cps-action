package productcode

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u *UpdateProductCodeRequest) Validate() error {

	u.ProductName = strings.TrimSpace(u.ProductName)
	return atLeastOne(u)
}

func validateOptionalProductCodes(value interface{}) error {
	if pc, ok := value.(*model.ProductCodes); ok {
		if *pc == (model.ProductCodes{}) {
			return nil
		}
		return pc.Validate()
	}
	return validation.NewError("invalid_product_codes", "ProductCodes is required and must be valid.")
}

func atLeastOne(u *UpdateProductCodeRequest) error {
	var validationErrors []error
	u.ProductName = strings.TrimSpace(u.ProductName)
	if (u.CBEProductCodes == (model.ProductCodes{})) && (u.CBEIFBProductCodes == (model.ProductCodes{})) && u.ProductName == "" {
		validationErrors = append(validationErrors, validation.NewError("missing_product_codes", "At least one of CBEProductCodes or CBEIFBProductCodes must be present"))
	}
	if u.CBEProductCodes != (model.ProductCodes{}) {
		e := validateOptionalProductCodes(u.CBEProductCodes)
		if e != nil {
			validationErrors = append(validationErrors, e)
		}
	}
	if u.CBEIFBProductCodes != (model.ProductCodes{}) {
		e := validateOptionalProductCodes(u.CBEIFBProductCodes)
		if e != nil {
			validationErrors = append(validationErrors, e)
		}
	}
	errString := ""
	for i := range validationErrors {
		errString += validationErrors[i].Error()
	}
	if len(errString) > 0 {
		return errors.New(errString)
	}
	return nil
}
