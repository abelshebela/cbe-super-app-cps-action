package budget_handler

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
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
	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "icons_image", 10<<20)
	if err != nil {
		h.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	allowedImageTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/bmp":     true,
		"image/svg+xml": true,
	}

	if !allowedImageTypes[fileHeader.Header.Get("Content-Type")] {
		h.logger.Errorf("invalid file type: %v", fileHeader.Header.Get("Content-Type"))
		common_util.SendErrorResponse(w, "INVALID_IMG_FORMAT", http.StatusBadRequest, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("Incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:          userContext.UserCode,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"icon_url": fileHeader},
		RequestAction:    entities.RequestCreateBudgetIcon,
		MakerActionTime:  time.Now(),
	}

	action, err := h.budgetService.CreateIcon(r.Context(), cpsAction)
	if err != nil {
		h.logger.Errorf("failed to create budget icon: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}

	_, err = common_util.StructToMap(action)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), 500, nil)
		return
	}
	common_util.BaseResponseMaker(map[string]interface{}{}, w, "Budget icon request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetFetchIcons(w http.ResponseWriter, r *http.Request) {

	filterParams := common_util.ExtractFilterParams(r)
	icons, err := h.budgetService.FetchIcons(r.Context(), filterParams)
	if err != nil {
		h.logger.Errorf("failed to fetch budget icons: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusBadRequest, nil)
		return
	}
	data, _ := common_util.StructToMap(icons)

	common_util.BaseResponseMaker(data, w, "Icons fetched successfully", 200)
}

func (h *BudgetHandler) BudgetUpdateIcon(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	file, fileHeader, err := common_util.ParseMultipartFormFile(r, "icons_image", 10<<20)
	if err != nil {
		h.logger.Errorf("error parsing file: %v", err)
		common_util.SendErrorResponse(w, common_util.MissingOrInvalidImage, 0, nil)
		return
	}
	defer file.Close()

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("Incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:          userContext.UserCode,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"icon_url": fileHeader},
		RequestAction:    entities.RequestUpdateBudgetIcon,
		MakerActionTime:  time.Now(),
	}

	_, err = h.budgetService.UpdateIcon(r.Context(), id, cpsAction)
	if err != nil {
		h.logger.Errorf("failed to update budger icon: %v", err)
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)

		return
	}

	common_util.BaseResponseMaker(map[string]interface{}{}, w, "Update request submitted for approval", 200)
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
	if ctx.IsIncomplete() {
		h.logger.Errorf("Incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:          ctx.UserID,
		MakerName:        ctx.FullName,
		MakerPhoneNumber: ctx.PhoneNumber,
		Department:       ctx.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"color": req.Color},
		RequestAction:    entities.RequestCreateBudgetColor,
		MakerActionTime:  time.Now(),
	}

	_, err := h.budgetService.CreateColor(r.Context(), req.Color, cpsAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}

	common_util.BaseResponseMaker(map[string]interface{}{}, w, "Color creation request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetFetchColors(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	if filterParams.Page < 1 || filterParams.PerPage < 1 {
		common_util.SendErrorResponse(w, "INVALID_INPUT_PARAMETERS", 0, nil)
		return
	}

	colors, err := h.budgetService.FetchColors(r.Context(), filterParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	data, _ := common_util.StructToMap(colors)

	common_util.BaseResponseMaker(data, w, "Colors fetched successfully", 200)
}

func (h *BudgetHandler) BudgetUpdateColor(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'id'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	var req budget.UpdateColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidJSONPayload, http.StatusBadRequest, nil)
		return
	}

	if err := req.Validate(); err != nil {
		common_util.SendErrorResponse(w, common_util.InvalidInput, http.StatusBadRequest, nil)
		return
	}

	userContext := ctx_util.ExtractUserContext(r)
	if userContext.IsIncomplete() {
		h.logger.Errorf("Incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction := entities.CPSAction{
		MakerID:          userContext.UserCode,
		MakerName:        userContext.FullName,
		MakerPhoneNumber: userContext.PhoneNumber,
		Department:       userContext.Department,
		ActionStatus:     entities.ActionPending,
		ActionType:       entities.ActionCreate,
		CurrentAction:    map[string]interface{}{"color": req.Color},
		RequestAction:    entities.RequestUpdateBudgetColor,
		MakerActionTime:  time.Now(),
		CreatedAt:        time.Now(),
		LastModifiedAt:   time.Now(),
	}

	_, err := h.budgetService.UpdateColor(r.Context(), id, req.Color, cpsAction)
	if err != nil {
		common_util.SendErrorResponse(w, err.Error(), http.StatusNotFound, nil)
		return
	}

	common_util.BaseResponseMaker(map[string]interface{}{}, w, "Color update request submitted for approval", 200)
}

func (h *BudgetHandler) BudgetCheckerApproval(w http.ResponseWriter, r *http.Request) {
	action_code, ok := common_util.GetParam(r, "action_code")
	if !ok {
		h.logger.Errorf("missing or invalid parameter 'action_code'")
		common_util.SendErrorResponse(w, common_util.InvalidInputParameters, 0, nil)
		return
	}

	ctx := ctx_util.ExtractUserContext(r)

	if ctx.IsIncomplete() {
		h.logger.Errorf("Incomplete user information")
		common_util.SendErrorResponse(w, common_util.IncompleteUserInfo, 0, nil)
		return
	}

	cpsAction := entities.CPSAction{
		ActionCode:         action_code,
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
