package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AdRepository interface {
	CreateOneAdvert(ctx context.Context, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	DeleteOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	Authorize(ctx context.Context, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error)
	Reject(ctx context.Context, cpsAction model.RejectCPSAction) (*model.CPSAction, error)
}
