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

// CreateBudgetCategory
//
//	@Summary		Create Budget Category
//	@Description	Create a new budget category with the provided information (multipart/form-data)
//	@Tags			Budget Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			name		formData	string	true	"Budget category name"				example(Monthly Groceries)
//	@Param			color		formData	string	false	"Hex color code"					example(#FF5733)
//	@Param			icon		formData	file	false	"Icon image file (png, jpg, etc)"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"Budget category created successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category [post]
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

// UpdateBudgetCategory
//
//	@Summary		Update Budget Category
//	@Description	Update a budget category by the provided ID (multipart/form-data)
//	@Tags			Budget Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string										true	"Budget category ID"
//	@Param			name	formData	string										false	"Budget category name"				example(Monthly Groceries)
//	@Param			color	formData	string										false	"Hex color code"					example(#FF5733)
//	@Param			icon	formData	file										false	"Icon image file (png, jpg, etc)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}		"Budget category update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}		"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}		"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/{id} [patch]
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
		b.logger.Errorf("failed to update budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryUpdateSubmittedForApproval, nil)
}

// GetBudgetCategoryByID
//
//	@Summary		Get Budget Category by ID
//	@Description	Retrieve a budget category by its ID
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Budget category ID"
//	@Success		200	{object}	localization.StandardResponse{data=budget_category.BudgetCategoryResponse}	"Budget category retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Budget category not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/{id} [get]
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

// GetAllBudgetCategories
//
//	@Summary		Get All Budget Categories
//	@Description	Retrieve all budget categories with optional filters
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			page		query	int		false	"Page number"
//	@Param			per_page	query	int		false	"Items per page"
//	@Param			search		query	string	false	"Searchable fields name"
//	@Param			enabled		query	bool	false	"Status filter (e.g.,true or false)"
//	@Param			name		query	string	false	"Status filter by  name"
//	@Success		200	{object}	localization.StandardResponse{data=[]budget_category.BudgetCategoryResponse}	"Budget categories retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category [get]
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

// DeleteBudgetCategory
//
//	@Summary		Delete Budget Category
//	@Description	Delete a budget category by the provided ID
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Budget category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Budget category delete request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/{id} [delete]
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

// EnableBudgetCategory
//
//	@Summary		Enable Budget Category
//	@Description	Enable a budget category by the provided ID
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Budget category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Budget category enable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/enable/{id} [patch]
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
		b.logger.Errorf("failed to enable budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryEnableSubmittedForApproval, nil)
}

// DisableBudgetCategory
//
//	@Summary		Disable Budget Category
//	@Description	Disable a budget category by the provided ID
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Budget category ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Budget category disable request submitted successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/disable/{id} [patch]
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
		b.logger.Errorf("failed to disnable budget category: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDisableSubmittedForApproval, nil)
}
