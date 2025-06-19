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
	// "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/service"

	common "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"

	inbound "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/port/inbound/service_details"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/utils"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
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

// func (h *HttpStore) ServiceFeeMaker(w http.ResponseWriter, r *http.Request) {

// 	var request dto.ServiceFeeMakerRequest

// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		h.handleError(w, fmt.Errorf("invalid request body"))
// 		return
// 	}

// 	update, err := h.Application.GetServiceDetailsByID(r.Context(), request.ServiceID)
// 	if err != nil {
// 		h.handleError(w, err)
// 		return
// 	}

// 	update.Tiers = request.Tries

// 	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
// 	if !ok {
// 		h.handleError(w, fmt.Errorf("unauthorized"))
// 		return
// 	}

// 	response, err := h.Application.ServiceFeeMaker(r.Context())

// }

func (h *HttpStore) ServiceFeeApprove(w http.ResponseWriter, r *http.Request) {}
func (h *HttpStore) ServiceFeeReject(w http.ResponseWriter, r *http.Request)  {}
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
func (h *HttpStore) UpdateCapMinAmountHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.SingleCapServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	serviceDetails, err := h.Application.GetServiceDetailsByID(r.Context(), req.ServiceId)
	if err != nil {
		h.handleError(w, err)
		return
	}

	serviceDetails.Cap.MinAmount = req.SingleCap

	_, err = h.Application.UpdateServiceCap(r.Context(), req.ServiceId, &serviceDetails.Cap)
	if err != nil {
		h.handleError(w, err)
		return
	}

	utils.WriteSuccessResponse(w, nil, "Minimum cap amount updated successfully")
}
func (h *HttpStore) ApproveServiceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.ApproveServiceDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.handleError(w, fmt.Errorf("invalid request body"))
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		h.handleError(w, fmt.Errorf("unauthorized"))
		return
	}
	req.CheckerID = claims.UserID

	err := h.Application.ApproveServiceDetails(r.Context(), &req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	action := "approved"
	if !req.Approve {
		action = "rejected"
	}
	utils.WriteSuccessResponse(w, nil, "Service details "+action+" successfully")
}

func (h *HttpStore) InitiateServiceFeeUpdate(w http.ResponseWriter, r *http.Request) {
	var req dto.ServiceFeeMakerRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind action data", err)
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	var cpsAction dto.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.MakerUser = dto.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}
	cpsAction.ActionData = req.Tries
	cpsAction.Department = department

	ctx := r.Context()
	updateFeeResponse, err := h.Application.InitiateServiceFeeUpdate(ctx, &cpsAction)

	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*dto.UpdateServiceDetailsResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           updateFeeResponse,
	}

	res.SendJSON()

}
func (h *HttpStore) ApproveServiceFeeUpdate(w http.ResponseWriter, r *http.Request) {
	var req dto.ActionData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind action data", err)
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	var cpsAction dto.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.CheckerUser = dto.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}

	cpsAction.ActionData = req.Tier
	cpsAction.Department = department

	ctx := r.Context()
	updateFeeResponse, err := h.Application.ApproveServiceFeeUpdate(ctx, &cpsAction)

	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*dto.UpdateServiceDetailsResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           updateFeeResponse,
	}

	res.SendJSON()
}

func (h *HttpStore) RejectServiceFeeUpdate(w http.ResponseWriter, r *http.Request) {
	var req dto.RejectCPSAction

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("failed to bind action data", err)
		err = fmt.Errorf("failed to bind error data %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid request",
		})
		middleware.ErrorHandler(w, err)
		return

	}

	var cpsAction dto.CPSAction

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction.CheckerUser = dto.User{
		UserCode:    user_code,
		FullName:    full_name,
		PhoneNumber: phone_number,
	}

	cpsAction.ActionData = req.Tier
	cpsAction.Department = department

	ctx := r.Context()
	authorizeFayda, err := h.Application.ApproveServiceFeeUpdate(ctx, &cpsAction)
	if err != nil {
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*dto.UpdateServiceDetailsResponse]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           authorizeFayda,
	}

	res.SendJSON()

}
