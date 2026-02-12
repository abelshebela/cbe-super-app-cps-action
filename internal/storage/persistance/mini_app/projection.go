package mini_app

import (
	// "cbe-super-app-cps-action/internal/constants/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"

	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MiniAppDocumentMapper(miniApp mini_model.MiniApp) *mini_model.MiniApp {
	if miniApp.ID == bson.NilObjectID {
		miniApp.ID = bson.NewObjectID()
	}
	if miniApp.CreatedAt.IsZero() {
		miniApp.CreatedAt = time.Now()
	}
	if miniApp.UpdatedAt.IsZero() {
		miniApp.UpdatedAt = time.Now()
	}
	if miniApp.IsDeleted {
		miniApp.IsDeleted = false
	}
	if miniApp.Enabled {
		miniApp.Enabled = false
	}
	return &miniApp
}

func MiniAppDocumentToBsonM(miniApp mini_model.MiniApp) bson.M {
	update := bson.M{
		"last_modified_at": time.Now(),
	}

	if miniApp.AppName != "" {
		update["app_name"] = miniApp.AppName
	}
	if miniApp.AppIcon != "" {
		update["app_icon"] = miniApp.AppIcon
	}
	if miniApp.BannerImage != "" {
		update["banner_image"] = miniApp.BannerImage
	}

	if miniApp.AppType != "" {
		update["app_type"] = miniApp.AppType
	}

	if miniApp.CategoryID != bson.NilObjectID {
		update["category_id"] = miniApp.CategoryID
	}

	if miniApp.MerchantID != bson.NilObjectID {
		update["merchant_id"] = miniApp.MerchantID
	}
	if miniApp.AppViewType != "" {
		update["app_view_type"] = miniApp.AppViewType
	}
	if miniApp.ServiceCode != "" {
		update["service_code"] = miniApp.ServiceCode
	}
	if miniApp.ServiceKey != "" {
		update["service_key"] = miniApp.ServiceKey
	}
	if miniApp.URL != "" {
		update["url"] = miniApp.URL
	}
	update["is_featured"] = miniApp.IsFeatured
	update["enabled"] = miniApp.Enabled

	return update
}

func MiniAppMerchantMapper(data mini_model.MiniAppMerchant) bson.M {
	result := bson.M{}

	if data.MerchantCode != "" {
		result["merchant_code"] = data.MerchantCode
	}
	if data.MerchantName != "" {
		result["merchant_name"] = data.MerchantName
	}
	if data.KYC != (types.KYC{}) {
		result["kyc"] = data.KYC
	}
	if data.BankAccountNumber != "" {
		result["bank_account_number"] = data.BankAccountNumber
	}

	if data.SettlementMethod != "" {
		result["settlement_method"] = data.SettlementMethod
	}

	result["enabled"] = data.Enabled
	result["is_deleted"] = data.IsDeleted

	return result
}

func ToMiniAppMerchantDomain(miniAppMerchant *mini_model.MiniAppMerchant) *mini_model.MiniAppMerchant {

	return &mini_model.MiniAppMerchant{
		ID:           miniAppMerchant.ID,
		MerchantCode: miniAppMerchant.MerchantCode,
		MerchantName: miniAppMerchant.MerchantName,
		KYC: types.KYC{
			Status: miniAppMerchant.KYC.Status,
			Representative: types.KYCInformation{
				Name:  miniAppMerchant.KYC.Representative.Name,
				Email: miniAppMerchant.KYC.Representative.Email,
				Phone: miniAppMerchant.KYC.Representative.Phone,
			},
		},
		BankAccountNumber: miniAppMerchant.BankAccountNumber,

		Enabled:   miniAppMerchant.Enabled,
		IsDeleted: miniAppMerchant.IsDeleted,
		CreatedAt: miniAppMerchant.CreatedAt,
		UpdatedAt: miniAppMerchant.UpdatedAt,
		DeletedAt: miniAppMerchant.DeletedAt,
	}
}
