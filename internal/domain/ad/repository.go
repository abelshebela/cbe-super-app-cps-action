package ad

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ADRepository interface {
	CreateAdvert(ctx context.Context, advert *Advert) (*Advert, error)
	UpdateAdvert(ctx context.Context, advert *Advert) (*Advert, error)
	DeleteAdvert(ctx context.Context, id string) (*Advert, error)
	EnableDisableAdvert(ctx context.Context, id string, enable bool) (*Advert, error)
	FetchAdvertByID(ctx context.Context, id string) (*Advert, error)
	FetchAdverts(ctx context.Context, filterParams *util_constant.Filter) (*utils.PaginatedResponse[[]*Advert], error)
	GetAdvertByTitle(ctx context.Context, title string) (*Advert, error)
}
