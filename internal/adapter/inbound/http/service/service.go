package service

import (
	"encoding/json"
	"net/http"

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
		local_utils.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	data, err := sa.appService.GetAllService(r.Context(), local_utils.Filter(*filterParams)) // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_utils.BaseResponseMaker(data, w, "Successfuly Service retrived", 200)
}
func (sa *serviceHandler) GetAllMinimumTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		local_utils.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	res, err := sa.appService.GetAllMinimumTransferCap(r.Context(), local_utils.Filter(*filterParams)) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfuly Service retrived", 200)
}
func (sa *serviceHandler) GetAllMaximumTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		local_utils.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	res, err := sa.appService.GetAllMaximumTransferCap(r.Context(), local_utils.Filter(*filterParams)) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfuly Service retrived", 200)
}
func (sa *serviceHandler) GetAllServiceFee(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		local_utils.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	res, err := sa.appService.GetAllServiceFee(r.Context(), local_utils.Filter(*filterParams)) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfuly Service retrived", 200)
}
func (sa *serviceHandler) GetAllTotalTransferCap(w http.ResponseWriter, r *http.Request) {
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		local_utils.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	res, err := sa.appService.GetAllTotalTransferCap(r.Context(), local_utils.Filter(*filterParams)) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(res, w, "Successfuly Service retrived", 200)
}
func (sa *serviceHandler) GetServiceFeeDetail(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	filterParams := local_utils.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		local_utils.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	res, err := sa.appService.GetServiceFeeDetail(r.Context(), service_id) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := local_utils.StructToMap(res)
	if err != nil {
		local_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}
	local_utils.BaseResponseMaker(data, w, "Successfuly Service retrived", 200)
}
func (sa *serviceHandler) UpdateServiceFee(w http.ResponseWriter, r *http.Request) {
	var req ServiceFeeDetailDTO
	service_id := chi.URLParam(r, "id")

	if service_id == "" {
		local_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	if req.Validate() != nil {
		local_utils.SendErrorResponse(w, req.Validate().Error(), 0, nil)
		return
	}

	err := sa.appService.UpdateServiceFee(r.Context(), string(service_id), req) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Successfuly Service tire updated", 200)
}
func (sa *serviceHandler) UpdateSingleMaxTransfer(w http.ResponseWriter, r *http.Request) {
	var req SingleMaxTransferRequest
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		local_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	err := sa.appService.UpdateSingleMaxTransfer(r.Context(), string(service_id), req) // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Successfuly Service single transfer updated", 200)
}
func (sa *serviceHandler) UpdateTotalMaxTransferCap(w http.ResponseWriter, r *http.Request) {
	var req TotalMaxTransferUpdateRequest
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		local_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	if req.Validate() != nil {
		local_utils.SendErrorResponse(w, req.Validate().Error(), 0, nil)
		return
	}

	err := sa.appService.UpdateTotalMaxTransferCap(r.Context(), string(service_id), req) // to be continued
	// res,err := sa.appService.  // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Successfuly Service single transfer updated", 200)
}
func (sa *serviceHandler) UpdateMinimumTransferCap(w http.ResponseWriter, r *http.Request) {
	var req MinimumTransferUpdateRequest
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		local_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	err := sa.appService.UpdateMinimumTransferCap(r.Context(), string(service_id), req) // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Successfuly Service min amount transfer updated", 200)
}
func (sa *serviceHandler) DeleteServiceFeeTire(w http.ResponseWriter, r *http.Request) {
	service_id := chi.URLParam(r, "id")
	if service_id == "" {
		local_utils.SendErrorResponse(w, "UNHANDLED_SERVER_ERROR", 0, nil)
		return
	}

	err := sa.appService.DeleteServiceFeeTire(r.Context(), string(service_id)) // to be continued
	if err != nil {
		local_utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	local_utils.BaseResponseMaker(map[string]interface{}{}, w, "Successfuly Delete service tier delete request", 200)
}
