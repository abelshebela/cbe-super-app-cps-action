package topupDto

import (
	"errors"
	"mime/multipart"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// const maxFileSize = 2 * 1024 * 1024

// var allowedMIMETypes = map[string]bool{
// 	"image/jpeg": true,
// 	"image/png":  true,
// 	"image/gif":  true,
// 	"image/webp": true,
// }

func (w TopupRequest) IsEmpty() bool {
	return strings.TrimSpace(w.Name) == "" &&
		strings.TrimSpace(w.Code) == "" &&
		w.Avatar == nil
}

func (w TopupRequest) Validate(isCreate bool) error {
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
		rules = append(rules, validation.Field(&w.Name, validation.By(validateString("name", w.Name, true, localization.ErrorTopupNameRequired.Code)), validation.By(utils.TrimWhiteSpace)))
	} else if w.Name != "" {
		rules = append(rules, validation.Field(&w.Name, validation.By(validateString("name", w.Name, false, localization.ErrorTopupNameRequired.Code))))
	}

	if isCreate {
		rules = append(rules, validation.Field(&w.Code, validation.By(validateString("code", w.Code, true, localization.ErrorTopupCodeRequired.Code)), validation.By(utils.TrimWhiteSpace)))
	} else if w.Code != "" {
		rules = append(rules, validation.Field(&w.Code, validation.By(validateString("code", w.Code, false, localization.ErrorTopupCodeRequired.Code))))
	}

	if len(rules) > 0 {
		if err := validation.ValidateStruct(&w, rules...); err != nil {
			return err
		}
	}
	if isCreate {
		rules = append(rules, validation.Field(&w.Avatar,
			validation.Required.Error(localization.ErrorTopupAvatarRequired.Code),
			validation.By(func(value interface{}) error { return validateImage(value) }),
		))
	} else if w.Avatar != nil {
		rules = append(rules, validation.Field(&w.Avatar,
			validation.By(func(value interface{}) error { return validateImage(value) }),
		))
	}

	if !(w.Self || w.Other || w.Agent) {
		return validation.NewError(localization.ErrorTopupServiceOption.Code, localization.ErrorTopupServiceOption.Message)
	}
	return nil
}
func (w TopupRequest) AggregatedValidate(isCreate bool) error {
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
		} else if err := validateImage(w.Avatar); err != nil {
			errs["avatar"] = err
		}
	} else if w.Avatar != nil {
		if err := validateImage(w.Avatar); err != nil {
			errs["avatar"] = err
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidImage
	}
	if !utils.IsValidImage(file) {
		return localization.ErrorMissingOrInvalidImage
	}

	if file.Size > (15 << 20) {
		return validation.NewError("logo", localization.MsgFileTooLarge)
	}

	return nil
}
