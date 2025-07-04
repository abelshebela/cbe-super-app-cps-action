package bulkservices_application

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type ApplicationAbstracts interface {
	FetchServices(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, *common.ErrorDefinition)
	EnableDisableServicesMaker(ctx context.Context, Service_id string, action bool, makerId string) (string, *common.ErrorDefinition)
	EnableDisableServicesChecker(ctx context.Context, Action_id string, action bool, checkerId string) *common.ErrorDefinition

	FetchCifs(ctx context.Context, cif string) ([]domain.LinkedAccount, *common.ErrorDefinition)
	RemoveCifMaker(ctx context.Context, id []string, makerId string) (string, *common.ErrorDefinition)
	RemoveCifChecker(ctx context.Context, Action_id string, action bool, checkerId string) *common.ErrorDefinition
}
type ApplicationStore struct {
	service domain.ServiceImpl
	Logger  utils.Logger
}

func NewAttachDetachChecker(service domain.ServiceImpl, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
		Logger:  logger,
	}
}
