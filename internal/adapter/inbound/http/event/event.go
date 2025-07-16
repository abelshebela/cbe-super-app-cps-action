package eventhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	// "fmt"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

func (h *HttpStore) MakerCreateEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.EventCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_REQUEST_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	UserID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	FullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	PhoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	Department, _ := r.Context().Value(constant.ContextKey("department")).(string)
	if UserID == "" || FullName == "" || PhoneNumber == "" || Department == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}
	fmt.Println(UserID, FullName, PhoneNumber, Department, "nodjghdjkfgheidgjfdk")
	request_id, err := h.Application.MakerCreateEvent(r.Context(), req, UserID, FullName, PhoneNumber)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	data := map[string]interface{}{"request_id": request_id}
	utils.BaseResponseMaker(data, w, "Event creation request submitted successfully", http.StatusOK)
}

func (h *HttpStore) CheckerEvent(w http.ResponseWriter, r *http.Request) {
	var req dto.EventCheckerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_REQUEST_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}

	UserID, _ := r.Context().Value(constant.ContextKey("user_id")).(string)
	FullName, _ := r.Context().Value(constant.ContextKey("full_name")).(string)
	PhoneNumber, _ := r.Context().Value(constant.ContextKey("phone_number")).(string)
	Department, _ := r.Context().Value(constant.ContextKey("department")).(string)
	if UserID == "" || FullName == "" || PhoneNumber == "" || Department == "" {
		utils.SendErrorResponse(w, "UNAUTHORIZED", http.StatusUnauthorized, nil)
		return
	}

	err := h.Application.CheckerCreateEvent(r.Context(), req.Request_Id, req.Action, UserID, FullName, PhoneNumber)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.BaseResponseMaker(map[string]interface{}{"status": "success"}, w, "Event request reviewed successfully", http.StatusOK)
}

func (h *HttpStore) FetchEventById(w http.ResponseWriter, r *http.Request) {
	var req dto.EventFetchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "INVALID_REQUEST_PAYLOAD", http.StatusBadRequest, nil)
		return
	}
	if err := req.Validate(); err != nil {
		utils.SendErrorResponse(w, "INVALID_INPUT", http.StatusBadRequest, map[string]interface{}{"errors": err})
		return
	}
	resp, err := h.Application.FetchEvent(r.Context(), req.RequestId)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	data, err := utils.StructToMap(resp)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.BaseResponseMaker(data, w, "Event successfully retrieved", http.StatusOK)
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
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	data, err := utils.StructToMap(resp)
	if err != nil {
		utils.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}
	utils.BaseResponseMaker(data, w, "Events successfully retrieved", http.StatusOK)
}
