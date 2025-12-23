package walletDto

import (
	"errors"
	"mime/multipart"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/pkgs/utils"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const specialChars = "`~!@#$%^&*()-_=+[]{}\\|;:'\",<.>/?"

func (w WalletRequest) IsEmpty() bool {
	return strings.TrimSpace(w.Name) == "" &&
		strings.TrimSpace(w.Code) == "" &&
		w.Avatar == nil
}

func (w WalletRequest) Validate(isCreate bool) error {
	errs := validation.Errors{}

	// --- Name ---
	if isCreate || w.Name != "" {
		if strings.TrimSpace(w.Name) == "" && isCreate {
			errs["name"] = localization.ErrorWalletNameRequired
		} else if strings.ContainsAny(w.Name, specialChars) {
			errs["name"] = localization.ErrorInvalidWalletName
		}
	}

	// --- Code ---
	if isCreate || w.Code != "" {
		if strings.TrimSpace(w.Code) == "" && isCreate {
			errs["code"] = localization.ErrorWalletCodeRequired
		} else if strings.ContainsAny(w.Code, specialChars) {
			errs["code"] = localization.ErrorInvalidWalletCode
		}
	}

	// --- Type ---
	if isCreate || w.Type != "" {
		if strings.TrimSpace(string(w.Type)) == "" && isCreate {
			errs["type"] = localization.ErrorWalletTypeRequired
		}
	}

	// --- Avatar ---
	if isCreate {
		if w.Avatar == nil {
			errs["avatar"] = localization.ErrorWalletAvatarRequired
		} else if err := validateAvatar(w.Avatar); err != nil {
			errs["avatar"] = err
		}
	} else if w.Avatar != nil {
		if err := validateAvatar(w.Avatar); err != nil {
			errs["avatar"] = err
		}
	}

	// --- Services check (cannot be set) ---
	// if w.Self != nil || w.Other != nil || w.Agent != nil {
	// 	errs["recharge_option"] = localization.ErrorWalletRechangeOption
	// }

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func validateAvatar(file *multipart.FileHeader) error {
	if file == nil {
		return localization.ErrorWalletAvatarInvalid
	}
	if !utils.IsValidImage(file) {
		return errors.New(localization.ErrorWalletAvatarInvalidType.Code)
	}
	if file.Size > (15 << 20) { // 15 MB limit
		return validation.NewError("avatar", localization.MsgFileTooLarge)
	}
	return nil
}
