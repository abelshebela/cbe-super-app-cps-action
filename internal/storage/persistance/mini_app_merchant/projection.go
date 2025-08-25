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
	// KYC is a struct, include if not zero value
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
	result["enabled"] = data.Enabled

	return result
}
func ToMiniAppMerchantDomain(miniAppMerchant *model.MiniAppMerchant) *model.MiniAppMerchant {

	miniApps := []model.MiniApps{}
	for _, app := range miniAppMerchant.MiniApps {

		res := model.MiniApps{
			ID:        app.ID,
			Enabled:   app.Enabled,
			IsDeleted: app.IsDeleted,
		}

		miniApps = append(miniApps, res)
	}
	return &model.MiniAppMerchant{
		ID:           miniAppMerchant.ID,
		Code:         miniAppMerchant.Code,
		MerchantName: miniAppMerchant.MerchantName,
		MerchantType: miniAppMerchant.MerchantType,
		KYC: model.KYC{
			Status: miniAppMerchant.KYC.Status,
			Representative: model.KYCInformation{
				Name:  miniAppMerchant.KYC.Representative.Name,
				Email: miniAppMerchant.KYC.Representative.Email,
				Phone: miniAppMerchant.KYC.Representative.Phone,
			},
		},
		BankAccountNumber: miniAppMerchant.BankAccountNumber,
		Branches:          convertBranchesModelToDomain(miniAppMerchant.Branches),
		Email:             miniAppMerchant.Email,
		PhoneNumber:       miniAppMerchant.PhoneNumber,
		MiniApps:          miniApps,
		Enabled:           miniAppMerchant.Enabled,
		IsDeleted:         miniAppMerchant.IsDeleted,
		CreatedAt:         miniAppMerchant.CreatedAt,
		LastModifiedAt:    miniAppMerchant.LastModifiedAt,
		DeletedAt:         miniAppMerchant.DeletedAt,
	}
}
func convertBranchesModelToDomain(branches []model.BranchInformation) []model.BranchInformation {
	result := make([]model.BranchInformation, len(branches))
	for i, b := range branches {
		result[i] = model.BranchInformation{
			BranchCode:          b.BranchCode,
			BranchName:          b.BranchName,
			BranchAddress:       b.BranchAddress,
			BranchOwner:         b.BranchOwner,
			BranchAccountNumber: b.BranchAccountNumber,
		}
	}
	return result
}
