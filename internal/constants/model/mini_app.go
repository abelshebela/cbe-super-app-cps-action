package model

import (
	// "cbe-super-app-cps-action/internal/constants/types"
	"time"

	// "go.mongodb.org/mongo-driver/v2/bson"
)

// type MiniApp struct {
// 	ID                  bson.ObjectID               `bson:"_id" json:"id"`
// 	CategoryID          bson.ObjectID               `bson:"category_id" json:"category_id"`
// 	ProductCodeID       bson.ObjectID               `bson:"product_code_id" json:"product_code_id"`
// 	MerchantID          bson.ObjectID               `bson:"merchant_id" json:"merchant_id"`
// 	AppCode             string                      `bson:"app_code" json:"app_code"`
// 	AppName             string                      `bson:"app_name" json:"app_name"`
// 	AppIcon             string                      `bson:"app_icon" json:"app_icon"`
// 	AppType             string                      `bson:"app_type" json:"app_type"`
// 	BannerImage         string                      `bson:"banner_image" json:"banner_image"`
// 	URL                 string                      `bson:"url" json:"url"`
// 	AppViewType         string                      `bson:"app_view_type" json:"app_view_type"`
// 	CommissionGLAccount string                      `bson:"commission_gl_account" json:"commission_gl_account"`
// 	IsFeatured          bool                        `bson:"is_featured" json:"is_featured"`
// 	Enabled             bool                        `bson:"enabled" json:"enabled"`
// 	IsDeleted           bool                        `bson:"is_deleted" json:"is_deleted"`
// 	CreatedAt           time.Time                   `bson:"created_at" json:"created_at"`
// 	UpdatedAt           time.Time                   `bson:"last_modified_at" json:"last_modified_at"`
// 	DeletedAt           *time.Time                  `bson:"deleted_at" json:"deleted_at"`
// 	Credential          types.CredentialInformation `bson:"credential" json:"credential"`
// }


type CredentialInformation struct {
	MerchantAppID string `json:"merchant_app_id"`
	FabricAppID   string `json:"fabric_app_id"`
	ShortCode     string `json:"short_code"`
	AppSecret     string `json:"app_secret"`
	PrivateKey    string `json:"private_key"`
	PublicKey     string `json:"public_key"`
}

type AppType string

const (
	Financial    AppType = "FINANCIAL"
	NonFinancial AppType = "NON_FINANCIAL"
)

type AppViewType string

const (
	AppViewTypeBoth AppViewType = "BOTH"
	AppViewTypeCB   AppViewType = "CB"
	AppViewTypeIFB  AppViewType = "IFB"
)
type MiniApp struct {
	ID          string                `json:"id"`
	CategoryID  string                `json:"category_id"`
	MerchantID  string                `json:"merchant_id"`
	ServiceID   string                `json:"service_id"`
	ServiceKey  string                `json:"service_key"`
	ServiceCode string                `json:"service_code"`
	AppName     string                `json:"app_name"`
	AppIcon     string                `json:"app_icon"`
	AppType     AppType               `json:"app_type"`
	BannerImage string                `json:"banner_image"`
	URL         string                `json:"url"`
	AppViewType AppViewType           `json:"app_view_type"`
	IsFeatured  bool                  `json:"is_featured"`
	Enabled     bool                  `json:"enabled"`
	IsDeleted   bool                  `json:"is_deleted"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"last_modified_at"`
	DeletedAt   *time.Time            `json:"deleted_at,omitempty"`
	Credential  CredentialInformation `json:"credential"`
	Caps        []Cap                 `json:"cap,omitempty"`
}

type MiniAppCategory struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	IsDeleted bool       `json:"is_deleted"`
	IsEnabled bool       `json:"is_enabled"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}



type SettlementMethod string

const (
	SettlementMethodDirect       SettlementMethod = "DIRECT"
	SettlementMethodGL           SettlementMethod = "GL"
	SettlementMethodMultiAccount SettlementMethod = "MULTI_ACCOUNT"
)

type MerchantType string

const (
	MerchantTypeMiniApp   MerchantType = "MINI_APP"
	MerchantTypeEcommerce MerchantType = "ECOMMERCE"
	MerchantTypeEvent     MerchantType = "EVENT"
	MerchantTypeLogistics MerchantType = "LOGISTICS"
)

type MiniAppMerchant struct {
	ID                string           `json:"id"`
	MerchantName      string           `json:"merchant_name"`
	MerchantCode      string           `json:"merchant_code"`
	SettlementMethod  SettlementMethod `json:"settlement_method"`
	// KYC               KYC              `json:"kyc,omitempty"`
	PhoneNumber       string            `json:"phone_number"`
	Email             string            `json:"email"`
	BankAccountNumber string           `json:"bank_account_number"`
	Enabled           bool             `json:"enabled"`
	IsDeleted         bool             `json:"is_deleted"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
	DeletedAt         *time.Time       `json:"deleted_at,omitempty"`
}



