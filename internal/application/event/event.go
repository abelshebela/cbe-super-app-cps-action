package event_application

import (
	"context"

	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	cps_entitites "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	evententity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationAbstracts interface {
	CreateEvent(ctx context.Context, event dto.EventRequest, maker domain.Maker) error
	UpdateEvent(ctx context.Context, id string, event dto.EventRequest, maker domain.Maker) error
	DeleteEvent(ctx context.Context, id string, maker domain.Maker) error
	EnableDisableEvent(ctx context.Context, id string, maker domain.Maker, enable bool) error

	FetchEventByID(ctx context.Context, id string) (*evententity.Event, error)
	FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*evententity.Event], error)
}
type ApplicationStore struct {
	service    domain.EventService
	cpsService cps_service.CPSActionService
	logger     utils.Logger
}

func NewEventApplication(service domain.EventService,
	cpsService cps_service.CPSActionService,
	logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service:    service,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (a *ApplicationStore) CreateEvent(ctx context.Context, event dto.EventRequest, maker domain.Maker) error {
	res, err := a.service.CreateEvent(ctx, event)
	if err != nil {
		return err
	}

	cpsAction := a.cpsService.BuildCPSAction(ctx, cpsactions.CreateCPSRequest{
		User: maker,
	})
	return nil
}

func (a *ApplicationStore) UpdateEvent(ctx context.Context, id string, event dto.EventRequest, maker domain.Maker) error {

	return nil
}
func (a *ApplicationStore) DeleteEvent(ctx context.Context, id string, maker domain.Maker) error {
	return nil
}
func (a *ApplicationStore) EnableDisableEvent(ctx context.Context, id string, maker domain.Maker, enable bool) error {
	return nil
}

func (a *ApplicationStore) FetchEventByID(ctx context.Context, id string) (*evententity.Event, error) {
	return nil, nil
}
func (a *ApplicationStore) FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*evententity.Event], error) {
	return nil, nil
}
