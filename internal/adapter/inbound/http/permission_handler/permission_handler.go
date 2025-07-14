package permission_handler

import (
	"encoding/json"
	"fmt"
	// "fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/permission"

	// ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	local_commen "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/common"
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

	maker_ID, full_name, phone_number, department, ok := extractUserContextClaims(r, w)
	if !ok {
		return
	}

	cpsAction := model.CPSAction{
		MakerID:          maker_ID,
		MakerName:        full_name,
		MakerPhoneNumber: phone_number,
		Department:       department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(entities.ActionCreate),
		RequestAction:    string(entities.RequestPermissionGroup),
	}

	cpsAction, err := h.permissionService.CreatePermissionGroup(request.GroupName, request.Role, request.PermissionCategoryLists, cpsAction)
	if err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", maker_ID)
	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS ")
	data, _ := util.StructToMap(cpsAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusCreated)
}

func (h *PermissionHandler) ApprovePermissionGroup(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	Checker_ID, full_name, phone_number, department, ok := extractUserContextClaims(r, w)

	fmt.Println("handler========================================")
	fmt.Println(Checker_ID, full_name, phone_number, department, ok)
	fmt.Println("handler========================================")

	if !ok {
		return
	}

	h.logger.Infof("[ApprovePermissionGroup] approving action %s", actionCode)

	checker := model.CPSAction{
		CheckerID:          Checker_ID,
		CheckerName:        full_name,
		CheckerPhoneNumber: phone_number,
		Department:         department,
	}

	approvedAction, err := h.permissionService.ApprovePermissionGroup(actionCode, checker)
	if err != nil {
		h.logger.Errorf("[ApprovePermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(approvedAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusCreated)
}

func (h *PermissionHandler) RejectPermissionGroup(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	checker_ID, full_name, phone_number, department, ok := extractUserContextClaims(r, w)
	if !ok {
		return
	}

	var payload struct {
		Reason string `json:"rejection_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.logger.Errorf("[RejectPermissionGroup] failed to decode request: %v", err)
		util.SendErrorResponse(w, "Invalid request payload", http.StatusBadRequest, nil)
		return
	}
	if payload.Reason == "" {
		util.SendErrorResponse(w, "Rejection reason is required", http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[RejectPermissionGroup] rejecting action %s by user %s", actionCode, checker_ID)

	checker := model.CPSAction{
		CheckerID:          checker_ID,
		CheckerName:        full_name,
		CheckerPhoneNumber: phone_number,
		Department:         department,
	}

	rejectedAction, err := h.permissionService.RejectPermissionGroup(actionCode, checker, payload.Reason)
	if err != nil {
		h.logger.Errorf("[RejectPermissionGroup] service error: %v", err)
		util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	def, _ := local_commen.GetSuccessResponseByCode("SUCCESS")
	data, _ := util.StructToMap(rejectedAction)
	util.BaseResponseMaker(data, w, def.Message, http.StatusCreated)
}

// Extracts and validates user context claims, returns user info and ok flag
func extractUserContextClaims(r *http.Request, w http.ResponseWriter) (userID, fullName, phoneNumber, department string, ok bool) {
	userID, ok1 := r.Context().Value(constant.ContextKey("user_code")).(string)
	fullName, ok2 := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, ok3 := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, ok4 := r.Context().Value(constant.ContextKey("department")).(string)

	if !ok1 || userID == "" {
		util.SendErrorResponse(w, "missing or invalid user_code in context", http.StatusBadRequest, nil)
		return
	}
	if !ok2 || fullName == "" {
		util.SendErrorResponse(w, "missing or invalid full_name in context", http.StatusBadRequest, nil)
		return
	}
	if !ok3 || phoneNumber == "" {
		util.SendErrorResponse(w, "missing or invalid phone_number in context", http.StatusBadRequest, nil)
		return
	}
	if !ok4 || department == "" {
		util.SendErrorResponse(w, "missing or invalid department in context", http.StatusBadRequest, nil)
		return
	}
	return userID, fullName, phoneNumber, department, true
}
