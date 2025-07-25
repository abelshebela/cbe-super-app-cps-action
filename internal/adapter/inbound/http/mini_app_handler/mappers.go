package miniapphandler

import (
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
		CommissionGLAccount: r.CommissionGLAccount,
		MerchantID:          r.MerchantID,
		IsEventMiniApp:      r.IsEventMiniApp,
		IsThreeClick:        r.IsThreeClick,
		AppType:             appType,
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

	// Map CredentialInformation
	creds := []struct {
		Environment   miniappentity.EnvironmentType
		MerchantAppID string
		FabricAppID   string
		ShortCode     string
		AppSecret     string
		PrivateKey    string
		PublicKey     string
	}{
		{miniappentity.UatEnvironment, r.UATMerchantAppID, r.UATFabricAppID, r.UATShortCode, r.UATAppSecret, r.UATPrivateKey, r.UATPublicKey},
		{miniappentity.ProductionEnvironment, r.ProdMerchantAppID, r.ProdFabricAppID, r.ProdShortCode, r.ProdAppSecret, r.ProdPrivateKey, r.ProdPublicKey},
		{miniappentity.TestEnvironment, r.TestMerchantAppID, r.TestFabricAppID, r.TestShortCode, r.TestAppSecret, r.TestPrivateKey, r.TestPublicKey},
		{miniappentity.DevEnvironment, r.DevMerchantAppID, r.DevFabricAppID, r.DevShortCode, r.DevAppSecret, r.DevPrivateKey, r.DevPublicKey},
	}
	for _, c := range creds {
		if c.MerchantAppID != "" || c.FabricAppID != "" || c.ShortCode != "" ||
			c.AppSecret != "" || c.PrivateKey != "" || c.PublicKey != "" {
			result.Credential = []miniappentity.CredentialInformation{
				{
					Environment:   c.Environment,
					MerchantAppID: c.MerchantAppID,
					FabricAppID:   c.FabricAppID,
					ShortCode:     c.ShortCode,
					AppSecret:     c.AppSecret,
					PrivateKey:    c.PrivateKey,
					PublicKey:     c.PublicKey,
				},
			}
			break // Only one credential set is allowed
		}
	}

	return &result, nil
}

func ToMiniAppResponse(m *miniappentity.MiniApp) MiniAppResponse {
	return MiniAppResponse{
		ID:                m.ID,
		AppName:           m.AppName,
		AppIcon:           m.AppIcon,
		CommisonGLAccount: m.CommissionGLAccount,
		AppType:           m.AppType,
		MerchantID:        m.MerchantID,
		ProductCode:       m.ProductCode,
		Credential:        m.Credential,
		IsEventMiniApp:    m.IsEventMiniApp,
		IsThreeClick:      m.IsThreeClick,
		Enabled:           m.Enabled,
		CreatedAt:         m.CreatedAt,
		LastModifiedAt:    m.LastModifiedAt,
	}
}
