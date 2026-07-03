package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

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
	accessList_cache "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/catch/access_list"
	service_cache "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/catch/service"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type servicesService struct {
	repo            storage.ServicesRepository
	ussdMerchant    storage.UssdMerchantRepository
	cps             service.CPSActionService
	core            coreio.CBECoreAPIInterface
	logger          utils.Logger
	serviceCache    service_cache.ServiceCatch
	accessListCache accessList_cache.AccessListCatch
}

func NewServicesService(repo storage.ServicesRepository, ussdMerchant storage.UssdMerchantRepository, cps service.CPSActionService, core coreio.CBECoreAPIInterface, serviceCache service_cache.ServiceCatch, accessListCache accessList_cache.AccessListCatch, logger utils.Logger) *servicesService {
	return &servicesService{
		repo:            repo,
		ussdMerchant:    ussdMerchant,
		cps:             cps,
		core:            core,
		serviceCache:    serviceCache,
		accessListCache: accessListCache,
		logger:          logger,
	}
}

func (s *servicesService) Create(ctx context.Context, req service_dto.CreateServiceRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("Service creating...")

	accessList, err := s.repo.FindServiceListByID(ctx, req.ServiceKeyId)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		log.Errorf("[servicesService][Create] error checking existing service for serviceKeyId=%s: %v", req.ServiceKeyId, err)
		return err
	}

	var service *service_dto.ServiceResponse
	if accessList != nil {
		service, err = s.repo.FindServiceByAccessListID(ctx, req.ServiceKeyId)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			log.Errorf("[servicesService][Create] error checking existing service for serviceKey=%s: %v", req.ServiceKey, err)
			return err
		}

		if accessList != nil && service != nil && strings.EqualFold(req.ServiceCode, service.ServiceCode) {
			log.Warnf("[servicesService][Create] duplicate service detected for serviceKey=%s", req.ServiceKey)
			return errors.New(localization.ErrorServiceExists.Code)
		}
	}

	if req.ProductGlAccount != "" {
		var accountDetail *model.AccountDetail
		var accountDetailForPl model.AccountDetail
		if !*req.IsPlAccount && req.ProductGlAccount != "" {
			accountDetail, err = s.ValidateAccountNumberWithExternalAPI(ctx, req.ProductGlAccount)
			if err != nil {
				log.Errorf("[servicesService][Authorize] account number validation failed for account number %s: %v", req.ProductGlAccount, err)
				return errors.New(localization.ErrorAccountNumberValidationFailed.Code)
			}
		}

		var accountID string
		if req.ProductGlAccount != "" {
			accountID, err = s.repo.CheckAccountNumberExistence(ctx, req.ProductGlAccount)
			if err != nil {
				log.Errorf("failed while checking account number existence: %v", err)
				return err
			}
		}

		if *req.IsPlAccount && accountDetail == nil {
			accountDetailForPl = model.AccountDetail{
				CustomerName:  accessList.ServiceName,
				AccountNumber: req.ProductGlAccount,
				AccountType:   "PL Account",
				Currency:      "ETB",
				CustomerID:    req.ProductGlAccount,
			}

			accountDetail = &accountDetailForPl
		}

		if req.ProductGlAccount != "" && accountDetail != nil && accountID == "" {
			log.Errorf("failed while inserting account number to ACCOUNTS: %v", err)
			accountID, err = s.repo.InsertAccountNumberToAccounts(ctx, *accountDetail)
			if err != nil {
				return err
			}
		}
	}

	mapped := core.MapToServiceModel(req, *accessList)

	return core.HandleCPSAction(ctx, s.cps, "", constants.RequestCreateService, mapped, nil, constants.ActionCreate)
}

func (s *servicesService) Update(ctx context.Context, id string, req service_dto.UpdateServiceRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[servicesService][Update] called with id=%s, req=%+v", id, req)
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[servicesService][Update] error fetching previous service by id=%s: %v", id, err)
		return err
	}

	serviceKeyId := service_dto.StringPointer(req.ServiceKeyId, prev.ServiceKeyId)

	accessList, err := s.repo.FindServiceListByID(ctx, serviceKeyId)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		log.Errorf("[servicesService][Update] error checking existing service for serviceKeyId=%s: %v", req.ServiceKeyId, err)
		return err
	}

	var service *service_dto.ServiceResponse
	if accessList != nil {
		if req.ServiceKeyId != nil {
			service, err = s.repo.FindServiceByAccessListID(ctx, *req.ServiceKeyId)
			if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
				log.Errorf("[servicesService][Create] error checking existing service for serviceKey=%s: %v", req.ServiceKey, err)
				return err
			}
		}

		if service != nil && service.ID != id && strings.EqualFold(*req.ServiceCode, service.ServiceCode) {
			log.Warnf("[servicesService][Create] duplicate service detected for serviceKey=%s", req.ServiceKey)
			return errors.New(localization.ErrorServiceExists.Code)
		}
	}

	if req.ProductGlAccount != nil {
		var accountDetail *model.AccountDetail
		var accountDetailForPl model.AccountDetail
		if !*req.IsPlAccount && *req.ProductGlAccount != "" {
			accountDetail, err = s.ValidateAccountNumberWithExternalAPI(ctx, *req.ProductGlAccount)
			if err != nil {
				log.Errorf("[servicesService][Authorize] account number validation failed for account number %s: %v", *req.ProductGlAccount, err)
				return errors.New(localization.ErrorAccountNumberValidationFailed.Code)
			}
		}

		var accountID string
		if *req.ProductGlAccount != "" {
			accountID, err = s.repo.CheckAccountNumberExistence(ctx, *req.ProductGlAccount)
			if err != nil {
				log.Errorf("failed while checking account number existence: %v", err)
				return errors.New(localization.ErrorUnexpectedError.Code)
			}
		}

		if *req.IsPlAccount && accountDetail == nil {
			accountDetailForPl = model.AccountDetail{
				CustomerName:  accessList.ServiceName,
				AccountNumber: *req.ProductGlAccount,
				AccountType:   "PL Account",
				Currency:      "ETB",
				CustomerID:    *req.ProductGlAccount,
			}

			accountDetail = &accountDetailForPl
		}

		if *req.ProductGlAccount != "" && accountDetail != nil && accountID == "" {
			log.Errorf("failed while inserting account number to ACCOUNTS: %v", err)
			accountID, err = s.repo.InsertAccountNumberToAccounts(ctx, *accountDetail)
			if err != nil {
				return errors.New(localization.ErrorUnexpectedError.Code)
			}
		}

	}

	mapped := core.MapToServiceUpdateModel(req, *prev)
	log.Infof("[servicesService][Update] mapped update model: %+v", mapped)
	err = core.HandleCPSAction(ctx, s.cps, id, constants.RequestUpdateService, mapped, prev, constants.ActionUpdate)
	if err != nil {
		log.Errorf("[servicesService][Update] HandleCPSAction failed: %v", err)
		return err
	}
	log.Infof("[servicesService][Update] service update successful for id=%s", id)
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
	log := local_util.LoggerFromCtx(ctx, s.logger)
	if id == "" {
		return localization.ErrorInvalidID
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[ServiceSVC][DeleteService] failed to delete error: %v", err)
		return err
	}

	// Check service
	wal, err := s.repo.FindWalletByServiceId(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return err
	}
	if wal {
		return errors.New("There is an active wallet connected with this service")
	}

	// Check Donation
	don, err := s.repo.FindDonationByServiceId(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return err
	}
	if don {
		return errors.New("There is an active donation connected with this service")
	}

	// Check ussd_merchants
	if s.ussdMerchant != nil {
		ussdMerchant, err := s.ussdMerchant.Find(ctx, bson.M{
			"service": id,
		})
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
		if strings.EqualFold(ussdMerchant.Service, id) {
			return errors.New("There is an active ussd merchant connected with this service")
		}
	}

	deleted := core.BuildServiceDeleteSnapshot(*prev)
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestDeleteService, deleted, prev, constants.ActionDelete)
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

	mapped := core.MapServiceListDtoUpdateToModel(*existing, req)
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
	log := local_util.LoggerFromCtx(ctx, s.logger)
	prev, err := s.repo.FindServiceListByID(ctx, id)
	if err != nil {
		if err.Error() == localization.ErrorAccessListNotFound.Code {
			return localization.ErrorAccessListNotFound
		}
		log.Errorf("[servicesService][DeleteServiceKey] error fetching access list by id=%s: %v", id, err)
		return err
	}

	// Check service
	service, err := s.repo.FindServiceByAccessListID(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return err
	}
	if service != nil {
		return errors.New("There is an active service with this access list")
	}

	// Check Access list by superapp role
	sar, err := s.repo.FindSupperAppRoleByAccessList(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return err
	}
	if sar {
		return errors.New("There is an active customer segmentation with this access list")
	}

	// Check Geographical Area
	geo, err := s.repo.FindGeographicalLocationByAccessList(ctx, id)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return err
	}
	if geo {
		return errors.New("There is an active geographical_location with this access list")
	}

	s.logger.Infof("[servicesService][DeleteServiceKey] Deleting service key with id=%s, found service list: %+v", id, prev)
	deleted := core.BuildServiceKeyDeleteSnapshot(*prev)
	return core.HandleCPSAction(ctx, s.cps, id, constants.RequestDeleteServiceList, deleted, prev, constants.ActionDelete)
}

func (s *servicesService) ValidateAccountNumberWithExternalAPI(ctx context.Context, accountNumber string) (*model.AccountDetail, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	response, err := s.core.NameLookup(ctx, coreio.NameLookupParam{AccountNumber: accountNumber})
	if err != nil {
		return nil, err
	}

	if !response.Success {
		var message string
		for _, msg := range response.Messages {
			message += msg
		}

		log.Warnf("(core) failed to get account details: %s", message)
		return nil, fmt.Errorf("account lookup failed: %s", message)
	}

	if response.Detail == nil {
		log.Errorf("account lookup successful but no account details found for account number %s", accountNumber)
		return nil, fmt.Errorf("no account details found for account number %s", accountNumber)
	}

	detail := response.Detail
	return &model.AccountDetail{
		AccountNumber:  detail.AccountNumber,
		CustomerName:   detail.AccountName,
		Restriction:    detail.RestrictionType,
		Currency:       detail.Currency,
		WorkingBalance: "",
		CustomerID:     detail.CustomerNumber,
		AccountType:    detail.RestrictionType,
	}, nil
}

func (s *servicesService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[servicesService][Authorize] Authorize called for action: %+v", action)

	var err error

	switch action.RequestAction {
	case string(constants.RequestDeleteService):
		err = s.repo.Delete(ctx, action.UniqueId, "")

		// Get data for cache
		serviceDoc, unmarshalErr := local_util.JsonUnmarshal[imodel.Service](action.CurrentAction)
		if unmarshalErr != nil {
			return nil, localization.ErrorInvalidActionData
		}

		var sources []string
		for _, c := range serviceDoc.Cap {
			sources = append(sources, string(c.Source))
		}
		sourceJoined := strings.Join(sources, ",")

		// Delete from cache
		s.serviceCache.Delete(ctx, service_cache.ServiceKey{
			Source:      service_cache.Source(sourceJoined),
			ServiceKey:  serviceDoc.ServiceKey,
			ServiceCode: serviceDoc.ServiceCode,
		})

	case string(constants.RequestDeleteServiceList):
		err = s.repo.Delete(ctx, "", action.UniqueId)

		// Get data for cache
		listDoc, unmarshalErr := local_util.JsonUnmarshal[imodel.ServiceKey](action.CurrentAction)
		if unmarshalErr != nil {
			return nil, localization.ErrorInvalidActionData
		}

		var source accessList_cache.Source
		if listDoc.IsSuperAppEnabled && listDoc.IsUSSDEnabled {
			source = accessList_cache.SourceBoth
		} else if listDoc.IsSuperAppEnabled {
			source = accessList_cache.SourceAPP
		} else if listDoc.IsUSSDEnabled {
			source = accessList_cache.SourceUSSD
		}

		// Delete from cache
		s.accessListCache.Delete(ctx, accessList_cache.AccessListKey{
			Source:     source,
			ServiceKey: listDoc.ServiceKey,
		})

	case string(constants.RequestCreateService):
		serviceDoc, unmarshalErr := local_util.JsonUnmarshal[imodel.Service](action.CurrentAction)
		if unmarshalErr != nil {
			return nil, localization.ErrorInvalidActionData
		}
		s.logger.Infof("[servicesService][Authorize] Authorizing create service with data: %+v", serviceDoc)
		err = s.repo.Create(ctx, serviceDoc.ProductGlAccount, serviceDoc)
		if err == nil {
			action.CurrentAction = serviceDoc
		}

		var values []service_cache.ServiceData

		for _, c := range serviceDoc.Cap {
			values = append(values, service_cache.ServiceData{
				ServiceID:          action.UniqueId,
				AccessListID:       serviceDoc.ServiceKeyId,
				ServiceKey:         serviceDoc.ServiceKey,
				ServiceCode:        serviceDoc.ServiceCode,
				MinimumFraudAmount: serviceDoc.MinimumFraudAmount,
				ProductGlAccount:   serviceDoc.ProductGlAccount,
				ServiceName:        serviceDoc.ServiceName,
				Source:             service_cache.Source(c.Source),
				Currency:           c.Currency,
				SingleTransferCap:  c.SingleCap,
				MinimumTransferCap: c.MinimumTransferCap,
			})

		}
		s.serviceCache.Set(ctx, values)

	case string(constants.RequestUpdateService):
		serviceDoc, unmarshalErr := local_util.JsonUnmarshal[imodel.Service](action.CurrentAction)
		if unmarshalErr != nil {
			return nil, localization.ErrorInvalidActionData
		}

		s.logger.Infof("[servicesService][Authorize] Authorizing update service with data: %+v", serviceDoc)
		err = s.repo.Update(ctx, action.UniqueId, serviceDoc, serviceDoc.ProductGlAccount)
		if err == nil {
			action.CurrentAction = serviceDoc
		}

		var values []service_cache.ServiceData

		for _, c := range serviceDoc.Cap {
			values = append(values, service_cache.ServiceData{
				ServiceID:          action.UniqueId,
				AccessListID:       serviceDoc.ServiceKeyId,
				ServiceKey:         serviceDoc.ServiceKey,
				ServiceCode:        serviceDoc.ServiceCode,
				MinimumFraudAmount: serviceDoc.MinimumFraudAmount,
				ProductGlAccount:   serviceDoc.ProductGlAccount,
				ServiceName:        serviceDoc.ServiceName,
				Source:             service_cache.Source(c.Source),
				Currency:           c.Currency,
				SingleTransferCap:  c.SingleCap,
				MinimumTransferCap: c.MinimumTransferCap,
			})

		}
		s.serviceCache.Update(ctx, values)

	case string(constants.RequestEnableService):
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, true)
	case string(constants.RequestDisableService):
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, false)
	case string(constants.RequestCreateServiceList):
		listDoc, unmarshalErr := local_util.JsonUnmarshal[imodel.ServiceKey](action.CurrentAction)
		if unmarshalErr != nil {
			return nil, localization.ErrorInvalidActionData
		}
		err = s.repo.CreateServiceKey(ctx, listDoc)
		if err == nil {
			action.CurrentAction = listDoc
		}

		// Set to cache
		var source accessList_cache.Source
		if listDoc.IsSuperAppEnabled && listDoc.IsUSSDEnabled {
			source = accessList_cache.SourceBoth
		} else if listDoc.IsSuperAppEnabled {
			source = accessList_cache.SourceAPP
		} else if listDoc.IsUSSDEnabled {
			source = accessList_cache.SourceUSSD
		}

		s.accessListCache.Set(ctx, accessList_cache.AccessListData{
			ID:          action.UniqueId,
			Source:      source,
			ServiceName: listDoc.ServiceName,
			ServiceKey:  listDoc.ServiceKey,
			AccountType: accessList_cache.AccountType(listDoc.AccountType),
		})

	case string(constants.RequestUpdateServiceList):
		listDoc, unmarshalErr := local_util.JsonUnmarshal[imodel.ServiceKey](action.CurrentAction)
		if unmarshalErr != nil {
			return nil, localization.ErrorInvalidActionData
		}
		prevListDoc, prevErr := local_util.JsonUnmarshal[imodel.ServiceKey](action.PreviousAction)
		if prevErr != nil {
			return nil, localization.ErrorInvalidActionData
		}
		err = s.repo.UpdateServiceKey(ctx, action.UniqueId, prevListDoc.ServiceKey, listDoc)
		if err == nil {
			action.CurrentAction = listDoc
		}

		// Set to cache
		var source accessList_cache.Source
		if listDoc.IsSuperAppEnabled && listDoc.IsUSSDEnabled {
			source = accessList_cache.SourceBoth
		} else if listDoc.IsSuperAppEnabled {
			source = accessList_cache.SourceAPP
		} else if listDoc.IsUSSDEnabled {
			source = accessList_cache.SourceUSSD
		}

		s.accessListCache.Update(ctx,
			accessList_cache.AccessListData{
				ID:          action.UniqueId,
				Source:      source,
				ServiceName: listDoc.ServiceName,
				ServiceKey:  listDoc.ServiceKey,
				AccountType: accessList_cache.AccountType(listDoc.AccountType),
			})

	case string(constants.RequestEnableServiceList):
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, true)
	case string(constants.RequestDisableServiceList):
		err = s.repo.EnableOrDisableServiceList(ctx, action.UniqueId, false)
	default:
		return nil, localization.ErrorInvalidRequest
	}

	if err != nil {
		log.Errorf("[servicesService][Authorize] error occurred: %v", err)
		return nil, err
	}

	return action, nil
}
