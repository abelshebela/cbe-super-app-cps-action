package bank_core

import (
	"cbe-super-app-cps-action/internal/constants"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func ParseMultipartFormFile(r *http.Request, key string, maxMemory int64, action string, logger utils.Logger) (multipart.File, *multipart.FileHeader, error) {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		return nil, nil, fmt.Errorf("%s", localization.MsgFileNotFound)
	}
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, nil, fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {

			if action == constants.CREATE {
				if err == http.ErrMissingFile {
					logger.Errorf("Error logo file is missing error: %v", err)
					return nil, nil, fmt.Errorf("%s", localization.MsgFileNotFound)
				}
			}
			logger.Infof("Logo not provided for update - skipping file update")
			return nil, nil, nil
			// else {
			// 	if err == http.ErrMissingFile {
			// 		logger.Infof("Logo is not provided in the update")
			// 	}
			// }
		}

		return nil, nil, fmt.Errorf("missing or invalid file for key '%s': %w", key, err)
	}

	return file, fileHeader, nil
}

func ValidateBankRequest(r *http.Request, data interface{}) localization.ResponseCode {
	var name, code, bic *string

	switch v := data.(type) {
	case *bank_dto.UpdateBankRequest:
		if v == nil {
			return localization.ErrorNoDataProvidedForBankUpdate
		}
		name, code, bic = &v.Name, &v.Code, &v.BIC
	case *bank_dto.CreateBankRequest:
		if v == nil {
			return localization.ErrorNoDataProvidedForBankUpdate
		}
		name, code, bic = &v.Name, &v.Code, &v.BIC
	default:
		return localization.ErrorInvalidBankRequest
	}

	if *name != "" && isInvalidFormat(name) {
		return localization.ErrorInvalidFormatForName
	}
	if *code != "" && isInvalidFormat(code) {
		return localization.ErrorInvalidFormatForCode
	}
	if *bic != "" && isInvalidFormat(bic) {
		return localization.ErrorInvalidFormatForBIC
	}

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
