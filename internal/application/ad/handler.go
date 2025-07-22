package ad

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADHandlers interface {
	CreateOneAdvert(ctx context.Context, adCpsReq model.CreateCPSAction) (*entities.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Advert], error)
	GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error)
	UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	DeleteOneAdvert(ctx context.Context, id string, adCpsReq model.CreateCPSAction) (*model.CPSAction, error)
}

type ADHandler struct {
	adDomain   service.AdvertService
	bucketName string
	minio      config.MinioClientInterface
	logger     utils.Logger
}

func InitADHandler(adDomain service.AdvertService, minioClinet config.MinioClientInterface, bucketName string, logger utils.Logger) ADHandlers {
	return ADHandler{
		adDomain:   adDomain,
		minio:      minioClinet,
		logger:     logger,
		bucketName: bucketName,
	}
}

func (a ADHandler) CreateOneAdvert(ctx context.Context, adCpsReq model.CreateCPSAction) (*entities.CPSAction, error) {
	return a.adDomain.CreateOneAdvert(ctx, adCpsReq)
}

func (a ADHandler) DeleteOneAdvert(ctx context.Context, id string, adCpsReq model.CreateCPSAction) (*model.CPSAction, error) {
	return a.adDomain.DeleteOneAdvert(ctx, id, adCpsReq)
}

func (a ADHandler) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.Advert], error) {
	adverts, err := a.adDomain.GetAllAdvert(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return adverts, nil
}

func (a ADHandler) GetOneAdvert(ctx context.Context, id string) (*entity.Advert, error) {
	advert, err := a.adDomain.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (a ADHandler) UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error) {
	advertCpsAction, err := a.adDomain.UpdateOneAdvert(ctx, id, cpsAction)
	if err != nil {
		return nil, err
	}

	return advertCpsAction, nil
}
