package hq

import (
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/hq"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"github.com/go-chi/chi/v5"
)

type HQHTTPHandler struct {
	handler hq.ApplicationAbstracts
}

func NewHQHTTPHandler(handler hq.ApplicationAbstracts) *HQHTTPHandler {
	return &HQHTTPHandler{
		handler: handler,
	}
}

func (h *HQHTTPHandler) GetHQ(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		utils.SendErrorResponse(w, "INVALID_ID", 0, nil)
		return
	}

	hqResp, err := h.handler.GetHQ(r.Context(), id)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hqResp)
}

func (h *HQHTTPHandler) UpdateBlockTimeRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateBlockTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	if err := request.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if userID == "" || fullName == "" || phoneNumber == "" || department == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	actionCode, err := h.handler.UpdateBlockTimeRequest(r.Context(), request, userID, phoneNumber, fullName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Update block time request submitted for approval",
		"data":    map[string]interface{}{"action_id": actionCode},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *HQHTTPHandler) UpdateArchiveTimeRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateArchiveTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	if err := request.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if userID == "" || fullName == "" || phoneNumber == "" || department == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	actionCode, err := h.handler.UpdateArchiveTimeRequest(r.Context(), request, userID, phoneNumber, fullName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Update archive time request submitted for approval",
		"data":    map[string]interface{}{"action_id": actionCode},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *HQHTTPHandler) UpdateBlockTime(w http.ResponseWriter, r *http.Request) {
	var request dto.ApproveRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}

	if err := request.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if userID == "" || fullName == "" || phoneNumber == "" || department == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	if err := h.handler.UpdateBlockTime(r.Context(), request, userID, phoneNumber, fullName); err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	action := "approved"
	if request.Decison == utils.DecisionDenied {
		action = "rejected"
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "update request " + action + " successfully Approved",
		"data":    map[string]interface{}{"action": action},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *HQHTTPHandler) UpdateArchiveTime(w http.ResponseWriter, r *http.Request) {
	var request dto.ApproveRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}

	if err := request.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	userID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	fullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department, _ := r.Context().Value(constant.ContextKey("department")).(string)

	if userID == "" || fullName == "" || phoneNumber == "" || department == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	if err := h.handler.UpdateArchiveTime(r.Context(), request, userID, phoneNumber, fullName); err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	action := "approved"
	if request.Decison == utils.DecisionDenied {
		action = "rejected"
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "update request " + action + " successfully Approved",
		"data":    map[string]interface{}{"action": action},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
