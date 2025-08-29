package miniappdto

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"
)

type MiniAppResponse struct {
	ID                string                      `json:"id"`
	AppName           string                      `json:"app_name"`
	AppIcon           string                      `json:"app_icon"`
	BannerImage       string                      `json:"banner_image"`
	CommisonGLAccount string                      `json:"commison_gl_account,omitempty"`
	AppType           constants.AppType               `json:"app_type"`
	MerchantID        string                      `json:"merchant_id"`
	AppViewType       constants.AppViewType           `json:"app_view_type"`
	URL               string                      `json:"url,omitempty"`
	Stage             constants.Stage                 `json:"stage"`
	ProductCode       []types.ProductCode         `json:"product_code"`
	Credential        types.CredentialInformation `json:"credential"`
	IsEventMiniApp    bool                        `json:"is_event_mini_app"`
	IsThreeClick      bool                        `json:"is_three_click"`
	Enabled           bool                        `json:"enabled"`
	CreatedAt         time.Time                   `json:"created_at"`
	LastModifiedAt    time.Time                   `json:"last_modified_at"`
}
