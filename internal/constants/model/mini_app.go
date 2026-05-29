package model

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniApp struct {
	ID                  bson.ObjectID               `bson:"_id" json:"id"`
	CategoryID          bson.ObjectID               `bson:"category_id" json:"category_id"`
	ProductCodeID       bson.ObjectID               `bson:"product_code_id" json:"product_code_id"`
	MerchantID          bson.ObjectID               `bson:"merchant_id" json:"merchant_id"`
	AppCode             string                      `bson:"app_code" json:"app_code"`
	AppName             string                      `bson:"app_name" json:"app_name"`
	AppIcon             string                      `bson:"app_icon" json:"app_icon"`
	AppType             string                      `bson:"app_type" json:"app_type"`
	BannerImage         string                      `bson:"banner_image" json:"banner_image"`
	URL                 string                      `bson:"url" json:"url"`
	AppViewType         string                      `bson:"app_view_type" json:"app_view_type"`
	CommissionGLAccount string                      `bson:"commission_gl_account" json:"commission_gl_account"`
	IsFeatured          bool                        `bson:"is_featured" json:"is_featured"`
	Enabled             bool                        `bson:"enabled" json:"enabled"`
	IsDeleted           bool                        `bson:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time                   `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time                   `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt           *time.Time                  `bson:"deleted_at" json:"deleted_at"`
	Credential          types.CredentialInformation `bson:"credential" json:"credential"`
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
