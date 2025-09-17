package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound"
	local_utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type serviceHandler struct {
	appService ApplicationService
	logger     utils.Logger
}

func NewServiceHandler(domain ApplicationService, logger utils.Logger) inbound.Service {
	return &serviceHandler{
		appService: domain,
		logger:     logger,
	}
}

func (sa *serviceHandler) GetAllService(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		sa.logger.Errorf("Invalid pagination parameters: page=%d, perPage=%d", filterParams.Page, filterParams.PerPage)
		local_utils.SendErrorResponse(w, local_utils.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	data, err := sa.appService.GetAllService(r.Context(), local_utils.Filter(*filterParams))
	if err != nil {
		sa.logger.Errorf("Failed to get all services: %v", err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	local_utils.BaseResponseMaker(data, w, "Successfully retrieved services", 200)
}

func (sa *serviceHandler) GetAllMinimumTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		sa.logger.Errorf("Invalid pagination parameters: page=%d, perPage=%d", filterParams.Page, filterParams.PerPage)
		local_utils.SendErrorResponse(w, local_utils.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	res, err := sa.appService.GetAllMinimumTransferCap(r.Context(), local_utils.Filter(*filterParams))
	if err != nil {
		sa.logger.Errorf("Failed to get minimum transfer caps: %v", err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfully retrieved minimum transfer caps", 200)
}

func (sa *serviceHandler) GetAllMaximumTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		sa.logger.Errorf("Invalid pagination parameters: page=%d, perPage=%d", filterParams.Page, filterParams.PerPage)
		local_utils.SendErrorResponse(w, local_utils.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	res, err := sa.appService.GetAllMaximumTransferCap(r.Context(), local_utils.Filter(*filterParams))
	if err != nil {
		sa.logger.Errorf("Failed to get maximum transfer caps: %v", err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfully retrieved maximum transfer caps", 200)
}

func (sa *serviceHandler) GetAllServiceFee(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		sa.logger.Errorf("Invalid pagination parameters: page=%d, perPage=%d", filterParams.Page, filterParams.PerPage)
		local_utils.SendErrorResponse(w, local_utils.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	res, err := sa.appService.GetAllServiceFee(r.Context(), local_utils.Filter(*filterParams))
	if err != nil {
		sa.logger.Errorf("Failed to get service fees: %v", err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfully retrieved service fees", 200)
}

func (sa *serviceHandler) GetAllTotalTransferCap(w http.ResponseWriter, r *http.Request) {
	res, err := sa.appService.GetAllTotalTransferCap(r.Context())
	if err != nil {
		sa.logger.Errorf("Failed to get total transfer cap: %v", err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfully retrieved total transfer cap", 200)
}

func (sa *serviceHandler) GetServiceFeeDetail(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		sa.logger.Errorf("Missing service ID parameter")
		local_utils.SendErrorResponse(w, local_utils.InvalidID, http.StatusBadRequest, nil)
		return
	}

	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		sa.logger.Errorf("Invalid pagination parameters: page=%d, perPage=%d", filterParams.Page, filterParams.PerPage)
		local_utils.SendErrorResponse(w, local_utils.InvalidInputParameters, http.StatusBadRequest, nil)
		return
	}

	res, err := sa.appService.GetServiceFeeDetail(r.Context(), service_id)
	if err != nil {
		sa.logger.Errorf("Failed to get service fee detail for ID %s: %v", service_id, err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	data, err := local_utils.StructToMap(res)
	if err != nil {
		sa.logger.Errorf("Failed to convert service fee detail to map: %v", err)
		local_utils.SendErrorResponse(w, local_utils.UnhandledServerError, http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(data, w, "Successfully retrieved service fee detail", 200)
}

func (sa *serviceHandler) UpdateServiceFee(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		sa.logger.Errorf("Missing service ID parameter")
		local_utils.SendErrorResponse(w, local_utils.InvalidID, http.StatusBadRequest, nil)
		return
	}

	var req ServiceFeeDetailDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sa.logger.Errorf("Failed to decode service fee update request: %v", err)
		local_utils.SendErrorResponse(w, local_utils.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		sa.logger.Errorf("Service fee validation failed: %v", err)
		local_utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	err := sa.appService.UpdateServiceFee(r.Context(), service_id, req)
	if err != nil {
		sa.logger.Errorf("Failed to update service fee for ID %s: %v", service_id, err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Service fee update request submitted successfully", 200)
}

func (sa *serviceHandler) UpdateSingleMaxTransfer(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		sa.logger.Errorf("Missing service ID parameter")
		local_utils.SendErrorResponse(w, local_utils.InvalidID, http.StatusBadRequest, nil)
		return
	}

	var req SingleMaxTransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sa.logger.Errorf("Failed to decode single max transfer request: %v", err)
		local_utils.SendErrorResponse(w, local_utils.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		sa.logger.Errorf("Single max transfer validation failed: %v", err)
		local_utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	err := sa.appService.UpdateSingleMaxTransfer(r.Context(), service_id, req)
	if err != nil {
		sa.logger.Errorf("Failed to update single max transfer for ID %s: %v", service_id, err)
		if strings.Contains(err.Error(), "SINGLE_MAX_TRANSFER_CANNOT_BE_LESS_OR_EQUAL_TO_MIN_AMOUNT") {
			local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": err.Error()}})
		} else if strings.Contains(err.Error(), "individual_single_cap") || strings.Contains(err.Error(), "corporate_single_cap") {
			local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": err.Error()}})
		} else {
			local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		}
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Single max transfer update request submitted successfully", 200)
}

func (sa *serviceHandler) UpdateTotalMaxTransferCap(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		sa.logger.Errorf("Missing service ID parameter")
		local_utils.SendErrorResponse(w, local_utils.InvalidID, http.StatusBadRequest, nil)
		return
	}

	var req TotalMaxTransferUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sa.logger.Errorf("Failed to decode total max transfer cap request: %v", err)
		local_utils.SendErrorResponse(w, local_utils.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		sa.logger.Errorf("Total max transfer cap validation failed: %v", err)
		local_utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	converter := local_utils.NewNumericConverter()
	totalCap, err := converter.IsPositiveNumber(req.TotalTransferLimit, "total_cap")
	if err != nil {
		sa.logger.Errorf("Failed to extract total cap value: %v", err)
		local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": "Invalid total cap value"}})
		return
	}

	err = sa.appService.UpdateTotalMaxTransferCap(r.Context(), service_id, totalCap)
	if err != nil {
		sa.logger.Errorf("Failed to update total max transfer cap for ID %s: %v", service_id, err)
		if strings.Contains(err.Error(), "total_cap") || strings.Contains(err.Error(), "exceed the new total cap") {
			local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": err.Error()}})
		} else if strings.Contains(err.Error(), "individual_single_cap") || strings.Contains(err.Error(), "individual_daily_cap") || strings.Contains(err.Error(), "corporate_single_cap") || strings.Contains(err.Error(), "corporate_daily_cap") {
			local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": err.Error()}})
		} else {
			local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		}
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Total max transfer cap update request submitted successfully", 200)
}

func (sa *serviceHandler) UpdateMinimumTransferCap(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		sa.logger.Errorf("Missing service ID parameter")
		local_utils.SendErrorResponse(w, local_utils.InvalidID, http.StatusBadRequest, nil)
		return
	}

	var req MinimumTransferUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sa.logger.Errorf("Failed to decode minimum transfer cap request: %v", err)
		local_utils.SendErrorResponse(w, local_utils.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		sa.logger.Errorf("Minimum transfer cap validation failed: %v", err)
		local_utils.SendErrorResponse(w, err, http.StatusBadRequest, nil)
		return
	}

	err := sa.appService.UpdateMinimumTransferCap(r.Context(), service_id, req)
	if err != nil {
		sa.logger.Errorf("Failed to update minimum transfer cap for ID %s: %v", service_id, err)
		if strings.Contains(err.Error(), "min_amount") {
			local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": err.Error()}})
		} else if strings.Contains(err.Error(), "cannot exceed total cap") {
			local_utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": map[string]string{"business_rule": err.Error()}})
		} else {
			local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		}
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Minimum transfer cap update request submitted successfully", 200)
}

func (sa *serviceHandler) DeleteServiceFeeTire(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		sa.logger.Errorf("Missing service ID parameter")
		local_utils.SendErrorResponse(w, local_utils.InvalidID, http.StatusBadRequest, nil)
		return
	}

	err := sa.appService.DeleteServiceFeeTire(r.Context(), service_id)
	if err != nil {
		sa.logger.Errorf("Failed to delete service fee tier for ID %s: %v", service_id, err)
		local_utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Service fee tier deletion request submitted successfully", 200)
}
