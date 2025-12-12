package productcode

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Validate ensures UpdateProductCodeRequest has at least one valid update.
func (u *UpdateProductCodeRequest) Validate() error {
	// Normalize spaces
	u.ProductName = local_util.ExtraSpaceRemover(u.ProductName)
	hasProductName := u.ProductName != ""
	hasCBE := u.CBEProductCodes != (model.ProductCodes{})
	hasIFB := u.CBEIFBProductCodes != (model.ProductCodes{})

	if !hasProductName && !hasCBE && !hasIFB {
		return validation.NewError(
			"missing_update_fields",
			localization.ErrorProductCodeUpdateRequestValidationErrorAtLeastOne.Code,
		)
	}

	errs := validation.Errors{}

	// Validate ProductName only if provided
	if hasProductName {
		if err := validation.Validate(u.ProductName,
			validation.RuneLength(2, 100).Error("Error Invalid Length"),
		); err != nil {
			errs["product_name"] = err
		}
	}

	// Validate nested ProductCodes only if provided
	if hasCBE {
		if err := u.CBEProductCodes.Validate(); err != nil {
			errs["cbe_product_codes"] = err
		}
	}
	if hasIFB {
		if err := u.CBEIFBProductCodes.Validate(); err != nil {
			errs["cbe_ifb_product_codes"] = err
		}
	}

	// Return collected errors if any
	if len(errs) > 0 {
		return errs
	}
	return nil
}
