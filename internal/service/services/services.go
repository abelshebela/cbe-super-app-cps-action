package services

import (
	"context"
	"database/sql"
	"errors"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	service_dto "cbe-super-app-cps-action/internal/constants/dto/services"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/services/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	coreio "github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type servicesService struct {
	repo   storage.ServicesRepository
	cps    service.CPSActionService
	core   coreio.CBECoreAPIInterface
	logger utils.Logger
}

func NewServicesService(repo storage.ServicesRepository, cps service.CPSActionService, core coreio.CBECoreAPIInterface, logger utils.Logger) *servicesService {
	return &servicesService{
		repo:   repo,
		cps:    cps,
		core:   core,
		logger: logger,
	}
}

func (s *servicesService) Create(ctx context.Context, req service_dto.CreateServiceRequest) error {
	s.logger.Infof("Service creating...")

	accessList, err := s.repo.FindServiceListByID(ctx, req.ServiceKeyId)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		s.logger.Errorf("[servicesService][Create] error checking existing service for serviceKeyId=%s: %v", req.ServiceKeyId, err)
		return err
	}

	if req.ProductGlAccount != "" {
		accountDetail, err := s.ValidateAccountNumberWithExternalAPI(ctx, req.ProductGlAccount)
		if err != nil {
			s.logger.Errorf("[servicesService][Authorize] account number validation failed for account number %s: %v", req.ProductGlAccount, err)
			return errors.New(localization.ErrorAccountNumberValidationFailed.Code)
		}

		err = s.repo.CheckAccountNumberExistence(ctx, *accountDetail, req.ProductGlAccount, req.ProductGlAccountCurrency)
		if err != nil {
			s.logger.Errorf("[servicesService][Authorize] failed to check account number error: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	mapped := core.MapToServiceModel(req, *accessList)

	return core.HandleCPSAction(ctx, s.cps, "", constants.RequestCreateService, mapped, nil, constants.ActionCreate)
}

func (s *servicesService) Update(ctx context.Context, id string, req service_dto.UpdateServiceRequest) error {
	s.logger.Infof("[servicesService][Update] called with id=%s, req=%+v", id, req)
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[servicesService][Update] error fetching previous service by id=%s: %v", id, err)
		return err
	}

	s.logger.Infof("[servicesService][Update] previous service: %+v", prev)
	serviceKeyId := service_dto.StringPointer(req.ServiceKeyId, prev.ServiceKeyId)
	s.logger.Infof("[servicesService][Update] resolved serviceKeyId: %s", serviceKeyId)

	if req.ProductGlAccount != nil {
		accountDetail, err := s.ValidateAccountNumberWithExternalAPI(ctx, *req.ProductGlAccount)
		if err != nil {
			s.logger.Errorf("[servicesService][Authorize] account number validation failed for account number %s: %v", *req.ProductGlAccount, err)
			return errors.New(localization.ErrorAccountNumberValidationFailed.Code)
		}

		err = s.repo.CheckAccountNumberExistence(ctx, *accountDetail, *req.ProductGlAccount, *req.ProductGlAccountCurrency)
		if err != nil {
			s.logger.Errorf("[servicesService][Authorize] failed to check account number error: %v", err)
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	mapped := core.MapToServiceUpdateModel(req, *prev)
	s.logger.Infof("[servicesService][Update] mapped update model: %+v", mapped)
	err = core.HandleCPSAction(ctx, s.cps, id, constants.RequestUpdateService, mapped, prev, constants.ActionUpdate)
	if err != nil {
		s.logger.Errorf("[servicesService][Update] HandleCPSAction failed: %v", err)
		return err
	}
	s.logger.Infof("[servicesService][Update] service update successful for id=%s", id)
	return nil
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
	payload := imodel.ServiceKey{IsEnabled: true, ServiceKey: prev.ServiceKey}
	return core.HandleCPSAction(ctx, s.cps, prev.ServiceKeyId, constants.RequestEnableService, payload, prev, constants.ActionUpdate)
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
	payload := imodel.ServiceKey{IsEnabled: false, ServiceKey: prev.ServiceKey}
	return core.HandleCPSAction(ctx, s.cps, prev.ServiceKeyId, constants.RequestDisableService, payload, prev, constants.ActionUpdate)
}

func (s *servicesService) DeleteServices(ctx context.Context, id string) error {
	if id == "" {
		return localization.ErrorInvalidID
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[ServiceSVC][DeleteService] failed to delete error: %v", err)
		return err
	}

	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestDeleteServiceList, nil, prev, constants.ActionDelete)
}

func (s *servicesService) GetAll(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]service_dto.ServiceResponse], error) {
	return s.repo.FindAllWithPagination(ctx, filter)
}

func (s *servicesService) CreateServiceList(ctx context.Context, req *service_dto.CreateServiceList) error {
	list, err := s.repo.FindServiceListByExactNameOrKey(ctx, req.ServiceName, req.ServiceKey)
	if err != nil && err.Error() != localization.ErrorAccessListNotFound.Code {
		return err
	}
	if list != nil {
		return errors.New(localization.ErrorServiceListAlreadyExists.Code)
	}

	mapped := core.MapServiceListDtoToModel(req)
	return core.HandleCPSAction(ctx, s.cps, "", constants.RequestCreateServiceList, mapped, nil, constants.ActionCreate)
}

func (s *servicesService) UpdateServiceList(ctx context.Context, id string, req *service_dto.UpdateServiceList) error {
	existing, err := s.repo.FindServiceListByID(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorAccessListNotFound.Code {
			return localization.ErrorAccessListNotFound
		}
		return err
	}

	// if strings.EqualFold(existing.ServiceName, req.ServiceName) || strings.EqualFold(existing.ServiceKey, req.ServiceKey) {
	// 	return errors.New(localization.ErrorNoChangesDetected.Code)
	// }

	mapped := core.MapServiceListDtoUpdateToModel(req)
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestUpdateServiceList, mapped, existing, constants.ActionUpdate)
}

func (s *servicesService) GetAllServiceList(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]imodel.ServiceKey], error) {
	return s.repo.FindAllServiceListWithPagination(ctx, filter)
}

func (s *servicesService) GetByID(ctx context.Context, id string) (*service_dto.ServiceResponse, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *servicesService) EnableOrDisableServiceList(ctx context.Context, id string, enable bool) error {
	if id == "" {
		return localization.ErrorInvalidID
	}
	prev, err := s.repo.FindServiceListByID(ctx, id)
	if err != nil {
		return err
	}
	if prev.IsEnabled == enable {
		if enable {
			return localization.ErrorAlreadyEnabled
		}
		return localization.ErrorAlreadyDisabled
	}
	payload := *prev
	payload.IsEnabled = enable

	var requestAction constants.RequestAction
	if enable {
		requestAction = constants.RequestEnableServiceList
	} else {
		requestAction = constants.RequestDisableServiceList
	}
	return core.HandleCPSAction(ctx, s.cps, id, requestAction, payload, prev, constants.ActionUpdate)
}

func (s *servicesService) DeleteServiceKey(ctx context.Context, id string) error {
	prev, err := s.repo.FindServiceListByID(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorAccessListNotFound.Code {
			return localization.ErrorAccessListNotFound
		}
		s.logger.Errorf("[servicesService][DeleteServiceKey] error fetching access list by id=%s: %v", id, err)
		return err
	}

	s.logger.Infof("[servicesService][DeleteServiceKey] Deleting service key with id=%s, found service list: %+v", id, prev)
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestDeleteServiceList, nil, prev, constants.ActionDelete)
}

func (s *servicesService) ValidateAccountNumberWithExternalAPI(ctx context.Context, accountNumber string) (*coreio.AccountLookupResult, error) {
	response, err := s.core.AccountLookup(coreio.AccountLookupParam{AccountNumber: accountNumber})
	if err != nil {
		return nil, err
	}

	if !response.Success {
		var message string
		for _, msg := range response.Messages {
			message += msg
		}

		s.logger.Warnf("(core) failed to get account details: %s", message)

		return nil, err
	}

	if response.Detail == nil {
		s.logger.Errorf("account lookup successful but no account details found for account number %s", accountNumber)
		return nil, err
	}

	return response, nil
}

func (s *servicesService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[servicesService][Authorize] Authorize called for action: %+v", action)

	serviceDoc, err := local_util.JsonUnmarshal[imodel.Service](action.CurrentAction)
	if err != nil {
		return nil, localization.ErrorInvalidActionData
	}

	switch action.RequestAction {
	case string(constants.RequestCreateService):
		s.logger.Infof("[servicesService][Authorize] Authorizing create service with data: %+v", serviceDoc)

		return nil, s.repo.Create(ctx, serviceDoc.ProductGlAccount, serviceDoc)
	case string(constants.RequestUpdateService):
		s.logger.Infof("[servicesService][Authorize] Authorizing update service with data: %+v", serviceDoc)
		return nil, s.repo.Update(ctx, action.UniqueId, serviceDoc, serviceDoc.ProductGlAccount)
	case string(constants.RequestEnableService):
		// err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, true)
	case string(constants.RequestDisableService):
		// err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, false)
	case string(constants.RequestDeleteService):
		err = s.repo.Delete(ctx, action.UniqueId, "")
	case string(constants.RequestCreateServiceList):
		listDoc, err := local_util.JsonUnmarshal[imodel.ServiceKey](action.CurrentAction)
		if err != nil {
			return nil, localization.ErrorInvalidActionData
		}

		err = s.repo.CreateServiceKey(ctx, listDoc)
	case string(constants.RequestUpdateServiceList):
		listDoc, err := local_util.JsonUnmarshal[imodel.ServiceKey](action.CurrentAction)
		if err != nil {
			return nil, localization.ErrorInvalidActionData
		}
		prevListDoc, err := local_util.JsonUnmarshal[imodel.ServiceKey](action.PreviousAction)
		if err != nil {
			return nil, localization.ErrorInvalidActionData
		}

		err = s.repo.UpdateServiceKey(ctx, action.UniqueId, prevListDoc.ServiceKey, listDoc)
	case string(constants.RequestDeleteServiceList):
		err = s.repo.Delete(ctx, "", action.UniqueId)
	// case string(constants.RequestDeleteServiceKey):
	// 	err = s.repo.DeleteServiceKey(ctx, action.UniqueId)
	case string(constants.RequestEnableServiceList):
		// listDoc, err := local_util.JsonUnmarshal[model.ServiceKey](action.PreviousAction)
		// if err != nil {
		// 	return nil, localization.ErrorInvalidActionData
		// }
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, true)
	case string(constants.RequestDisableServiceList):
		// listDoc, err := local_util.JsonUnmarshal[model.ServiceKey](action.PreviousAction)
		// if err != nil {
		// 	return nil, localization.ErrorInvalidActionData
		// }

		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, false)

	default:
		return nil, localization.ErrorInvalidRequest
	}

	if err != nil {
		s.logger.Errorf("[servicesService][Authorize] error occurred: %v", err)
		return nil, err
	}

	action.CurrentAction = serviceDoc
	return action, nil
}
