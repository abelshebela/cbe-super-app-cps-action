package hq

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/hq"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	// common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

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
	data, err := utils.StructToMap(hqResp)
	if err != nil {
		utils.SendErrorResponse(w, fmt.Errorf("unhandled Server Error"), 500, nil)
		return
	}
	utils.BaseResponseMaker(data, w, "Successfly fetched", 200)
}

func (h *HQHTTPHandler) GetAllHQ(w http.ResponseWriter, r *http.Request) {
	filterParams := utils.ExtractFilterParams(r)
	ctx := r.Context()

	hqResp, err := h.handler.GetHQDetail(ctx, filterParams)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	data, _ := utils.StructToMap(hqResp)
	utils.BaseResponseMaker(data, w, "HQs fetched successfully", http.StatusOK)
}

func (h *HQHTTPHandler) GetBlockTime(w http.ResponseWriter, r *http.Request) {

	resp, err := h.handler.GetBlockTime(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	data, _ := utils.StructToMap(resp)
	utils.BaseResponseMaker(data, w, "Block time fetched successfully", http.StatusOK)
}

func (h *HQHTTPHandler) GetArchiveTime(w http.ResponseWriter, r *http.Request) {
	resp, err := h.handler.GetArchiveTime(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	data, _ := utils.StructToMap(resp)
	utils.BaseResponseMaker(data, w, "Archive time fetched successfully", http.StatusOK)
}

func (h *HQHTTPHandler) GetPasswordExpiry(w http.ResponseWriter, r *http.Request) {
	resp, err := h.handler.GetPasswordExpiry(r.Context())
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	data, _ := utils.StructToMap(resp)
	utils.BaseResponseMaker(data, w, "Password expiry fetched successfully", http.StatusOK)
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
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}
	_, err := h.handler.UpdateBlockTimeRequest(r.Context(), request, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, 0, "Update block time request submitted for approval")
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
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}
	action, err := h.handler.UpdateArchiveTimeRequest(r.Context(), request, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, map[string]string{"action_code": action.ActionCode}, "Update archive time request submitted for approval")
}

func (h *HQHTTPHandler) UpdatePasswordExpiryRequest(w http.ResponseWriter, r *http.Request) {
	var request dto.UpdatePasswordExpiryRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.SendErrorResponse(w, "INVALID_JSON_PAYLOAD", 0, nil)
		return
	}
	if err := request.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}
	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		utils.SendErrorResponse(w, utils.IncompleteUserInfo, 0, nil)
		return
	}
	action, err := h.handler.UpdatePasswordExpiryRequest(r.Context(), request, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}
	utils.WriteSuccessResponse(w, map[string]string{"action_code": action.ActionCode}, "Update password expiry request submitted for approval")
}
