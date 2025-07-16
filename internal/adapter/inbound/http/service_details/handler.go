package service_details_inbound

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	service_details_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/service_details"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"

	common "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"

	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/service_details"
	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type HttpStore struct {
	Application service_details_app.ApplicationAbstracts
	logger      shared_utils.Logger
}

type CurrentUser struct {
	Department  string
	UserCode    string
	FullName    string
	PhoneNumber string
}

func NewHttpServiceDetails(app service_details_app.ApplicationAbstracts, logger shared_utils.Logger) inbound.ServiceDetailsInbound {
	return &HttpStore{
		Application: app,
		logger:      logger,
	}
}

func (a *HttpStore) buildUserContext(r *http.Request) (*CurrentUser, error) {
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		return nil, fmt.Errorf(common_util.IncompleteUserInfo)
	}

	return &CurrentUser{
		Department:  userContext.Department,
		UserCode:    userContext.UserCode,
		FullName:    userContext.FullName,
		PhoneNumber: userContext.PhoneNumber,
	}, nil
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

	makerUser := contexts.ExtractUserContext(r)
	update.MakerID = makerUser.UserID
	response, err := h.Application.UpdateServiceDetailsRequest(r.Context(), update.ID, &update)
	if err != nil {
		h.handleError(w, err)
		return
	}

	utils.WriteSuccessResponse(w, response, "Update request submitted for approval")
}

func (h *HttpStore) UpdateServiceDetailsChecker(w http.ResponseWriter, r *http.Request) {
	var request dto.ApproveServiceDetailsRequest
	action_code := chi.URLParam(r, "action_code")

	request.ActionID = action_code
	if r.Method == "GET" {
		request.Approve = true
	} else {
		request.Approve = false
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			h.handleError(w, fmt.Errorf("invalid request body"))
			return
		}
	}
	checkerUser := contexts.ExtractUserContext(r)
	request.CheckerID = checkerUser.UserID
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
	curUser, curErr := h.buildUserContext(r)

	if curErr != nil {
		middleware.ErrorHandler(w, curErr)
		return
	}

	cpsAction.MakerUser = dto.User{
		UserCode:    curUser.UserCode,
		FullName:    curUser.FullName,
		PhoneNumber: curUser.PhoneNumber,
	}
	cpsAction.ActionData = req.Tries
	cpsAction.Department = curUser.Department

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
	curUser, curErr := h.buildUserContext(r)

	if curErr != nil {
		middleware.ErrorHandler(w, curErr)
		return
	}

	cpsAction.CheckerUser = dto.User{
		UserCode:    curUser.UserCode,
		FullName:    curUser.FullName,
		PhoneNumber: curUser.PhoneNumber,
	}

	cpsAction.Department = curUser.Department
	cpsAction.ActionData = req.Tier

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
	curUser, curErr := h.buildUserContext(r)

	if curErr != nil {
		middleware.ErrorHandler(w, curErr)
		return
	}

	cpsAction.CheckerUser = dto.User{
		UserCode:    curUser.UserCode,
		FullName:    curUser.FullName,
		PhoneNumber: curUser.PhoneNumber,
	}

	cpsAction.ActionData = req.Tier
	cpsAction.Department = curUser.Department

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
