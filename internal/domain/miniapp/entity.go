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

type ProductCode struct {
	ID          string
	BranchType  BranchType
	ProductCode string
}

type AppType struct {
	UAT        string
	Production string
	Test       string
	Dev        string
}

type CredentialInformation struct {
	ID            string
	Environment   EnvironmentType
	MerchantAppID string
	FabricAppID   string
	ShortCode     string
	AppSecret     string
	PrivateKey    string
	PublicKey     string
}

type MiniApp struct {
	ID                string
	AppName           string
	AppIcon           string
	CommisonGLAccount string
	AppType           AppType
	MerchantID        string
	ProductCode       []ProductCode
	Credential        []CredentialInformation
	IsEventMiniApp    bool
	IsThreeClick      bool
	Enabled           bool
	IsDeleted         bool
	CreatedAt         time.Time
	LastModifiedAt    time.Time
	DeletedAt         time.Time
}
