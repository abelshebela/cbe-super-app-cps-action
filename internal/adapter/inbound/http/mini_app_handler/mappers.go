package miniapphandler

import (
	"strings"

	miniappentity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
)

// ToMiniAppCreateRequest maps MiniAppRequest to MiniAppCreateRequest
func (r MiniAppRequest) ToMiniAppCreateRequest(isCreate bool) (*miniappentity.MiniAppCreateRequest, error) {
	// Run existing validation to avoid duplication
	if err := r.Validate(isCreate); err != nil {
		return nil, err
	}

	appType, err := r.GetAppType(isCreate)

	if err != nil {
		return nil, err
	}
	// Initialize result
	result := miniappentity.MiniAppCreateRequest{
		AppName:             r.AppName,
		AppIcon:             r.AppIcon,
		BannerImage:         r.BannerImage,
		CommissionGLAccount: r.CommissionGLAccount,
		MerchantID:          r.MerchantID,
		IsEventMiniApp:      r.IsEventMiniApp,
		IsThreeClick:        r.IsThreeClick,
		AppType:             appType,
		AppViewType:         miniappentity.AppViewType(strings.ToUpper(r.AppViewType)),
		URL:                 r.URL,
		MPAASID:             r.MPAASID,
		Stage:               miniappentity.StageUat,
	}

	// Map ProductCode
	result.ProductCode = make([]miniappentity.ProductCode, 0, 2)
	products := []struct {
		BranchType     miniappentity.BranchType
		ProductCode    string
		VATCode        string
		ServiceFeeCode string
	}{
		{miniappentity.IFB, r.IFBProductCode, r.IFBVATCode, r.IFBServiceFeeCode},
		{miniappentity.CB, r.CBProductCode, r.CBVATCode, r.CBServiceFeeCode},
	}
	for _, p := range products {
		if p.ProductCode != "" || p.VATCode != "" || p.ServiceFeeCode != "" {
			result.ProductCode = append(result.ProductCode, miniappentity.ProductCode{
				BranchType:     p.BranchType,
				ProductCode:    p.ProductCode,
				VATCode:        p.VATCode,
				ServiceFeeCode: p.ServiceFeeCode,
			})
		}
	}

	return &result, nil
}

func ToMiniAppResponse(m *miniappentity.MiniApp) MiniAppResponse {
	return MiniAppResponse{
		ID:                m.ID,
		AppName:           m.AppName,
		AppIcon:           m.AppIcon,
		BannerImage:       m.BannerImage,
		CommisonGLAccount: m.CommissionGLAccount,
		AppType:           m.AppType,
		MerchantID:        m.MerchantID,
		ProductCode:       m.ProductCode,
		Credential:        m.Credential,
		URL:               m.URL,
		MPAASID:           m.MPAASID,
		AppViewType:       m.AppViewType,
		Stage:             m.Stage,
		IsEventMiniApp:    m.IsEventMiniApp,
		IsThreeClick:      m.IsThreeClick,
		Enabled:           m.Enabled,
		CreatedAt:         m.CreatedAt,
		LastModifiedAt:    m.LastModifiedAt,
	}
}
