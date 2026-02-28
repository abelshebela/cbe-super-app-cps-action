package event_merchant_dto

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"errors"
	"strings"
)

func (req CreateEventMerchantRequest) Validate() error {
	if strings.TrimSpace(req.MerchantID) == "" {
		return errors.New(localization.ErrorEventMerchantInvalidMerchantID.Code)
	}
	if strings.TrimSpace(req.MerchantType) == "" {
		return errors.New(localization.ErrorEventMerchantInvalidMerchantType.Code)
	}
	if strings.TrimSpace(req.SettlementMethod) == "" {
		return errors.New(localization.ErrorEventMerchantInvalidSettlementMethod.Code)
	}
	if strings.TrimSpace(req.MerchantName) == "" {
		return errors.New(localization.ErrorEventMerchantInvalidMerchantName.Code)
	}
	if strings.TrimSpace(req.BankAccountNumber) == "" {
		return errors.New(localization.ErrorEventMerchantInvalidBankAccountNumber.Code)
	}
	if req.IsEventMerchant == nil || *req.IsEventMerchant == false {
		return errors.New(localization.ErrorEventMerchantInvalidIsEventMerchant.Code)
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
	return nil
}

func (req UpdateEventMerchantRequest) Validate() error {
	if strings.TrimSpace(req.MerchantID) != "" && len(strings.TrimSpace(req.MerchantID)) == 0 {
		return errors.New(localization.ErrorEventMerchantInvalidMerchantID.Code)
	}
	if strings.TrimSpace(req.MerchantType) != "" && len(strings.TrimSpace(req.MerchantType)) == 0 {
		return errors.New(localization.ErrorEventMerchantInvalidMerchantType.Code)
	}
	if strings.TrimSpace(req.SettlementMethod) != "" && len(strings.TrimSpace(req.SettlementMethod)) == 0 {
		return errors.New(localization.ErrorEventMerchantInvalidSettlementMethod.Code)
	}
	if strings.TrimSpace(req.MerchantName) != "" && len(strings.TrimSpace(req.MerchantName)) == 0 {
		return errors.New(localization.ErrorEventMerchantInvalidMerchantName.Code)
	}
	if strings.TrimSpace(req.BankAccountNumber) != "" && len(strings.TrimSpace(req.BankAccountNumber)) == 0 {
		return errors.New(localization.ErrorEventMerchantInvalidBankAccountNumber.Code)
	}
	if req.IsEventMerchant == nil || *req.IsEventMerchant == false {
		return errors.New(localization.ErrorEventMerchantInvalidIsEventMerchant.Code)
	}
	// if strings.TrimSpace(req.Email) != "" {
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
	return nil
}
