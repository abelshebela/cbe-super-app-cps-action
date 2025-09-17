package service_details

import (
	inbound "cbe-super-app-cps-action/internal/constants/interfaces/service_details"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"
dto "cbe-super-app-cps-action/internal/constants/dto/service_details"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"cbe-super-app-cps-action/internal/handlers/rest/http/service_details/core"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)


type serviceAdapter struct {
	logger     utils.Logger
	serviceApp service.ServiceService
}
func InitServiceAdapter(serviceApp service.ServiceService, logger utils.Logger) inbound.ServiceAdapter {
	return &serviceAdapter{
		logger:     logger,
		serviceApp: serviceApp,
	}
}

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
func (s *serviceAdapter) GetAllTotalTransferCap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalTransferCaps, err := s.serviceApp.GetAllTotalTransferCap(ctx)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, totalTransferCaps)
}

func (s *serviceAdapter) GetServiceFeeDetail(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}
	ctx := r.Context()

	serviceFeeDetails, err := s.serviceApp.GetServiceFeeDetail(ctx,service_id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessServiceFeeDetailFetched, serviceFeeDetails)
}

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

func (s *serviceAdapter) UpdateTotalMaxTransferCap(w http.ResponseWriter, r *http.Request) {


	var req dto.TotalMaxTransferUpdateRequest
	if !core.DecodeJSONBody(w, r, &req, s.logger) {
		return 
	}

	if err := req.Validate(); err != nil {
		s.logger.Errorf("Failed to decode total max transfer request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	ctx := r.Context()

	err:= s.serviceApp.UpdateTotalMaxTransferCap(ctx, req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCreatedForUpdateTotalMaxTransferCap, nil)
}

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
	}

	if err := req.Validate(); err != nil {
		s.logger.Errorf("Failed to decode total minimum transfer request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}
	ctx := r.Context()

	err:= s.serviceApp.UpdateMinimumTransferCap(ctx, service_id, req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCPSActionCreatedForUpdateMinimumTransferCap, nil)

}

func (s *serviceAdapter) DeleteServiceFeeTire(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		s.logger.Errorf("service ID is required to fetch one")
		localization.SendErrorResponse(w, localization.ErrorServiceDetailIDRequired, nil, nil)
		return
	}
	ctx:= r.Context()
	err := s.serviceApp.DeleteServiceFeeTire(ctx, service_id)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDeleteRequestCreated, nil)
}