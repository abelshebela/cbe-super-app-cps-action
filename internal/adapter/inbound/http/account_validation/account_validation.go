package accountvalidation_inbound

import (
	"encoding/json"
	"net/http"

	accountvalidation_app "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) FetchAccountValidation(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.SendErrorResponse(w, "INVALID_ID", 0, nil)
		return
	}

	resp, err := h.Application.GetAccountValidation(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Account fetched successfully",
		"data":    resp,
	})
}

func (h *HttpStore) UpdateAccountValidationMaker(w http.ResponseWriter, r *http.Request) {
	var req accountvalidation_app.UpdateAccountValidationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	resp, err := h.Application.UpdateAccountValidationRequest(r.Context(), req.ID, req.Validation, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	h.sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Update request submitted for approval",
		"data":    resp,
	})
}

func (h *HttpStore) UpdateAccountValidationChecker(w http.ResponseWriter, r *http.Request) {
	var req accountvalidation_app.ApproveRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	err := h.Application.UpdateAccountValidation(r.Context(), req.ActionID, req.Approve, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		errMsg := err.Error()
		utils.SendErrorResponse(w, errMsg, 0, nil)
		return

	}

	action := "approved"
	if !req.Approve {
		action = "rejected"
	}
	h.sendSuccessResponse(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "update request " + action + " successfully Approved",
		"data":    action,
	})
}

func (h *HttpStore) sendSuccessResponse(w http.ResponseWriter, status int, data interface{}) {
	resp := struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: status,
		Data:   data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
