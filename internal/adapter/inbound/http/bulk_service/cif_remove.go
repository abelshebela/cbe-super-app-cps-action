package bulkservices_inbound

import (
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) SearchAccountByAccountNumber(w http.ResponseWriter, r *http.Request) {
	accountNumber, ok := common_util.GetParam(r, "account_number")
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}
	account, err := h.Application.FetchAccounts(r.Context(), accountNumber)
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

	actionCode, err := h.Application.RemoveCifMaker(r.Context(), req.Cif, maker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	res := dto.CifRemoveMakerResponse{
		ActionCode: actionCode,
		UserID:     req.UserID,
	}
	utils.WriteSuccessResponse(w, res, "successfully retrived")
}

func (h *HttpStore) RemoveCifCheckerApprove(w http.ResponseWriter, r *http.Request) {
	h.RemoveCifChecker(w, r, dto.CifRemoveCheckerRequest{
		ServiceAction: true,
	})
}

func (h *HttpStore) RemoveCifCheckerReject(w http.ResponseWriter, r *http.Request) {
	var req dto.CifRemoveCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)
		return
	}
	req.ServiceAction = false
	h.RemoveCifChecker(w, r, req)
}

func (h *HttpStore) RemoveCifChecker(w http.ResponseWriter, r *http.Request, req dto.CifRemoveCheckerRequest) {

	actionID, ok := common_util.GetParam(r, "action_code")
	if !ok {
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	checker, err := h.extractUserIDFromContext(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}
	cpsAction, err := h.Application.RemoveCifChecker(r.Context(), actionID, req.ServiceAction, req.RejectReason, checker)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.BaseResponseMaker(cpsAction, w, "action successfully completed", 200)
}
