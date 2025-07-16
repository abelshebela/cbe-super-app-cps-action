package event

import (
	"context"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type EventRepository interface {
	CreateCpsAction(ctx context.Context, Action entity.CPSAction) (*entity.CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action entity.CPSAction) error
	FetchCpsActionByID(ctx context.Context, Action_Id string) (*entity.CPSAction, error)
	CreateEvent(ctx context.Context, event Event) (*Event, error)
	FetchEventById(ctx context.Context, event_id string) (*Event, error)
	FetchEvent(ctx context.Context, limit, offset int) ([]*Event, error)
}
