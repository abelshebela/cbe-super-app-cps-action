package budget_category

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/budget_category"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	inbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/inbound/budget_category"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/utils/common"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/utils/file"
	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetCategoryAdapter struct {
	BudgetCategoryHandler budget_category.BudgetCategoryApplictionService
	logger                utils.Logger
	fileService           file.FileService
}

func InitBudgetCategoryAdapter(
	budgetCategoryHandler budget_category.BudgetCategoryApplictionService,
	logger utils.Logger,
	fileService file.FileService,
) inbound.BudgetCategoryInbound {
	return BudgetCategoryAdapter{
		BudgetCategoryHandler: budgetCategoryHandler,
		logger:                logger,
		fileService:           fileService,
	}
}

func (b BudgetCategoryAdapter) getUserFromContext(ctx context.Context) (action.User, error) {
	userCode, ok := ctx.Value(constant.ContextKey("user_id")).(string)
	if !ok || userCode == "" {
		return action.User{}, errors.New("user_id not found in context")
	}

	fullName, _ := ctx.Value(constant.ContextKey("full_name")).(string)
	phoneNumber, _ := ctx.Value(constant.ContextKey("phone_number")).(string)

	return action.User{
		UserCode:    userCode,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}, nil
}

func (b BudgetCategoryAdapter) CreateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		b.logger.Errorf("failed to parse multipart form: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	file, fileHeader, err := r.FormFile("icon")
	if err != nil {
		b.logger.Errorf("failed to get file: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}
	defer file.Close()

	uploadResult, err := b.fileService.UploadImage(r.Context(), file, fileHeader, "budget-category-icons")
	if err != nil {
		b.logger.Errorf("upload failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	req := dto.CreateBudgetCategoryRequest{
		Name:        r.FormValue("name"),
		Icon:        uploadResult.Key,
		BucketName:  uploadResult.Bucket,
		ObjectName:  uploadResult.ObjectName,
		Description: r.FormValue("description"),
	}

	if err := req.Validate(); err != nil {
		b.logger.Errorf("validation failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	user, err := b.getUserFromContext(r.Context())
	if err != nil {
		b.logger.Errorf("user context error: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	actionId, err := b.BudgetCategoryHandler.CreateBudgetCategoryAction(r.Context(), req, user)
	if err != nil {
		b.logger.Errorf("create action error: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	common.SuccessResponse(w, http.StatusCreated, "SUCCESS", map[string]string{"action_id": actionId})
}

func (b BudgetCategoryAdapter) ApproveBudgetCategoryActionHTTP(w http.ResponseWriter, r *http.Request) {
	b.BudgetCategoryHandler.ApproveBudgetCategoryActionHTTP(w, r)
}

func (b BudgetCategoryAdapter) UpdateBudgetCategory(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		b.logger.Errorf("failed to parse multipart form: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	var req dto.UpdateBudgetCategoryRequest
	req.ID = chi.URLParam(r, "budget_category_id")
	req.Name = r.FormValue("name")
	req.Description = r.FormValue("description")

	file, fileHeader, err := r.FormFile("icon")
	if err == nil {
		defer file.Close()

		uploadResult, err := b.fileService.UploadImage(r.Context(), file, fileHeader, "budget-category-icons")
		if err != nil {
			b.logger.Errorf("upload failed: %v", err)
			common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
			return
		}

		req.Icon = uploadResult.Key
		req.BucketName = uploadResult.Bucket
		req.ObjectName = uploadResult.ObjectName
	}

	if err := req.Validate(); err != nil {
		b.logger.Errorf("validation failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	user, err := b.getUserFromContext(r.Context())
	if err != nil {
		b.logger.Errorf("user context error: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	actionId, err := b.BudgetCategoryHandler.UpdateBudgetCategoryAction(r.Context(), req, user)
	if err != nil {
		b.logger.Errorf("failed to update budget category action: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	common.SuccessResponse(w, http.StatusOK, "SUCCESS", map[string]string{"action_id": actionId})
}

func (b BudgetCategoryAdapter) DeleteBudgetCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.DeleteBudgetCategoryRequest
	req.ID = chi.URLParam(r, "budget_category_id")

	if err := req.Validate(); err != nil {
		b.logger.Errorf("delete validation failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	user, err := b.getUserFromContext(r.Context())
	if err != nil {
		b.logger.Errorf("user context error: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	actionId, err := b.BudgetCategoryHandler.DeleteBudgetCategoryAction(r.Context(), req, user)
	if err != nil {
		b.logger.Errorf("delete action failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	common.SuccessResponse(w, http.StatusOK, "SUCCESS", map[string]string{"action_id": actionId})
}

func (b BudgetCategoryAdapter) GetBudgetCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.GetBudgetCategoryRequest
	req.ID = chi.URLParam(r, "budget_category_id")

	if err := req.Validate(); err != nil {
		b.logger.Errorf("get validation failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	result, err := b.BudgetCategoryHandler.GetBudgetCategory(r.Context(), req)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("get handler failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	if err == mongo.ErrNoDocuments {
		common.ErrorResponse(w, http.StatusNotFound, "NOT_FOUND", map[string]interface{}{"error": "category not found"})
		return
	}

	common.SuccessResponse(w, http.StatusOK, "SUCCESS", map[string]interface{}{"category": result})
}

func (b BudgetCategoryAdapter) GetAllBudgetCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.GetAllBudgetCategoryRequest
	req.Page = 1
	req.Limit = 10

	req.Search = r.URL.Query().Get("search")
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = int64(page)
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			req.Limit = int64(limit)
		}
	}

	if err := req.Validate(); err != nil {
		b.logger.Errorf("list validation failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	results, err := b.BudgetCategoryHandler.GetAllBudgetCategory(r.Context(), req)
	if err != nil && err != mongo.ErrNoDocuments {
		b.logger.Errorf("get all handler failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	common.SuccessResponse(w, http.StatusOK, "SUCCESS", map[string]interface{}{"categories": results})
}
