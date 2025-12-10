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
		b.logger.Errorf("[CreateBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[CreateBudgetCategory] request sent successfully")
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
		localization.SendBadRequestResponse(w, err.Error())
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
		b.logger.Errorf("[UpdateBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[UpdateBudgetCategory] request sent successfully for id: %s", id)
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
		b.logger.Errorf("[GetBudgetCategoryByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[GetBudgetCategoryByID] budget category retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryFetched, budgetCategory)
}

func (b *budgetCategoryAdapter) GetAllBudgetCategories(w http.ResponseWriter, r *http.Request) {
	filterParams := common_util.ExtractFilterParams(r)

	budgetCategories, err := b.budgetCategoryApplication.FetchBudgetCategory(r.Context(), filterParams)
	if err != nil {
		b.logger.Errorf("[GetAllBudgetCategories] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[GetAllBudgetCategories] retrieved %d budget categories", len(budgetCategories.Data))
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
		b.logger.Errorf("[DeleteBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[DeleteBudgetCategory] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDeleteSubmittedForApproval, nil)
}

func (b *budgetCategoryAdapter) EnableBudgetCategory(w http.ResponseWriter, r *http.Request) {
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

	err := b.budgetCategoryApplication.EnableOrDisableBudgetCategory(r.Context(), id, true)
	if err != nil {
		b.logger.Errorf("[EnableBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[EnableBudgetCategory] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryEnableSubmittedForApproval, nil)
}
func (b *budgetCategoryAdapter) DisableBudgetCategory(w http.ResponseWriter, r *http.Request) {
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

	err := b.budgetCategoryApplication.EnableOrDisableBudgetCategory(r.Context(), id, false)
	if err != nil {
		b.logger.Errorf("[DisableBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[DisableBudgetCategory] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDisableSubmittedForApproval, nil)
}
