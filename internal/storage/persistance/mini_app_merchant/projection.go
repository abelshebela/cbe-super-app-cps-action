package mini_app_merchant

import (
shared_type "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

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
	if data.KYC != (shared_type.KYC{}) {
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
		KYC: shared_type.KYC{
			Status: miniAppMerchant.KYC.Status,
			Representative: shared_type.KYCInformation{
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
func convertBranchesModelToDomain(branches []shared_type.BranchInformation) []shared_type.BranchInformation {
	result := make([]shared_type.BranchInformation, len(branches))
	for i, b := range branches {
		result[i] = shared_type.BranchInformation{
			BranchCode:          b.BranchCode,
			BranchName:          b.BranchName,
			BranchAddress:       b.BranchAddress,
			BranchOwner:         b.BranchOwner,
			BranchAccountNumber: b.BranchAccountNumber,
		}
	}
	return result
}
