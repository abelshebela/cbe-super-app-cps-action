package miniapp

import "time"

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

// type ProductCode struct {
// 	ID          string
// 	BranchType  BranchType
// 	ProductCode string
// }

type ProductCode struct {
	ID          string     `bson:"id"`
	BranchType  BranchType `bson:"branch_type"`
	ProductCode string     `bson:"product_code"`
}

type AppType struct {
	UAT        string
	Production string
	Test       string
	Dev        string
}

// type CredentialInformation struct {
// 	ID            string
// 	Environment   EnvironmentType
// 	MerchantAppID string
// 	FabricAppID   string
// 	ShortCode     string
// 	AppSecret     string
// 	PrivateKey    string
// 	PublicKey     string
// }

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
	ID                string                  `bson:"id"`
	AppName           string                  `bson:"app_name"`
	AppIcon           string                  `bson:"app_icon"`
	CommisonGLAccount string                  `bson:"commison_gl_account"`
	AppType           AppType                 `bson:"app_type"`
	MerchantID        string                  `bson:"merchant_id"`
	ProductCode       []ProductCode           `bson:"product_code"`
	Credential        []CredentialInformation `bson:"credential"`
	IsEventMiniApp    bool                    `bson:"is_event_mini_app"`
	IsThreeClick      bool                    `bson:"is_three_click"`
	Enabled           bool                    `bson:"enabled"`
	IsDeleted         bool                    `bson:"is_deleted"`
	CreatedAt         time.Time               `bson:"created_at"`
	LastModifiedAt    time.Time               `bson:"last_modified_at"`
	DeletedAt         time.Time               `bson:"deleted_at"`
}

// type MiniApp struct {
// 	ID                string
// 	AppName           string
// 	AppIcon           string
// 	CommisonGLAccount string
// 	AppType           AppType
// 	MerchantID        string
// 	ProductCode       []ProductCode
// 	Credential        []CredentialInformation
// 	IsEventMiniApp    bool
// 	IsThreeClick      bool
// 	Enabled           bool
// 	IsDeleted         bool
// 	CreatedAt         time.Time
// 	LastModifiedAt    time.Time
// 	DeletedAt         time.Time
// }
