package permission_handler

import (
	"encoding/json"

	// "fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	// common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/permission"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/entities"

	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// local_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
)

type PermissionHandler struct {
	permissionService permission.PermissionService
	logger            utils.Logger
}

func NewPermissionHTTPHandler(service permission.PermissionService, logger utils.Logger) inbound.PermissionPortHandler {
	return &PermissionHandler{
		permissionService: service,
		logger:            logger,
	}
}

func (h *PermissionHandler) CreatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var request permission.CreatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreatePermissionGroup] failed to decode request: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}
	role := strings.ToUpper(request.Role)
	cpsAction := model.CPSAction{
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(entities.ActionCreate),
		RequestAction:    string(constants.RequestCreatePermissionGroup),
	}

	_, err := h.permissionService.CreatePermissionGroup(r.Context(), "", request.GroupName, role, request.PermissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	util.WriteSuccessResponse(w, 0, "group permission Request  created successfully")
}

func (h *PermissionHandler) GetPermissionGroups(w http.ResponseWriter, r *http.Request) {

	filterparams := util.ExtractFilterParams(r)
	permissionGroups, err := h.permissionService.GetPermissionGroups(r.Context(), filterparams)
	if err != nil {
		h.logger.Errorf("[GetPermissionGroups] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	util.WriteSuccessResponse(w, permissionGroups, "Permission groups fetched successfully")
}

func (h *PermissionHandler) GetPermissionGroup(w http.ResponseWriter, r *http.Request) {

	permissionGroup, err := h.permissionService.GetPermissionGroup(chi.URLParam(r, "group_name"))
	if err != nil {
		h.logger.Errorf("[GetPermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	util.WriteSuccessResponse(w, permissionGroup, "Permission group fetched successfully")
}

func (h *PermissionHandler) UpdatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var request permission.UpdatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] failed to decode request: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UpdatePermissionGroup] validation failed: %v", err)
		util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	oldGroupName := chi.URLParam(r, "group_name")
	if oldGroupName == "" {
		util.SendErrorResponse(w, "Group name is required", http.StatusNotFound, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		util.SendErrorResponse(w, util.IncompleteUserInfo, 0, nil)
		return
	}

	// Extract values from optional fields, using old values as defaults if not provided
	groupName := oldGroupName
	if request.GroupName != nil {
		groupName = *request.GroupName
	}

	// Get current permission group to use as defaults for missing fields
	currentGroup, err := h.permissionService.GetPermissionGroup(oldGroupName)
	if err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] failed to get current permission group: %v", err)
		util.SendErrorResponse(w, "Failed to get current permission group", 0, nil)
		return
	}

	// Use current values as defaults for missing fields
	var role string = currentGroup.Role
	if request.Role != nil {
		role = strings.ToUpper(*request.Role)
	}

	var permissionCategoryLists []string
	if request.PermissionCategoryLists != nil {
		permissionCategoryLists = *request.PermissionCategoryLists
		h.logger.Infof("[UpdatePermissionGroup] Using provided permission categories: %+v", permissionCategoryLists)
	} else {
		h.logger.Infof("[UpdatePermissionGroup] No permission categories provided, will keep existing ones unchanged")
		permissionCategoryLists = []string{}
	}

	cpsAction := model.CPSAction{
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(entities.ActionUpdate),
		RequestAction:    string(constants.RequestUpdatePermissionGroup),
	}

	_, err = h.permissionService.UpdatePermissionGroupRequest(r.Context(), oldGroupName, groupName, role, permissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[UpdatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	util.WriteSuccessResponse(w, 0, "group permission Request updated successfully")
}

func (h *PermissionHandler) GetAllPermissionCategoriesWithPermissions(w http.ResponseWriter, r *http.Request) {
	categories, err := h.permissionService.GetAllPermissionCategoriesWithPermissions(r.Context())
	if err != nil {
		h.logger.Errorf("[GetAllPermissionCategoriesWithPermissions] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	util.WriteSuccessResponse(w, categories, "Permission categories fetched successfully")
}
