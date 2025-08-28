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

	if err := h.svc.CreateUserRequest(r.Context(), req); err != nil {
		h.logger.Errorf("[CreateUserRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserCreationRequestSubmitted, nil)
}

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

func (h *handler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	userCode := strings.TrimSpace(chi.URLParam(r, "user_code"))
	if userCode == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorUserCodeRequired.Code)
		return
	}

	user, err := h.svc.FetchUserByUserCode(r.Context(), userCode)
	if err != nil {
		h.logger.Errorf("[FetchUserByUserCode] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUserRetrieved, user)
}

func (h *handler) GetAllCPSUsers(w http.ResponseWriter, r *http.Request) {
	filter := local_util.ExtractFilterParams(r)
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PerPage < 1 {
		filter.PerPage = 10
	}
	if filter.PerPage > 100 {
		filter.PerPage = 100
	}

	users, err := h.svc.GetAllCPSUsers(r.Context(), filter)
	if err != nil {
		h.logger.Errorf("[GetAllCPSUsers] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessCpsUsersRetrieved, users)
}

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
