package ad

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	
)

type ADRepo interface {
	CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Advert], error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	Reject(ctx context.Context, cpsAction model.RejectCPSAction) (*model.CPSAction, error)
}
