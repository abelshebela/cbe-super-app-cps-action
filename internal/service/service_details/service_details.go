package service_details

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/service_details/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceDetails struct {
	logger      utils.Logger
	serviceRepo storage.ServiceDetailsRepository
	hqRepo      storage.HQRepository
	cpsService  service.CPSActionService
}

func NewServiceDetailsService(client *mongo.Client, ServiceDetailsRepo storage.ServiceDetailsRepository, hqRepo storage.HQRepository, cpsAction service.CPSActionService, logger utils.Logger) service.ServiceService {
	return &ServiceDetails{
		logger:      logger,
		serviceRepo: ServiceDetailsRepo,
		hqRepo:      hqRepo,
		cpsService:  cpsAction,
	}
}
func (s *ServiceDetails) GetAllService(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ServiceDetails], error) {
	projection := bson.M{}
	return s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
}
func (s *ServiceDetails) GetAllMinimumTransferCap(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.MinimumTransferCapResponse], error) {
	projection := bson.M{}

	response, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		return nil, err
	}

	// Map the response to DTOs
	mappedData := core.MapSliceToMinimumTransferCapResponse(response.Data)

	return &types.PaginatedResponse[[]*dto.MinimumTransferCapResponse]{
		Data: mappedData,
		Meta: response.Meta,
	}, nil
}
func (s *ServiceDetails) GetAllMaximumTransferCap(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.MaximumTransferCapResponse], error) {
	projection := bson.M{}

	response, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		return nil, err
	}
	mappedData := core.MapSliceToMaximumTransferCapResponse(response.Data)

	return &types.PaginatedResponse[[]*dto.MaximumTransferCapResponse]{
		Data: mappedData,
		Meta: response.Meta,
	}, nil
}
func (s *ServiceDetails) GetAllServiceFee(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.ServiceFeeResponse], error) {
	projection := bson.M{}

	response, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		return nil, err
	}

	// Map the response to DTOs
	mappedData := core.MapSliceToServiceFeeResponse(response.Data)

	return &types.PaginatedResponse[[]*dto.ServiceFeeResponse]{
		Data: mappedData,
		Meta: response.Meta,
	}, nil
}
func (s *ServiceDetails) GetAllTotalTransferCap(ctx context.Context) (*dto.TotalTransferCapResponse, error) {
	projection := bson.M{}

	hq, err := s.hqRepo.Find(ctx, projection)
	if err != nil {
		return nil, err
	}

	// Map the response to DTO
	return core.MapToTotalTransferCapResponse(hq), nil
}
func (s *ServiceDetails) GetServiceFeeDetail(ctx context.Context, id string) (*dto.ServiceFeeDetailResponse, error) {
	projection := bson.M{}

	service, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		return nil, err
	}

	// Map the response to DTO
	return core.MapToServiceFeeDetailResponse(service), nil
}
func (s *ServiceDetails) UpdateServiceFee(ctx context.Context, id string, req dto.ServiceFeeDetailDTO) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	projection := bson.M{
		"tire": 1,
	}
	serviceDetail, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		return err
	}
	cpsAction := lib.CpsModelBuilder(id, makerData, serviceDetail, req, string(constants.RequestUpdateServiceFee), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}
func (s *ServiceDetails) UpdateSingleMaxTransfer(ctx context.Context, id string, req dto.SingleMaxTransferRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	projection := bson.M{
		"cap": 1,
	}
	prev, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		s.logger.Errorf("error fetching previous service data for update single max transfer: %v", err)
		return err
	}
	if req.CDailyCap <= prev.Cap.MinAmount || req.CSingleCap <= prev.Cap.MinAmount || req.IDailyCap <= prev.Cap.MinAmount || req.ISingleCap <= prev.Cap.MinAmount {
		return errors.New(localization.ErrorSingleMaxTransferCannotBeLessOrEqualToMinAmount.Code)
	}
	projection = bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	hq, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for updat single transfer cap: %v", err)
		return err
	}

	if req.CDailyCap > hq.TotalCap || req.IDailyCap > hq.TotalCap {
		return errors.New(localization.ErrorSingleTransferCanNotBeGreaterThanCap.Code)
	}
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, req, string(constants.RequestUpdateServiceSingle), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}
func (s *ServiceDetails) UpdateTotalMaxTransferCap(ctx context.Context, req dto.TotalMaxTransferUpdateRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	err := core.ValidateTotalCapAgainstServices(ctx, s.serviceRepo, s.logger, req.TotalTransferLimit)
	if err != nil {
		s.logger.Errorf("Total cap validation failed: %v", err)
		return err
	}
	projection := bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	prev, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for update total max transfer cap: %v", err)
		return err
	}

	cpsAction := lib.CpsModelBuilder(prev.ID.String(), makerData, prev, req, string(constants.RequestUpdateServiceTotal), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil
}
func (s *ServiceDetails) UpdateMinimumTransferCap(ctx context.Context, id string, req dto.MinimumTransferUpdateRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	projection := bson.M{
		"cap": 1,
	}
	prev, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		s.logger.Errorf("error fetching previous service data for update minimum transfer cap: %v", err)
		return err
	}

	newMinAmount := req.Minimum

	updatedCap := model.Cap{
		ISingleCap:         prev.Cap.ISingleCap,
		IDailyCap:          prev.Cap.IDailyCap,
		CorporateSingleCap: prev.Cap.CorporateSingleCap,
		CorporateDailyCap:  prev.Cap.CorporateDailyCap,
		MinAmount:          newMinAmount,
	}
	if updatedCap.MinAmount >= updatedCap.ISingleCap ||
		updatedCap.MinAmount >= updatedCap.IDailyCap ||
		updatedCap.MinAmount >= updatedCap.CorporateSingleCap ||
		updatedCap.MinAmount >= updatedCap.CorporateDailyCap {
		return errors.New(localization.ErrorMinAmountCanNotBeGreaterThanCap.Code)
	}
	projection = bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	hq, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for update minimum: %v", err)
		return err
	}

	if req.Minimum >= hq.TotalCap {
		return errors.New(localization.ErrorMinAmountCanNotBeGreaterThanTotal.Code)
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, prev, req, string(constants.RequestUpdateServiceMinCap), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil
}
func (s *ServiceDetails) DeleteServiceFeeTire(ctx context.Context, id string) error {
	s.logger.Infof("Deleting service fee tire for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	projection := bson.M{
		"tier": 1,
	}
	prev, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		s.logger.Errorf("error fetching previous service data for delete service fee tire: %v", err)
		return err
	}
	req := map[string]interface{}{"id": id, "IsDeleted": true}
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, req, string(constants.RequestDeleteServiceFee), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	s.logger.Infof("Successfully created CPS action for delete service fee tire, id: %s", id)
	return nil

}

func (s *ServiceDetails) applyServiceUpdate(ctx context.Context, cpsAction *model.CPSAction) error {

	serviceID := cpsAction.UniqueId
	if serviceID == "" {
		return errors.New("service ID is required (UniqueId field is empty)")
	}

	existingService, err := s.serviceRepo.FindByID(ctx, bson.M{}, serviceID)
	if err != nil {
		s.logger.Errorf("Failed to fetch existing service details: %v", err)
		return err
	}

	updatedService, err := core.MapServiceDetailsForUpdate(existingService, cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("Failed to map service details for update: %v", err)
		return err
	}

	return s.serviceRepo.Update(ctx, updatedService.ID.Hex(), updatedService)
}

func (s *ServiceDetails) applyTotalCapUpdate(ctx context.Context, cpsAction *model.CPSAction) error {
	projection := bson.M{}
	hq, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for updattotal maximum cap: %v", err)
		return err
	}

	totalCap, err := core.TotalCapMapper(hq, cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("Failed to map service details for update: %v", err)
		return err
	}

	return s.hqRepo.Update(ctx, "total_cap", totalCap.TotalCap, time.Now())
}

func (s *ServiceDetails) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	fmt.Println("ServiceDetails Authorize called")
	if cpsAction.ActionStatus != constants.Approved {
		s.logger.Errorf("Tried to authorize service action without cps action approval")
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	s.logger.Infof("Authorizing service action, action: %s", cpsAction.RequestAction)

	var err error
	switch cpsAction.RequestAction {
	case string(constants.RequestUpdateServiceDetails),
		string(constants.RequestUpdateServiceFee),
		string(constants.RequestUpdateServiceSingle),
		string(constants.RequestUpdateServiceMinCap),
		string(constants.RequestDeleteServiceFee):

		err = s.applyServiceUpdate(ctx, cpsAction)

	case string(constants.RequestUpdateServiceTotal):
		err = s.applyTotalCapUpdate(ctx, cpsAction)

	default:
		s.logger.Errorf("Unsupported action requested: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	if err != nil {
		s.logger.Errorf("Failed to process service action: %s, error: %v", cpsAction.RequestAction, err)
		return nil, err
	}

	s.logger.Infof("Service action authorization completed: %s", cpsAction.RequestAction)
	return cpsAction, nil
}
