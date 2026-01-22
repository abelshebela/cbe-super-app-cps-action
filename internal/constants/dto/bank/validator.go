package bank_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c CreateBankRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error(localization.MsgBankNameRequired)),
		validation.Field(&c.Type,
			validation.Required.Error(localization.MsgBankTypeRequired),
			validation.In("BANK", "WALLET", "MFI").Error(localization.MsgInvalidRequestBankType),
		),
		validation.Field(&c.BICCode, validation.Required.Error(localization.MsgBankBICRequired)),
		validation.Field(&c.Logo, validation.By(func(value interface{}) error { return validateLogo(value) })),
	)
}

func (u UpdateBankRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.NilOrNotEmpty,
		),
		validation.Field(&u.Logo, validation.By(func(value interface{}) error {
			if value == nil {
				return nil
			}
			return validateLogo(value)
		})),
		// Code and BIC are optional on update — allow empty values by not validating them here
	)
}

func (u UpdateLogo) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Logo, validation.By(func(value interface{}) error { return validateLogo(value) })),
	)
}

func validateLogo(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorBankImageMissingOrInvalid
	}

	if !utils.IsValidImage(file) {
		return errors.New(localization.ErrorBankImageMissingOrInvalid.Code)
	}

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}
