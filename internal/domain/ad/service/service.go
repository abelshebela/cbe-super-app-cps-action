package service

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/repository"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADDomain struct {
	ADRepo repository.Repository
	logger utils.Logger
}

type AdvertService interface {
	CreateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	UpdateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	DeleteOneAdvert(ctx context.Context, cpsAction entity.CPSAction) error
	Authorize(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
	Reject(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error)
}

var _ AdvertService = (*ADDomain)(nil)

func InitADDomian(adRepo repository.Repository, logger utils.Logger) AdvertService {
	return &ADDomain{
		ADRepo: adRepo,
		logger: logger,
	}
}

func (a *ADDomain) CreateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	cpsActionRes, err := a.ADRepo.CreateOneAdvert(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (a *ADDomain) DeleteOneAdvert(ctx context.Context, cpsAction entity.CPSAction) error {
	if err := a.ADRepo.DeleteOneAdvert(ctx, cpsAction); err != nil {
		return err
	}
	return nil
}

func (a *ADDomain) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*entity.AdvertResponse, error) {
	advertRes, err := a.ADRepo.GetAllAdvert(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return advertRes, nil
}

func (a *ADDomain) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	advertRes, err := a.ADRepo.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	return advertRes, nil
}

func (a *ADDomain) UpdateOneAdvert(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	advertAction, err := a.ADRepo.UpdateOneAdvert(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return advertAction, nil
}

func (a *ADDomain) Authorize(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	cpsActionRes, err := a.ADRepo.Authorize(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}

func (a *ADDomain) Reject(ctx context.Context, cpsAction entity.CPSAction) (*entity.CPSAction, error) {
	cpsActionRes, err := a.ADRepo.Reject(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return cpsActionRes, nil
}
