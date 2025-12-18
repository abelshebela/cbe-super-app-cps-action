package productcode

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

// Validate ensures UpdateProductCodeRequest has at least one valid update.
func (u *UpdateProductCodeRequest) Validate() error {
	// Normalize spaces
	u.ProductName = utils.ExtraSpaceRemover(u.ProductName)
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
		if err := validatecode(&u.CBEProductCodes); err != nil {
			errs["cbe_product_codes"] = err
		}
	}
	if hasIFB {
		if err := validatecode(&u.CBEIFBProductCodes); err != nil {
			errs["cbe_ifb_product_codes"] = err
		}
	}

	// Return collected errors if any
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validatecode(pc *model.ProductCodes) error {
	pc.PRD = utils.ExtraSpaceRemover(pc.PRD)
	pc.VATPRD = utils.ExtraSpaceRemover(pc.VATPRD)
	pc.SFPRD = utils.ExtraSpaceRemover(pc.SFPRD)
	pc.TRXN = utils.ExtraSpaceRemover(pc.TRXN)

	err := validation.ValidateStruct(pc,
		validation.Field(&pc.PRD, validation.By(utils.NoSpecialChars)),
		validation.Field(&pc.VATPRD, validation.By(utils.NoSpecialChars)),
		validation.Field(&pc.SFPRD, validation.By(utils.NoSpecialChars)),
		validation.Field(&pc.TRXN, validation.By(utils.NoSpecialChars)),
	)

	if err != nil {
		return fmt.Errorf("%s: %v", localization.ErrorProductCodesValidationError.Code, err)
	}
	return nil
}
