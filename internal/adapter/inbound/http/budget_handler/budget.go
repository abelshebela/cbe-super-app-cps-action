package budget_handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/middleware"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget"
)

type BudgetHandler struct {
	budgetService budget.BudgetService
	logger        utils.Logger
}

func NewBudgetHTTPHandler(budgetService budget.BudgetService, logger utils.Logger) inbound.BudgetPortHandler {
	return &BudgetHandler{
		budgetService: budgetService,
		logger:        logger,
	}
}

func (h *BudgetHandler) CreateBudgetIcon(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Errorf("failed to parse form data: %v", err)
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	file, fileHeader, err := r.FormFile("icons_image")
	if err != nil {
		err = fmt.Errorf("failed to read icons image: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid icons image",
		})
		h.logger.Errorf("icons_image error: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}
	defer file.Close()

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction := entities.CPSAction{
		MakerID:          user_code,
		MakerName:        full_name,
		MakerPhoneNumber: phone_number,
		Department:       department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"icon_url": fileHeader},
		RequestAction:    entities.RequestBudgetIcon,
		MakerActionTime:  time.Now(),
	}

	action, err := h.budgetService.CreateIcon(r.Context(), cpsAction)
	if err != nil {
		h.logger.Errorf("failed to create budget icon: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}

	res := common.Response[*entities.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           action,
	}

	res.SendJSON()
}

func (h *BudgetHandler) BudgetFetchIcons(w http.ResponseWriter, r *http.Request) {
	icons, err := h.budgetService.FetchIcons(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := common.Response[[]*entities.Icon]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           icons,
	}

	res.SendJSON()
}

func (h *BudgetHandler) BudgetUpdateIcon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Errorf("failed to parse form data: %v", err)
		err = fmt.Errorf("failed to parse multipart form: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "invalid multipart form",
		})
		middleware.ErrorHandler(w, err)
		return
	}

	file, fileHeader, err := r.FormFile("icons_image")
	if err != nil {
		err = fmt.Errorf("failed to read icons image: %w", constant.ErrorDefinition{
			Code:    http.StatusBadRequest,
			Message: "missing or invalid icons image",
		})
		h.logger.Errorf("icons_image error: %v", err)
		middleware.ErrorHandler(w, err)
		return
	}
	defer file.Close()

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction := entities.CPSAction{
		MakerID:          user_code,
		MakerName:        full_name,
		MakerPhoneNumber: phone_number,
		Department:       department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"icon_url": fileHeader},
		RequestAction:    entities.RequestBudgetIcon,
		MakerActionTime:  time.Now(),
	}

	action, err := h.budgetService.UpdateIcon(r.Context(), id, cpsAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := common.Response[*entities.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           action,
	}

	res.SendJSON()
}

func (h *BudgetHandler) BudgetCreateColor(w http.ResponseWriter, r *http.Request) {
	var req budget.BudgetCreateColor

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction := entities.CPSAction{
		MakerID:          user_code,
		MakerName:        full_name,
		MakerPhoneNumber: phone_number,
		Department:       department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"color": req.Color},
		RequestAction:    entities.RequestBudgetColor,
		MakerActionTime:  time.Now(),
	}

	action, err := h.budgetService.CreateColor(r.Context(), req.Color, cpsAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := common.Response[*entities.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           action,
	}

	res.SendJSON()
}

func (h *BudgetHandler) BudgetFetchColors(w http.ResponseWriter, r *http.Request) {
	colors, err := h.budgetService.FetchColors(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := common.Response[[]*entities.Color]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           colors,
	}

	res.SendJSON()
}

func (h *BudgetHandler) BudgetUpdateColor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req budget.UpdateColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction := entities.CPSAction{
		MakerID:          user_code,
		MakerName:        full_name,
		MakerPhoneNumber: phone_number,
		Department:       department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"color": req.Color},
		RequestAction:    entities.RequestBudgetColor,
		MakerActionTime:  time.Now(),
	}

	action, err := h.budgetService.UpdateColor(r.Context(), id, req.Color, cpsAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := common.Response[*entities.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           action,
	}

	res.SendJSON()
}

func (h *BudgetHandler) BudgetCheckerApproval(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "action_code")

	user_code := r.Context().Value(constant.ContextKey("user_code")).(string)
	full_name := r.Context().Value(constant.ContextKey("full_name")).(string)
	phone_number := r.Context().Value(constant.ContextKey("phone_number")).(string)
	department := r.Context().Value(constant.ContextKey("department")).(string)

	cpsAction := entities.CPSAction{
		ActionCode:         code,
		CheckerID:          user_code,
		CheckerName:        full_name,
		CheckerPhoneNumber: phone_number,
		Department:         department,
		CheckerActionTime:  time.Now(),
	}

	action, err := h.budgetService.ApproveAction(r.Context(), cpsAction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := common.Response[*entities.CPSAction]{
		ResponseWriter: w,
		Status:         http.StatusOK,
		Data:           action,
	}

	res.SendJSON()
}
