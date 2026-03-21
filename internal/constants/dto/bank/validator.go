package bank_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"mime/multipart"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var bankNameRegex = regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)

func (c CreateBankRequest) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Name,
			validation.Required.Error(localization.MsgBankNameRequired),
			validation.Match(bankNameRegex).Error("Bank name must not contain special characters"),
		),
		validation.Field(&c.BICCode, validation.Required.Error(localization.MsgBankBICRequired)),
		validation.Field(&c.Logo, validation.Required.Error(localization.MsgBankLogoRequired), validation.By(func(value interface{}) error { return validateLogo(value) })),
		validation.Field(&c.HasAlphaNumeric,
			validation.Required.Error(localization.MsgBankHasAlphaNumericRequired),
		),
		validation.Field(&c.AccountLength,
			validation.Required.Error(localization.MsgBankAccountLengthRequired),
			validation.Min(1).Error("Account length must be positive"),
		),
	)
}

func (u UpdateBankRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name,
			validation.When(u.Name != "",
				validation.Match(bankNameRegex).Error("Bank name must not contain special characters"),
			),
		),
		validation.Field(&u.BICCode,
			validation.When(u.BICCode != "", validation.Length(1, 0)),
		),
		validation.Field(&u.Logo, validation.By(func(value interface{}) error {
			file, ok := value.(*multipart.FileHeader)
			if !ok || file == nil {
				return nil
			}
			return validateLogo(file)
		})),
		validation.Field(&u.HasAlphaNumeric,
			validation.When(u.HasAlphaNumeric != nil && *u.HasAlphaNumeric, validation.In(true, false)),
		),
		validation.Field(&u.AccountLength,
			validation.When(u.AccountLength != 0,
				validation.Min(1).Error("Account length must be positive"),
			),
		),
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
