package miniappdto

import (
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type MiniAppResponse struct {
	ID                bson.ObjectID               `json:"id" bson:"_id"`
	AppName           string                      `json:"app_name" bson:"app_name"`
	AppIcon           string                      `json:"app_icon" bson:"app_icon"`
	BannerImage       string                      `json:"banner_image" bson:"banner_image"`
	CommisonGLAccount string                      `json:"commison_gl_account" bson:"commison_gl_account"`
	AppType           constants.AppType           `json:"app_type" bson:"app_type"`
	AppViewType       constants.AppViewType       `json:"app_view_type" bson:"app_view_type"`
	URL               string                      `json:"url" bson:"url"`
	Stage             constants.Stage             `json:"stage" bson:"stage"`
	ProductCode       []types.ProductCode         `json:"product_code" bson:"product_code"`
	Credential        types.CredentialInformation `json:"credential" bson:"credential"`
	IsEventMiniApp    bool                        `json:"is_event_mini_app" bson:"is_event_mini_app"`
	IsThreeClick      bool                        `json:"is_three_click" bson:"is_three_click"`
	Enabled           bool                        `json:"enabled" bson:"enabled"`
	CreatedAt         time.Time                   `json:"created_at" bson:"created_at"`
	LastModifiedAt    time.Time                   `json:"last_modified_at" bson:"last_modified_at"`
	Merchant          Merchant                    `json:"merchant" bson:"merchant"`
}

type Merchant struct {
	ID           bson.ObjectID `json:"id" bson:"_id"`
	MerchantName string        `json:"merchant_name" bson:"merchant_name"`
}
