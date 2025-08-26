package walletDto

import (
	"errors"
	"mime/multipart"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const maxFileSize = 2 * 1024 * 1024

var allowedMIMETypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

func (w WalletRequest) IsEmpty() bool {
	return strings.TrimSpace(w.Name) == "" &&
		strings.TrimSpace(w.Code) == "" &&
		w.Avatar == nil
}

func (w WalletRequest) Validate(isCreate bool) error {
	if !isCreate && w.IsEmpty() {
		return nil
	}

	var rules []*validation.FieldRules

	validateString := func(fieldName, value string, required bool, errCode string) validation.RuleFunc {
		return func(_ interface{}) error {
			if required && strings.TrimSpace(value) == "" {
				return errors.New(errCode)
			}
			return nil
		}
	}

 	if isCreate {
		rules = append(rules, validation.Field(&w.Name, validation.By(validateString("name", w.Name, true, localization.ErrorWalletNameRequired.Code))))
	} else if w.Name != "" {
		rules = append(rules, validation.Field(&w.Name, validation.By(validateString("name", w.Name, false, localization.ErrorWalletNameRequired.Code))))
	}

 	if isCreate {
		rules = append(rules, validation.Field(&w.Code, validation.By(validateString("code", w.Code, true, localization.ErrorWalletCodeRequired.Code))))
	} else if w.Code != "" {
		rules = append(rules, validation.Field(&w.Code, validation.By(validateString("code", w.Code, false, localization.ErrorWalletCodeRequired.Code))))
	}

 	if isCreate {
		rules = append(rules, validation.Field(&w.Avatar,
			validation.Required.Error(localization.ErrorWalletAvatarRequired.Code),
			validation.By(func(value interface{}) error { return validateAvatar(value) }),
		))
	} else if w.Avatar != nil {
		rules = append(rules, validation.Field(&w.Avatar,
			validation.By(func(value interface{}) error { return validateAvatar(value) }),
		))
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&w, rules...); err != nil {
			return err
		}
	}

	return nil
}

func validateAvatar(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return errors.New(localization.ErrorWalletAvatarInvalid.Code)
	}

	if file.Size > maxFileSize {
		return errors.New(localization.ErrorWalletAvatarTooLarge.Code)
	}

	ct := file.Header.Get("Content-Type")
	if ct == "" || !allowedMIMETypes[ct] {
		return errors.New(localization.ErrorWalletAvatarInvalidType.Code)
	}
	return nil
}
