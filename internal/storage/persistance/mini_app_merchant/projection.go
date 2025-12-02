package mini_app_merchant

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"

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
	if data.KYC != (types.KYC{}) {
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
	// mini apps are not embedded in merchant

	// MUST include enabled and is_deleted
	result["enabled"] = data.Enabled
	result["is_deleted"] = data.IsDeleted

	return result
}

func ToMiniAppMerchantDomain(miniAppMerchant *model.MiniAppMerchant) *model.MiniAppMerchant {

	return &model.MiniAppMerchant{
		ID:           miniAppMerchant.ID,
		Code:         miniAppMerchant.Code,
		MerchantName: miniAppMerchant.MerchantName,
		MerchantType: miniAppMerchant.MerchantType,
		KYC: types.KYC{
			Status: miniAppMerchant.KYC.Status,
			Representative: types.KYCInformation{
				Name:  miniAppMerchant.KYC.Representative.Name,
				Email: miniAppMerchant.KYC.Representative.Email,
				Phone: miniAppMerchant.KYC.Representative.Phone,
			},
		},
		BankAccountNumber: miniAppMerchant.BankAccountNumber,
		Branches:          convertBranchesModelToDomain(miniAppMerchant.Branches),
		Email:             miniAppMerchant.Email,
		PhoneNumber:       miniAppMerchant.PhoneNumber,

		Enabled:        miniAppMerchant.Enabled,
		IsDeleted:      miniAppMerchant.IsDeleted,
		CreatedAt:      miniAppMerchant.CreatedAt,
		LastModifiedAt: miniAppMerchant.LastModifiedAt,
		DeletedAt:      miniAppMerchant.DeletedAt,
	}
}
func convertBranchesModelToDomain(branches []types.BranchInformation) []types.BranchInformation {
	result := make([]types.BranchInformation, len(branches))
	for i, b := range branches {
		result[i] = types.BranchInformation{
			BranchCode:          b.BranchCode,
			BranchName:          b.BranchName,
			BranchAddress:       b.BranchAddress,
			BranchOwner:         b.BranchOwner,
			BranchAccountNumber: b.BranchAccountNumber,
		}
	}
	return result
}
