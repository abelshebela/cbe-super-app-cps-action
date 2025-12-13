package topupDto

import (
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
	var rules []*validation.FieldRules

	if isCreate {
		rules = append(rules,
			validation.Field(
				&w.Name,
				validation.Required.Error(localization.ErrorTopupNameRequired.Code),
				validation.By(utils.TrimWhiteSpace),
			),
		)
	} else {
		rules = append(rules,
			validation.Field(
				&w.Name,
				validation.By(utils.TrimWhiteSpace),
			),
		)
	}

	if isCreate {
		rules = append(rules,
			validation.Field(
				&w.Code,
				validation.Required.Error(localization.ErrorTopupCodeRequired.Code),
				validation.By(utils.TrimWhiteSpace),
			),
		)
	} else {
		rules = append(rules,
			validation.Field(
				&w.Code,
				validation.By(utils.TrimWhiteSpace),
			),
		)
	}

	if isCreate {
		rules = append(rules,
			validation.Field(
				&w.Avatar,
				validation.Required.Error(localization.ErrorTopupAvatarRequired.Code),
				validation.By(validateImage),
			),
		)
	} else if w.Avatar != nil {
		rules = append(rules,
			validation.Field(
				&w.Avatar,
				validation.By(validateImage),
			),
		)
	}

	if err := validation.ValidateStruct(&w, rules...); err != nil {
		return err
	}

	if isCreate {
		if !(w.Self || w.Other || w.Agent) {
			return validation.NewError(
				localization.ErrorTopupServiceOption.Code,
				localization.ErrorTopupServiceOption.Message,
			)
		}
	}

	return nil
}

func (w TopupRequest) AggregatedValidate(isCreate bool) error {
	errs := validation.Errors{}
	if strings.TrimSpace(w.Name) == "" {
		errs["name"] = localization.ErrorWalletNameRequired
	}
	if strings.TrimSpace(w.Code) == "" {
		errs["code"] = localization.ErrorWalletCodeRequired
	}
	if !(w.Self || w.Other || w.Agent) {
		errs["recharge_option"] = localization.ErrorWalletRechangeOption
	}
	if isCreate {
		if w.Avatar == nil {
			errs["avatar"] = localization.ErrorWalletAvatarRequired
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
