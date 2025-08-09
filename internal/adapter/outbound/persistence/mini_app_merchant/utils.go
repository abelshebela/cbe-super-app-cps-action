package miniappmerchant

import (
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func buildUpdateSet(m *entities.MiniAppMerchant) bson.M {
	set := bson.M{}

	if m.MerchantType != "" {
		set["merchant_type"] = m.MerchantType
	}
	if m.MerchantName != "" {
		set["merchant_name"] = m.MerchantName
	}
	// Update KYC fields if present
	if m.KYC.Status != "" {
		set["kyc.status"] = m.KYC.Status
	}
	if m.KYC.Representative.Name != "" {
		set["kyc.representative.name"] = m.KYC.Representative.Name
	}
	if m.KYC.Representative.Email != "" {
		set["kyc.representative.email"] = m.KYC.Representative.Email
	}
	if m.KYC.Representative.Phone != "" {
		set["kyc.representative.phone"] = m.KYC.Representative.Phone
	}
	if m.PhoneNumber != "" {
		set["phone_number"] = m.PhoneNumber
	}
	if m.Email != "" {
		set["email"] = m.Email
	}
	if m.BankAccountNumber != "" {
		set["bank_account_number"] = m.BankAccountNumber
	}
	if len(m.Branches) > 0 {
		set["branches"] = m.Branches
	}
	set["enabled"] = m.Enabled

	return set
}
