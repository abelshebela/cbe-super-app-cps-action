package accountvalidation_inbound

import (
	"encoding/json"
	"net/http"

	accountvalidation_app "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) FetchAccountValidation(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Validation ID required")
		return
	}

	validation, err := h.Application.GetAccountValidation(r.Context(), id)
	if err != nil {
		switch err.Error() {
		case "validation rule ID cannot be empty":
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		case "validation rule not found":
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		default:
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch validation rule")
		}
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

	claims, _ := r.Context().Value("claims").(middleware.UserPayload)

	actionID, err := h.Application.UpdateAccountValidationRequest(r.Context(), req.ID, req.Validation, claims.UserID)
	if err != nil {
		switch err.Error() {
		case "validation rule ID cannot be empty",
			"maker ID cannot be empty",
			"validation failed: identifier cannot be empty",
			"validation failed: min length cannot exceed max length",
			"validation failed: service ID cannot be empty":
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		case "validation rule not found":
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		case "validation rule already has a pending action":
			utils.WriteErrorResponse(w, http.StatusConflict, err.Error())
		default:
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Update request failed")
		}
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

	err := h.Application.UpdateAccountValidation(r.Context(), req.ActionID, req.Approve, claims.UserID)
	if err != nil {
		switch err.Error() {
		case "action ID cannot be empty",
			"checker ID cannot be empty":
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		case "action not found":
			utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		case "action is not pending":
			utils.WriteErrorResponse(w, http.StatusConflict, err.Error())
		case "validation failed: identifier cannot be empty",
			"validation failed: min length cannot exceed max length",
			"validation failed: service ID cannot be empty":
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		default:
			utils.WriteErrorResponse(w, http.StatusInternalServerError, "Approval process failed")
		}
		return
	}

	action := "approved"
	if !req.Approve {
		action = "rejected"
	}
	utils.WriteSuccessResponse(w, nil, "Update request "+action+" successfully")
}
