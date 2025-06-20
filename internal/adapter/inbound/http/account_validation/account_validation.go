package accountvalidation_inbound

import (
	"encoding/json"
	"net/http"

	accountvalidation_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/utils"
)

var errorMap = map[string]int{
	"validation rule ID cannot be empty":          http.StatusBadRequest,
	"maker ID cannot be empty":                    http.StatusBadRequest,
	"validation failed: identifier cannot be empty": http.StatusBadRequest,
	"validation failed: min length cannot exceed max length": http.StatusBadRequest,
	"validation failed: service ID cannot be empty": http.StatusBadRequest,
	"validation rule not found":                    http.StatusNotFound,
	"validation rule already has a pending action": http.StatusConflict,
	"action ID cannot be empty":                    http.StatusBadRequest,
	"checker ID cannot be empty":                   http.StatusBadRequest,
	"action not found":                             http.StatusNotFound,
	"action is not pending":                        http.StatusConflict,
}

func (h *HttpStore) FetchAccountValidation(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation ID required")
		return
	}

	validation, err := h.Application.GetAccountValidation(r.Context(), id)
	if err != nil {
		errMsg := err.Error()
		status, ok := errorMap[errMsg]
		if !ok {
			status = http.StatusInternalServerError
			errMsg = "Failed to fetch validation rule"
		}
		utils.WriteErrorResponse(w, status, errMsg)
		return
	}

	utils.WriteSuccessResponse(w, validation, "Successfully retrieved validation rule")
}

func (h *HttpStore) UpdateAccountValidationMaker(w http.ResponseWriter, r *http.Request) {
	var req accountvalidation_app.UpdateAccountValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	actionID, err := h.Application.UpdateAccountValidationRequest(r.Context(), req.ID, req.Validation, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		errMsg := err.Error()
		status, ok := errorMap[errMsg]
		if !ok {
			status = http.StatusInternalServerError
			errMsg = "Update request failed"
		}
		utils.WriteErrorResponse(w, status, errMsg)
		return
	}

	resp := accountvalidation_app.UpdateAccountValidationResponse{ActionID: actionID}
	utils.WriteSuccessResponse(w, resp, "Update request submitted for approval")
}

func (h *HttpStore) UpdateAccountValidationChecker(w http.ResponseWriter, r *http.Request) {
	var req accountvalidation_app.ApproveRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	claims, _ := r.Context().Value("claims").(middleware.UserPayload)

	err := h.Application.UpdateAccountValidation(r.Context(), req.ActionID, req.Approve, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		errMsg := err.Error()
		status, ok := errorMap[errMsg]
		if !ok {
			status = http.StatusInternalServerError
			errMsg = "Approval process failed"
		}
		utils.WriteErrorResponse(w, status, errMsg)
		return
	}

	action := "approved"
	if !req.Approve {
		action = "rejected"
	}
	utils.WriteSuccessResponse(w, nil, "Update request "+action+" successfully")
}