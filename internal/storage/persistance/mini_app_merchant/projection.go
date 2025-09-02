package mini_app_merchant

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MiniAppMerchantMapper maps MiniAppMerchant model to BSON for database operations
func MiniAppMerchantMapper(data model.MiniAppMerchant) bson.M {
	result := bson.M{}

	if data.Code != "" {
		result["merchant_code"] = data.Code
	}
	if data.MerchantName != "" {
		result["merchant_name"] = data.MerchantName
	}
	if data.MerchantType != "" {
		result["merchant_type"] = data.MerchantType
	}
	if data.KYC != (model.KYC{}) {
		result["kyc"] = data.KYC
	}
	if data.BankAccountNumber != "" {
		result["bank_account_number"] = data.BankAccountNumber
	}
	if len(data.Branches) > 0 {
		result["branches"] = data.Branches
	}
	if data.Email != "" {
		result["email"] = data.Email
	}
	if data.PhoneNumber != "" {
		result["phone_number"] = data.PhoneNumber
	}
	if len(data.MiniApps) > 0 {
		result["mini_apps"] = data.MiniApps
	}

	// Only include enabled/is_deleted if explicitly set (optional)
	// Comment these lines if you don't want to overwrite them every update
	result["enabled"] = data.Enabled
	result["is_deleted"] = data.IsDeleted

	return result
}
