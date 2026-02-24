package core

import (
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func LogisticsMerchantMapper(merchant local_model.LogisticsMerchant) bson.M {
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
	return result
}
