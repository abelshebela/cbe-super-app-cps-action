package bpsmakerhandler

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/bps_user"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"net/http"

	common_utils "cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BPSUserHandler struct {
	Service service.BPSUserService
	logger  utils.Logger
}

func InitBPSUserMakerHandler(service service.BPSUserService, logger utils.Logger) bps_user.BPSUserHandler {
	return BPSUserHandler{
		Service: service,
		logger:  logger,
	}
}

func (h BPSUserHandler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w,localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	user, err := h.Service.FetchUserByUserCode(r.Context(), userCode)
	if err != nil {
		h.logger.Errorf("FetchUserRequest failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessUserRetrieved, user)
}

func (h BPSUserHandler) GetAllBPSUsers(w http.ResponseWriter, r *http.Request) {
	filterParams := common_utils.ExtractFilterParams(r)
	
	users, err := h.Service.GetAllBPSUsers(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("GetAllBPSUser request failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w,localization.SuccessUserRetrieved,users)
}

func (h BPSUserHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}

	err := h.Service.UpdateBpsUser(r.Context(), userCode, false)
	if err != nil {
		h.logger.Errorf("Disable user request failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w,localization.SuccessBankDisableRequestSent,map[string]string{})
}

func (h BPSUserHandler) EnableUser(w http.ResponseWriter, r *http.Request) {
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		localization.SendErrorResponse(w, localization.ErrorUserCodeRequired, nil, nil)
		return
	}
	
	err := h.Service.UpdateBpsUser(r.Context(), userCode,true)
	if err != nil {
		h.logger.Errorf("Enable user request failed: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	
	localization.SendSuccessResponse(w,localization.SuccessBankEnableRequestSent,map[string]string{})
}
