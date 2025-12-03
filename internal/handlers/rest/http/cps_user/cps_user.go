package cps_user

import (
	"encoding/json"
	"net/http"
	"strings"

	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	localization "cbe-super-app-cps-action/internal/constants/localization"

	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type handler struct {
	svc    service.CPSUserService
	logger utils.Logger
}

func InitCPSUserHandler(svc service.CPSUserService, logger utils.Logger) *handler {
	return &handler{svc: svc, logger: logger}
}

// CreateUserRequest creates a new CPS user request
//
//	@Summary		Create CPS user request
//	@Description	Creates a new CPS user request that requires approval
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			request	body		cpsuser.CreateUserRequest				true	"CPS user creation request"
//	@Success		201		{object}	localization.StandardResponse{data=nil}	"CPS user creation request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/create [post]
func (h *handler) CreateUserRequest(w http.ResponseWriter, r *http.Request) {
	var req cpsuser.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[CreateUserRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}
	// Normalize and validate request
	req.Normalize()
	if err := req.Validate(); err != nil {
		h.logger.Errorf("[CreateUserRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	formattedPhone, err := local_util.ValidateAndNormalizePhoneNumber(req.PhoneNumber)
	if err != nil {
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	req.PhoneNumber = formattedPhone
	if err := h.svc.CreateUserRequest(r.Context(), req); err != nil {
		h.logger.Errorf("[CreateUserRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserCreationRequestSubmitted, nil)
}

// UpdateUserRequest updates an existing CPS user request
//
//	@Summary		Update CPS user request
//	@Description	Updates an existing CPS user request that requires approval
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Param			request		body		cpsuser.UpdateUserRequest				true	"CPS user update request"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user update request submitted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - Invalid input"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/update/{user_code} [patch]
func (h *handler) UpdateUserRequest(w http.ResponseWriter, r *http.Request) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	var req cpsuser.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[UpdateUserRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.MsgInvalidInput)
		return
	}

	// ensure the request carries the target user code from the path
	req.UserCode = userCode

	req.Normalize()

	if err := req.Validate(); err != nil {
		h.logger.Errorf("[UpdateUserRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := h.svc.UpdateUserRequest(r.Context(), req); err != nil {
		h.logger.Errorf("[UpdateUserRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserUpdateRequestSubmitted, nil)
}

// FetchUserByUserCode retrieves a CPS user by user code
//
//	@Summary		Get CPS user by user code
//	@Description	Retrieves a specific CPS user by their user code
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//
// // @Success 200 {object} localization.StandardResponse{data=cps_user_resp} "CPS user retrieved successfully"
//
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/{user_code} [get]
func (h *handler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	// Build detailed response in the service layer
	user, err := h.svc.GetCpsUserDetail(r.Context(), userCode)
	if err != nil {
		h.logger.Errorf("[FetchUserByUserCode] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserRetrieved, user)
}

// GetAllCPSUsers retrieves all CPS users with pagination
//
//	@Summary		Get all CPS users
//	@Description	Retrieves a paginated list of all CPS users with optional filtering
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int										false	"Page number"		default(1)
//	@Param			per_page	query		int										false	"Items per page"	default(10)
//	@Param			search		query		string									false	"Search term"
//
// // @Success 200 {object} localization.StandardResponse{data=cps_users_paginated_resp} "CPS users retrieved successfully"
//
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users [get]
func (h *handler) GetAllCPSUsers(w http.ResponseWriter, r *http.Request) {
	filterParasm := local_util.ExtractFilterParams(r)

	users, err := h.svc.GetAllCPSUsers(r.Context(), filterParasm)
	if err != nil {
		h.logger.Errorf("[GetAllCPSUsers] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUsersRetrieved, users)
}

// DeleteUserRequest deletes a CPS user request
//
//	@Summary		Delete CPS user request
//	@Description	Deletes a CPS user request that requires approval
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user deleted successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/delete/{user_code} [delete]
func (h *handler) DeleteUserRequest(w http.ResponseWriter, r *http.Request) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	if err := h.svc.DeleteUserRequest(r.Context(), userCode); err != nil {
		h.logger.Errorf("[DeleteUserRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserDeleted, nil)
}

// DisableUser disables a CPS user
//
//	@Summary		Disable CPS user
//	@Description	Disables a CPS user account
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user disabled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/disable/{user_code} [post]
func (h *handler) DisableUser(w http.ResponseWriter, r *http.Request) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	if err := h.svc.DisableUser(r.Context(), userCode); err != nil {
		h.logger.Errorf("[DisableUser] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserDisabled, nil)
}

// EnableUser enables a CPS user
//
//	@Summary		Enable CPS user
//	@Description	Enables a CPS user account
//	@Tags			CPS Users
//	@Accept			json
//	@Produce		json
//	@Param			user_code	path		string									true	"User Code"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"CPS user enabled successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request - User code required"
//	@Failure		404			{object}	localization.StandardResponse{data=nil}	"User not found"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/cps_users/enable/{user_code} [post]
func (h *handler) EnableUser(w http.ResponseWriter, r *http.Request) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	if err := h.svc.EnableUser(r.Context(), userCode); err != nil {
		h.logger.Errorf("[EnableUser] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SucessCpsUserEnabled, nil)
}
