package permission

import (
	"cbe-super-app-cps-action/internal/constants/dto/permission"
	permission_int "cbe-super-app-cps-action/internal/constants/interfaces/permission"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type PermissionHandler struct {
	PermissionService service.PermissionService
	logger            utils.Logger
}

func InitPermissionHandler(svc service.PermissionService, logger utils.Logger) permission_int.PermissionHandler {
	return &PermissionHandler{
		PermissionService: svc,
		logger:            logger,
	}
}

func (h *PermissionHandler) CreatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var request permission.CreatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreatePermissionGroup] failed to decode request: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToDecodeRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := common_utils.ExtractUserContext(r)
	if common_utils.IsIncomplete(userContext) {
		localization.SendErrorByCodeResponse(w, localization.ErrorIncompleteUserInfo.Code)
		return
	}

	err := h.PermissionService.CreatePermissionGroup(r.Context(), request)
	if err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupRequestCreated, nil)
}

func (h *PermissionHandler) GetPermissionGroups(w http.ResponseWriter, r *http.Request) {
	filterparams := common_utils.ExtractFilterParams(r)
	permissionGroups, err := h.PermissionService.GetPermissionGroups(r.Context(), filterparams)
	if err != nil {
		h.logger.Errorf("[GetPermissionGroups] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupsFetched, permissionGroups)
}

func (h *PermissionHandler) GetPermissionGroup(w http.ResponseWriter, r *http.Request) {
	permissionGroup, err := h.PermissionService.GetPermissionGroup(chi.URLParam(r, "group_name"))
	if err != nil {
		h.logger.Errorf("[GetPermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupFetched, permissionGroup)
}

func (h *PermissionHandler) UpdatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var request permission.UpdatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UpdatePermissionGroup] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	oldGroupName := chi.URLParam(r, "group_name")
	if oldGroupName == "" {
		localization.SendErrorResponse(w, localization.ErrorGroupNameRequired, nil, nil)
		return
	}

	userContext := common_utils.ExtractUserContext(r)
	if common_utils.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	// Set the old group name from URL parameter
	request.OldGroupName = oldGroupName

	err := h.PermissionService.UpdatePermissionGroup(r.Context(), request)
	if err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[UpdatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupRequestUpdated, nil)
}

func (h *PermissionHandler) GetAllPermissionCategoriesWithPermissions(w http.ResponseWriter, r *http.Request) {
	categories, err := h.PermissionService.GetAllPermissionCategoriesWithPermissions(r.Context())
	if err != nil {
		h.logger.Errorf("[GetAllPermissionCategoriesWithPermissions] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessPermissionCategoriesFetched, categories)
}
