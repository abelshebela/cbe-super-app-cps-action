package bpsmakerhandler

import (
	"net/http"

	bpsapp "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/bps_user"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/bps_user"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	local_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BPSUserHandler struct {
	Service bpsapp.ApplicationService
	logger  utils.Logger
}

func InitBPSUserMakerHandler(service bpsapp.ApplicationService, logger utils.Logger) inbound.BPSUserHandler {
	return BPSUserHandler{
		Service: service,
		logger:  logger,
	}
}

func (h BPSUserHandler) GetPendingUserActions(w http.ResponseWriter, r *http.Request) {
	pendingUserAction, err := h.Service.GetPendingUserActions(r.Context())
	if err != nil {
		h.logger.Errorf("GetPendingUserAction failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	local_util.BaseResponseMaker(pendingUserAction, w, "Pending users fetched successfully", 200)
}

func (h BPSUserHandler) FetchUserByUserCode(w http.ResponseWriter, r *http.Request) {
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		local_util.SendErrorResponse(w, "USER_CODE_IS_REQUIRED", 0, nil)
		return
	}

	user, err := h.Service.FetchUserByUserCode(r.Context(), userCode)
	if err != nil {
		h.logger.Errorf("FetchUserRequest failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, err := local_util.StructToMap(user)
	if err != nil {
		h.logger.Errorf("failed to convert data to map: %v", err)
		local_util.SendErrorResponse(w, "Failed to convert data to map", http.StatusInternalServerError, nil)
		return
	}
	local_util.BaseResponseMaker(data, w, "BPS user fetched successfully", 200)
}

func (h BPSUserHandler) GetAllBPSUsers(w http.ResponseWriter, r *http.Request) {
	// filterParams := local_util.ExtractFilterParams(r)
	filterParams := local_util.ExtractMongoFilterParams(r)

	users, err := h.Service.GetAllBPSUsers(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("GetAllBPSUser request failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	doc, _ := local_util.StructToMap(users)
	local_util.BaseResponseMaker(doc, w, "BPS users fetched successfully", 200)
}

func (h BPSUserHandler) DisableUser(w http.ResponseWriter, r *http.Request) {
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		local_util.SendErrorResponse(w, "USER_CODE_IS_REQUIRED", 0, nil)
		return
	}

	// Extract user context from request
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	action, err := h.Service.DisableUser(r.Context(), userCode, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		h.logger.Errorf("Disable user request failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]string{"action_code": action.ActionCode}
	local_util.BaseResponseMaker(response, w, "User disable request submitted successfully", 200)
}

func (h BPSUserHandler) EnableUser(w http.ResponseWriter, r *http.Request) {
	userCode := chi.URLParam(r, "user_code")
	if userCode == "" {
		local_util.SendErrorResponse(w, "USER_CODE_IS_REQUIRED", 0, nil)
		return
	}

	// Extract user context from request
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	action, err := h.Service.EnableUser(r.Context(), userCode, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		h.logger.Errorf("Enable user request failed: %v", err)
		local_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]string{"action_code": action.ActionCode}
	local_util.BaseResponseMaker(response, w, "User enable request submitted successfully", 200)
}
