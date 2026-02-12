package core

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func EventMerchantMapper(merchant model.EventMerchant) bson.M {
	result := bson.M{}

	if merchant.MerchantID != "" {
		result["merchant_id"] = merchant.MerchantID
	}
	if merchant.MerchantType != "" {
		result["merchant_type"] = merchant.MerchantType
	}
	if merchant.SettlementMethod != "" {
		result["settlement_method"] = merchant.SettlementMethod
	}
	if merchant.MerchantName != "" {
		result["merchant_name"] = merchant.MerchantName
	}
	if merchant.BankAccountNumber != "" {
		result["bank_account_number"] = merchant.BankAccountNumber
	}
	if merchant.Email != "" {
		result["email"] = merchant.Email
	}
	if merchant.PhoneNumber != "" {
		result["phone_number"] = merchant.PhoneNumber
	}
	return result
}
