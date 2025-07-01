package eventhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cbe-super-app-cps-action/internal/application/dto"
	"cbe-super-app-cps-action/internal/application/middleware"
	"cbe-super-app-cps-action/pkgs/utils"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func (h *HttpStore) MakerCreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.EventCreateRequest
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
	maker_id := claims.UserID
	maker_phone := claims.PhoneNumber

	request_id, err := h.Application.MakerCreateEvent(r.Context(), req, maker_id, "", maker_phone)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	res := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: request_id,
	}
	utils.WriteSuccessResponse(w, res, "successfull")
}
func (h *HttpStore) CheckerEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.EventCheckerRequest
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
	checker_phone := claims.PhoneNumber
	err := h.Application.CheckerCreateEvent(r.Context(), req.Request_Id, req.Action, checker_id, "", checker_phone)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	utils.WriteSuccessResponse(w, nil, "success")
}

func (h *HttpStore) FetchEventById(w http.ResponseWriter, r *http.Request) {
	var req dto.EventFetchRequest
	resp, err := h.Application.FetchEvent(r.Context(), req.RequestId)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "bad request")
		return
	}
	utils.WriteSuccessResponse(w, resp, "successfully retrived")
}
func (h *HttpStore) FetchEvent(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = DefaultPage
	}
	pageSize, _ := strconv.Atoi(q.Get("page_size"))
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset := (page - 1) * pageSize
	limit := pageSize
	resp, err := h.Application.FetchAllEvents(r.Context(), limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNoContent, "")
		return
	}
	utils.WriteSuccessResponse(w, resp, "successfull")
}
