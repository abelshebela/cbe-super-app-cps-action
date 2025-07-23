package event_outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"
)

type Outbound interface {
	CreateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action action.CPSAction) error
	FetchCpsActionById(ctx context.Context, Action_Id string) (*entities.CPSAction, error)
	CreateEvent(ctx context.Context, event event.Event) (event.Event, error)

	FetchEventById(ctx context.Context, event_id string) (event.Event, error)
	FetchEvent(ctx context.Context, limit, offset int) ([]event.Event, error)
}
