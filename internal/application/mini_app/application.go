package miniapp_application

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type ApplicationAbstracts interface {
	MakerCreateMiniApp(ctx context.Context, miniApp dto.MiniAppCreateRequest, maker model.User) (*entities.CPSAction, error)
	MakerUpdateMiniApp(ctx context.Context, req dto.MiniAppCreateRequest, maker model.User) (*entities.CPSAction, error)
	MakerDeleteMiniApp(ctx context.Context, maker *model.User, id string) (*entities.CPSAction, error)
	ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*model.MiniApp], error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}

type ApplicationStore struct {
	service domain.MiniAppService
	Logger  utils.Logger
}

func NewApplicationService(service domain.MiniAppService, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
		Logger:  logger,
	}
}

func (a *ApplicationStore) MakerCreateMiniApp(ctx context.Context, miniApp dto.MiniAppCreateRequest, maker model.User) (*entities.CPSAction, error) {
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
	action, err := a.service.CreateMiniAppAction(ctx, data, maker)
	if err != nil {
		a.Logger.Errorf("[mini_app.MakerCreateMiniApp] ", err.Error())
		return nil, err
	}
	return action, nil
}

func (a *ApplicationStore) MakerUpdateMiniApp(ctx context.Context, req dto.MiniAppCreateRequest, maker model.User) (*entities.CPSAction, error) {
	data := domain.MiniApp{
		AppName:           req.AppName,
		AppIcon:           req.AppIcon,
		CommisonGLAccount: req.CommisonGLAccount,
		AppType: domain.AppType{
			UAT:        req.AppType.UAT,
			Production: req.AppType.Production,
			Test:       req.AppType.Test,
			Dev:        req.AppType.Dev,
		},
		MerchantID: req.MerchantID,
		ProductCode: func() []domain.ProductCode {
			var productCodes []domain.ProductCode
			for _, pc := range req.ProductCode {
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
			for _, cred := range req.Credential {
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
		IsEventMiniApp: req.IsEventMiniApp,
		IsThreeClick:   req.IsThreeClick,
		Enabled:        req.Enabled,
		LastModifiedAt: time.Now(),
	}
	update, err := a.service.UpdateMiniAppAction(ctx, data, maker, req.ID)
	if err != nil {
		a.Logger.Errorf("[mini_app.MakerUpdateMiniApp] %v", err)
		return nil, err
	}
	return update, nil
}

func (a *ApplicationStore) MakerDeleteMiniApp(ctx context.Context, maker *model.User, id string) (*entities.CPSAction, error) {
	cpsAction, err := a.service.DeleteMiniAppAction(ctx, *maker, id)
	if err != nil {
		a.Logger.Errorf("[mini_app.MakerDeleteMiniApp] %v", err)
		return nil, err
	}
	return cpsAction, nil
}

func (a *ApplicationStore) ListMiniApp(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*model.MiniApp], error) {
	list, err := a.service.ListMiniApp(ctx, filterParam)
	if err != nil {
		a.Logger.Errorf("[mini_app.ListMiniApp] %v", err)
		return nil, err
	}
	return list, nil
}

func (a *ApplicationStore) DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error) {
	detail, err := a.service.DetailMiniAppByID(ctx, id)
	if err != nil {
		a.Logger.Errorf("[mini_app.DetailMiniAppByID] %v", err)
		return model.MiniApp{}, err
	}
	return detail, nil
}
