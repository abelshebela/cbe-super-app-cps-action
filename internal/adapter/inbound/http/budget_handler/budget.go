package budget_handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
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
		common_util.SendErrorResponse(w, "Failed to parse multipart form", http.StatusBadRequest, nil)
		return
	}

	file, fileHeader, err := r.FormFile("icons_image")
	if err != nil {
		h.logger.Errorf("failed to read icons image: %c", err)
		common_util.SendErrorResponse(w, "Missing or invalid icons image", http.StatusBadRequest, nil)
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
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	data, err := common_util.StructToMap(action)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}
	common_util.BaseResponseMaker(data, w, "Budget icon request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetFetchIcons(w http.ResponseWriter, r *http.Request) {
	icons, err := h.budgetService.FetchIcons(r.Context())
	if err != nil {
		h.logger.Errorf("failed to fetch budget icons: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusInternalServerError, nil)
		return
	}

	data := map[string]interface{}{"icons": icons}
	common_util.BaseResponseMaker(data, w, "Icons fetched successfully", 200)
}

func (h *BudgetHandler) BudgetUpdateIcon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		// h.logger.Errorf("failed to parse form data: %v", err)
		// common_util.SendErrorResponse(w, "Failed to parse multipart form", http.StatusBadRequest, nil)
		return
	}

	file, fileHeader, err := r.FormFile("icons_image")
	if err != nil {
		h.logger.Errorf("icons_image error: %v", err)
		common_util.SendErrorResponse(w, "Missing or invalid icons image", http.StatusBadRequest, nil)
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
		h.logger.Errorf("failed to update budger icon: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)

		return
	}

	data, err := common_util.StructToMap(action)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
	}

	common_util.BaseResponseMaker(data, w, "Update request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetCreateColor(w http.ResponseWriter, r *http.Request) {
	var req budget.BudgetCreateColor

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidInput, http.StatusBadRequest, nil)
		return
	}

	ctx := ctx_util.ExtractUserContext(r)

	cpsAction := entities.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"color": req.Color},
		RequestAction:    entities.RequestBudgetColor,
		MakerActionTime:  time.Now(),
	}

	action, err := h.budgetService.CreateColor(r.Context(), req.Color, cpsAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}

	data, err := common_util.StructToMap(action)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(data, w, "Color creattion request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetFetchColors(w http.ResponseWriter, r *http.Request) {
	colors, err := h.budgetService.FetchColors(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data := map[string]interface{}{"colors": colors}
	common_util.BaseResponseMaker(data, w, "Colors fetched successfully", 200)
}

func (h *BudgetHandler) BudgetUpdateColor(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req budget.UpdateColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidInput, http.StatusBadRequest, nil)
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
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}

	data, err := common_util.StructToMap(action)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(data, w, "Color update request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetCheckerApproval(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "action_code")

	ctx := ctx_util.ExtractUserContext(r)

	cpsAction := entities.CPSAction{
		ActionCode:         code,
		CheckerID:          ctx.UserID,
		CheckerName:        ctx.FullName,
		CheckerPhoneNumber: ctx.PhoneNumber,
		Department:         ctx.Department,
		CheckerActionTime:  time.Now(),
	}

	action, err := h.budgetService.ApproveAction(r.Context(), cpsAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}

	data, err := common_util.StructToMap(action)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}

	common_util.BaseResponseMaker(data, w, "Action approved successfully", 200)
}
