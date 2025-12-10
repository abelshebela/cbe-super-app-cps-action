package mini_app

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MiniAppDocumentMapper(miniApp model.MiniApp) *model.MiniApp {
	if miniApp.ID == bson.NilObjectID {
		miniApp.ID = bson.NewObjectID()
	}
	if miniApp.CreatedAt.IsZero() {
		miniApp.CreatedAt = time.Now()
	}
	if miniApp.LastModifiedAt.IsZero() {
		miniApp.LastModifiedAt = time.Now()
	}
	if miniApp.IsDeleted {
		miniApp.IsDeleted = false
	}
	if miniApp.Enabled {
		miniApp.Enabled = false
	}
	return &miniApp
}

func MiniAppDocumentToBsonM(miniApp model.MiniApp) bson.M {
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

	if len(miniApp.ProductCode) > 0 {
		productCodes := make([]bson.M, len(miniApp.ProductCode))
		for i, pc := range miniApp.ProductCode {
			productCodes[i] = bson.M{
				"id":               pc.ID,
				"branch_type":      pc.BranchType,
				"product_code":     pc.ProductCode,
				"vat_code":         pc.VATCode,
				"service_fee_code": pc.ServiceFeeCode,
			}
		}
		update["product_code"] = productCodes
	}

	cred := miniApp.Credential
	if cred.Environment != "" || cred.MerchantAppID != "" || cred.FabricAppID != "" ||
		cred.ShortCode != "" || cred.AppSecret != "" || cred.PrivateKey != "" ||
		cred.PublicKey != "" || cred.MiniAppCode != "" || cred.Signature != "" ||
		!cred.Timestamp.IsZero() {

		credMap := bson.M{"id": bson.NewObjectID()}
		if cred.Environment != "" {
			credMap["environment"] = cred.Environment
		}
		if cred.MerchantAppID != "" {
			credMap["merchant_appid"] = cred.MerchantAppID
		}
		if cred.FabricAppID != "" {
			credMap["fabric_appid"] = cred.FabricAppID
		}
		if cred.ShortCode != "" {
			credMap["short_code"] = cred.ShortCode
		}
		if cred.AppSecret != "" {
			credMap["app_secret"] = cred.AppSecret
		}
		if cred.PrivateKey != "" {
			credMap["private_key"] = cred.PrivateKey
		}
		if cred.PublicKey != "" {
			credMap["public_key"] = cred.PublicKey
		}
		if cred.MiniAppCode != "" {
			credMap["miniapp_code"] = cred.MiniAppCode
		}
		if cred.Signature != "" {
			credMap["signature"] = cred.Signature
		}
		if !cred.Timestamp.IsZero() {
			credMap["timestamp"] = cred.Timestamp
		}
		update["credential"] = credMap
	}

	if miniApp.MerchantID != bson.NilObjectID {
		update["merchant_id"] = miniApp.MerchantID
	}
	if miniApp.AppViewType != "" {
		update["app_view_type"] = miniApp.AppViewType
	}
	if miniApp.URL != "" {
		update["url"] = miniApp.URL
	}
	if miniApp.Stage != "" {
		update["stage"] = miniApp.Stage
	}

	if miniApp.AppMode != "" {
		update["app_mode"] = miniApp.AppMode
	}

	return update
}
