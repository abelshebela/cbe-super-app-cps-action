package bulkservices_inbound

import (
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/application/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/pkgs/utils"
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

	makerID, err := h.extractUserIDFromContext(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	actionID, err := h.Application.RemoveCifMaker(r.Context(), req.Cif, makerID)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	res := dto.CifRemoveMakerResponse{
		Action_Id: actionID,
	}
	utils.WriteSuccessResponse(w, res, "successfully retrived")
}
func (h *HttpStore) RemoveCifChecker(w http.ResponseWriter, r *http.Request) {
	var req dto.CifRemoveCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)

		return
	}
	checkerID, Cerr := h.extractUserIDFromContext(r)
	if Cerr != nil {
		utils.SendErrorResponse(w, Cerr.Error(), http.StatusUnauthorized, nil)
		return
	}
	err := h.Application.RemoveCifChecker(r.Context(), req.Action_Id, req.ServiceAction, checkerID)
	if err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, nil, "action successfull")
}
