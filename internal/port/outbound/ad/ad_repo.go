package ad

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type ADRepo interface {
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Advert], error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	HandleAdvertCreate(ctx context.Context, cpsRes *entities.CPSAction) (*entities.CPSAction, error)
	HandleAdvertUpdate(ctx context.Context, cpsRes *entities.CPSAction) (*entities.CPSAction, error)
	HandleAdvertDelete(ctx context.Context, cpsRes *entities.CPSAction) (*entities.CPSAction, error)
	
}
