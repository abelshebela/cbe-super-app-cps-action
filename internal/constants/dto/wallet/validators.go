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
	if !(w.Self || w.Other || w.Agent) {
		return validation.NewError(localization.ErrorWalletRechangeOption.Code, localization.ErrorWalletRechangeOption.Message)
	}

	return nil
}
func (w WalletRequest) AggregatedValidate(isCreate bool) error {
	errs := validation.Errors{}
	if strings.TrimSpace(w.Name) == "" {
		errs["name"] = errors.New(localization.ErrorWalletNameRequired.Code)
	}
	if strings.TrimSpace(w.Code) == "" {
		errs["code"] = errors.New(localization.ErrorWalletCodeRequired.Code)
	}
	if !(w.Self || w.Other || w.Agent) {
		errs["recharge_option"] = errors.New(localization.ErrorWalletRechangeOption.Code)
	}
	if isCreate {
		if w.Avatar == nil {
			errs["avatar"] = errors.New(localization.ErrorWalletAvatarRequired.Code)
		} else if err := validateAvatar(w.Avatar); err != nil {
			errs["avatar"] = err
		}
	} else if w.Avatar != nil {
		if err := validateAvatar(w.Avatar); err != nil {
			errs["avatar"] = err
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateAvatar(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return errors.New(localization.ErrorWalletAvatarInvalid.Code)
	}
	if !isImageFormat(file) {
		return errors.New(localization.ErrorWalletAvatarInvalidType.Code)
	}

	if file.Size > maxFileSize {
		return errors.New(localization.ErrorWalletAvatarTooLarge.Code)
	}

	// this is causing issue with mobile and front end upload
	// ct := file.Header.Get("Content-Type")
	// if ct == "" || !allowedMIMETypes[ct] {
	// 	return errors.New(localization.ErrorWalletAvatarInvalidType.Code)
	// }
	return nil
}

// use this for image validation this works with the mobile and the frontend
func isImageFormat(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}
	contentType := fileHeader.Header.Get("Content-Type")

	// Acceptable image formats
	ext := strings.ToLower(strings.TrimPrefix(strings.ToLower(fileHeader.Filename[strings.LastIndex(fileHeader.Filename, "."):]), "."))
	switch ext {
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
		contentType = "image/png"
	case "gif":
		contentType = "image/gif"
	}

	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		return true
	default:
		return false
	}
}
