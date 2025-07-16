package event

import "context"

type EventService interface {
	CreateEventRequest(ctx context.Context, event Event, ticket Ticket, makerID, makerName, makerPhone string) (string, error)
	ApproveEventRequest(ctx context.Context, action_id string, action_taken bool, checkerID, checkerName, checkerPhone string) error
	FetchEvent(ctx context.Context, limit, offset int) ([]Event, error)
	FetchEventById(ctx context.Context, event_id string) (Event, error)
}

type Service struct {
	Repository Repository
}

func NewService(repository Repository) (EventService, error) {
	return &Service{
		Repository: repository,
	}, nil
}
