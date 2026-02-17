package services

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"

	service_dto "cbe-super-app-cps-action/internal/constants/dto/services"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/services/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type servicesService struct {
	repo   storage.ServicesRepository
	cps    service.CPSActionService
	logger utils.Logger
}

func NewServicesService(repo storage.ServicesRepository, cps service.CPSActionService, logger utils.Logger) *servicesService {
	return &servicesService{repo: repo, cps: cps, logger: logger}
}

func (s *servicesService) Create(ctx context.Context, req service_dto.CreateServiceRequest) error {
	exist, err := s.repo.CheckServiceExistence(ctx, req.ServiceCode, req.ServiceKey, req.ServiceName)
	if err != nil {
		return err
	}
	if exist {
		return errors.New(localization.ErrorServiceExists.Code)
	}

	// for _, service := range req.ServiceList {
	// 	if service.ServiceKey == nil || service.ServiceName == nil {
	// 		break
	// 	}
	// 	listExists, err := s.repo.CheckServiceExistence(ctx, "", *service.ServiceKey, *service.ServiceName)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if listExists {
	// 		return errors.New(localization.ErrorChildServiceExists.Code)
	// 	}
	// }

	mapped := core.MapToServiceModel(req)

	return core.HandleCPSAction(ctx, s.cps, "", constants.RequestCreateService, mapped, nil, constants.ActionCreate)
}

func (s *servicesService) Update(ctx context.Context, id string, req service_dto.UpdateServiceRequest) error {
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	searchCode := service_dto.StringPointer(req.ServiceCode, prev.ServiceCode)
	searchName := service_dto.StringPointer(req.ServiceName, prev.ServiceName)

	filterParam := types.Filter{
		Search: searchCode,
	}
	services, err := s.repo.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		return err
	}

	for _, svc := range services.Data {
		if svc.ID.Hex() != id && (svc.ServiceCode == searchCode || svc.ServiceName == searchName) {
			return errors.New(localization.ErrorServiceExists.Code)
		}
	}
	mapped := core.MapToServiceUpdateModel(req, *prev)
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestUpdateService, mapped, prev, constants.ActionUpdate)
}

func (s *servicesService) Enable(ctx context.Context, id string) error {
	if id == "" {
		return localization.ErrorInvalidID
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if prev.Enabled {
		return localization.ErrorAlreadyEnabled
	}
	payload := model.Service{Enabled: true}
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestEnableService, payload, prev, constants.ActionUpdate)
}

func (s *servicesService) Disable(ctx context.Context, id string) error {
	if id == "" {
		return localization.ErrorInvalidID
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if !prev.Enabled {
		return localization.ErrorAlreadyDisabled
	}
	payload := model.Service{Enabled: false}
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestDisableService, payload, prev, constants.ActionUpdate)
}

func (s *servicesService) GetAll(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]model.Service], error) {
	return s.repo.FindAllWithPagination(ctx, filter)
}

func (s *servicesService) GetAllServiceList(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]model.ServiceList], error) {
	return s.repo.FindAllServiceListWithPagination(ctx, filter)
}

func (s *servicesService) GetByID(ctx context.Context, id string) (*model.Service, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *servicesService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	serviceDoc, err := local_util.JsonUnmarshal[model.Service](action.CurrentAction)
	if err != nil {
		return nil, localization.ErrorInvalidActionData
	}

	switch action.RequestAction {
	case string(constants.RequestCreateService):
		err = s.repo.Create(ctx, serviceDoc)
	case string(constants.RequestUpdateService):
		err = s.repo.Update(ctx, action.UniqueId, serviceDoc)
	case string(constants.RequestEnableService):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
	case string(constants.RequestDisableService):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
	default:
		return nil, localization.ErrorInvalidRequest
	}
	if err != nil {
		return nil, err
	}
	action.CurrentAction = serviceDoc
	return action, nil
}
