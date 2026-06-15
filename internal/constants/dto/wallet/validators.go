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
		strings.TrimSpace(w.UniqueCode) == "" &&
		w.Self == nil &&
		w.Other == nil &&
		w.Agent == nil &&
		w.Avatar == nil
}

func (w WalletRequest) Validate(isCreate bool) error {
	errs := validation.Errors{}

	if isCreate || w.Name != "" {
		if strings.TrimSpace(w.Name) == "" && isCreate {
			errs["name"] = localization.ErrorWalletNameRequired
		}
	}
	// --- Code ---
	if isCreate || w.UniqueCode != "" {
		trimmed := strings.TrimSpace(w.UniqueCode)
		if trimmed == "" && isCreate {
			errs["unique_code"] = localization.ErrorWalletCodeRequired
		} else if len(trimmed) > 10 || !isAlpha(trimmed) {
			errs["unique_code"] = localization.ErrorInvalidWalletCode
		}
		w.UniqueCode = strings.ToUpper(trimmed)
	}
	if isCreate {
		if w.Self == nil && w.Other == nil && w.Agent == nil {
			errs["type"] = localization.ErrorWalletTypeRequired
		}
		if w.SelfServiceID == "" && w.OtherServiceID == "" && w.AgentServiceID == "" {
			errs["type"] = localization.ErrorWalletAgentServiceIDRequired
			errs["type"] = localization.ErrorWalletSelfServiceIDRequired
			errs["type"] = localization.ErrorWalletOtherServiceIDRequired
		}
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

	if w.Self != nil && *w.Self {
		if w.SelfServiceID == "" {
			return localization.ErrorWalletSelfServiceIDRequired
		}
	}
	if w.Other != nil && *w.Other {
		if w.OtherServiceID == "" {
			return localization.ErrorWalletOtherServiceIDRequired
		}
	}
	if w.Agent != nil && *w.Agent {
		if w.AgentServiceID == "" {
			return localization.ErrorWalletAgentServiceIDRequired
		}
	}

	if len(errs) > 0 {
		return errs
	}
	w.Name = strings.TrimSpace(w.Name)
	w.UniqueCode = strings.TrimSpace(w.UniqueCode)

	return nil
}

func validateAvatar(file *multipart.FileHeader) error {
	if file == nil {
		return localization.ErrorWalletAvatarInvalid
	}
	if !utils.IsValidImage(file) {
		return errors.New(localization.ErrorWalletAvatarInvalidType.Code)
	}
	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}
	return nil
}

// Helper: check if string is all alphabetic
func isAlpha(s string) bool {
	for _, r := range s {
		if (r < 'A' || r > 'Z') && (r < 'a' || r > 'z') {
			return false
		}
	}
	return true
}
