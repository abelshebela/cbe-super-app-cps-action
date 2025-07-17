package event

import (
	"context"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type EventRepository interface {
	CreateCpsAction(ctx context.Context, action entity.CPSAction) (*entity.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action entity.CPSAction) error
	FetchCpsActionByID(ctx context.Context, actionID string) (*entity.CPSAction, error)
	CreateEvent(ctx context.Context, event Event) (*Event, error)
	FetchEventByID(ctx context.Context, eventID string) (*Event, error)
	FetchEvent(ctx context.Context, limit, offset int) ([]*Event, error)
	CPSActionExists(ctx context.Context, cpsReq entity.CreateCPSAction) error
}
