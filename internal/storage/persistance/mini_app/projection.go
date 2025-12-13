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

	if miniApp.MerchantID != bson.NilObjectID {
		update["merchant_id"] = miniApp.MerchantID
	}
	if miniApp.AppViewType != "" {
		update["app_view_type"] = miniApp.AppViewType
	}
	if miniApp.URL != "" {
		update["url"] = miniApp.URL
	}
	if miniApp.CommissionGLAccount != "" {
		update["commission_gl_account"] = miniApp.CommissionGLAccount
	}
	update["is_featured"] = miniApp.IsFeatured
	update["enabled"] = miniApp.Enabled
	if miniApp.AppCode != "" {
		update["app_code"] = miniApp.AppCode
	}

	return update
}
