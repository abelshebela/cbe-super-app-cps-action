package dto

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
	ID          string     `json:"id"`
	BranchType  BranchType `json:"branch_type"`
	ProductCode string     `json:"product_code"`
}

type AppType struct {
	UAT        string `json:"uat"`
	Production string `json:"production"`
	Test       string `json:"test"`
	Dev        string `json:"dev"`
}

type CredentialInformation struct {
	Environment   EnvironmentType `json:"environment"`
	MerchantAppID string          `json:"merchant_app_id"`
	FabricAppID   string          `json:"fabric_app_id"`
	ShortCode     string          `json:"short_code"`
	AppSecret     string          `json:"app_secret"`
	PrivateKey    string          `json:"private_key"`
	PublicKey     string          `json:"public_key"`
}

type MiniAppCreateRequest struct {
	ID                string                  `json:"id"`
	AppName           string                  `json:"app_name"`
	AppIcon           string                  `json:"app_icon"`
	CommisonGLAccount string                  `json:"commison_gl_account"`
	AppType           AppType                 `json:"app_type"`
	MerchantID        string                  `json:"merchant_id"`
	ProductCode       []ProductCode           `json:"product_code"`
	Credential        []CredentialInformation `json:"credential"`
	IsEventMiniApp    bool                    `json:"is_event_mini_app"`
	IsThreeClick      bool                    `json:"is_three_click"`
	Enabled           bool                    `json:"enabled"`
}

type MiniAppCheckerRequest struct {
	Action_id string `json:"action_id"`
	Action    bool   `json:"action"`
}
