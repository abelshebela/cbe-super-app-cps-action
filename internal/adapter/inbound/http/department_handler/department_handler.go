// Package department_handler provides HTTP handlers for department-related operations.
package department_handler

import (
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/department"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/department"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/go-chi/chi/v5"
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

// extractUserContext extracts and validates user context from the request
func (h *DepartmentHandler) extractUserContext(w http.ResponseWriter, r *http.Request) (ctx_util.UserContext, bool) {
	ctx := ctx_util.ExtractUserContext(r)
	if ctx.IsIncomplete() {
		h.logger.Errorf("incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, http.StatusBadRequest, nil)
		return ctx, false
	}
	return ctx, true
}

// validateDepartmentID validates the department ID from URL parameters
func (h *DepartmentHandler) validateDepartmentID(w http.ResponseWriter, r *http.Request) (string, bool) {
	departmentID := chi.URLParam(r, "id")
	if departmentID == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, http.StatusBadRequest, nil)
		return "", false
	}
	return departmentID, true
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var request department.CreateDepartmentRequest
	curCtx := r.Context()

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreateDepartment] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	// Input validation
	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreateDepartment] validation failed: %v", err)
		common_util.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	// Extract and validate user context
	userCtx, ok := h.extractUserContext(w, r)
	if !ok {
		return
	}

	// Create maker user for application layer
	maker := department.Maker{
		UserCode:    userCtx.UserID,
		FullName:    userCtx.FullName,
		PhoneNumber: userCtx.PhoneNumber,
		Department:  userCtx.Department,
	}

	// Pass to application layer - let it handle CPS action creation
	if err := h.departmentService.CreateDepartment(curCtx, request, maker); err != nil {
		h.logger.Errorf("[CreateDepartment] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[CreateDepartment] request sent successfully by user: %s", userCtx.UserID)
	common_util.WriteSuccessResponse(w, nil, RequestSentSuccesfully)
}

func (h *DepartmentHandler) UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	var request department.DepartmentUpdateCPSActionRequest
	curCtx := r.Context()

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UpdateDepartmentRequest] failed to decode request: %v", err)
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	// Input validation
	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UpdateDepartmentRequest] validation failed: %v", err)
		common_util.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	// Extract and validate user context
	userCtx, ok := h.extractUserContext(w, r)
	if !ok {
		return
	}

	// Validate department ID
	departmentID, ok := h.validateDepartmentID(w, r)
	if !ok {
		return
	}
	maker := department.Maker{
		UserCode:    userCtx.UserID,
		FullName:    userCtx.FullName,
		PhoneNumber: userCtx.PhoneNumber,
		Department:  userCtx.Department,
	}

	// Pass to application layer - let it handle CPS action creation
	if err := h.departmentService.UpdateDepartment(curCtx, departmentID, request, maker); err != nil {
		h.logger.Errorf("[UpdateDepartmentRequest] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[UpdateDepartmentRequest] update request sent successfully by user: %s for department_id: %s", userCtx.UserID, departmentID)
	common_util.WriteSuccessResponse(w, nil, RequestSentSuccesfully)
}

func (h *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	departments, err := h.departmentService.GetAllDepartments(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("[GetAllDepartments] failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	common_util.WriteSuccessResponse(w, departments, "Departments fetched successfully")
}

func (h *DepartmentHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	department, err := h.departmentService.GetDepartmentByID(r.Context(), id)
	if err != nil {
		h.logger.Errorf("[GetDepartmentByID] failed: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	common_util.WriteSuccessResponse(w, department, "Department fetched successfully")
}

func (h *DepartmentHandler) EnableDepartment(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := h.extractUserContext(w, r)
	if !ok {
		return
	}

	departmentID, ok := h.validateDepartmentID(w, r)
	if !ok {
		return
	}

	// _, err := h.departmentService.GetDepartmentByID(r.Context(), departmentID)
	// if err != nil {
	// 	h.logger.Errorf("[EnableDepartment] department not found: %v", err)
	// 	common_util.SendErrorResponse(w, "NOT_FOUND", http.StatusNotFound, nil)
	// 	return
	// }

	// Create maker user for application layer
	maker := department.Maker{
		UserCode:    userCtx.UserID,
		FullName:    userCtx.FullName,
		PhoneNumber: userCtx.PhoneNumber,
		Department:  userCtx.Department,
	}

	// Pass to application layer - let it handle CPS action creation
	if err := h.departmentService.EnableDepartment(r.Context(), departmentID, maker); err != nil {
		h.logger.Errorf("[EnableDepartment] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[EnableDepartment] enable request sent successfully by user: %s for department_id: %s", userCtx.UserID, departmentID)
	common_util.WriteSuccessResponse(w, nil, RequestSentSuccesfully)
}

func (h *DepartmentHandler) DisableDepartment(w http.ResponseWriter, r *http.Request) {
	// Extract and validate user context
	userCtx, ok := h.extractUserContext(w, r)
	if !ok {
		return
	}

	// Validate department ID
	departmentID, ok := h.validateDepartmentID(w, r)
	if !ok {
		return
	}

	// Check if department exists
	_, err := h.departmentService.GetDepartmentByID(r.Context(), departmentID)
	if err != nil {
		h.logger.Errorf("[DisableDepartment] department not found: %v", err)
		common_util.SendErrorResponse(w, "NOT_FOUND", http.StatusNotFound, nil)
		return
	}

	// Create maker user for application layer
	maker := department.Maker{
		UserCode:    userCtx.UserID,
		FullName:    userCtx.FullName,
		PhoneNumber: userCtx.PhoneNumber,
		Department:  userCtx.Department,
	}

	// Pass to application layer - let it handle CPS action creation
	if err := h.departmentService.DisableDepartment(r.Context(), departmentID, maker); err != nil {
		h.logger.Errorf("[DisableDepartment] service error: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	h.logger.Infof("[DisableDepartment] disable request sent successfully by user: %s for department_id: %s", userCtx.UserID, departmentID)
	common_util.WriteSuccessResponse(w, nil, RequestSentSuccesfully)
}
