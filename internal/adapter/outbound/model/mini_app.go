package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
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
	UAT        AppType = "UAT"
	Production AppType = "PRODUCATION"
	Test       AppType = "TEST"
	Dev        AppType = "DEV"
)

type ProductCode struct {
	ID             string     `bson:"id"`
	BranchType     BranchType `bson:"branch_type"`
	ProductCode    string     `bson:"product_code"`
	VATCode        string     `bson:"vat_code"`
	ServiceFeeCode string     `bson:"service_fee_code"`
}

type CredentialInformation struct {
	ID            bson.ObjectID   `bson:"id" json:"id,omitempty"`
	Environment   EnvironmentType `bson:"environment" json:"environment"`
	MerchantAppID string          `bson:"merchant_app_id" json:"merchant_app_id"`
	FabricAppID   string          `bson:"fabric_app_id" json:"fabric_app_id"`
	ShortCode     string          `bson:"short_code" json:"short_code"`
	AppSecret     string          `bson:"app_secret" json:"app_secret"`
	PrivateKey    string          `bson:"private_key" json:"private_key"`
	PublicKey     string          `bson:"public_key" json:"public_key"`
	Timestamp     time.Time       `bson:"timestamp" json:"timestamp"`
	Signature     string          `bson:"signature" json:"-"`
	MiniAppCode   string          `bson:"mini_app_code" json:"mini_app_code"`
}

type MiniApp struct {
	ID                bson.ObjectID           `bson:"_id"`
	AppName           string                  `bson:"app_name"`
	AppIcon           string                  `bson:"app_icon"`
	CommisonGLAccount string                  `bson:"commison_gl_account"`
	AppType           string                  `bson:"app_type"`
	MerchantID        string                  `bson:"merchant_id"`
	ProductCode       []ProductCode           `bson:"product_code"`
	Credential        []CredentialInformation `bson:"credential"`
	AppViewType       string                  `bson:"app_view_type"`
	URL               string                  `bson:"url"`
	MPAASID           string                  `bson:"mpaas_id"`
	Stage             string                  `bson:"stage"`
	IsEventMiniApp    bool                    `bson:"is_event_mini_app"`
	IsThreeClick      bool                    `bson:"is_three_click"`
	Enabled           bool                    `bson:"enabled"`
	IsDeleted         bool                    `bson:"is_deleted"`
	CreatedAt         time.Time               `bson:"created_at"`
	LastModifiedAt    time.Time               `bson:"last_modified_at"`
	DeletedAt         time.Time               `bson:"deleted_at"`
}
