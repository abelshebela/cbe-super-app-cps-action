package bulkservices_inbound

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

func (h *HttpStore) extractUserIDFromContext(r *http.Request) (string, error) {
	userCtx := ctx_util.ExtractUserContext(r)
	if userCtx.IsIncomplete() {
		return "", fmt.Errorf(utils.IncompleteUserInfo)
	}
	return userCtx.UserID, nil
}

func (h *HttpStore) buildServiceActionMessage(serviceID string, isEnable bool, verb string) string {
	action := "enable"
	if !isEnable {
		action = "disable"
	}
	return fmt.Sprintf("Action %s on Service id :%s %s successfully", action, serviceID, verb)
}

func (h *HttpStore) FetchServices(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := utils.ExtractPaginator(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	resp, common_error := h.Application.FetchServices(r.Context(), int(offset), int(limit))
	if common_error != nil {
		utils.SendErrorResponse(w, common_error.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, resp, "successfully retrived")
}
func (h *HttpStore) EnableDisableServicesMaker(w http.ResponseWriter, r *http.Request) {
	var req dto.EnableDisableServiceMakerDtoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)
		return
	}

	makerID, err := h.extractUserIDFromContext(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	request_id, common_error := h.Application.EnableDisableServicesMaker(r.Context(), req.ServiceId, req.ServiceAction, makerID)
	if common_error != nil {
		utils.SendErrorResponse(w, common_error.Error(), 0, nil)
		return
	}

	resp := dto.EnableDisableServiceMakerDtoResponse{ActionId: request_id}
	msg := h.buildServiceActionMessage(req.ServiceId, req.ServiceAction, "requested")
	utils.WriteSuccessResponse(w, resp, msg)
}

func (h *HttpStore) EnableDisableServicesChecker(w http.ResponseWriter, r *http.Request) {
	var req dto.EnableDisableServiceCheckerDtoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, utils.InvalidInput, 0, nil)
		return
	}

	checkerID, err := h.extractUserIDFromContext(r)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusUnauthorized, nil)
		return
	}

	common_error := h.Application.EnableDisableServicesChecker(r.Context(), req.Action_Id, req.ServiceAction, checkerID)
	if common_error != nil {
		utils.SendErrorResponse(w, common_error.Error(), 0, nil)
		return
	}

	msg := h.buildServiceActionMessage(req.Action_Id, req.ServiceAction, "implemented")
	utils.WriteSuccessResponse(w, nil, msg)
}
