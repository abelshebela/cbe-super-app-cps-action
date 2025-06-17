package event_application

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/event"
)

type ApplicationAbstracts interface {
	MakerCreateEvent(ctx context.Context, event dto.EventCreateRequest, makerId, makerName, makerPhone string) (string, error)
	CheckerCreateEvent(ctx context.Context, actionId string, action bool, checkerId, checkerName, checkerPhone string) error
	FetchAllEvents(ctx context.Context, limit, offset int) ([]dto.EventResponse, error)
	FetchEvent(ctx context.Context, event_id string) (dto.EventDTO, error)
}
type ApplicationStore struct {
	service domain.EventService
}

func NewEventApplication(service domain.EventService) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
	}
}
