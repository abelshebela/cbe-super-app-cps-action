package event

import (
	"context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"


)

type EventRepository interface {
	CreateEvent(ctx context.Context, event Event) (*Event, error)
	UpdateEvent(ctx context.Context, event Event) (*Event, error)
	DeleteEvent(ctx context.Context, id string) (*Event, error)
	EnableDisableEvent(ctx context.Context, id string, enable bool) (*Event, error)

	FetchEventByID(ctx context.Context, id string) (*Event, error)
	FetchEvent(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Event], error)
}
