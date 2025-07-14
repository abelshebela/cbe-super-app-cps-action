package bulkservices_inbound

import (
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) SearchAccountByCif(w http.ResponseWriter, r *http.Request) {
	var req dto.CifSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)
		return
	}
	account, err := h.Application.FetchCifs(r.Context(), req.Cif)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, account, "successfully retrived")
}
func (h *HttpStore) RemoveCifMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.CifRemoveMakerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)
		return
	}

	maker, err := h.extractUserIDFromContext(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	actionID, err := h.Application.RemoveCifMaker(r.Context(), req.Cif, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	res := dto.CifRemoveMakerResponse{
		ActionID: actionID,
		UserID:   req.UserID,
	}
	utils.WriteSuccessResponse(w, res, "successfully retrived")
}
func (h *HttpStore) RemoveCifChecker(w http.ResponseWriter, r *http.Request) {
	var req dto.CifRemoveCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)

		return
	}
	checker, err := h.extractUserIDFromContext(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}
	err = h.Application.RemoveCifChecker(r.Context(), req.ActionID, req.ServiceAction, checker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, nil, "action successfully completed")
}
