// Package department_handler provides HTTP handlers for department-related operations.
package department_handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/department"
	cpsconstants "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func (h *DepartmentHandler) createCPSActionMaker(ctx ctx_util.UserContext) cpsactions.CPSAction {
	cpsAction := cpsactions.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     cpsconstants.ActionPending,
		ActionType:       cpsconstants.ActionCreate,
		RequestAction:    cpsconstants.RequestCreateDepartment,
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
	createdActionCode, err := h.departmentService.CreateCPSAction(curCtx, request.Department, request.PortalCards, request.PermissionGroups, cpsAction)
	if err != nil {
		h.logger.Errorf("[CreateDepartment] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[CreateDepartment] request sent successfully by user: %s with action_code: %s", userID, createdActionCode)

	common_util.BaseResponseMaker(map[string]string{"action_code": createdActionCode}, w, RequestSentSuccesfully, 200)
}

func (h *DepartmentHandler) UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	var request department.UpdateDepartmentRequest
	curCtx := h.getContext(r)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UpdateDepartmentRequest] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := validatePatchUpdateDepartmentRequest(request); err != nil {
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

	// Get department_id from URL params (PATCH route uses 'id')
	departmentID := chi.URLParam(r, "id")

	// Create CPS action for update
	cpsAction := cpsactions.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     cpsconstants.ActionPending,
		ActionType:       cpsconstants.ActionUpdate,
		RequestAction:    cpsconstants.RequestUpdateDepartment,
		CurrentAction: map[string]interface{}{
			"department_id":     departmentID,
			"department":        request.Department,
			"portal_cards":      request.PortalCards,
			"permission_groups": request.PermissionGroups,
		},
	}

	createdAction, err := h.departmentService.CreateDepartmentUpdateCPSAction(curCtx, cpsAction)
	if err != nil {
		h.logger.Errorf("[UpdateDepartmentRequest] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[UpdateDepartmentRequest] update request sent successfully by user: %s for department_id: %s, action_code: %s", userID, departmentID, createdAction.ActionCode)
	common_util.BaseResponseMaker(map[string]string{"action_code": createdAction.ActionCode}, w, RequestSentSuccesfully, 200)
}

// validatePatchUpdateDepartmentRequest validates only fields that are present for PATCH semantics
func validatePatchUpdateDepartmentRequest(req department.UpdateDepartmentRequest) error {
	var rules []error
	if req.Department != "" {
		if err := validation.Validate(req.Department, validation.Required); err != nil {
			rules = append(rules, err)
		}
	}
	if req.PortalCards != nil {
		if err := validation.Validate(req.PortalCards, validation.Each(validation.Required)); err != nil {
			rules = append(rules, err)
		}
	}
	if req.PermissionGroups != nil {
		if err := validation.Validate(req.PermissionGroups, validation.Each(validation.Required)); err != nil {
			rules = append(rules, err)
		}
	}
	if len(rules) > 0 {
		return fmt.Errorf("%v", rules)
	}
	return nil
}

func (h *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	departments, err := h.departmentService.GetAllDepartments(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("[GetAllDepartments] failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	common_util.BaseResponseMaker(departments, w, "Departments fetched successfully", http.StatusOK)
}
