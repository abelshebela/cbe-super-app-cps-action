package logistics_merchant_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"strings"
)

func Validation(v any) error {
	switch req := v.(type) {
	case CreateLogisticsMerchantRequest:
		if strings.TrimSpace(req.MerchantID) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidMerchantID.Message)
		}
		if strings.TrimSpace(req.MerchantType) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidMerchantType.Message)
		}
		if strings.TrimSpace(req.SettlementMethod) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidSettlementMethod.Message)
		}
		if strings.TrimSpace(req.MerchantName) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidMerchantName.Message)
		}
		if strings.TrimSpace(req.BankAccountNumber) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidBankAccountNumber.Message)
		}
		// if strings.TrimSpace(req.Email) == "" {
		// 	return errors.New(localization.ErrorLogisticMerchantInvalidEmail.Message)
		// } else {
		// 	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		// 		return errors.New(localization.ErrorLogisticMerchantInvalidEmail.Message)
		// 	}
		// }
		// if strings.TrimSpace(req.PhoneNumber) == "" {

		// 	return errors.New(localization.ErrorLogisticMerchantInvalidPhoneNumber.Message)
		// }
	case UpdateLogisticsMerchantRequest:
		// if strings.TrimSpace(req.Email) != "" {
		// 	return errors.New(localization.ErrorLogisticMerchantInvalidEmail.Message)
		// } else {
		// 	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		// 		return errors.New(localization.ErrorLogisticMerchantInvalidEmail.Message)
		// 	}
		// }
		// if strings.TrimSpace(req.PhoneNumber) != "" {
		// 	length := len(req.PhoneNumber)
		// 	if length < 10 || length > 13 {
		// 		return errors.New(localization.ErrorLogisticMerchantInvalidPhoneNumber.Message)
		// 	}
		// }
	default:
		return errors.New("invalid request type for validation")
	}
	return nil
}
