package miniappcore

import (
	"cbe-super-app-cps-action/internal/constants"
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/pkgs/utils"
	"errors"
	"log"
	"net/http"
	"strings"
)

func ParseMiniAppRequestFromMultipartForm(r *http.Request, isCreate bool) (miniappdto.MiniAppRequest, error) {
	var req miniappdto.MiniAppRequest

	file, fileHeader, err := utils.ParseMultipartFormFile(r, "app_icon", 2<<20)
	if err != nil {
		if err.Error() != localization.ErrorMissingFile.Code || isCreate {
			log.Println("Failed to upload app_icon: " + err.Error())
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}
	}
	if file != nil {
		defer file.Close()
	}
	req.AppIcon = fileHeader

	bannerFile, bannerFileHeader, err := utils.ParseMultipartFormFile(r, "banner_image", 2<<20)
	if err != nil {
		if err.Error() != localization.ErrorMissingFile.Code || isCreate {
			log.Println("Failed to upload banner_image: " + err.Error())
			return req, errors.New(localization.ErrorInvalidFileUpload.Code)
		}
	}
	if bannerFile != nil {
		defer bannerFile.Close()
	}
	req.BannerImage = bannerFileHeader

	get := func(key string) string {
		return strings.TrimSpace(r.FormValue(key))
	}

	req.AppName = get("app_name")
	req.CommissionGLAccount = get("commission_gl_account")
	req.MerchantID = get("merchant_id")
	req.URL = get("url")
	req.AppViewType = get("app_view_type")
	req.IFBProductCode = get("ifb_product_code")
	req.IFBVATCode = get("ifb_vat_code")
	req.IFBServiceFeeCode = get("ifb_service_fee_code")
	req.CBProductCode = get("cb_product_code")
	req.CBVATCode = get("cb_vat_code")
	req.CBServiceFeeCode = get("cb_service_fee_code")

	isEventMiniApp, err := parseBool(get("is_event_mini_app"))
	if err != nil {
		return req, errors.New(localization.ErrorInvalidBooleanFormat.Code)
	}
	req.IsEventMiniApp = isEventMiniApp

	isThreeClick, err := parseBool(get("is_three_click"))
	if err != nil {
		return req, errors.New(localization.ErrorInvalidBooleanFormat.Code)
	}
	req.IsThreeClick = isThreeClick

	return req, nil
}

func parseBool(value string) (bool, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "true":
		return true, nil
	case "false", "":
		return false, nil
	default:
		return false, errors.New(localization.ErrorInvalidBooleanFormat.Code)
	}
}

func ToMiniAppCreateRequest(r miniappdto.MiniAppRequest, isCreate bool) (*miniappdto.MiniAppCreateRequest, error) {

	result := miniappdto.MiniAppCreateRequest{
		AppName:             r.AppName,
		AppIcon:             r.AppIcon,
		BannerImage:         r.BannerImage,
		CommissionGLAccount: r.CommissionGLAccount,
		MerchantID:          r.MerchantID,
		IsEventMiniApp:      r.IsEventMiniApp,
		IsThreeClick:        r.IsThreeClick,
		AppType:             constants.URL,
		AppViewType:         constants.AppViewType(strings.ToUpper(r.AppViewType)),
		URL:                 r.URL,
		Stage:               constants.StageUat,
	}

	// Map ProductCode
	result.ProductCode = make([]types.ProductCode, 0, 2)
	products := []struct {
		BranchType     constants.BranchType
		ProductCode    string
		VATCode        string
		ServiceFeeCode string
	}{
		{constants.IFB, r.IFBProductCode, r.IFBVATCode, r.IFBServiceFeeCode},
		{constants.CB, r.CBProductCode, r.CBVATCode, r.CBServiceFeeCode},
	}
	for _, p := range products {
		if p.ProductCode != "" || p.VATCode != "" || p.ServiceFeeCode != "" {
			result.ProductCode = append(result.ProductCode, types.ProductCode{
				BranchType:     p.BranchType,
				ProductCode:    p.ProductCode,
				VATCode:        p.VATCode,
				ServiceFeeCode: p.ServiceFeeCode,
			})
		}
	}

	return &result, nil
}

func ToMiniAppResponse(m *model.MiniApp) miniappdto.MiniAppResponse {
	return miniappdto.MiniAppResponse{
		ID:                m.ID.Hex(),
		AppName:           m.AppName,
		AppIcon:           m.AppIcon,
		BannerImage:       m.BannerImage,
		CommisonGLAccount: m.CommissionGLAccount,
		AppType:           m.AppType,
		MerchantID:        m.MerchantID,
		ProductCode:       m.ProductCode,
		Credential:        m.Credential,
		URL:               m.URL,
		AppViewType:       m.AppViewType,
		Stage:             m.Stage,
		IsEventMiniApp:    m.IsEventMiniApp,
		IsThreeClick:      m.IsThreeClick,
		Enabled:           m.Enabled,
		CreatedAt:         m.CreatedAt,
		LastModifiedAt:    m.LastModifiedAt,
	}
}
