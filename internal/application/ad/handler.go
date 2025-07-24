package ad

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/service"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ADHandlers interface {
	CreateOneAdvert(ctx context.Context, adCpsReq model.CreateCPSAction) (*entities.CPSAction, error)
	GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.AdvertResponse], error)
	GetOneAdvert(ctx context.Context, id string) (*entity.AdvertResponse, error)
	UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error)
	DeleteOneAdvert(ctx context.Context, id string, adCpsReq model.CreateCPSAction) (*entities.CPSAction, error)
	EnableOrDisableAdvert(ctx context.Context, id string, requestAction cps_const.RequestAction, cpsReq model.CreateCPSAction) (*entities.CPSAction, error)
}

type ADHandler struct {
	adDomain   service.AdvertService
	bucketName string
	minio      config.MinioClientInterface
	logger     utils.Logger
	cpsService cps_service.CPSActionService
}

func InitADHandler(adDomain service.AdvertService, minioClinet config.MinioClientInterface, bucketName string,
	cpsService cps_service.CPSActionService,
	logger utils.Logger) ADHandlers {
	return ADHandler{
		adDomain:   adDomain,
		minio:      minioClinet,
		logger:     logger,
		bucketName: bucketName,
		cpsService: cpsService,
	}
}

func (a ADHandler) CreateOneAdvert(ctx context.Context, adCpsReq model.CreateCPSAction) (*entities.CPSAction, error) {
	action, err := a.adDomain.CreateOneAdvert(ctx, adCpsReq)

	if err != nil {
		return nil, err
	}

	// check pending action
	_, err = a.cpsService.CPSActionExists(ctx, entities.CheckCPSAction{
		UserCode:      action.MakerID,
		FullName:      action.MakerName,
		Department:    action.Department,
		PhoneNumber:   action.MakerPhoneNumber,
		RequestAction: string(action.RequestAction),
	})

	if err != nil {
		return nil, err
	}

	return a.cpsService.CreateCPSAction(ctx, action)

}

func (a ADHandler) DeleteOneAdvert(ctx context.Context, id string, adCpsReq model.CreateCPSAction) (*entities.CPSAction, error) {

	action, err := a.adDomain.DeleteOneAdvert(ctx, id, adCpsReq)
	if err != nil {
		return nil, err
	}

	_, err = a.cpsService.CPSActionExists(ctx, entities.CheckCPSAction{
		UserCode:      action.MakerID,
		FullName:      action.MakerName,
		Department:    action.Department,
		PhoneNumber:   action.MakerPhoneNumber,
		RequestAction: string(action.RequestAction),
	})

	if err != nil {
		return nil, err
	}

	return a.cpsService.CreateCPSAction(ctx, action)
}

func (a ADHandler) GetAllAdvert(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.AdvertResponse], error) {
	adverts, err := a.adDomain.GetAllAdvert(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return adverts, nil
}

func (a ADHandler) GetOneAdvert(ctx context.Context, id string) (*entity.AdvertResponse, error) {
	advert, err := a.adDomain.GetOneAdvert(ctx, id)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (a ADHandler) UpdateOneAdvert(ctx context.Context, id string, cpsAction model.CreateCPSAction) (*entities.CPSAction, error) {
	action, err := a.adDomain.UpdateOneAdvert(ctx, id, cpsAction)

	if err != nil {
		return nil, err
	}

	// check pending action
	_, err = a.cpsService.CPSActionExists(ctx, entities.CheckCPSAction{
		UserCode:      action.MakerID,
		FullName:      action.MakerName,
		Department:    action.Department,
		PhoneNumber:   action.MakerPhoneNumber,
		RequestAction: string(action.RequestAction),
	})

	if err != nil {
		return nil, err
	}

	return a.cpsService.CreateCPSAction(ctx, action)
}

func (a ADHandler) EnableOrDisableAdvert(ctx context.Context, id string, requestAction cps_const.RequestAction, cpsReq model.CreateCPSAction) (*entities.CPSAction, error) {
	action, err := a.adDomain.EnableOrDisableAdvert(ctx, id, requestAction, cpsReq)
	if err != nil {
		return nil, err
	}

	// check pending action
	_, err = a.cpsService.CPSActionExists(ctx, entities.CheckCPSAction{
		UserCode:      action.MakerID,
		FullName:      action.MakerName,
		Department:    action.Department,
		PhoneNumber:   action.MakerPhoneNumber,
		RequestAction: string(action.RequestAction),
	})

	if err != nil {
		return nil, err
	}
	return a.cpsService.CreateCPSAction(ctx, action)
}
