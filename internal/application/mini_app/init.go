package miniapp_application

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/miniapp"
)

type ApplicationAbstracts interface {
	MakerCreateMiniApp(ctx context.Context, miniApp dto.MiniAppCreateRequest, makerId string) (string, error)
	CheckerCreateMiniApp(ctx context.Context, actionId string, action bool, checkerId string) error
}
type ApplicationStore struct {
	service domain.MiniAppService
}

func NewAttachDetachChecker(service domain.MiniAppService) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
	}
}
