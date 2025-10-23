package budget_category

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/budget_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/budget_category/core"
	"cbe-super-app-cps-action/internal/service"
	common_util "cbe-super-app-cps-action/pkgs/utils"
	"net/http"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type budgetCategoryAdapter struct {
	budgetCategoryApplication service.BudgetCategoryService
	logger                    utils.Logger
}

func InitBudgetCategoryAdapter(budgetCategoryApplication service.BudgetCategoryService, logger utils.Logger) budget_category.BudgetCategoryPortHandler {
	return &budgetCategoryAdapter{
		logger:                    logger,
		budgetCategoryApplication: budgetCategoryApplication,
	}
}

func (b *budgetCategoryAdapter) CreateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		b.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		b.logger.Errorf("request validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err = b.budgetCategoryApplication.CreateBudgetCategory(r.Context(), req)
	if err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryRequestSubmittedForApproval, nil)
}

func (b *budgetCategoryAdapter) UpdateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	id := core.ExtractIDFromURL(r)
	if id == "" {
		b.logger.Errorf("budget category ID is required")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	req, err := core.ParseUpdateRequestFromMultipartForm(r)
	if err != nil {
		b.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		b.logger.Errorf("request validation failed: %v", err)
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err = b.budgetCategoryApplication.UpdateBudgetCategory(r.Context(), id, req)
	if err != nil {
		b.logger.Errorf("failed to update budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryUpdateSubmittedForApproval, nil)
}

func (b *budgetCategoryAdapter) GetBudgetCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	budgetCategory, err := b.budgetCategoryApplication.FetchBudgetCategoryByID(r.Context(), id)
	if err != nil {
		b.logger.Errorf("failed to fetch budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryFetched, budgetCategory)
}

func (b *budgetCategoryAdapter) GetAllBudgetCategories(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	budgetCategories, err := b.budgetCategoryApplication.FetchBudgetCategory(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("failed to fetch budget categories: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoriesFetched, budgetCategories)
}

func (b *budgetCategoryAdapter) DeleteBudgetCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := b.budgetCategoryApplication.DeleteBudgetCategory(r.Context(), id)
	if err != nil {
		b.logger.Errorf("failed to delete budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDeleteSubmittedForApproval, nil)
}

func (b *budgetCategoryAdapter) EnableBudgetCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	Enabled:= true
	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := b.budgetCategoryApplication.EnableOrDisableBudgetCategory(r.Context(), id, Enabled)
	if err != nil {
		b.logger.Errorf("failed to enable budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	action := "enabled"

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryToggleSubmittedForApproval, map[string]string{"action": action})
}
func (b *budgetCategoryAdapter) DisableBudgetCategory(w http.ResponseWriter, r *http.Request) {
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	Enabled:= false
	userContext := common_util.ExtractUserContext(r)
	if common_util.IsIncomplete(userContext) {
		b.logger.Errorf("Incomplete user information")
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := b.budgetCategoryApplication.EnableOrDisableBudgetCategory(r.Context(), id, Enabled)
	if err != nil {
		b.logger.Errorf("failed to disnable budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	action := "disabled"

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryToggleSubmittedForApproval, map[string]string{"action": action})
}

