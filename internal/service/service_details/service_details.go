package service_details

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"

	// "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/service_details/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllService", "ServiceDetails", "GetAllService")
	defer span.End()

	projection := bson.M{}
	result, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch services", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}
func (s *ServiceDetails) GetAllMinimumTransferCap(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.MinimumTransferCapResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllMinimumTransferCap", "ServiceDetails", "GetAllMinimumTransferCap")
	defer span.End()

	projection := bson.M{}

	response, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch minimum transfer caps", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
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
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllMaximumTransferCap", "ServiceDetails", "GetAllMaximumTransferCap")
	defer span.End()

	projection := bson.M{}

	response, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch maximum transfer caps", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	mappedData := core.MapSliceToMaximumTransferCapResponse(response.Data)

	return &types.PaginatedResponse[[]*dto.MaximumTransferCapResponse]{
		Data: mappedData,
		Meta: response.Meta,
	}, nil
}
func (s *ServiceDetails) GetAllServiceFee(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dto.ServiceFeeResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllServiceFee", "ServiceDetails", "GetAllServiceFee")
	defer span.End()

	projection := bson.M{}

	response, err := s.serviceRepo.FindAllWithPagination(ctx, projection, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch service fees", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
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
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllTotalTransferCap", "ServiceDetails", "GetAllTotalTransferCap")
	defer span.End()

	projection := bson.M{}

	hq, err := s.hqRepo.Find(ctx, projection)
	if err != nil {
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}

	// Map the response to DTO
	return core.MapToTotalTransferCapResponse(hq), nil
}
func (s *ServiceDetails) GetServiceFeeDetail(ctx context.Context, id string) (*dto.ServiceFeeDetailResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetServiceFeeDetail", "ServiceDetails", "GetServiceFeeDetail")
	defer span.End()

	projection := bson.M{}

	service, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		span.AddEvent("Failed to fetch service fee detail", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	// Map the response to DTO
	return core.MapToServiceFeeDetailResponse(service), nil
}
func (s *ServiceDetails) UpdateServiceFee(ctx context.Context, id string, req dto.ServiceFeeDetailDTO) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateServiceFee", "ServiceDetails", "UpdateServiceFee")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	projection := bson.M{
		"tiers": 1,
	}

	serviceDetail, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		span.AddEvent("Failed to find service detail", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	serviceData := *serviceDetail
	serviceMap := core.ServiceMapper(&serviceData, req)
	cpsAction := lib.CpsModelBuilder(id, makerData, serviceDetail, serviceMap, string(constants.RequestUpdateServiceFee), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}
func (s *ServiceDetails) UpdateSingleMaxTransfer(ctx context.Context, id string, req dto.SingleMaxTransferRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateSingleMaxTransfer", "ServiceDetails", "UpdateSingleMaxTransfer")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	projection := bson.M{
		"cap": 1,
	}
	prev, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		s.logger.Errorf("error fetching previous service data for update single max transfer: %v", err)
		span.AddEvent("Failed to find service detail", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if req.CDailyCap <= prev.Cap.MinAmount || req.CSingleCap <= prev.Cap.MinAmount || req.IDailyCap <= prev.Cap.MinAmount || req.ISingleCap <= prev.Cap.MinAmount {
		span.AddEvent("Single max transfer cannot be less or equal to min amount", trace.WithAttributes(
			attribute.String("error", localization.ErrorSingleMaxTransferCannotBeLessOrEqualToMinAmount.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorSingleMaxTransferCannotBeLessOrEqualToMinAmount.Code)
	}
	projection = bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	hq, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for updat single transfer cap: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if req.CDailyCap > hq.TotalCap || req.IDailyCap > hq.TotalCap {
		span.AddEvent("Single transfer cannot be greater than cap", trace.WithAttributes(
			attribute.String("error", localization.ErrorSingleTransferCanNotBeGreaterThanCap.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorSingleTransferCanNotBeGreaterThanCap.Code)
	}
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, req, string(constants.RequestUpdateServiceSingle), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}
func (s *ServiceDetails) UpdateTotalMaxTransferCap(ctx context.Context, req dto.TotalMaxTransferUpdateRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateTotalMaxTransferCap", "ServiceDetails", "UpdateTotalMaxTransferCap")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	err := core.ValidateTotalCapAgainstServices(ctx, s.serviceRepo, s.logger, req.TotalTransferLimit)
	if err != nil {
		s.logger.Errorf("Total cap validation failed: %v", err)
		span.AddEvent("Total cap validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	projection := bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	prev, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for update total max transfer cap: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	cpsAction := lib.CpsModelBuilder(prev.ID.String(), makerData, prev, req, string(constants.RequestUpdateServiceTotal), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	return nil
}
func (s *ServiceDetails) UpdateMinimumTransferCap(ctx context.Context, id string, req dto.MinimumTransferUpdateRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateMinimumTransferCap", "ServiceDetails", "UpdateMinimumTransferCap")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	projection := bson.M{
		"cap": 1,
	}
	prev, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		s.logger.Errorf("error fetching previous service data for update minimum transfer cap: %v", err)
		span.AddEvent("Failed to find service detail", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	newMinAmount := req.Minimum

	updatedCap := types.Cap{
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
		span.AddEvent("Min amount cannot be greater than cap", trace.WithAttributes(
			attribute.String("error", localization.ErrorMinAmountCanNotBeGreaterThanCap.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorMinAmountCanNotBeGreaterThanCap.Code)
	}
	projection = bson.M{
		"total_cap": 1,
		"_id":       1,
	}
	hq, err := s.hqRepo.Find(ctx, nil, projection)
	if err != nil {
		s.logger.Errorf("error fetching HQ data for update minimum: %v", err)
		span.AddEvent("Failed to fetch HQ", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if req.Minimum >= hq.TotalCap {
		span.AddEvent("Min amount cannot be greater than total", trace.WithAttributes(
			attribute.String("error", localization.ErrorMinAmountCanNotBeGreaterThanTotal.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorMinAmountCanNotBeGreaterThanTotal.Code)
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, prev, req, string(constants.RequestUpdateServiceMinCap), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	return nil
}
func (s *ServiceDetails) DeleteServiceFeeTire(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteServiceFeeTire", "ServiceDetails", "DeleteServiceFeeTire")
	defer span.End()

	s.logger.Infof("Deleting service fee tire for id: %s", id)
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	projection := bson.M{
		"tiers": 1,
	}
	prev, err := s.serviceRepo.FindByID(ctx, projection, id)
	if err != nil {
		s.logger.Errorf("error fetching previous service data for delete service fee tire: %v", err)
		span.AddEvent("Failed to find service detail", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	req := map[string]interface{}{"id": id, "IsDeleted": true}
	cpsAction := lib.CpsModelBuilder(id, makerData, prev, req, string(constants.RequestDeleteServiceFee), constants.UPDATE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
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

	return s.serviceRepo.Update(ctx, serviceID, updatedService)
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
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "ServiceDetails", "Authorize")
	defer span.End()

	fmt.Println("ServiceDetails Authorize called")
	if cpsAction.ActionStatus != constants.Approved {
		s.logger.Errorf("Tried to authorize service action without cps action approval")
		span.AddEvent("CPS action status invalid", trace.WithAttributes(
			attribute.String("error", localization.ErrorCPSActionStatusInvalid.Code),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
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
		if err != nil {
			span.AddEvent("Failed to apply service update", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateServiceTotal):
		err = s.applyTotalCapUpdate(ctx, cpsAction)
		if err != nil {
			span.AddEvent("Failed to apply total cap update", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, err
		}

	default:
		s.logger.Errorf("Unsupported action requested: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	s.logger.Infof("Service action authorization completed: %s", cpsAction.RequestAction)
	return cpsAction, nil
}
