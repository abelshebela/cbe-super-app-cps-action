package event_merchant_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"strings"
)

func Validation(v any) error {
	switch req := v.(type) {
	case CreateEventMerchantRequest:
		if strings.TrimSpace(req.MerchantID) == "" {
			return errors.New(localization.ErrorEventMerchantInvalidMerchantID.Message)
		}
		if strings.TrimSpace(req.MerchantType) == "" {
			return errors.New(localization.ErrorEventMerchantInvalidMerchantType.Message)
		}
		if strings.TrimSpace(req.SettlementMethod) == "" {
			return errors.New(localization.ErrorEventMerchantInvalidSettlementMethod.Message)
		}
		if strings.TrimSpace(req.MerchantName) == "" {
			return errors.New(localization.ErrorEventMerchantInvalidMerchantName.Message)
		}
		if strings.TrimSpace(req.BankAccountNumber) == "" {
			return errors.New(localization.ErrorEventMerchantInvalidBankAccountNumber.Message)
		}
		// if strings.TrimSpace(req.Email) == "" {
		// 	return errors.New(localization.ErrorEventMerchantInvalidEmail.Message)
		// } else {
		// 	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		// 		return errors.New(localization.ErrorEventMerchantInvalidEmail.Message)
		// 	}
		// }
		// if strings.TrimSpace(req.PhoneNumber) == "" {

		// 	return errors.New(localization.ErrorEventMerchantInvalidPhoneNumber.Message)
		// }
	case UpdateEventMerchantRequest:
		// if strings.TrimSpace(req.Email) != "" {
		// 	return errors.New(localization.ErrorEventMerchantInvalidEmail.Message)
		// } else {
		// 	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		// 		return errors.New(localization.ErrorEventMerchantInvalidEmail.Message)
		// 	}
		// }
		// if strings.TrimSpace(req.PhoneNumber) != "" {
		// 	length := len(req.PhoneNumber)
		// 	if length < 10 || length > 13 {
		// 		return errors.New(localization.ErrorEventMerchantInvalidPhoneNumber.Message)
		// 	}
		// }
	default:
		return errors.New("invalid request type for validation")
	}
	return nil
}
