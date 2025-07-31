// Package miniapphandler provides DTOs and validation functions for miniapp requests
package miniapphandler

import (
	"errors"
	"mime/multipart"
	"net/url"
	"regexp"
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
	AppViewType         string                `form:"app_view_type"`

	// App types
	AppType string `form:"app_type"`

	URL     string `form:"url"`
	MPAASID string `form:"mpaas_id"`

	// IFB product codes
	IFBProductCode    string `form:"ifb_product_code"`
	IFBVATCode        string `form:"ifb_vat_code"`
	IFBServiceFeeCode string `form:"ifb_service_fee_code"`

	// CB product codes
	CBProductCode    string `form:"cb_product_code"`
	CBVATCode        string `form:"cb_vat_code"`
	CBServiceFeeCode string `form:"cb_service_fee_code"`
}

// MiniAppResponse struct
type MiniAppResponse struct {
	ID                string                                `json:"id"`
	AppName           string                                `json:"app_name"`
	AppIcon           string                                `json:"app_icon"`
	CommisonGLAccount string                                `json:"commison_gl_account,omitempty"`
	AppType           miniappentity.AppType                 `json:"app_type"`
	MerchantID        string                                `json:"merchant_id"`
	AppViewType       miniappentity.AppViewType             `json:"app_view_type"`
	URL               string                                `json:"url,omitempty"`
	MPAASID           string                                `json:"mpaas_id,omitempty"`
	Stage             miniappentity.Stage                   `json:"stage"`
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
	var fieldRules []*validation.FieldRules
	if isCreate {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppName, validation.Required.Error("app_name is required")),
			validation.Field(&r.MerchantID, validation.Required.Error("merchant_id is required")),
			validation.Field(&r.AppIcon, validation.Required, validation.By(validateFile)),
			validation.Field(&r.AppViewType, validation.Required.Error("app_view_type is required")),
			validation.Field(&r.AppType, validation.Required.Error("app_type is required")),
		}
	} else {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppIcon, validation.By(validateFile)),
		}
	}

	err := validation.ValidateStruct(&r, fieldRules...)
	if err != nil {
		return err
	}

	return validation.Validate(&r,
		validation.By(validateAppType(r)),
		validation.By(validateProductCodes(r, isCreate)),
		validation.By(validateExclusiveAppFlags(r)),
		validation.By(validateAppViewType(r)),
	)
}

// validateFile validates the file upload
func validateFile(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return nil
	}
	if file.Size > (2 << 20) {
		return errors.New("FILE_TOO_LARGE")
	}
	return nil
}

// validateAppType ensures exactly one app type is provided and validates URL if AppType is "URL"
func validateAppType(r MiniAppRequest) validation.RuleFunc {
	return func(value any) error {
		if strings.ToUpper(r.AppType) == "URL" {
			if r.URL == "" {
				return errors.New("URL_REQUIRED")
			}
			// Validate if URL is valid
			if err := validateURL(r.URL); err != nil {
				return err
			}
		}
		return nil
	}
}

// validURL is a stricter regex for basic URL validation as a fallback
var validURL = regexp.MustCompile(`^https?://([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(/[\w-./?%&=]*)?$`)

// validateURL checks if the provided string is a valid URL
func validateURL(raw string) error {
	parsedURL, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New("INVALID_URL")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("INVALID_URL")
	}

	if parsedURL.Host == "" {
		return errors.New("INVALID_URL")
	}

	if !validURL.MatchString(raw) {
		return errors.New("INVALID_URL")
	}

	return nil
}

// validateProductCodes ensures product codes are valid and complete for each branch
func validateProductCodes(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value any) error {
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
				return errors.New("INCOMPLETE_BRANCH_PRODUCT_CODES")
			}
			if hasAll {
				validBranchCount++
			}
		}

		if isCreate && validBranchCount == 0 {
			return errors.New("NO_PRODUCT_CODES_PROVIDED")
		}

		if r.AppViewType == "BOTH" {
			if validBranchCount != 2 {
				return errors.New("BOTH_PRODUCT_CODES_REQUIRED")
			}
		}

		return nil
	}
}

// validateAppViewType ensures app_view_type is one of the allowed values
func validateAppViewType(r MiniAppRequest) validation.RuleFunc {
	return func(value any) error {
		viewType := strings.TrimSpace(r.AppViewType)
		if viewType == "" {
			return errors.New("APP_VIEW_TYPE_INVALID_OR_MISSING")
		}

		// Validate against allowed enum values
		switch viewType {
		case string(miniappentity.AppViewTypeBoth), string(miniappentity.AppViewTypeCB), string(miniappentity.AppViewTypeIFB):
			return nil

		default:
			return errors.New("INVALID_APP_VIEW_TYPE")
		}
	}
}

func (r *MiniAppRequest) GetAppType(isCreate bool) (miniappentity.AppType, error) {
	var selected miniappentity.AppType
	count := 0

	if strings.ToUpper(r.AppType) == "URL" {
		selected = miniappentity.URL
		count++
	}
	if strings.ToUpper(r.AppType) == "MPAASID" {
		selected = miniappentity.MPAASID
		count++
	}

	if count == 0 && isCreate {
		return "", errors.New("APP_TYPE_MISSING")
	}

	return selected, nil
}

func validateExclusiveAppFlags(r MiniAppRequest) validation.RuleFunc {
	return func(value any) error {
		if r.IsEventMiniApp && r.IsThreeClick {
			return errors.New("EXCLUSIVE_APP_FLAGS")
		}
		return nil
	}
}

func IsEmptyUpdate(dto *miniappentity.MiniAppCreateRequest) bool {
	return dto.AppName == "" &&
		dto.AppIcon == nil &&
		dto.CommissionGLAccount == "" &&
		dto.AppType == "" &&
		len(dto.ProductCode) == 0 &&
		len(dto.Credential) == 0 &&
		!dto.IsEventMiniApp &&
		!dto.IsThreeClick
}
