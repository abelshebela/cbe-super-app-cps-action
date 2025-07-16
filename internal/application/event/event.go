package event_application

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
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

func (a *ApplicationStore) MakerCreateEvent(ctx context.Context, event dto.EventCreateRequest, makerId, makerName, makerPhone string) (string, error) {
	ticket_req := domain.Ticket{}
	event_req := domain.Event{}

	request_id, err := a.service.CreateEventRequest(ctx, event_req, ticket_req, makerId, makerName, makerPhone)
	if err != nil {
		return "", err
	}
	return request_id, nil
}
func (a *ApplicationStore) CheckerCreateEvent(ctx context.Context, actionId string, action bool, checkerId, checkerName, checkerPhone string) error {
	return a.service.ApproveEventRequest(ctx, actionId, action, checkerId, checkerName, checkerPhone)
}

func (a *ApplicationStore) FetchAllEvents(ctx context.Context, limit, offset int) ([]dto.EventResponse, error) {
	data, err := a.service.FetchEvent(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	var result []dto.EventResponse
	for range /* _, d := */ data {
		result = append(result, dto.EventResponse{})
	}
	return result, nil
}
func (a *ApplicationStore) FetchEvent(ctx context.Context, event_id string) (dto.EventDTO, error) {
	data, err := a.service.FetchEventByID(ctx, event_id)
	if err != nil {
		return dto.EventDTO{}, err
	}
	return dto.EventDTO{
		EventID:   data.ID,
		EventCode: data.Code,
		EventName: data.Name,
	}, nil
}
