package miniapphandler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	miniappentity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EnvironmentType enum
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

// MiniAppRequest DTO for both create and update (PATCH) requests
type MiniAppRequest struct {
	AppName             string                `form:"app_name"`
	AppIcon             *multipart.FileHeader `form:"app_icon"`
	CommissionGLAccount string                `form:"commission_gl_account"`
	MerchantID          string                `form:"merchant_id"`
	IsEventMiniApp      bool                  `form:"is_event_mini_app"`
	IsThreeClick        bool                  `form:"is_three_click"`

	// App types by environment
	AppTypeUAT  string `form:"app_type_uat"`
	AppTypeProd string `form:"app_type_production"`
	AppTypeTest string `form:"app_type_test"`
	AppTypeDev  string `form:"app_type_dev"`

	// IFB product codes
	IFBProductCode    string `form:"ifb_product_code"`
	IFBVATCode        string `form:"ifb_vat_code"`
	IFBServiceFeeCode string `form:"ifb_service_fee_code"`

	// CB product codes
	CBProductCode    string `form:"cb_product_code"`
	CBVATCode        string `form:"cb_vat_code"`
	CBServiceFeeCode string `form:"cb_service_fee_code"`

	// UAT credentials
	UATMerchantAppID string `form:"uat_merchant_app_id"`
	UATFabricAppID   string `form:"uat_fabric_app_id"`
	UATShortCode     string `form:"uat_short_code"`
	UATAppSecret     string `form:"uat_app_secret"`
	UATPrivateKey    string `form:"uat_private_key"`
	UATPublicKey     string `form:"uat_public_key"`

	// Prod credentials
	ProdMerchantAppID string `form:"prod_merchant_app_id"`
	ProdFabricAppID   string `form:"prod_fabric_app_id"`
	ProdShortCode     string `form:"prod_short_code"`
	ProdAppSecret     string `form:"prod_app_secret"`
	ProdPrivateKey    string `form:"prod_private_key"`
	ProdPublicKey     string `form:"prod_public_key"`

	// Test credentials
	TestMerchantAppID string `form:"test_merchant_app_id"`
	TestFabricAppID   string `form:"test_fabric_app_id"`
	TestShortCode     string `form:"test_short_code"`
	TestAppSecret     string `form:"test_app_secret"`
	TestPrivateKey    string `form:"test_private_key"`
	TestPublicKey     string `form:"test_public_key"`

	// Dev credentials
	DevMerchantAppID string `form:"dev_merchant_app_id"`
	DevFabricAppID   string `form:"dev_fabric_app_id"`
	DevShortCode     string `form:"dev_short_code"`
	DevAppSecret     string `form:"dev_app_secret"`
	DevPrivateKey    string `form:"dev_private_key"`
	DevPublicKey     string `form:"dev_public_key"`
}

// MiniAppResponse struct
type MiniAppResponse struct {
	ID                string                                `json:"id"`
	AppName           string                                `json:"app_name"`
	AppIcon           string                                `json:"app_icon"`
	CommisonGLAccount string                                `json:"commison_gl_account"`
	AppType           miniappentity.AppType                 `json:"app_type"`
	MerchantID        string                                `json:"merchant_id"`
	ProductCode       []miniappentity.ProductCode           `json:"product_code"`
	Credential        []miniappentity.CredentialInformation `json:"credential"`
	IsEventMiniApp    bool                                  `json:"is_event_mini_app"`
	IsThreeClick      bool                                  `json:"is_three_click"`
	Enabled           bool                                  `json:"enabled"`
	CreatedAt         time.Time                             `json:"created_at"`
	LastModifiedAt    time.Time                             `json:"last_modified_at"`
}

// Validate validates the MiniAppRequest for create or update (PATCH)
func (r MiniAppRequest) Validate(isCreate bool) error {
	// Field-specific validations
	var fieldRules []*validation.FieldRules
	if isCreate {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppName, validation.Required.Error("app_name is required")),
			validation.Field(&r.CommissionGLAccount, validation.Required.Error("commission_gl_account is required")),
			validation.Field(&r.MerchantID, validation.Required.Error("merchant_id is required")),
			validation.Field(&r.AppIcon, validation.Required, validation.By(validateFile)),
		}
	} else {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppIcon, validation.By(validateFile)),
		}
	}

	// Validate struct-level fields
	err := validation.ValidateStruct(&r, fieldRules...)
	if err != nil {
		return err
	}

	// Struct-level custom validations
	return validation.Validate(&r,
		validation.By(validateAppType(r, isCreate)),
		validation.By(validateProductCodes(r, isCreate)),
		validation.By(validateCredentials(r, isCreate)),
		validation.By(validateExclusiveAppFlags(r)),
	)
}

// validateFile validates the file upload
func validateFile(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return nil // Handled by Required for create
	}
	if file.Size > (2 << 20) {
		return errors.New("file size should be less than 2MB")
	}
	return nil
}

// validateAppType ensures exactly one app type is provided
func validateAppType(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value interface{}) error {
		types := []string{r.AppTypeUAT, r.AppTypeProd, r.AppTypeTest, r.AppTypeDev}
		count := 0
		for _, t := range types {
			if strings.TrimSpace(t) != "" {
				count++
			}
		}

		if isCreate && count == 0 {
			return errors.New("one app type (app_type_uat, app_type_production, app_type_test, app_type_dev) must be provided")
		}
		if count > 1 {
			return errors.New("only one app type (app_type_uat, app_type_production, app_type_test, app_type_dev) can be provided")
		}
		return nil
	}
}

// validateProductCodes ensures product codes are valid and complete for each branch
func validateProductCodes(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value interface{}) error {
		products := []struct {
			BranchType     BranchType
			ProductCode    string
			VATCode        string
			ServiceFeeCode string
		}{
			{IFB, r.IFBProductCode, r.IFBVATCode, r.IFBServiceFeeCode},
			{CB, r.CBProductCode, r.CBVATCode, r.CBServiceFeeCode},
		}

		validBranchCount := 0

		for _, p := range products {
			hasAny := p.ProductCode != "" || p.VATCode != "" || p.ServiceFeeCode != ""
			hasAll := p.ProductCode != "" && p.VATCode != "" && p.ServiceFeeCode != ""

			if hasAny && !hasAll {
				return fmt.Errorf("all fields for branch_type %s must be provided if one is set", p.BranchType)
			}
			if hasAll {
				validBranchCount++
			}
		}

		if validBranchCount == 0 {
			return errors.New("at least one complete set of product codes (IFB or CB) must be provided")
		}

		return nil
	}
}

// validateCredentials ensures exactly one credential environment is provided and valid
func validateCredentials(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value interface{}) error {
		creds := []struct {
			Environment   EnvironmentType
			MerchantAppID string
			FabricAppID   string
			ShortCode     string
			AppSecret     string
			PrivateKey    string
			PublicKey     string
		}{
			{UatEnvironment, r.UATMerchantAppID, r.UATFabricAppID, r.UATShortCode, r.UATAppSecret, r.UATPrivateKey, r.UATPublicKey},
			{ProductionEnvironment, r.ProdMerchantAppID, r.ProdFabricAppID, r.ProdShortCode, r.ProdAppSecret, r.ProdPrivateKey, r.ProdPublicKey},
			{TestEnvironment, r.TestMerchantAppID, r.TestFabricAppID, r.TestShortCode, r.TestAppSecret, r.TestPrivateKey, r.TestPublicKey},
			{DevEnvironment, r.DevMerchantAppID, r.DevFabricAppID, r.DevShortCode, r.DevAppSecret, r.DevPrivateKey, r.DevPublicKey},
		}

		seenEnvs := make(map[EnvironmentType]bool)
		for _, c := range creds {
			if c.MerchantAppID == "" && c.FabricAppID == "" && c.ShortCode == "" &&
				c.AppSecret == "" && c.PrivateKey == "" && c.PublicKey == "" {
				continue
			}
			if c.MerchantAppID == "" || c.FabricAppID == "" || c.ShortCode == "" ||
				c.AppSecret == "" || c.PrivateKey == "" || c.PublicKey == "" {
				return fmt.Errorf("all credential fields for environment %s must be provided", c.Environment)
			}
			seenEnvs[c.Environment] = true
		}

		if isCreate && len(seenEnvs) != 1 {
			return errors.New("exactly one credential environment must be provided")
		}
		if len(seenEnvs) > 1 {
			return errors.New("only one credential environment can be provided")
		}
		return nil
	}
}

func (r *MiniAppRequest) GetAppType() (miniappentity.AppType, error) {
	var selected miniappentity.AppType
	count := 0

	if r.AppTypeUAT != "" {
		selected = miniappentity.UAT
		count++
	}
	if r.AppTypeProd != "" {
		selected = miniappentity.Production
		count++
	}
	if r.AppTypeTest != "" {
		selected = miniappentity.Test
		count++
	}
	if r.AppTypeDev != "" {
		selected = miniappentity.Dev
		count++
	}

	if count == 0 {
		return "", errors.New("one app type must be provided")
	}
	if count > 1 {
		return "", errors.New("only one app type is allowed")
	}
	return selected, nil
}

func validateExclusiveAppFlags(r MiniAppRequest) validation.RuleFunc {
	return func(value any) error {
		if r.IsEventMiniApp && r.IsThreeClick {
			return errors.New("only one of is_event_mini_app or is_three_click can be true")
		}
		return nil
	}
}
