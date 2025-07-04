package miniapp_application

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
)

type ApplicationAbstracts interface {
	MakerCreateMiniApp(ctx context.Context, miniApp dto.MiniAppCreateRequest, makerId string) (string, *common.ErrorDefinition)
	CheckerCreateMiniApp(ctx context.Context, actionId string, action bool, checkerId string) *common.ErrorDefinition
}
type ApplicationStore struct {
	service domain.MiniAppService
	Logger  utils.Logger
}

func NewAttachDetachChecker(service domain.MiniAppService, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
		Logger:  logger,
	}
}
