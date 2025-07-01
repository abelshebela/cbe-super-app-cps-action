package bulkservices_inbound

import (
	"encoding/json"
	"net/http"

	"cbe-super-app-cps-action/internal/application/dto"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) SearchAccountByCif(w http.ResponseWriter, r *http.Request) {
	var req dto.CifSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	account, err := h.Application.FetchCifs(r.Context(), req.Cif)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	utils.WriteSuccessResponse(w, account, "successfully retrived")
}
func (h *HttpStore) RemoveCifMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.CifRemoveMakerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	makerId := claims.UserID
	action_id, err := h.Application.RemoveCifMaker(r.Context(), req.Cif, makerId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	res := dto.CifRemoveMakerResponse{
		Action_Id: action_id,
	}
	utils.WriteSuccessResponse(w, res, "successfully retrived")
}
func (h *HttpStore) RemoveCifChecker(w http.ResponseWriter, r *http.Request) {
	var req dto.CifRemoveCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	checker_id := claims.UserID
	err := h.Application.RemoveCifChecker(r.Context(), req.Action_Id, req.ServiceAction, checker_id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}
	utils.WriteSuccessResponse(w, nil, "action successfull")
}
