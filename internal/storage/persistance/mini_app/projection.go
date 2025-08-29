package mini_app

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MiniAppMapper(doc *model.MiniApp) model.MiniApp {
	return *doc
}

func MiniAppDocumentMapper(miniApp model.MiniApp) *model.MiniApp {
	if miniApp.ID == bson.NilObjectID {
		miniApp.ID = bson.NewObjectID()
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
	if miniApp.CommissionGLAccount != "" {
		update["commison_gl_account"] = miniApp.CommissionGLAccount
	}
	if miniApp.AppType != "" {
		update["app_type"] = miniApp.AppType
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
	if !reflect.DeepEqual(miniApp.Credential, types.CredentialInformation{}) {
		cred := bson.M{
			"id": bson.NewObjectID(),
		}
		if miniApp.Credential.Environment != "" {
			cred["environment"] = miniApp.Credential.Environment
		}
		if miniApp.Credential.MerchantAppID != "" {
			cred["merchant_appid"] = miniApp.Credential.MerchantAppID
		}
		if miniApp.Credential.FabricAppID != "" {
			cred["fabric_appid"] = miniApp.Credential.FabricAppID
		}
		if miniApp.Credential.ShortCode != "" {
			cred["short_code"] = miniApp.Credential.ShortCode
		}
		if miniApp.Credential.AppSecret != "" {
			cred["app_secret"] = miniApp.Credential.AppSecret
		}
		if miniApp.Credential.PrivateKey != "" {
			cred["private_key"] = miniApp.Credential.PrivateKey
		}
		if miniApp.Credential.PublicKey != "" {
			cred["public_key"] = miniApp.Credential.PublicKey
		}
		if miniApp.Credential.MiniAppCode != "" {
			cred["miniapp_code"] = miniApp.Credential.MiniAppCode
		}
		if miniApp.Credential.Signature != "" {
			cred["signature"] = miniApp.Credential.Signature
		}
		if !miniApp.Credential.Timestamp.IsZero() {
			cred["timestamp"] = miniApp.Credential.Timestamp
		}
		update["credential"] = cred
	}
	if miniApp.MerchantID != "" {
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
	if miniApp.IsEventMiniApp {
		update["is_event_mini_app"] = miniApp.IsEventMiniApp
	}
	if miniApp.IsThreeClick {
		update["is_three_click"] = miniApp.IsThreeClick
	}

	return update
}
