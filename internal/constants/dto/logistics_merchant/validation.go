package logistics_merchant_dto

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"regexp"
	"strings"
)

func Validation(v any) error {
	allowedCharsRegex := regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)
	switch req := v.(type) {
	case CreateLogisticsMerchantRequest:
		if strings.TrimSpace(req.MerchantID) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidMerchantID.Message)
		}
		if strings.TrimSpace(req.SettlementMethod) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidSettlementMethod.Message)
		}
		if !allowedCharsRegex.MatchString(req.SettlementMethod) {
			return errors.New("Settlement method must not contain special characters")
		}
		if strings.TrimSpace(req.MerchantName) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidMerchantName.Message)
		}
		if !allowedCharsRegex.MatchString(req.MerchantName) {
			return errors.New("Merchant name must not contain special characters")
		}
		if strings.TrimSpace(req.BankAccountNumber) == "" {
			return errors.New(localization.ErrorLogisticMerchantInvalidBankAccountNumber.Message)
		}
		if !allowedCharsRegex.MatchString(req.BankAccountNumber) {
			return errors.New("Bank account number must not contain special characters")
		}
		if req.IsLogisticsMerchant == nil || !*req.IsLogisticsMerchant {
			return errors.New(localization.ErrorLogisticMerchantInvalidIsLogisticsMerchant.Message)
		}

	case UpdateLogisticsMerchantRequest:
		if !allowedCharsRegex.MatchString(req.MerchantName) {
			return errors.New("Merchant name must not contain special characters")
		}
		if !allowedCharsRegex.MatchString(req.BankAccountNumber) {
			return errors.New("Bank account number must not contain special characters")
		}
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
		// if req.IsLogisticsMerchant == nil || *req.IsLogisticsMerchant == false {
		// 	return errors.New(localization.ErrorLogisticMerchantInvalidIsLogisticsMerchant.Message)
		// }
	case EnableOrDisableLogisticsMerchantsRequest:
		if len(req.MerchantIDs) == 0 {
			return errors.New(localization.ErrorInvalidInputParameters.Code)
		}
		for _, id := range req.MerchantIDs {
			if strings.TrimSpace(id) == "" {
				return errors.New(localization.ErrorInvalidInputParameters.Code)
			}
		}
	default:
		return errors.New("invalid request type for validation")
	}
	return nil
}
