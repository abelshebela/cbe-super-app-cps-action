package budget

import (
	budget_dto "cbe-super-app-cps-action/internal/constants/dto/budget"
	"cbe-super-app-cps-action/internal/constants/interfaces/budget"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/handlers/rest/http/budget/core"
	"cbe-super-app-cps-action/internal/service"
	common_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type budget_icon_resp *model.Icon
type budget_icons_paginated_resp *types.PaginatedResponse[[]*model.Icon]
type budget_color_resp *model.Color
type budget_colors_paginated_resp *types.PaginatedResponse[[]*model.Color]

type budgetAdapter struct {
	budgetApplication service.BudgetService
	logger            utils.Logger
}

func InitBudgetAdapter(budgetApplication service.BudgetService, logger utils.Logger) budget.BudgetPortHandler {
	return &budgetAdapter{
		logger:            logger,
		budgetApplication: budgetApplication,
	}
}

// CreateBudgetIcon creates a new budget icon
// @Summary Create budget icon
// @Description Creates a new budget icon by uploading an image file
// @Tags Budget Icons
// @Accept multipart/form-data
// @Produce json
// @Param icons_image formData file true "Budget icon image file"
// @Success 200 {object} localization.StandardResponse{data=nil} "Budget icon creation request submitted for approval"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Invalid file or missing file"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/icons [post]
func (b *budgetAdapter) CreateBudgetIcon(w http.ResponseWriter, r *http.Request) {
	// file, fileHeader, err := core.ParseMultipartFormFile(r, "icons_image", 10<<20)
	// if err != nil {
	// 	b.logger.Errorf("error parsing file: %v", err)
	// 	localization.SendErrorByCodeResponse(w, err.Error())
	// 	return
	// }

	file, fileHeader, err := core.ParseMultipartFormFile(r, "icons_image", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorBankImageMissingOrInvalid, nil, nil)
		return
	}
	defer file.Close()

	// if err := core.FileValidator(w, file, *fileHeader, b.logger); err != nil {
	// 	localization.SendErrorByCodeResponse(w, err.Error())
	// 	return
	// }

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err = b.budgetApplication.CreateBudgetIcon(r.Context(), fileHeader, &file)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetIconRequestSubmittedForApproval, nil)
}

// BudgetFetchIcons retrieves all budget icons with pagination
// @Summary Get all budget icons
// @Description Retrieves a paginated list of all budget icons with optional filtering
// @Tags Budget Icons
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=budget_icons_paginated_resp} "Budget icons retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/icons [get]
func (b *budgetAdapter) BudgetFetchIcons(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	icons, err := b.budgetApplication.BudgetFetchIcons(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget icons: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetIconsFetched, icons)
}

// BudgetUpdateIcon updates an existing budget icon
// @Summary Update budget icon
// @Description Updates an existing budget icon by uploading a new image file
// @Tags Budget Icons
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Icon ID"
// @Param icons_image formData file true "Updated budget icon image file"
// @Success 200 {object} localization.StandardResponse{data=nil} "Budget icon update request submitted for approval"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Invalid file or missing file"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Icon not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/icons/{id} [put]
func (b *budgetAdapter) BudgetUpdateIcon(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	file, fileHeader, err := core.ParseMultipartFormFile(r, "icons_image", 10<<20)
	if err != nil {
		b.logger.Errorf("error parsing file: %v", err)
		localization.SendErrorResponse(w, localization.ErrorMissingOrInvalidImage, nil, nil)
		return
	}

	// if err := core.FileValidator(w, file, *fileHeader, b.logger); err != nil {
	// 	localization.SendErrorByCodeResponse(w, err.Error())
	// 	return
	// }

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err = b.budgetApplication.BudgetUpdateIcon(r.Context(), id, fileHeader, &file)
	if err != nil {
		b.logger.Errorf("failed to update budger icon: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())

		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetIconRequestSubmittedForApproval, nil)
}

// BudgetCreateColor creates a new budget color
// @Summary Create budget color
// @Description Creates a new budget color
// @Tags Budget Colors
// @Accept json
// @Produce json
// @Param request body budget_dto.BudgetCreateColor true "Budget color creation request"
// @Success 200 {object} localization.StandardResponse{data=nil} "Budget color creation request submitted for approval"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Invalid input"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/colors [post]
func (b *budgetAdapter) BudgetCreateColor(w http.ResponseWriter, r *http.Request) {
	var req budget_dto.BudgetCreateColor

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		b.logger.Errorf("error decoding request body: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	color := &model.Color{
		Color: req.Color,
	}
	err := b.budgetApplication.BudgetCreateColor(r.Context(), color)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetColorRequestSubmittedForApproval, nil)

}

// BudgetFetchColors retrieves all budget colors with pagination
// @Summary Get all budget colors
// @Description Retrieves a paginated list of all budget colors with optional filtering
// @Tags Budget Colors
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} localization.StandardResponse{data=budget_colors_paginated_resp} "Budget colors retrieved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/colors [get]
func (b *budgetAdapter) BudgetFetchColors(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)
	if filterParams.Page < 1 {
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	colors, err := b.budgetApplication.BudgetFetchColors(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget colors: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetColorsFetched, colors)
}

// BudgetUpdateColor updates an existing budget color
// @Summary Update budget color
// @Description Updates an existing budget color
// @Tags Budget Colors
// @Accept json
// @Produce json
// @Param id path string true "Color ID"
// @Param request body budget_dto.UpdateColorRequest true "Budget color update request"
// @Success 200 {object} localization.StandardResponse{data=nil} "Budget color update request submitted for approval"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Invalid input"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Color not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/colors/{id} [put]
func (b *budgetAdapter) BudgetUpdateColor(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	var req budget_dto.UpdateColorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidJSONPayload, nil, nil)
		return
	}

	if err := req.Validate(); err != nil {
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	color := &model.Color{
		Color: req.Color,
	}

	err := b.budgetApplication.BudgetUpdateColor(r.Context(), id, color)
	if err != nil {
		b.logger.Errorf("failed to update budget color: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetColorUpdateSubmittedForApproval, nil)
}

// BudgetCheckerApproval approves a budget action
// @Summary Approve budget action
// @Description Approves a budget action by action code
// @Tags Budget Actions
// @Accept json
// @Produce json
// @Param action_code path string true "Action Code"
// @Success 200 {object} localization.StandardResponse{data=nil} "Budget action approved successfully"
// @Failure 400 {object} localization.StandardResponse{data=nil} "Bad request - Invalid action code"
// @Failure 401 {object} localization.StandardResponse{data=nil} "Unauthorized"
// @Failure 403 {object} localization.StandardResponse{data=nil} "Forbidden"
// @Failure 404 {object} localization.StandardResponse{data=nil} "Action not found"
// @Failure 500 {object} localization.StandardResponse{data=nil} "Internal server error"
// @Security BearerAuth
// @Router /budgets/actions/approve/{action_code} [post]
func (b *budgetAdapter) BudgetCheckerApproval(w http.ResponseWriter, r *http.Request) {
	actionCode, ok := common_util.GetParam(r, "action_code")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'action_code'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := b.budgetApplication.BudgetCheckerApproval(r.Context(), actionCode)
	if err != nil {
		b.logger.Errorf("failed to approve action: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCheckerActionApproved, nil)
}
