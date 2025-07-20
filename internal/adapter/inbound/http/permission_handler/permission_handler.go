package permission_handler

import (
	"encoding/json"

	// "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	constants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

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

	cpsAction := model.CPSAction{
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(entities.ActionCreate),
		RequestAction:    string(constants.RequestCreatePermissionGroup),
	}

	cpsAction, err := h.permissionService.CreatePermissionGroup("", request.GroupName, request.Role, request.PermissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	common_util.WriteSuccessResponse(w, cpsAction, "group permission Request  created successfully")
}

func (h *PermissionHandler) GetPermissionGroups(w http.ResponseWriter, r *http.Request) {
	permissionGroups, err := h.permissionService.GetPermissionGroups()
	if err != nil {
		h.logger.Errorf("[GetPermissionGroups] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	common_util.WriteSuccessResponse(w, permissionGroups, "Permission groups fetched successfully")
}

func (h *PermissionHandler) GetPermissionGroup(w http.ResponseWriter, r *http.Request) {
	permissionGroup, err := h.permissionService.GetPermissionGroup(chi.URLParam(r, "group_name"))
	if err != nil {
		h.logger.Errorf("[GetPermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, permissionGroup, "Permission group fetched successfully")
}

func (h *PermissionHandler) UpdatePermissionGroup(w http.ResponseWriter, r *http.Request) {
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

	cpsAction := model.CPSAction{
		MakerID:          userContext.UserID,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(entities.ActionUpdate),
		RequestAction:    string(constants.RequestCreatePermissionGroup),
	}

	cpsAction, err := h.permissionService.CreatePermissionGroup(oldGroupName, request.GroupName, request.Role, request.PermissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[UpdatePermissionGroup] request sent successfully by user: %s", userContext.UserID)

	common_util.WriteSuccessResponse(w, cpsAction, "group permission Request updated successfully")
}

// func (h *PermissionHandler) ApprovePermissionGroup(w http.ResponseWriter, r *http.Request) {
// 	actionCode := chi.URLParam(r, "action_code")
// 	checkerUser := contexts.ExtractUserContext(r)

// 	h.logger.Infof("[ApprovePermissionGroup] approving action %s", actionCode)

// 	checker := model.CPSAction{
// 		CheckerID:          checkerUser.UserID,
// 		CheckerName:        checkerUser.FullName,
// 		CheckerPhoneNumber: checkerUser.PhoneNumber,
// 		Department:         checkerUser.Department,
// 	}

// 	approvedAction, err := h.permissionService.ApprovePermissionGroup(actionCode, checker)
// 	if err != nil {
// 		h.logger.Errorf("[ApprovePermissionGroup] service error: %v", err)
// 		util.SendErrorResponse(w, err.Error(), 0, nil)
// 		return
// 	}
// 	common_util.WriteSuccessResponse(w, approvedAction, "Action Approved")
// }

// func (h *PermissionHandler) RejectPermissionGroup(w http.ResponseWriter, r *http.Request) {
// 	actionCode := chi.URLParam(r, "action_code")
// 	checker_ID, full_name, phone_number, department, ok := extractUserContextClaims(r, w)
// 	if !ok {
// 		return
// 	}

// 	var payload struct {
// 		Reason string `json:"rejection_reason"`
// 	}

// 	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
// 		h.logger.Errorf("[RejectPermissionGroup] failed to decode request: %v", err)
// 		util.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest, nil)
// 		return
// 	}
// 	if payload.Reason == "" {
// 		util.SendErrorResponse(w, "Rejection reason is required", http.StatusBadRequest, nil)
// 		return
// 	}

// 	h.logger.Infof("[RejectPermissionGroup] rejecting action %s by user %s", actionCode, checker_ID)

// 	checker := model.CPSAction{
// 		CheckerID:          checker_ID,
// 		CheckerName:        full_name,
// 		CheckerPhoneNumber: phone_number,
// 		Department:         department,
// 	}

// 	rejectedAction, err := h.permissionService.RejectPermissionGroup(actionCode, checker, payload.Reason)
// 	if err != nil {
// 		h.logger.Errorf("[RejectPermissionGroup] service error: %v", err)
// 		util.SendErrorResponse(w, err.Error(), 0, nil)
// 		return
// 	}
// 	common_util.WriteSuccessResponse(w, rejectedAction, "Action Rejected")
// }

// // Extracts and validates user context claims, returns user info and ok flag
// func extractUserContextClaims(r *http.Request, w http.ResponseWriter) (userID, fullName, phoneNumber, department string, ok bool) {
// 	userID, ok1 := r.Context().Value(constant.ContextKey("user_code")).(string)
// 	fullName, ok2 := r.Context().Value(constant.ContextKey("full_name")).(string)
// 	phoneNumber, ok3 := r.Context().Value(constant.ContextKey("phone_number")).(string)
// 	department, ok4 := r.Context().Value(constant.ContextKey("department")).(string)

// 	if !ok1 || userID == "" {
// 		fmt.Println("[extractUserContextClaims] missing or invalid user_code in context")
// 		util.SendErrorResponse(w, "missing or invalid user_code in context", http.StatusBadRequest, nil)
// 		return
// 	}
// 	if !ok2 || fullName == "" {
// 		fmt.Println("[extractUserContextClaims] missing or invalid full_name in context")
// 		util.SendErrorResponse(w, "missing or invalid full_name in context", http.StatusBadRequest, nil)
// 		return
// 	}
// 	if !ok3 || phoneNumber == "" {
// 		fmt.Println("[extractUserContextClaims] missing or invalid phone_number in context")
// 		util.SendErrorResponse(w, "missing or invalid phone_number in context", http.StatusBadRequest, nil)
// 		return
// 	}
// 	if !ok4 || department == "" {
// 		fmt.Println("[extractUserContextClaims] missing or invalid department in context")
// 		util.SendErrorResponse(w, "missing or invalid department in context", http.StatusBadRequest, nil)
// 		return
// 	}
// 	return userID, fullName, phoneNumber, department, true
// }
