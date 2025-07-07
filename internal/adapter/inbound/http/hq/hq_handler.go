package hq

import (
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/hq"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

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

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	actionCode, err := h.handler.UpdateBlockTimeRequest(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"action_code": actionCode})
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

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	actionCode, err := h.handler.UpdateArchiveTimeRequest(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"action_code": actionCode})
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

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	if err := h.handler.UpdateBlockTime(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName); err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
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

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.SendErrorResponse(w, "UNAUTHORIZED", 0, nil)
		return
	}

	if err := h.handler.UpdateArchiveTime(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName); err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
}
