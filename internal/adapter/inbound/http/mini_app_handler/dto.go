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
	BannerImage         *multipart.FileHeader `form:"banner_image"`
	CommissionGLAccount string                `form:"commission_gl_account"`
	MerchantID          string                `form:"merchant_id"`
	IsEventMiniApp      bool                  `form:"is_event_mini_app"`
	IsThreeClick        bool                  `form:"is_three_click"`
	AppViewType         string                `form:"app_view_type"`

	URL string `form:"url"`

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
	ID                string                              `json:"id"`
	AppName           string                              `json:"app_name"`
	AppIcon           string                              `json:"app_icon"`
	BannerImage       string                              `json:"banner_image"`
	CommisonGLAccount string                              `json:"commison_gl_account,omitempty"`
	AppType           miniappentity.AppType               `json:"app_type"`
	MerchantID        string                              `json:"merchant_id"`
	AppViewType       miniappentity.AppViewType           `json:"app_view_type"`
	URL               string                              `json:"url,omitempty"`
	Stage             miniappentity.Stage                 `json:"stage"`
	ProductCode       []miniappentity.ProductCode         `json:"product_code"`
	Credential        miniappentity.CredentialInformation `json:"credential"`
	IsEventMiniApp    bool                                `json:"is_event_mini_app"`
	IsThreeClick      bool                                `json:"is_three_click"`
	Enabled           bool                                `json:"enabled"`
	CreatedAt         time.Time                           `json:"created_at"`
	LastModifiedAt    time.Time                           `json:"last_modified_at"`
}

// Validate validates the MiniAppRequest for create or update (PATCH)
func (r MiniAppRequest) Validate(isCreate bool) error {
	var fieldRules []*validation.FieldRules
	if isCreate {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppName,
				validation.Required.Error("app_name is required"),
				validation.Match(regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)).Error("app_name must not contain special characters"),
			),
			validation.Field(&r.MerchantID, validation.Required.Error("merchant_id is required")),
			validation.Field(&r.AppIcon, validation.Required, validation.By(validateFile)),
			validation.Field(&r.AppViewType, validation.Required.Error("app_view_type is required")),
			validation.Field(&r.URL, validation.Required.Error("url is required")),
			validation.Field(&r.BannerImage, validation.Required.Error("banner_image is required"), validation.By(validateFile)),
		}
	} else {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppIcon, validation.By(validateFile)),
			validation.Field(&r.BannerImage, validation.By(validateFile)),
		}
	}

	err := validation.ValidateStruct(&r, fieldRules...)
	if err != nil {
		return err
	}

	return nil
}

const MaxAvatarSize = 2 * 1024 * 1024

// isImageFormat checks if the content type is an allowed image
func isImageFormat(fileHeader *multipart.FileHeader) bool {
	if fileHeader == nil {
		return false
	}
	contentType := fileHeader.Header.Get("Content-Type")
	switch contentType {
	case "image/jpeg", "image/png", "image/gif":
		return true
	default:
		return false
	}
}

// validateAvatarFile checks the size and format of the uploaded image
func validateFile(value any) error {
	fileHeader, ok := value.(*multipart.FileHeader)
	if !ok || fileHeader == nil {
		return nil // Nothing to validate
	}

	if fileHeader.Size > MaxAvatarSize {
		return validation.NewError("file_too_large", "image must not exceed 2MB")
	}

	if !isImageFormat(fileHeader) {
		return validation.NewError("invalid_image_format", "image must be a valid image (jpeg, png, gif)")
	}

	return nil
}

// validateAppType ensures exactly one app type is provided and validates URL if AppType is "URL"
func validateAppType(r MiniAppRequest) validation.RuleFunc {
	return func(value any) error {
		if err := validateURL(r.URL); err != nil {
			return err
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

func validateAppViewType(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value any) error {
		viewType := strings.TrimSpace(r.AppViewType)
		if viewType == "" {
			if isCreate {
				return errors.New("APP_VIEW_TYPE_INVALID_OR_MISSING")
			}
			return nil
		}

		// Validate allowed enum values
		switch viewType {
		case string(miniappentity.AppViewTypeBoth),
			string(miniappentity.AppViewTypeCB),
			string(miniappentity.AppViewTypeIFB):
			return nil
		default:
			return errors.New("INVALID_APP_VIEW_TYPE")
		}
	}
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
		!dto.IsEventMiniApp &&
		!dto.IsThreeClick &&
		dto.BannerImage == nil
}
