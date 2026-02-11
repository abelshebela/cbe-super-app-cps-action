package ecommercemerchant

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MiniAppMerchantMapper(data model.EcommerceMerchant) bson.M {
	update := bson.M{}

	if data.MerchantName != "" {
		update["merchant_name"] = data.MerchantName
	}
	if data.BankAccountNumber != "" {
		update["bank_account_number"] = data.BankAccountNumber
	}
	if data.SettlementMethod != "" {
		update["settlement_method"] = data.SettlementMethod
	}
	if data.Email != "" {
		update["email"] = data.Email
	}
	if data.PhoneNumber != "" {
		update["phone_number"] = data.PhoneNumber
	}

	if len(data.Branches) > 0 {
		update["branches"] = data.Branches
	}

	update["updated_at"] = time.Now()

	return update
}

func ToMiniAppMerchantDomain(miniAppMerchant *model.MiniAppMerchant) *model.MiniAppMerchant {

	return &model.MiniAppMerchant{
		ID:           miniAppMerchant.ID,
		MerchantCode: miniAppMerchant.MerchantCode,
		MerchantName: miniAppMerchant.MerchantName,
		// KYC: types.KYC{
		// 	Status: miniAppMerchant.KYC.Status,
		// 	Representative: types.KYCInformation{
		// 		Name:  miniAppMerchant.KYC.Representative.Name,
		// 		Email: miniAppMerchant.KYC.Representative.Email,
		// 		Phone: miniAppMerchant.KYC.Representative.Phone,
		// 	},
		// },
		BankAccountNumber: miniAppMerchant.BankAccountNumber,

		Enabled:   miniAppMerchant.Enabled,
		IsDeleted: miniAppMerchant.IsDeleted,
		CreatedAt: miniAppMerchant.CreatedAt,
		UpdatedAt: miniAppMerchant.UpdatedAt,
		DeletedAt: miniAppMerchant.DeletedAt,
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
