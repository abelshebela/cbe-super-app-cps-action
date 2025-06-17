package event_outbound

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/event"
)

type Outbound interface {
	CreateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action action.CPSAction) error
	FetchCpsActionById(ctx context.Context, Action_Id string) (entity.CPSAction, error)
	CreateEvent(ctx context.Context, event event.Event) (event.Event, error)

	FetchEventById(ctx context.Context, event_id string) (event.Event, error)
	FetchEvent(ctx context.Context, limit, offset int) ([]event.Event, error)
}
