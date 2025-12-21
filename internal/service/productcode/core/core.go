package core

import (
	"cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func NonEmptyProductCodes(request, existing model.ProductCodes) model.ProductCodes {
	return model.ProductCodes{
		PRD:    utils.NonEmptyString(request.PRD, existing.PRD),
		VATPRD: utils.NonEmptyString(request.VATPRD, existing.VATPRD),
		SFPRD:  utils.NonEmptyString(request.SFPRD, existing.SFPRD),
		TRXN:   utils.NonEmptyString(request.TRXN, existing.TRXN),
	}
}
