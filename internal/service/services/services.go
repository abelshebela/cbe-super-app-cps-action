package services

import (
	"context"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
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

// Public API (raises CPS actions)
func (s *servicesService) Create(ctx context.Context, req model.Services) error {
	if err := core.ValidateCreate(req, s.repo); err != nil {
		return err
	}
	// new entity, uniqueId empty

	return core.HandleCPSAction(ctx, s.cps, "", constants.RequestCreateService, req, nil, constants.ActionCreate)
}

func (s *servicesService) Update(ctx context.Context, id string, req model.Services) error {
	// fetch existing and validate uniqueness if code or name changes
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if req.ServiceCode != "" {
		res, err := s.repo.FindAllWithPagination(ctx, types.Filter{Filters: map[string]interface{}{"service_code": req.ServiceCode}})
		if err != nil {
			return err
		}
		for _, v := range res.Data {
			if v.ID.Hex() != id {
				return localization.ErrorServiceExists
			}
		}
	}
	if req.ServiceName != "" {
		res, err := s.repo.FindAllWithPagination(ctx, types.Filter{Filters: map[string]interface{}{"service_name": req.ServiceName}})
		if err != nil {
			return err
		}
		for _, v := range res.Data {
			if v.ID.Hex() != id {
				return localization.ErrorServiceExists
			}
		}
	}
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestUpdateService, req, prev, constants.ActionUpdate)
}

func (s *servicesService) Enable(ctx context.Context, id string) error {
	if id == "" {
		return localization.ErrorInvalidID
	}
	// carry minimal payload
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if prev.Enabled {
		return localization.ErrorAlreadyEnabled
	}
	payload := model.Services{Enabled: true}
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
	payload := model.Services{Enabled: false}
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestDisableService, payload, prev, constants.ActionUpdate)
}

// Reads (direct)
func (s *servicesService) GetAll(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.Services], error) {
	return s.repo.FindAllWithPagination(ctx, filter)
}

func (s *servicesService) GetByID(ctx context.Context, id string) (*model.Services, error) {
	return s.repo.FindByID(ctx, id)
}

// Authorize applies changes on approval
func (s *servicesService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	serviceDoc, err := local_util.JsonUnmarshal[model.Services](action.CurrentAction)
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

// local CPS wrapper functions removed; using services/core helper directly above
