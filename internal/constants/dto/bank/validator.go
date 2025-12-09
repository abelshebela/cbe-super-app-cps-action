package bank_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"mime/multipart"
	"regexp"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (c CreateBankRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name, validation.Required.Error(localization.MsgBankNameRequired), validation.Length(3, 25), validation.Match(regexp.MustCompile(`^[a-zA-Z0-9 ]+$`)).Error("Name must not be contain special char")),
		validation.Field(&c.Code, validation.Required.Error(localization.MsgBankCodeRequired)),
		validation.Field(&c.BIC, validation.Required.Error(localization.MsgBankBICRequired)),
		validation.Field(&c.AccountLength, validation.Required.Error(localization.MsgBankAccountLengthRequired)),
		validation.Field(&c.Logo, validation.By(func(value interface{}) error { return validateLogo(value) })),
	)
}

func (u UpdateBankRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.NilOrNotEmpty,
			validation.Length(3, 50),
			is.Alpha,
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

	if file.Size > (15 << 20) {
		return validation.NewError("logo", localization.MsgFileTooLarge)
	}

	return nil
}
