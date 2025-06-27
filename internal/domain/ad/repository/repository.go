package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/ad/entity"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
)

type Repository interface {
	CreateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	UpdateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	DeleteOneAdvert(ctx context.Context, cpsAction entity.CPSAction) error
	Authorize(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	Reject(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
}
