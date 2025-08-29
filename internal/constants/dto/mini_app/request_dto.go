package miniappdto

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
	"mime/multipart"
)

type EnvironmentType string

const (
	UatEnvironment        EnvironmentType = "UAT"
	DevEnvironment        EnvironmentType = "DEV"
	TestEnvironment       EnvironmentType = "TEST"
	ProductionEnvironment EnvironmentType = "PRODUCTION"
)

// BranchType enum
type BranchType string

const (
	IFB BranchType = "IFB"
	CB  BranchType = "CB"
)

type MiniAppRequest struct {
	AppName             string                `form:"app_name"`
	AppIcon             *multipart.FileHeader `form:"app_icon"`
	BannerImage         *multipart.FileHeader `form:"banner_image"`
	CommissionGLAccount string                `form:"commission_gl_account"`
	MerchantID          string                `form:"merchant_id"`
	IsEventMiniApp      bool                  `form:"is_event_mini_app"`
	IsThreeClick        bool                  `form:"is_three_click"`
	AppViewType         string                `form:"app_view_type"`

	URL string `form:"url"`

	// IFB product codes
	IFBProductCode    string `form:"ifb_product_code"`
	IFBVATCode        string `form:"ifb_vat_code"`
	IFBServiceFeeCode string `form:"ifb_service_fee_code"`

	// CB product codes
	CBProductCode    string `form:"cb_product_code"`
	CBVATCode        string `form:"cb_vat_code"`
	CBServiceFeeCode string `form:"cb_service_fee_code"`
}
type MiniAppCreateRequest struct {
	ID                  string                      `json:"id"`
	AppName             string                      `json:"app_name"`
	AppIcon             *multipart.FileHeader       `json:"app_icon"`
	BannerImage         *multipart.FileHeader       `json:"banner_image"`
	CommissionGLAccount string                      `json:"commison_gl_account"`
	AppType             constants.AppType           `json:"app_type"`
	MerchantID          string                      `json:"merchant_id"`
	URL                 string                      `json:"url"`
	Stage               constants.Stage             `json:"stage"`
	AppViewType         constants.AppViewType       `json:"app_view_type"`
	ProductCode         []types.ProductCode         `json:"product_code"`
	Credential          types.CredentialInformation `json:"credential"`
	IsEventMiniApp      bool                        `json:"is_event_mini_app"`
	IsThreeClick        bool                        `json:"is_three_click"`
}
