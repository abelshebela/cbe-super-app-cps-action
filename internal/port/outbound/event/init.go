package event_outbound

import (
	"context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/event"

)

type EventRepository interface {
	CreateEvent(ctx context.Context, event event.Event) (*event.Event, error)
	UpdateEvent(ctx context.Context, event event.Event) (*event.Event, error)
	DeleteEvent(ctx context.Context, id string) (*event.Event, error)
	EnableDisableEvent(ctx context.Context, id string, enable bool) (*event.Event, error)
	EventNameExists(ctx context.Context, name string, id *string) (bool, error)

	FetchEventByID(ctx context.Context, id string) (*event.Event, error)
	FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*event.Event], error)
}