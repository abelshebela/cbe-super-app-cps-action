package model

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type ProductCode struct {
	ID                 string       `json:"id" bson:"_id"`
	ProductName        string       `json:"product_name" bson:"product_name"`
	CBEProductCodes    ProductCodes `json:"cbe_product_codes" bson:"cbe_product_codes"`
	CBEIFBProductCodes ProductCodes `json:"cbe_ifb_product_codes" bson:"cbe_ifb_product_codes"`
	CreatedAt          time.Time    `json:"created_at" bson:"created_at"`
	LastUpdatedAt      time.Time    `json:"last_updated_at" bson:"last_updated_at"`
}

type ProductCodes struct {
	PRD    string `json:"prd" bson:"prd"`
	VATPRD string `json:"vat_prd" bson:"vatprd"`
	SFPRD  string `json:"sf_prd" bson:"sfprd"`
	TRXN   string `json:"trxn" bson:"trxn"`
}

func NonEmptyProductCodes(request, existing ProductCodes) ProductCodes {
	return ProductCodes{
		PRD:    utils.NonEmptyString(request.PRD, existing.PRD),
		VATPRD: utils.NonEmptyString(request.VATPRD, existing.VATPRD),
		SFPRD:  utils.NonEmptyString(request.SFPRD, existing.SFPRD),
		TRXN:   utils.NonEmptyString(request.TRXN, existing.TRXN),
	}
}

func (pc *ProductCodes) Validate() error {
	pc.PRD = local_util.ExtraSpaceRemover(pc.PRD)
	pc.VATPRD = local_util.ExtraSpaceRemover(pc.VATPRD)
	pc.SFPRD = local_util.ExtraSpaceRemover(pc.SFPRD)

	err := validation.ValidateStruct(&pc,
		validation.Field(&pc.PRD,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&pc.VATPRD,
			validation.By(utils.NoSpecialChars),
		),
		validation.Field(&pc.SFPRD,
			validation.By(utils.NoSpecialChars),
		),
	)
	if err!= nil{
		err = fmt.Errorf(localization.ErrorProductCodesValidationError.Code)
	}
	return err
}
