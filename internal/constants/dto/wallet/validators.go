package walletDto

import (
	"errors"
	"mime/multipart"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

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
		rules = append(rules, validation.Field(&w.Name, validation.By(validateString("name", w.Name, true, localization.ErrorWalletNameRequired.Code)), validation.By(utils.TrimWhiteSpace)))
	} else if w.Name != "" {
		rules = append(rules, validation.Field(&w.Name, validation.By(validateString("name", w.Name, false, localization.ErrorWalletNameRequired.Code))))
	}

	if isCreate {
		rules = append(rules, validation.Field(&w.Code, validation.By(validateString("code", w.Code, true, localization.ErrorWalletCodeRequired.Code)), validation.By(utils.TrimWhiteSpace)))
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

	if isCreate {
		rules = append(rules, validation.Field(&w.Type, validation.By(validateString("type", string(w.Type), true, localization.ErrorWalletTypeRequired.Code))))
	} else if string(w.Type) != "" {
		rules = append(rules, validation.Field(&w.Type, validation.By(validateString("type", string(w.Type), false, localization.ErrorWalletTypeRequired.Code))))
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

const specialChars = "`~!@#$%^&*()-_=+[]{}\\|;:'\",<.>/?"

func (w WalletRequest) AggregatedValidate(isCreate bool) error {
	errs := validation.Errors{}

	if isCreate {
		if strings.TrimSpace(w.Name) == "" {
			errs["name"] = localization.ErrorWalletNameRequired
		} else {
			if strings.ContainsAny(w.Name, specialChars) {
				errs["name"] = localization.ErrorInvalidWalletName
			}
		}
		if strings.TrimSpace(w.Code) == "" {
			errs["code"] = localization.ErrorWalletCodeRequired
		} else {
			if strings.ContainsAny(w.Code, specialChars) {
				errs["code"] = localization.ErrorInvalidWalletCode
			}
		}
		if !(w.Self || w.Other || w.Agent) {
			errs["recharge_option"] = localization.ErrorWalletRechangeOption
		}
		if w.Avatar == nil {
			errs["avatar"] = localization.ErrorWalletAvatarRequired
		} else if err := validateAvatar(w.Avatar); err != nil {
			errs["avatar"] = err
		}
		if strings.TrimSpace(string(w.Type)) == "" {
			errs["type"] = localization.ErrorWalletTypeRequired
		}
	} else {
		if w.Name != "" && strings.TrimSpace(w.Name) == "" {
			errs["name"] = localization.ErrorWalletNameRequired
		} else if w.Name != "" {
			if strings.ContainsAny(w.Name, specialChars) {
				errs["name"] = localization.ErrorInvalidWalletName
			}
		}
		if w.Code != "" && strings.TrimSpace(w.Code) == "" {
			errs["code"] = localization.ErrorWalletCodeRequired
		} else if w.Code != "" {
			if strings.ContainsAny(w.Code, specialChars) {
				errs["code"] = localization.ErrorInvalidWalletCode
			}
		}
		if w.Avatar != nil {
			if err := validateAvatar(w.Avatar); err != nil {
				errs["avatar"] = err
			}
		}
		if string(w.Type) != "" && strings.TrimSpace(string(w.Type)) == "" {
			errs["type"] = localization.ErrorWalletTypeRequired
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
		return localization.ErrorWalletAvatarInvalid
	}
	if !utils.IsValidImage(file) {
		return errors.New(localization.ErrorWalletAvatarInvalidType.Code)
	}
	if file.Size > (15 << 20) {
		return validation.NewError("logo", localization.MsgFileTooLarge)
	}

	return nil
}
