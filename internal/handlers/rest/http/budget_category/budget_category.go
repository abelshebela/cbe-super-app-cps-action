package budget_category

import (
	"cbe-super-app-cps-action/internal/constants/interfaces/budget_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	core "cbe-super-app-cps-action/internal/handlers/rest/http/budget_category/core"
	"cbe-super-app-cps-action/internal/service"
	common_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"net/http"
	"strings"

	"cbe-super-app-cps-action/internal/constants/types"

	constants "cbe-super-app-cps-action/internal/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
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
//	@Param			name	formData	string									true	"Budget category name"	example(Monthly Groceries)
//	@Param			color	formData	string									false	"Hex color code"		example(#FF5733)
//	@Param			icon	formData	file									false	"Icon image file (png, jpg, etc)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Budget category created successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category [post]
func (b *budgetCategoryAdapter) CreateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "createBudgetCategory", "handler", "budgetCategory")
	defer span.End()

	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	req, err := core.ParseRequestFromMultipartForm(r, true)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
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
	req.Name = strings.TrimSpace(req.Name)
	err = b.budgetCategoryApplication.CreateBudgetCategory(ctx, req)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[CreateBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		b.logger.Infof("[CreateBudgetCategory] request sent successfully")
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryCreateRequestSubmittedForApproval, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryCreatedSP, nil)

	}
}

// UpdateBudgetCategory
//
//	@Summary		Update Budget Category
//	@Description	Update a budget category by the provided ID (multipart/form-data)
//	@Tags			Budget Category
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		path		string									true	"Budget category ID"
//	@Param			name	formData	string									false	"Budget category name"	example(Monthly Groceries)
//	@Param			color	formData	string									false	"Hex color code"		example(#FF5733)
//	@Param			icon	formData	file									false	"Icon image file (png, jpg, etc)"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Budget category update request submitted successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/{id} [patch]
func (b *budgetCategoryAdapter) UpdateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "updateBudgetCategory", "handler", "budgetCategory")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

	id := core.ExtractIDFromURL(r)
	if id == "" {
		b.logger.Errorf("budget category ID is required")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	req, err := core.ParseUpdateRequestFromMultipartForm(r)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("failed to parse request from multipart form: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	if err := req.Validate(); err != nil {
		span.RecordError(err)
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

	span.SetAttributes(attribute.String("budget_category.id", id))
	req.Name = strings.TrimSpace(req.Name)
	err = b.budgetCategoryApplication.UpdateBudgetCategory(ctx, id, req)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[UpdateBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {
		b.logger.Infof("[UpdateBudgetCategory] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryUpdateSubmittedForApproval, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryUpdatedSP, nil)

	}
}

// GetBudgetCategoryByID
//
//	@Summary		Get Budget Category by ID
//	@Description	Retrieve a budget category by its ID
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string																		true	"Budget category ID"
//	@Success		200	{object}	localization.StandardResponse{data=budget_category.BudgetCategoryResponse}	"Budget category retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}										"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}										"Budget category not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}										"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category/{id} [get]
func (b *budgetCategoryAdapter) GetBudgetCategoryByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "getBudgetCategoryById", "handler", "budgetCategory")
	defer span.End()
	id, ok := common_util.GetParam(r, "id")
	if !ok {
		b.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorInvalidInputParameter, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("budget_category.id", id))

	budgetCategory, err := b.budgetCategoryApplication.FetchBudgetCategoryByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[GetBudgetCategoryByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	b.logger.Infof("[GetBudgetCategoryByID] budget category retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryFetched, budgetCategory)
}

// GetAllBudgetCategories
//
//	@Summary		Get All Budget Categories
//	@Description	Retrieve all budget categories with optional filters
//	@Tags			Budget Category
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																				false	"Page number"
//	@Param			per_page	query		int																				false	"Items per page"
//	@Param			search		query		string																			false	"Searchable fields name"
//	@Param			enabled		query		bool																			false	"Status filter (e.g.,true or false)"
//	@Param			name		query		string																			false	"Status filter by  name"
//	@Success		200			{object}	localization.StandardResponse{data=[]budget_category.BudgetCategoryResponse}	"Budget categories retrieved successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}											"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}											"Internal server error"
//	@Security		BearerAuth
//	@Router			/budget-category [get]
func (b *budgetCategoryAdapter) GetAllBudgetCategories(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "getAllBudgetCategories", "handler", "budgetCategory")
	defer span.End()
	filterParams := common_util.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := common_util.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := common_util.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	budgetCategories, err := b.budgetCategoryApplication.FetchBudgetCategory(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[GetAllBudgetCategories] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("budget_category.count", len(budgetCategories.Data)))
	b.logger.Infof("[GetAllBudgetCategories] retrieved %d budget categories", len(budgetCategories.Data))
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
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "deleteBudgetCategory", "handler", "budgetCategory")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

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

	span.SetAttributes(attribute.String("budget_category.id", id))

	err := b.budgetCategoryApplication.DeleteBudgetCategory(ctx, id)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[DeleteBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {

		b.logger.Infof("[DeleteBudgetCategory] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDeleteSubmittedForApproval, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDeletedSP, nil)

	}
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
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "enableBudgetCategory", "handler", "budgetCategory")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

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

	span.SetAttributes(attribute.String("budget_category.id", id))

	err := b.budgetCategoryApplication.EnableOrDisableBudgetCategory(ctx, id, true)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[EnableBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if md.IsMakerOnly {
		b.logger.Infof("[EnableBudgetCategory] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryEnableSubmittedForApproval, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryEnabledSP, nil)

	}
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
	ctx, span := common_util.TraceLogger(r.Context(), "handler", "disableBudgetCategory", "handler", "budgetCategory")
	defer span.End()
	md := &types.ContextMetadata{}
	ctx = context.WithValue(ctx, constants.ContextKeyMetadata, md)

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

	span.SetAttributes(attribute.String("budget_category.id", id))

	err := b.budgetCategoryApplication.EnableOrDisableBudgetCategory(ctx, id, false)
	if err != nil {
		span.RecordError(err)
		b.logger.Errorf("[DisableBudgetCategory] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	if md.IsMakerOnly {

		b.logger.Infof("[DisableBudgetCategory] request sent successfully for id: %s", id)
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDisableSubmittedForApproval, nil)
	} else {
		localization.SendSuccessResponse(w, localization.SuccessBudgetCategoryDisabledSP, nil)

	}
}
