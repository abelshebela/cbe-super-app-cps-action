package miniapp

import (
	"mime/multipart"
	"time"
)

type EnvironmentType string

const (
	UatEnvironment        EnvironmentType = "UAT"
	DevEnvironment        EnvironmentType = "DEV"
	TestEnvironment       EnvironmentType = "TEST"
	ProductionEnvironment EnvironmentType = "PRODUCTION"
)

type BranchType string

const (
	IFB BranchType = "IFB"
	CB  BranchType = "CB"
)

type AppType string

const (
	URL AppType = "URL"
)

type Stage string

const (
	StageUat Stage = "UAT"
)

type AppViewType string

const (
	AppViewTypeBoth AppViewType = "BOTH"
	AppViewTypeCB   AppViewType = "CB"
	AppViewTypeIFB  AppViewType = "IFB"
)

type ProductCode struct {
	ID             string     `bson:"id" json:"id,omitempty"`
	BranchType     BranchType `bson:"branch_type" json:"branch_type"`
	ProductCode    string     `bson:"product_code" json:"product_code"`
	VATCode        string     `bson:"vat_code" json:"vat_code"`
	ServiceFeeCode string     `bson:"service_fee_code" json:"service_fee_code"`
}

type CredentialInformation struct {
	ID            string          `bson:"id" json:"id,omitempty"`
	Environment   EnvironmentType `bson:"environment" json:"environment"`
	MerchantAppID string          `bson:"merchant_app_id" json:"merchant_app_id"`
	FabricAppID   string          `bson:"fabric_app_id" json:"fabric_app_id"`
	ShortCode     string          `bson:"short_code" json:"short_code"`
	MiniAppCode   string          `bson:"mini_app_code" json:"mini_app_code"`
	AppSecret     string          `bson:"app_secret" json:"app_secret"`
	PrivateKey    string          `bson:"private_key" json:"-"`
	PublicKey     string          `bson:"public_key" json:"public_key"`
	Signature     string          `bson:"signature" json:"-"`
	Timestamp     time.Time       `bson:"timestamp" json:"-"`
}

type MiniApp struct {
	ID                  string                `bson:"_id" json:"id,omitempty"`
	AppName             string                `bson:"app_name" json:"app_name"`
	AppIcon             string                `bson:"app_icon" json:"app_icon"`
	BannerImage         string                `bson:"banner_image" json:"banner_image"`
	CommissionGLAccount string                `bson:"commison_gl_account" json:"commison_gl_account,omitempty"`
	AppType             AppType               `bson:"app_type" json:"app_type"`
	MerchantID          string                `bson:"merchant_id" json:"merchant_id"`
	ProductCode         []ProductCode         `bson:"product_code" json:"product_code"`
	Credential          CredentialInformation `bson:"credential" json:"credential,omitempty"`
	URL                 string                `bson:"url" json:"url,omitempty"`
	AppViewType         AppViewType           `bson:"app_view_type" json:"app_view_type"`
	Stage               Stage                 `bson:"stage" json:"stage"`
	IsEventMiniApp      bool                  `bson:"is_event_mini_app" json:"is_event_mini_app"`
	IsThreeClick        bool                  `bson:"is_three_click" json:"is_three_click"`
	Enabled             bool                  `bson:"enabled" json:"enabled"`
	IsDeleted           bool                  `bson:"is_deleted" json:"is_deleted"`
	CreatedAt           time.Time             `bson:"created_at" json:"created_at"`
	LastModifiedAt      time.Time             `bson:"last_modified_at" json:"last_modified_at"`
	DeletedAt           time.Time             `bson:"deleted_at" json:"deleted_at"`
}

type MiniAppCreateRequest struct {
	ID                  string                `json:"id"`
	AppName             string                `json:"app_name"`
	AppIcon             *multipart.FileHeader `json:"app_icon"`
	BannerImage         *multipart.FileHeader `json:"banner_image"`
	CommissionGLAccount string                `json:"commison_gl_account"`
	AppType             AppType               `json:"app_type"`
	MerchantID          string                `json:"merchant_id"`
	URL                 string                `json:"url"`
	Stage               Stage                 `json:"stage"`
	AppViewType         AppViewType           `json:"app_view_type"`
	ProductCode         []ProductCode         `json:"product_code"`
	Credential          CredentialInformation `json:"credential"`
	IsEventMiniApp      bool                  `json:"is_event_mini_app"`
	IsThreeClick        bool                  `json:"is_three_click"`
}
