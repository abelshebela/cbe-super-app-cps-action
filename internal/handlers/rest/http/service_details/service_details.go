package service_details

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/service_details"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/service_details/core"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type serviceAdapter struct {
	logger     utils.Logger
	serviceApp service.ServiceService
}

type paginatedServiceDetails types.PaginatedResponse[[]*model.ServiceDetails]
type paginatedMinimumTransferCapResponse types.PaginatedResponse[[]*dto.MinimumTransferCapResponse]
type paginatedMaximumTransferCapResponse types.PaginatedResponse[[]*dto.MaximumTransferCapResponse]
type paginatedServiceFeeResponse types.PaginatedResponse[[]*dto.ServiceFeeResponse]

func InitServiceAdapter(serviceApp service.ServiceService, logger utils.Logger) inbound.ServiceAdapter {
	return &serviceAdapter{
		logger:     logger,
		serviceApp: serviceApp,
	}
}

// Get All Services
// @Summary Get All Services
// @Description Retrieves all services with pagination
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} localization.StandardResponse{data=paginatedServiceDetails}
// @Failure 400,500 {object} localization.StandardResponse{data=nil}
// @Router /service [get]
func (s *serviceAdapter) GetAllService(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	services, err := s.serviceApp.GetAllService(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, services)
}

// Get All Minimum Transfer Caps
// @Summary Get All Minimum Transfer Caps
// @Description Retrieves all minimum transfer caps with pagination
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} localization.StandardResponse{data=paginatedMinimumTransferCapResponse}
// @Failure 400,500 {object} localization.StandardResponse{data=nil}
// @Router /service/minimum [get]
func (s *serviceAdapter) GetAllMinimumTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	minTransferCaps, err := s.serviceApp.GetAllMinimumTransferCap(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, minTransferCaps)
}

// Get All Maximum Transfer Caps
// @Summary Get All Maximum Transfer Caps
// @Description Retrieves all maximum transfer caps with pagination
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} localization.StandardResponse{data=paginatedMaximumTransferCapResponse}
// @Failure 400,500 {object} localization.StandardResponse{data=nil}
// @Router /service/maximum [get]
func (s *serviceAdapter) GetAllMaximumTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	maxTransferCaps, err := s.serviceApp.GetAllMaximumTransferCap(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, maxTransferCaps)
}

// Get All Service Fees
// @Summary Get All Service Fees
// @Description Retrieves all service fees with pagination
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} localization.StandardResponse{data=paginatedServiceFeeResponse}
// @Failure 400,500 {object} localization.StandardResponse{data=nil}
// @Router /service/service_fee [get]
func (s *serviceAdapter) GetAllServiceFee(w http.ResponseWriter, r *http.Request) {
	filterParams := local_util.ExtractFilterParams(r)
	if filterParams.Page < 0 || filterParams.PerPage < 0 {
		localization.SendErrorResponse(w, localization.ErrorInvalidRequest, nil, nil)
		return
	}
	ctx := r.Context()

	serviceFees, err := s.serviceApp.GetAllServiceFee(ctx, filterParams)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, serviceFees)
}

// Get All Total Transfer Caps
// @Summary Get All Total Transfer Caps
// @Description Retrieves all total transfer caps
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Success 200 {object} localization.StandardResponse{data=dto.TotalTransferCapResponse}
// @Failure 400,500 {object} localization.StandardResponse{data=nil}
// @Router /service/total/transfer_cap [get]
func (s *serviceAdapter) GetAllTotalTransferCap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalTransferCaps, err := s.serviceApp.GetAllTotalTransferCap(ctx)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, totalTransferCaps)
}

// Get Service Fee Detail
// @Summary Get Service Fee Detail
// @Description Retrieves service fee detail by service ID
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} localization.StandardResponse{data=dto.ServiceFeeDetailResponse}
// @Failure 400,404,500 {object} localization.StandardResponse{data=nil}
// @Router /service/service_fee/detail/{id} [get]
func (s *serviceAdapter) GetServiceFeeDetail(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}
	ctx := r.Context()

	serviceFeeDetails, err := s.serviceApp.GetServiceFeeDetail(ctx, service_id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, serviceFeeDetails)
}

// Update Service Fee
// @Summary Update Service Fee
// @Description Updates service fee detail by service ID
// @Tags ServiceDetails
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param body body dto.ServiceFeeDetailDTO true "Service Fee Detail DTO"
// @Success 200 {object} localization.StandardResponse{data=nil}
// @Failure 400,404,422,500 {object} localization.StandardResponse{data=nil}
// @Router /service/service_fee/update/{id} [put]
func (s *serviceAdapter) UpdateServiceFee(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}

	var req dto.ServiceFeeDetailDTO
	if !core.DecodeJSONBody(w, r, &req, s.logger) {
		return
	}

	if err := req.Validate(); err != nil {
		s.logger.Errorf("Service fee validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	ctx := r.Context()

	err := s.serviceApp.UpdateServiceFee(ctx, service_id, req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCreatedForUpdateServiceFee, nil)
}

// Update Single Max Transfer
// @Summary Update Single Max Transfer
// @Description Updates single max transfer by service ID
// @Tags ServiceDetails
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param body body dto.SingleMaxTransferRequest true "Single Max Transfer DTO"
// @Success 200 {object} localization.StandardResponse{data=nil}
// @Failure 400,404,422,500 {object} localization.StandardResponse{data=nil}
// @Router /service/single_transfer_max/update/id} [put]
func (s *serviceAdapter) UpdateSingleMaxTransfer(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}

	var req dto.SingleMaxTransferRequest
	if !core.DecodeJSONBody(w, r, &req, s.logger) {
		return
		return
	}

	if err := req.Validate(); err != nil {
		s.logger.Errorf("Failed to decode single max transfer request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	ctx := r.Context()

	err := s.serviceApp.UpdateSingleMaxTransfer(ctx, service_id, req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCreatedForUpdateSingleMaxTransfer, nil)

}

// Update Total Max Transfer Cap
// @Summary Update Total Max Transfer Cap
// @Description Updates total max transfer cap
// @Tags ServiceDetails
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.TotalMaxTransferUpdateRequest true "Total Max Transfer Update DTO"
// @Success 200 {object} localization.StandardResponse{data=nil}
// @Failure 400,404,422,500 {object} localization.StandardResponse{data=nil}
// @Router /service/total_transfer_max/update [put]
func (s *serviceAdapter) UpdateTotalMaxTransferCap(w http.ResponseWriter, r *http.Request) {

	var req dto.TotalMaxTransferUpdateRequest
	if !core.DecodeJSONBody(w, r, &req, s.logger) {
		return
		return
	}

	if err := req.Validate(); err != nil {
		s.logger.Errorf("Failed to decode total max transfer request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	ctx := r.Context()

	err := s.serviceApp.UpdateTotalMaxTransferCap(ctx, req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCreatedForUpdateTotalMaxTransferCap, nil)
}

// Update Minimum Transfer Cap
// @Summary Update Minimum Transfer Cap
// @Description Updates minimum transfer cap by service ID
// @Tags ServiceDetails
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Service ID"
// @Param body body dto.MinimumTransferUpdateRequest true "Minimum Transfer Update DTO"
// @Success 200 {object} localization.StandardResponse{data=nil}
// @Failure 400,404,422,500 {object} localization.StandardResponse{data=nil}
// @Router /service/minimum_transfer/update/{id} [put]
func (s *serviceAdapter) UpdateMinimumTransferCap(w http.ResponseWriter, r *http.Request) {

	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}

	var req dto.MinimumTransferUpdateRequest
	if !core.DecodeJSONBody(w, r, &req, s.logger) {
		return
		return
	}

	if err := req.Validate(); err != nil {
		s.logger.Errorf("Failed to decode total minimum transfer request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	ctx := r.Context()

	err := s.serviceApp.UpdateMinimumTransferCap(ctx, service_id, req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCreatedForUpdateMinimumTransferCap, nil)

}

// Delete Service Fee Tire
// @Summary Delete Service Fee Tire
// @Description Deletes service fee tire by service ID
// @Tags ServiceDetails
// @Security BearerAuth
// @Produce json
// @Param id path string true "Service ID"
// @Success 200 {object} localization.StandardResponse{data=nil}
// @Failure 400,404,500 {object} localization.StandardResponse{data=nil}
// @Router /service/service_fee/delete/{id} [delete]
func (s *serviceAdapter) DeleteServiceFeeTire(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}
	ctx := r.Context()
	err := s.serviceApp.DeleteServiceFeeTire(ctx, service_id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDeleteRequestCreated, nil)
}
