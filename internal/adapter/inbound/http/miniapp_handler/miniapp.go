package miniapphandler

import (
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/middleware"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) MakerCreateMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	maker_id := claims.UserID
	// Process the request using application logic
	response, err := h.Application.MakerCreateMiniApp(r.Context(), req, maker_id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "Failed to create mini app")
		return
	}
	utils.WriteSuccessResponse(w, response, "successful")
}
func (h *HttpStore) CheckerMiniApp(w http.ResponseWriter, r *http.Request) {
	var req dto.MiniAppCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		utils.WriteErrorResponse(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	checker_id := claims.UserID
	err := h.Application.CheckerCreateMiniApp(r.Context(), req.Action_id, req.Action, checker_id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	utils.WriteSuccessResponse(w, nil, "successful")
}
