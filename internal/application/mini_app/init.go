package miniapp_application

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
)

type ApplicationAbstracts interface {
	MakerCreateMiniApp(ctx context.Context, miniApp dto.MiniAppCreateRequest, maker model.User) (string, *common.ErrorDefinition)
	CheckerCreateMiniApp(ctx context.Context, actionId string, action bool, checker model.User) *common.ErrorDefinition
	MakerUpdateMiniApp(ctx context.Context, req dto.MiniAppCreateRequest, maker model.User) (string, error)
	MakerDeleteMiniApp(ctx context.Context, maker *model.User) (string, error)
	ListMiniApp(ctx context.Context) ([]*domain.MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (domain.MiniApp, error)
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
