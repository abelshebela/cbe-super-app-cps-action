package hq

import (
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/application/dto"
	"cbe-super-app-cps-action/internal/application/hq"
	"cbe-super-app-cps-action/internal/application/middleware"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type HQHTTPHandler struct {
	handler hq.ApplicationAbstracts
	logger  utils.Logger
}

func NewHQHTTPHandler(handler hq.ApplicationAbstracts, logger utils.Logger) *HQHTTPHandler {
	return &HQHTTPHandler{
		handler: handler,
		logger:  logger,
	}
}

func (h *HQHTTPHandler) GetHQ(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		h.logger.Errorf("HQ ID is empty")
		http.Error(w, "HQ ID is required", http.StatusBadRequest)
		return
	}

	hq, err := h.handler.GetHQ(r.Context(), id)
	if err != nil {
		h.logger.Errorf("failed to get HQ: %v", err)
		http.Error(w, "Failed to get HQ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hq)
}

func (h *HQHTTPHandler) UpdateBlockTimeRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateBlockTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	actionCode, err := h.handler.UpdateBlockTimeRequest(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		h.logger.Errorf("failed to update block time request: %v", err)
		http.Error(w, "Failed to update block time request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"action_code": actionCode})
}

func (h *HQHTTPHandler) UpdateArchiveTimeRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdateArchiveTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	actionCode, err := h.handler.UpdateArchiveTimeRequest(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName)
	if err != nil {
		h.logger.Errorf("failed to update archive time request: %v", err)
		http.Error(w, "Failed to update archive time request", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"action_code": actionCode})
}

func (h *HQHTTPHandler) UpdateBlockTime(w http.ResponseWriter, r *http.Request) {
	var request dto.ApproveRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.handler.UpdateBlockTime(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName); err != nil {
		h.logger.Errorf("failed to update block time: %v", err)
		http.Error(w, "Failed to update block time", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *HQHTTPHandler) UpdateArchiveTime(w http.ResponseWriter, r *http.Request) {
	var request dto.ApproveRejectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.handler.UpdateArchiveTime(r.Context(), request, claims.UserID, claims.PhoneNumber, claims.FullName); err != nil {
		h.logger.Errorf("failed to update archive time: %v", err)
		http.Error(w, "Failed to update archive time", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
