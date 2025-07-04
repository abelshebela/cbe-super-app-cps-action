// Package department_handler provides HTTP handlers for department-related operations.
package department_handler

import (
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
	cur_ctx := r.Context()

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
	if err := h.departmentService.CreateDepartment(cur_ctx, request.Department, request.PortalCards, cpsAction); err != nil {
		h.logger.Errorf("[CreateDepartment] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	// writing response
	h.logger.Infof("[CreateDepartment] request sent successfully by user: %s", userID)
	common_util.WriteSuccessResponse(w, nil, RequestSentSuccesfully)
}

func (h *DepartmentHandler) ApproveDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	actionCode := chi.URLParam(r, "action_code")
	ctx := ctx_util.ExtractUserContext(r)
	cur_ctx := r.Context()

	// context validation
	if ctx.IsIncomplete() {
		h.logger.Errorf("[ApproveRequest] incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[ApproveRequest] user %s is approving action %s", ctx.Department, actionCode)

	cpsAction, err := h.departmentService.ValidateActionRequest(cur_ctx, actionCode, ctx.Department)
	if err != nil {
		h.logger.Errorf("[ApproveRequest] validation failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	h.logger.Infof("[ApproveRequest] action validated: %+v", cpsAction)
	cpsAction.ActionType = entities.ActionUpdate

	if serviceErr := h.departmentService.ApproveActionByType(cur_ctx, *cpsAction); serviceErr != nil {
		h.logger.Errorf("[ApproveRequest] service error: %v", serviceErr)
		common_util.SendErrorResponse(w, serviceErr.Error(), http.StatusInternalServerError, nil)
		return
	}

	checker := h.createCPSActionChecker(ctx)

	if err := h.departmentService.ApproveActionRequest(cur_ctx, actionCode, checker); err != nil {
		h.logger.Errorf("[ApproveRequest] failed to approve action request: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.logger.Infof("[ApproveRequest] action %s approved by user %s", actionCode, ctx.UserID)
	common_util.WriteSuccessResponse(w, nil, ActionApproved)
}