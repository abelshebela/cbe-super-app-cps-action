package model

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniApp struct {
	ID                  bson.ObjectID               `bson:"_id" json:"id"`
	CategoryID          bson.ObjectID               `bson:"category_id" json:"category_id"`
	AppName             string                      `bson:"app_name" json:"app_name"`
	AppIcon             string                      `bson:"app_icon" json:"app_icon"`
	AppCode             string                      `bson:"app_code" json:"app_code"`
	BannerImage         string                      `bson:"banner_image" json:"banner_image"`
	CommissionGLAccount string                      `bson:"commison_gl_account" json:"commison_gl_account,omitempty"`
	AppType             constants.AppType           `bson:"app_type" json:"app_type"`
	MerchantID          string                      `bson:"merchant_id" json:"merchant_id"`
	ProductCode         []types.ProductCode         `bson:"product_code" json:"product_code"`
	Credential          types.CredentialInformation `bson:"credential" json:"credential,omitempty"`
	URL                 string                      `bson:"url" json:"url,omitempty"`
	AppViewType         constants.AppViewType       `bson:"app_view_type" json:"app_view_type"`
	Stage               constants.Stage             `bson:"stage" json:"stage"`
	IsEventMiniApp      bool                        `bson:"is_event_mini_app" json:"is_event_mini_app"`
	IsThreeClick        bool                        `bson:"is_three_click" json:"is_three_click"`
	Enabled             bool                        `bson:"enabled" json:"enabled"`
	IsDeleted           bool                        `bson:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time                   `bson:"created_at" json:"created_at"`
	LastModifiedAt      time.Time                   `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt           time.Time                   `bson:"deleted_at" json:"deleted_at"`
}

type MiniAppCategory struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	Icon      string        `bson:"icon" json:"icon"`
	IsDeleted bool          `bson:"is_deleted" json:"is_deleted"`
	IsEnabled bool          `bson:"is_enabled" json:"is_enabled"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}
