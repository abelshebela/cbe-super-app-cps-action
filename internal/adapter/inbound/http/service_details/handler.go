package service_details_inbound

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	service_details_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/service_details"
	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/service_details"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/utils"
)

type HttpStore struct {
	Application service_details_app.ApplicationAbstracts
	logger      shared_utils.Logger
}

func NewHttpServiceDetails(app service_details_app.ApplicationAbstracts, logger shared_utils.Logger) inbound.ServiceDetailsInbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

func (h *HttpStore) handleError(w http.ResponseWriter, err error) {
	h.logger.Errorf("Error occurred: %v", err)

	statusCode := http.StatusInternalServerError
	message := "Internal server error"

	if err.Error() == "service ID cannot be empty" ||
		err.Error() == "maker ID cannot be empty" ||
		err.Error() == "checker ID cannot be empty" ||
		err.Error() == "action ID cannot be empty" {
		statusCode = http.StatusBadRequest
		message = err.Error()
	} else if err.Error() == "service not found" {
		statusCode = http.StatusNotFound
		message = err.Error()
	} else if err.Error() == "service already has a pending action" {
		statusCode = http.StatusConflict
		message = err.Error()
	}

	utils.WriteErrorResponse(w, statusCode, message)
}

func (h *HttpStore) GetAllServiceDetails(w http.ResponseWriter, r *http.Request) {
	services, err := h.Application.GetAllServiceDetails(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	utils.WriteSuccessResponse(w, services, "Successfully retrieved all service details")
}

func (h *HttpStore) GetServiceDetailsByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.handleError(w, fmt.Errorf("service ID cannot be empty"))
		return
	}

	service, err := h.Application.GetServiceDetailsByID(r.Context(), id)
	if err != nil {
		h.handleError(w, err)
		return
	}

	utils.WriteSuccessResponse(w, service, "Successfully retrieved service details")
}

func (h *HttpStore) UpdateServiceDetailsMaker(w http.ResponseWriter, r *http.Request) {
	var update dto.UpdateServiceDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		h.handleError(w, fmt.Errorf("unauthorized"))
		return
	}

	update.MakerID = claims.UserID
	response, err := h.Application.UpdateServiceDetailsRequest(r.Context(), update.ID, &update)
	if err != nil {
		h.handleError(w, err)
		return
	}

	utils.WriteSuccessResponse(w, response, "Update request submitted for approval")
}

func (h *HttpStore) UpdateServiceDetailsChecker(w http.ResponseWriter, r *http.Request) {
	var request dto.ApproveServiceDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		h.handleError(w, fmt.Errorf("unauthorized"))
		return
	}

	request.CheckerID = claims.UserID
	err := h.Application.UpdateServiceDetails(r.Context(), &request)
	if err != nil {
		h.handleError(w, err)
		return
	}

	action := "approved"
	if !request.Approve {
		action = "rejected"
	}
	utils.WriteSuccessResponse(w, nil, "Update request "+action+" successfully")
}

func (h *HttpStore) ServiceDetailsDailyCapMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.DailyCapServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	update, err := h.Application.GetServiceDetailsByID(r.Context(), req.ServiceId)
	if err != nil {
		h.handleError(w, err)
		return
	}
	update.Cap.DailyCap = req.DailyCap
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		h.handleError(w, fmt.Errorf("unauthorized"))
		return
	}
	res, err := h.Application.UpdateServiceDetailsRequest(r.Context(), claims.UserID, &dto.UpdateServiceDetailsRequest{
		ID:                 update.ID,
		ServiceID:          update.ServiceID,
		ServiceCode:        update.ServiceCode,
		ServiceName:        update.ServiceName,
		ServiceType:        update.ServiceType,
		Key:                update.Key,
		Cap:                update.Cap,
		CBEProductCodes:    update.CBEProductCodes,
		CBEIFBProductCodes: update.CBEIFBProductCodes,
		AboveAmount:        update.AboveAmount,
		AboveServiceFee:    update.AboveServiceFee,
		PaymentType:        update.PaymentType,
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	response := res.ActionID
	utils.WriteSuccessResponse(w, response, "success")
	//&update.MakerID = claims.UserID
}
func (h *HttpStore) ServiceDetailsSingleCapMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.SingleCapServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	update, err := h.Application.GetServiceDetailsByID(r.Context(), req.ServiceId)
	if err != nil {
		h.handleError(w, err)
		return
	}
	update.Cap.SingleCap = req.SingleCap
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		h.handleError(w, fmt.Errorf("unauthorized"))
		return
	}
	res, err := h.Application.UpdateServiceDetailsRequest(r.Context(), claims.UserID, &dto.UpdateServiceDetailsRequest{
		ID:                 update.ID,
		ServiceID:          update.ServiceID,
		ServiceCode:        update.ServiceCode,
		ServiceName:        update.ServiceName,
		ServiceType:        update.ServiceType,
		Key:                update.Key,
		Cap:                update.Cap,
		CBEProductCodes:    update.CBEProductCodes,
		CBEIFBProductCodes: update.CBEIFBProductCodes,
		AboveAmount:        update.AboveAmount,
		AboveServiceFee:    update.AboveServiceFee,
		PaymentType:        update.PaymentType,
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	response := res.ActionID
	utils.WriteSuccessResponse(w, response, "success")
}

func (h *HttpStore) TotalTransferCapMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.TotalTransferCapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	update, err := h.Application.GetServiceDetailsByID(r.Context(), req.ServiceId)
	if err != nil {
		h.handleError(w, err)
		return
	}
	update.Cap.MaxAmount = req.TotalCap
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		h.handleError(w, fmt.Errorf("unauthorized"))
		return
	}
	res, err := h.Application.UpdateServiceDetailsRequest(r.Context(), claims.UserID, &dto.UpdateServiceDetailsRequest{
		ID:                 update.ID,
		ServiceID:          update.ServiceID,
		ServiceCode:        update.ServiceCode,
		ServiceName:        update.ServiceName,
		ServiceType:        update.ServiceType,
		Key:                update.Key,
		Cap:                update.Cap,
		CBEProductCodes:    update.CBEProductCodes,
		CBEIFBProductCodes: update.CBEIFBProductCodes,
		AboveAmount:        update.AboveAmount,
		AboveServiceFee:    update.AboveServiceFee,
		PaymentType:        update.PaymentType,
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	response := res.ActionID
	utils.WriteSuccessResponse(w, response, "success")
}
