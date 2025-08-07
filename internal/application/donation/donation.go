package donation

import (
	"context"
	"mime/multipart"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DonationAbstract interface {
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest, maker cps_entities.User) error
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest, maker cps_entities.User) error
	DeleteDonationCategory(ctx context.Context, id string, maker cps_entities.User) error
}

type DonationService interface {
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) (*dto.DonationCategoryResponse, error)
	UploadIcon(ctx context.Context, icon *multipart.FileHeader) (string, error)
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error)
}

type DonationStore struct {
	service    DonationService
	cpsService cps_service.CPSActionService
	logger     utils.Logger
}

func NewDonationApplication(service DonationService, cpsService cps_service.CPSActionService, logger utils.Logger) DonationAbstract {
	return &DonationStore{
		service:    service,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (a *DonationStore) handleCPSAction(ctx context.Context, maker cps_entities.User, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cps_entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	_, err := a.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest, maker cps_entities.User) error {
	iconURL, err := a.service.UploadIcon(ctx, donation.Icon)
	if err != nil {
		a.logger.Errorf("failed to upload icon: %v", err)
		return err
	}

	cpsRequest := dto.DonationCategoryCPSRequest{
		CategoryName: donation.CategoryName,
		Icon:         iconURL,
	}

	if err := a.handleCPSAction(ctx, maker, cps_const.RequestCreateDonationCategory, cpsRequest, nil, cps_const.ActionCreate); err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error) {
	paginatedResponse, err := a.service.FetchDonationCategory(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("failed to fetch donation categories: %v", err)
		return nil, err
	}
	return paginatedResponse, nil
}

func (a *DonationStore) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	donationCategory, err := a.service.FetchDonationCategoryByID(ctx, id)
	if err != nil {
		a.logger.Errorf("failed to fetch donation category by ID: %v", err)
		return nil, err
	}
	return donationCategory, nil
}

func (a *DonationStore) UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest, maker cps_entities.User) error {
	var iconURL string
	var err error

	if donation.Icon != nil {
		iconURL, err = a.service.UploadIcon(ctx, donation.Icon)
		if err != nil {
			a.logger.Errorf("failed to upload icon: %v", err)
			return err
		}
	}

	cpsRequest := dto.DonationCategoryCPSRequest{
		CategoryName: donation.CategoryName,
		Icon:         iconURL,
	}

	prevData := map[string]interface{}{"id": id}
	if err := a.handleCPSAction(ctx, maker, cps_const.RequestUpdateDonationCategory, cpsRequest, prevData, cps_const.ActionUpdate); err != nil {
		return err
	}
	return nil
}

func (a *DonationStore) DeleteDonationCategory(ctx context.Context, id string, maker cps_entities.User) error {
	return nil
}
