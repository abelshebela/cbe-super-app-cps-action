package core

import (
	"cbe-super-app-cps-action/internal/constants"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"go.mongodb.org/mongo-driver/v2/bson"

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
	if strings.EqualFold(strings.TrimSpace(existing.AccountNumber), normalizedAcct) {
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	return nil
}

func ExistingIdentifierForUpdate(existing imodel.UssdMerchant, id string, req ussd_merchant_dto.UpdateUssdMerchantRequest) error {
	if &existing == nil {
		return nil
	}

	normalizedEmail := strings.TrimSpace(strings.ToLower(req.Email))
	normalizedPhone := local_util.FormatPhoneNumber(req.PhoneNumber)
	normalizedAcct := strings.TrimSpace(req.AccountNumber)

	existingID := local_util.FirstHex24(existing.ID.String())
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

func ModelToBson(data *imodel.UssdMerchant) bson.M {
	update := bson.M{}

	if data.Name != "" {
		update["name"] = data.Name
	}
	if data.Logo != "" {
		update["logo"] = data.Logo
	}
	if data.SettlementMethod != "" {
		update["settlement_method"] = data.SettlementMethod
	}
	if data.PhoneNumber != "" {
		update["phone_number"] = data.PhoneNumber
	}
	if data.Email != "" {
		update["email"] = data.Email
	}
	if data.Service != "" {
		update["service"] = data.Service
	}
	if data.AccountNumber != "" {
		update["account_number"] = data.AccountNumber
	}
	return update
}

func CreateCredentials(ussd_marchant *imodel.UssdMerchant, cfg config.VaultConfig) error {

	cred, _, err := local_util.LocalEncryptPassword(ussd_marchant.PhoneNumber+local_util.UniqueIdGenerator(), constants.Cred, constants.Empty, constants.Empty, &cfg)
	if err != nil {
		return err
	}

	ussd_marchant.Credential = cred
	return nil
}
