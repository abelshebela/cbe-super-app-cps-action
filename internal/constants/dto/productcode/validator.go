package productcode

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (u *UpdateProductCodeRequest) Validate() error {
	u.ProductName = local_util.ExtraSpaceRemover(u.ProductName)
	return atLeastOne(u)
}

func validateOptionalProductCodes(value interface{}) error {
	if pc, ok := value.(*model.ProductCodes); ok {
		if *pc == (model.ProductCodes{}) {
			return nil
		}
		return pc.Validate()
	}
	return fmt.Errorf(localization.ErrorProductCodeValidationError.Code)
}

func atLeastOne(u *UpdateProductCodeRequest) error {
	var validationErrors []error
	u.ProductName = local_util.ExtraSpaceRemover(u.ProductName)
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
	if len(validationErrors) > 0 {
		return fmt.Errorf(localization.ErrorProductCodeUpdateRequestValidationErrorAtLeastOne.Code)
	}
	return nil
}
