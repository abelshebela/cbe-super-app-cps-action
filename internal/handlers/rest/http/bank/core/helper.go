package bank_core

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64, action string, logger utils.Logger) (multipart.File, *multipart.FileHeader, error) {

	file, fileHeader, err := local_util.ParseMultipartFormFile(r, "logo", int64(constants.MaxMemoryForUpload))
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {

			if action == constants.CREATE {
				if err == http.ErrMissingFile {
					logger.Errorf("[BankHelper][ParseFile] logo missing: %v", err)
					return nil, nil, fmt.Errorf("%s", localization.MsgFileNotFound)
				}
			}
			logger.Infof("[BankHelper][ParseFile] no logo, skipping")
			return nil, nil, nil
		}

		return nil, nil, fmt.Errorf("missing or invalid file for key '%s': %w", key, err)
	}

	return file, fileHeader, nil
}

func ValidateBankRequest(r *http.Request, data interface{}) localization.ResponseCode {
	var bic *string

	switch v := data.(type) {
	case *bank_dto.UpdateBankRequest:
		if v == nil {
			return localization.ErrorNoDataProvidedForBankUpdate
		}
		if err := v.Validate(); err != nil {
			return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
		}
		bic = &v.BICCode
	case *bank_dto.CreateBankRequest:
		if v == nil {
			return localization.ErrorNoDataProvidedForBankUpdate
		}
		if err := v.Validate(); err != nil {
			return localization.ErrorToResponseCode(err.Error(), 400, err.Error())
		}
		bic = &v.BICCode
	default:
		return localization.ErrorInvalidBankRequest
	}

	// if *name != "" && isInvalidFormat(name) {
	// 	return localization.ErrorInvalidFormatForName
	// }

	if *bic != "" && isInvalidFormat(bic) {
		return localization.ErrorInvalidFormatForBIC
	}

	// if bankType != nil && *bankType != "" {
	// 	switch *bankType {
	// 	case constants.Bank, constants.Wallet, constants.MFI:
	// 		// valid type
	// 	default:
	// 		return localization.ErrorInvalidBankRequest
	// 	}
	// }

	return localization.ResponseCode{}
}

// Helper function to check for invalid formats
func isInvalidFormat(s *string) bool {
	if s == nil {
		return false
	}
	str := strings.TrimSpace(*s)
	if str == "" {
		return true
	}
	for _, r := range str {
		if !(r == ' ' || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')) {
			return true
		}
	}
	return false
}
