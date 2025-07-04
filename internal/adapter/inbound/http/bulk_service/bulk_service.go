package bulkservices_inbound

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) FetchServices(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.ExtractPaginator(r)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "")
		return
	}
	resp, common_error := h.Application.FetchServices(r.Context(), int(offset), int(limit))
	if common_error != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	utils.WriteSuccessResponse(w, resp, "successfully retrived")

}
func (h *HttpStore) EnableDisableServicesMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.EnableDisableServiceMakerDtoRequest

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
	// TODO: Add further logic to handle the request
	request_id, common_error := h.Application.EnableDisableServicesMaker(r.Context(), req.ServiceId, req.ServiceAction, makerId)
	if common_error != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "")
		return
	}
	action := "enable"
	if !req.ServiceAction {
		action = "disable"
	}
	resp := dto.EnableDisableServiceMakerDtoResponse{
		ActionId: request_id,
	}
	msg := fmt.Sprintf("Action %v  on Service id :%v requested successfully", action, req.ServiceId)
	utils.WriteSuccessResponse(w, resp, msg)
}

func (h *HttpStore) EnableDisableServicesChecker(w http.ResponseWriter, r *http.Request) {
	var req dto.EnableDisableServiceCheckerDtoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}
	claims, ok := r.Context().Value("claims").(middleware.UserPayload)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	checkerId := claims.UserID
	// TODO: Add further logic to handle the request
	common_error := h.Application.EnableDisableServicesChecker(r.Context(), req.Action_Id, req.ServiceAction, checkerId)
	if common_error != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "")
		return
	}
	action := "enable"
	if !req.ServiceAction {
		action = "disable"
	}
	msg := fmt.Sprintf("Action %v  on Service id :%v implemented successfully", action, req.Action_Id)
	utils.WriteSuccessResponse(w, nil, msg)
}
