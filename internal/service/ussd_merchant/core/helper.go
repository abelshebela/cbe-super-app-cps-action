package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"
)

// ExistingIdentifier checks whether any identifier (email, phone_number, account_number)
// in the request already exists in the provided result set. Returns a specific
// localization error code when a duplicate is found, or nil otherwise.
func ExistingIdentifier(existing *imodel.UssdMerchant, req ussd_merchant_dto.CreateUssdMerchantRequest) error {
	if existing == nil {
		return nil
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	normalizedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	normalizedAcct := strings.TrimSpace(req.AccountNumber)

	// Email duplicate
	if normalizedEmail != "" && strings.EqualFold(strings.TrimSpace(existing.Email), normalizedEmail) {
		return errors.New(localization.ErrorEmailAlreadyExist.Code)
	}
	// Phone duplicate
	if normalizedPhone != "" {
		// Normalize stored phone too, just in case
		storedPhone := local_util.FormatPhoneNumber(existing.PhoneNumber)
		if storedPhone == normalizedPhone {
			return errors.New(localization.ErrorPhonenumberAlreadyExist.Code)
		}
	}
	// Account number duplicate

	if normalizedAcct != "" && strings.EqualFold(strings.TrimSpace(existing.AccountNumber), normalizedAcct) {
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	return nil
}

func ExistingIdentifierForUpdate(existing imodel.UssdMerchant, id string, req ussd_merchant_dto.UpdateUssdMerchantRequest) error {
	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	normalizedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	normalizedAcct := strings.TrimSpace(req.AccountNumber)

	existingID := existing.ID
	// Email duplicate
	if normalizedEmail != "" && strings.EqualFold(existing.Email, normalizedEmail) && existingID != id {
		return errors.New(localization.ErrorEmailAlreadyExist.Code)
	}
	// Phone duplicate
	if normalizedPhone != "" && existingID != id {
		// Normalize stored phone too, just in case
		storedPhone := local_util.FormatPhoneNumber(existing.PhoneNumber)
		if storedPhone == normalizedPhone {
			return errors.New(localization.ErrorPhonenumberAlreadyExist.Code)
		}
	}
	// Account number duplicate
	if normalizedAcct != "" && strings.EqualFold(strings.TrimSpace(existing.AccountNumber), normalizedAcct) && existingID != id {
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	return nil
}

func UssdMerchant(req ussd_merchant_dto.CreateUssdMerchantRequest) imodel.UssdMerchant {
	return imodel.UssdMerchant{
		MerchantCode:     local_util.UniqueIdGenerator(),
		Name:             req.Name,
		SettlementMethod: constants.SettlementMethod(req.SettlementMethod),
		PhoneNumber:      req.PhoneNumber,
		AccountNumber:    req.AccountNumber,
		Email:            req.Email,
		Service:          req.Service,
	}
}

func UssdMerchantUpdate(req ussd_merchant_dto.UpdateUssdMerchantRequest, ussdMerchant *ussd_merchant_dto.UssdMerchantResponse) imodel.UssdMerchant {

	return imodel.UssdMerchant{
		MerchantCode:     local_util.UniqueIdGenerator(),
		Name:             req.Name,
		SettlementMethod: constants.SettlementMethod(req.SettlementMethod),
		PhoneNumber:      req.PhoneNumber,
		AccountNumber:    req.AccountNumber,
		Email:            req.Email,
		Service:          req.Service,
	}
}

func CreateCredentials(merchant *imodel.UssdMerchant, cfg config.VaultConfig) error {
	cred, _, err := local_util.LocalEncryptPassword(merchant.PhoneNumber+local_util.UniqueIdGenerator(), constants.Cred, constants.Empty, constants.Empty, &cfg)
	if err != nil {
		return err
	}

	merchant.Credential = cred
	return nil
}
