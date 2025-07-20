package hq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/hq"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

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
		utils.SendErrorResponse(w, fmt.Errorf("Unhandled Server Error"), 500, nil)
		return
	}
	utils.BaseResponseMaker(data, w, "Successfly fetched", 200)
}

func (h *HQHTTPHandler) GetAllHQ(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page := constant.DefaultPage
	if pageInt, err := strconv.Atoi(query.Get("page")); err == nil && pageInt > 0 {
		page = pageInt
	}

	per_page := constant.DefaultPerPage
	if perPageInt, err := strconv.Atoi(query.Get("per_page")); err == nil &&
		perPageInt <= 10 && perPageInt > 0 {
		per_page = perPageInt
	}

	search := query.Get("search")
	filter := query.Get("filter")

	filterParams := &constant.Filter{
		Page:    page,
		PerPage: per_page,
		Search:  search,
		Filters: filter,
	}

	ctx := r.Context()
	customers, err := h.handler.GetHQDetail(ctx, filterParams)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)

		return
	}

	data, _ := utils.StructToMap(customers)
	utils.BaseResponseMaker(data, w, "Successfuly HQ data fetched", http.StatusAccepted)

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
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	actionCode, err := h.handler.UpdateBlockTimeRequest(r.Context(), request, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, actionCode, "Update block time request submitted for approval")
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
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}
	actionCode, err := h.handler.UpdateArchiveTimeRequest(r.Context(), request, userContext.UserID, userContext.PhoneNumber, userContext.FullName, userContext.Department)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), 0, nil)
		return
	}

	common_util.WriteSuccessResponse(w, actionCode, "Update archive time request submitted for approval")

}
