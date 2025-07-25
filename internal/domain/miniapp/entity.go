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

type ProductCode struct {
	ID             string     `bson:"id"`
	BranchType     BranchType `bson:"branch_type"`
	ProductCode    string     `bson:"product_code"`
	VATCode        string     `bson:"vat_code"`
	ServiceFeeCode string     `bson:"service_fee_code"`
}

type AppType string

const (
	UAT        AppType = "UAT"
	Production AppType = "PRODUCATION"
	Test       AppType = "TEST"
	Dev        AppType = "DEV"
)

type CredentialInformation struct {
	ID            string          `bson:"id"`
	Environment   EnvironmentType `bson:"environment"`
	MerchantAppID string          `bson:"merchant_app_id"`
	FabricAppID   string          `bson:"fabric_app_id"`
	ShortCode     string          `bson:"short_code"`
	AppSecret     string          `bson:"app_secret"`
	PrivateKey    string          `bson:"private_key"`
	PublicKey     string          `bson:"public_key"`
}

type MiniApp struct {
	ID                  string                  `bson:"_id"`
	AppName             string                  `bson:"app_name"`
	AppIcon             string                  `bson:"app_icon"`
	CommissionGLAccount string                  `bson:"commison_gl_account"`
	AppType             AppType                 `bson:"app_type"`
	MerchantID          string                  `bson:"merchant_id"`
	ProductCode         []ProductCode           `bson:"product_code"`
	Credential          []CredentialInformation `bson:"credential"`
	IsEventMiniApp      bool                    `bson:"is_event_mini_app"`
	IsThreeClick        bool                    `bson:"is_three_click"`
	Enabled             bool                    `bson:"enabled"`
	IsDeleted           bool                    `bson:"is_deleted"`
	CreatedAt           time.Time               `bson:"created_at"`
	LastModifiedAt      time.Time               `bson:"last_modified_at"`
	DeletedAt           time.Time               `bson:"deleted_at"`
}

type MiniAppCreateRequest struct {
	ID                  string                  `json:"id"`
	AppName             string                  `json:"app_name"`
	AppIcon             *multipart.FileHeader   `json:"app_icon"`
	CommissionGLAccount string                  `json:"commison_gl_account"`
	AppType             AppType                 `json:"app_type"`
	MerchantID          string                  `json:"merchant_id"`
	ProductCode         []ProductCode           `json:"product_code"`
	Credential          []CredentialInformation `json:"credential"`
	IsEventMiniApp      bool                    `json:"is_event_mini_app"`
	IsThreeClick        bool                    `json:"is_three_click"`
}
