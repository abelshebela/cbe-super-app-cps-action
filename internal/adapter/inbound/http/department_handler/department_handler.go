// Package department_handler provides HTTP handlers for department-related operations.
package department_handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/department"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/go-chi/chi/v5"
)

const (
	ActionApproved         = "Action approved"
	RequestSentSuccesfully = "request sent successfully"
)

type DepartmentHandler struct {
	departmentService department.DepartmentService
	logger            utils.Logger
}

func NewDepartmentHTTPHandler(service department.DepartmentService, logger utils.Logger) inbound.DepartmentPortHandler {
	return &DepartmentHandler{
		departmentService: service,
		logger:            logger,
	}
}

func (h *DepartmentHandler) createCPSActionMaker(ctx ctx_util.UserContext) entities.CPSAction {
	cpsAction := entities.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		RequestAction:    entities.RequestDepartment,
	}

	return cpsAction
}

func (h *DepartmentHandler) getContext(r *http.Request) context.Context {
	return r.Context()
}

func (h *DepartmentHandler) createCPSActionChecker(ctx ctx_util.UserContext) entities.CPSAction {
	cpsAction := entities.CPSAction{
		CheckerID:          ctx.UserID,
		CheckerName:        ctx.FullName,
		CheckerPhoneNumber: ctx.PhoneNumber,
	}

	return cpsAction
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var request department.CreateDepartmentRequest
	curCtx := h.getContext(r)

	// this is decoding
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreateDepartment] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)

		return
	}

	// Input validation
	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreateDepartment] validation failed: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidInput, http.StatusBadRequest, nil)
		return
	}

	ctx := ctx_util.ExtractUserContext(r)
	userID := ctx.UserID

	// context validation
	if ctx.IsIncomplete() {
		h.logger.Errorf("[CreateDepartment] incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusBadRequest, nil)
		return
	}

	cpsAction := h.createCPSActionMaker(ctx)

	// creating department
	createdAction, err := h.departmentService.CreateCPSAction(curCtx, request.Department, request.PortalCards, cpsAction)
	if err != nil {
		h.logger.Errorf("[CreateDepartment] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[CreateDepartment] request sent successfully by user: %s with action_code: %s", userID, createdAction.ActionCode)

	data, err := common_util.StructToMap(createdAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.BaseResponseMaker(data, w, RequestSentSuccesfully, 200)
}

func (h *DepartmentHandler) ApproveDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	ctx := ctx_util.ExtractUserContext(r)
	curCtx := h.getContext(r)

	// context validation
	if ctx.IsIncomplete() {
		h.logger.Errorf("[ApproveRequest] incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusBadRequest, nil)
		return
	}

	cpsAction, err := h.departmentService.ValidateActionRequest(curCtx, actionCode, ctx.Department)
	if err != nil {
		h.logger.Errorf("[ApproveRequest] validation failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	if cpsAction == nil {
		common_util.SendErrorResponse(w, common_util.AccountNotFound, http.StatusNotFound, nil)
		return
	}

	actionCopy := *cpsAction

	if serviceErr := h.departmentService.ApproveActionByType(curCtx, actionCopy); serviceErr != nil {
		h.logger.Errorf("[ApproveRequest] service error: %v", serviceErr)
		common_util.SendErrorResponse(w, serviceErr.Error(), http.StatusInternalServerError, nil)
		return
	}

	checker := h.createCPSActionChecker(ctx)

	if err := h.departmentService.ApproveActionRequest(curCtx, actionCode, checker); err != nil {
		h.logger.Errorf("[ApproveRequest] failed to approve action request: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.logger.Infof("[ApproveRequest] action %s approved by user %s", actionCode, ctx.UserID)
	common_util.BaseResponseMaker(nil, w, ActionApproved, 200)
}

func (h *DepartmentHandler) UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	var request department.UpdateDepartmentRequest
	curCtx := h.getContext(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UpdateDepartmentRequest] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UpdateDepartmentRequest] validation failed: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidInput, http.StatusBadRequest, nil)
		return
	}

	ctx := ctx_util.ExtractUserContext(r)
	userID := ctx.UserID

	if ctx.IsIncomplete() {
		h.logger.Errorf("[UpdateDepartmentRequest] incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusBadRequest, nil)
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionUpdate,
		RequestAction:    entities.RequestDepartment,
		CurrentAction: map[string]any{
			"department_code": request.DepartmentCode,
			"department":      request.Department,
			"portal_cards":    request.PortalCards,
		},
	}

	createdAction, err := h.departmentService.CreateDepartmentUpdateCPSAction(curCtx, cpsAction)
	if err != nil {
		h.logger.Errorf("[UpdateDepartmentRequest] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[UpdateDepartmentRequest] update request sent successfully by user: %s with action_code: %s", userID, createdAction.ActionCode)

	data, err := common_util.StructToMap(createdAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, RequestSentSuccesfully, 200)
}

func (h *DepartmentHandler) RejectDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	ctx := ctx_util.ExtractUserContext(r)
	cur_ctx := r.Context()

	// context validation
	if ctx.IsIncomplete() {
		h.logger.Errorf("[RejectDepartmentRequest] incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusBadRequest, nil)
		return
	}

	var req RejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("[RejectDepartmentRequest] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}
	if req.RejectionReason == "" {
		h.logger.Warnf("[RejectDepartmentRequest] rejection reason required")
		common_util.SendErrorResponse(w, "rejection reason required", http.StatusBadRequest, nil)
		return
	}

	// Validate action exists and is pending
	cpsAction, err := h.departmentService.ValidateActionRequest(cur_ctx, actionCode, ctx.Department)
	if err != nil {
		h.logger.Errorf("[RejectDepartmentRequest] validation failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	if cpsAction == nil {
		common_util.SendErrorResponse(w, common_util.AccountNotFound, http.StatusNotFound, nil)
		return
	}

	// Prepare checker action with rejection reason
	reason := req.RejectionReason
	checker := entities.CPSAction{
		CheckerID:          ctx.UserID,
		CheckerName:        ctx.FullName,
		CheckerPhoneNumber: ctx.PhoneNumber,
		RejectionReason:    &reason,
	}

	// Call application/service layer
	if cpsAction.ActionType == entities.ActionUpdate {
		cpsAction.CheckerID = ctx.UserID
		cpsAction.CheckerName = ctx.FullName
		cpsAction.CheckerPhoneNumber = ctx.PhoneNumber
		cpsAction.RejectionReason = &reason
		if err := h.departmentService.RejectDepartmentUpdate(cur_ctx, *cpsAction); err != nil {
			h.logger.Errorf("[RejectDepartmentRequest] failed to reject department update: %v", err)
			common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
			return
		}
	} else {
		if err := h.departmentService.RejectActionRequest(cur_ctx, actionCode, checker); err != nil {
			h.logger.Errorf("[RejectDepartmentRequest] failed to reject action request: %v", err)
			common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
			return
		}
	}

	h.logger.Infof("[RejectDepartmentRequest] action %s rejected by user %s", actionCode, ctx.UserID)
	common_util.BaseResponseMaker(nil, w, "Action rejected", 200)
}

func (h *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	departments, err := h.departmentService.GetAllDepartments(r.Context())
	if err != nil {
		h.logger.Errorf("[GetAllDepartments] failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	data, err := common_util.StructToMap(departments)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Departments fetched successfully", http.StatusOK)
}
