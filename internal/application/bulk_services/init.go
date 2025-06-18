package bulkservices_application

import (
	"context"

	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type ApplicationAbstracts interface {
	FetchServices(ctx context.Context, offset, limit int) ([]domain.ServiceDetails, error)
	EnableDisableServicesMaker(ctx context.Context, Service_id string, action bool, makerId string) (string, error)
	EnableDisableServicesChecker(ctx context.Context, Action_id string, action bool, checkerId string) error

	FetchCifs(ctx context.Context, cif string) ([]domain.LinkedAccount, error)
	RemoveCifMaker(ctx context.Context, id []string, makerId string) (string, error)
	RemoveCifChecker(ctx context.Context, Action_id string, action bool, checkerId string) error
}
type ApplicationStore struct {
	service domain.ServiceImpl
}

func NewAttachDetachChecker(service domain.ServiceImpl) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
	}
}
