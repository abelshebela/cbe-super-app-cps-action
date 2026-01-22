package miniappdto

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"mime/multipart"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (r MiniAppRequest) Validate(isCreate bool) error {
	var fieldRules []*validation.FieldRules
	if isCreate {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppName,
				validation.Required.Error(localization.ErrorMiniAppNameRequired.Code),
				// validation.Match(regexp.MustCompile(`^[a-zA-Z0-9 _-]+$`)).Error(localization.ErrorInvalidAppNameFormat.Code),
			),
			validation.Field(&r.MerchantID, validation.Required.Error(localization.ErrorMiniAppMerchantIDRequired.Code)),
			validation.Field(&r.AppIcon, validation.Required.Error(localization.ErrorAppIconRequired.Code), validation.By(validateImage)),
			validation.Field(&r.AppViewType, validation.Required.Error(localization.ErrorAppViewTypeRequired.Code), validation.By(ValidateAppViewType(r, isCreate))),
			validation.Field(&r.URL, validation.Required.Error(localization.ErrorMiniAppURLRequired.Code), validation.By(validateURL)),
			validation.Field(&r.BannerImage, validation.By(validateImage)),
		}
	} else {
		fieldRules = []*validation.FieldRules{
			validation.Field(&r.AppIcon, validation.By(validateImage)),
			validation.Field(&r.BannerImage, validation.By(validateImage)),
		}
	}

	return validation.ValidateStruct(&r, fieldRules...)
}

func validateImage(value interface{}) error {
	file, ok := value.(*multipart.FileHeader)
	if !ok || file == nil {
		return localization.ErrorMissingOrInvalidImage
	}
	if !utils.IsValidImage(file) {
		return localization.ErrorMissingOrInvalidImage
	}

	if file.Size > (10 << 20) {
		return validation.NewError("logo", "file size exceeds 10MB limit")
	}

	return nil
}

// validateURL checks if the provided string is a valid URL
func validateURL(value any) error {
	raw, ok := value.(string)
	if !ok || raw == "" {
		return nil // Nothing to validate
	}

	parsedURL, err := url.ParseRequestURI(raw)
	if err != nil {
		return errors.New(localization.ErrorInvalidURL.Code)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New(localization.ErrorInvalidURL.Code)
	}

	if parsedURL.Host == "" {
		return errors.New(localization.ErrorInvalidURL.Code)
	}

	validURL := regexp.MustCompile(`^https?://([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}(/[\w-./?%&=]*)?$`)
	if !validURL.MatchString(raw) {
		return errors.New(localization.ErrorInvalidURL.Code)
	}

	return nil
}

// validateProductCodes ensures product codes are valid and complete for each branch
func validateProductCodes(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value any) error {
		products := []struct {
			BranchType     constants.BranchType
			ProductCode    string
			VATCode        string
			ServiceFeeCode string
		}{
			{constants.IFB, r.IFBProductCode, r.IFBVATCode, r.IFBServiceFeeCode},
			{constants.CB, r.CBProductCode, r.CBVATCode, r.CBServiceFeeCode},
		}

		validBranchCount := 0

		for _, p := range products {
			hasAny := p.ProductCode != "" || p.VATCode != "" || p.ServiceFeeCode != ""
			hasAll := p.ProductCode != "" && p.VATCode != "" && p.ServiceFeeCode != ""

			if hasAny && !hasAll {
				return errors.New(localization.ErrorIncompleteBranchProductCodes.Code)
			}
			if hasAll {
				validBranchCount++
			}
		}

		if isCreate && validBranchCount == 0 {
			return errors.New(localization.ErrorNoProductCodesProvided.Code)
		}

		if r.AppViewType == string(constants.AppViewTypeBoth) {
			if validBranchCount != 2 {
				return errors.New(localization.ErrorBothProductCodesRequired.Code)
			}
		}

		return nil
	}
}

// ValidateAppViewType ensures the app view type is valid
func ValidateAppViewType(r MiniAppRequest, isCreate bool) validation.RuleFunc {
	return func(value any) error {
		viewType := strings.TrimSpace(r.AppViewType)
		if viewType == "" {
			if isCreate {
				return errors.New(localization.ErrorAppViewTypeInvalidOrMissing.Code)
			}
			return nil
		}

		switch viewType {
		case string(constants.AppViewTypeBoth), string(constants.AppViewTypeCB), string(constants.AppViewTypeIFB):
			return nil
		default:
			return errors.New(localization.ErrorInvalidAppViewType.Code)
		}
	}
}

// ValidateExclusiveAppFlags ensures IsEventMiniApp and IsThreeClick are not both true
func ValidateExclusiveAppFlags(r MiniAppRequest) validation.RuleFunc {
	return func(value any) error {
		if r.IsEventMiniApp && r.IsThreeClick {
			return errors.New(localization.ErrorExclusiveAppFlags.Code)
		}
		return nil
	}
}

// IsEmptyUpdate checks if the MiniAppCreateRequest is empty
func IsEmptyUpdate(dto *MiniAppCreateRequest) bool {
	return dto.AppName == "" &&
		dto.AppIcon == nil &&
		dto.CommissionGLAccount == "" &&
		dto.AppType == "" &&
		len(dto.ProductCode) == 0 &&
		!dto.IsEventMiniApp &&
		!dto.IsThreeClick &&
		dto.BannerImage == nil &&
		dto.MerchantID == "" &&
		dto.AppViewType == "" &&
		dto.URL == "" &&
		dto.Stage == "" &&
		reflect.DeepEqual(dto.Credential, types.CredentialInformation{})
}
