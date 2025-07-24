package dto

import (
	"fmt"
	"mime/multipart"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
	AppIcon           *multipart.FileHeader   `json:"app_icon"`
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

func (r MiniAppCreateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.AppName, validation.Required.Error("app_name is required")),
		validation.Field(&r.AppIcon, validation.By(validateFile)),
		validation.Field(&r.CommisonGLAccount, validation.Required.Error("commison_gl_account is required")),
		validation.Field(&r.MerchantID, validation.Required.Error("merchant_id is required")),
		validation.Field(&r.ProductCode, validation.Required.Error("product_code is required")),
		validation.Field(&r.Credential, validation.Required.Error("credential is required")),
	)
}

func validateFile(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok {
		return fmt.Errorf("invalid file")
	}
	if file.Size > (2 << 20) { // 2MB
		return fmt.Errorf("file size should be less than 2MB")
	}
	return nil
}
