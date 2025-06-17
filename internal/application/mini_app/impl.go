package miniapp_application

import (
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/miniapp"
)

func (a *ApplicationStore) MakerCreateMiniApp(ctx context.Context, miniApp dto.MiniAppCreateRequest, makerId string) (string, error) {
	data := domain.MiniApp{
		AppName:           miniApp.AppName,
		AppIcon:           miniApp.AppIcon,
		CommisonGLAccount: miniApp.CommisonGLAccount,
		AppType: domain.AppType{
			UAT:        miniApp.AppType.UAT,
			Production: miniApp.AppType.Production,
			Test:       miniApp.AppType.Test,
			Dev:        miniApp.AppType.Dev,
		},
		MerchantID: miniApp.MerchantID,
		ProductCode: func() []domain.ProductCode {
			var productCodes []domain.ProductCode
			for _, pc := range miniApp.ProductCode {
				productCodes = append(productCodes, domain.ProductCode{
					ID:          pc.ID,
					BranchType:  domain.BranchType(pc.BranchType),
					ProductCode: pc.ProductCode,
				})
			}
			return productCodes
		}(),
		Credential: func() []domain.CredentialInformation {
			var credentials []domain.CredentialInformation
			for _, cred := range miniApp.Credential {
				credentials = append(credentials, domain.CredentialInformation{
					Environment:   domain.EnvironmentType(cred.Environment),
					MerchantAppID: cred.MerchantAppID,
					FabricAppID:   cred.FabricAppID,
					ShortCode:     cred.ShortCode,
					AppSecret:     cred.AppSecret,
					PrivateKey:    cred.PrivateKey,
					PublicKey:     cred.PublicKey,
				})
			}
			return credentials
		}(),
		IsEventMiniApp: miniApp.IsEventMiniApp,
		IsThreeClick:   miniApp.IsThreeClick,
		Enabled:        miniApp.Enabled,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
		DeletedAt:      time.Time{},
	}
	action_id, err := a.service.CreateMiniAppAction(ctx, data, makerId)
	if err != nil {
		return "", err
	}
	return action_id, nil
}

func (a *ApplicationStore) CheckerCreateMiniApp(ctx context.Context, actionId string, action bool, checkerId string) error {
	return a.service.CheckMiniApp(ctx, actionId, action, checkerId)
}
