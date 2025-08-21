package ad

import (
	"context"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ADRepository interface {
	CreateAdvert(ctx context.Context, advert *entity.Advert) (*entity.Advert, error)
	UpdateAdvert(ctx context.Context, advert *entity.Advert) (*entity.Advert, error)
	DeleteAdvert(ctx context.Context, id string) (*entity.Advert, error)
	EnableDisableAdvert(ctx context.Context, id string, enable bool) (*entity.Advert, error)
	FetchAdvertByID(ctx context.Context, id string) (*entity.Advert, error)
	FetchAdverts(ctx context.Context, filterParams *util_constant.Filter) (*utils.PaginatedResponse[[]*entity.Advert], error)
	GetAdvertByTitle(ctx context.Context, title string) (*entity.Advert, error)
}
