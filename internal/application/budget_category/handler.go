package budget_category

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	budgetCategory "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/utils/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetCategoryApplictionService interface {
	CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error

	ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action.User) error
	CreateBudgetCategoryAction(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest, maker action.User) (string, error)
	UpdateBudgetCategoryAction(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest, maker action.User) (string, error)
	DeleteBudgetCategoryAction(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest, maker action.User) (string, error)
	GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*budgetCategory.BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*budgetCategory.BudgetCategory, error)
	ApproveBudgetCategoryActionHTTP(w http.ResponseWriter, r *http.Request)
}

type BudgetCategoryHandler struct {
	service budgetCategory.BudgetCategoryService
	logger  utils.Logger
}

func InitBudgetCategoryHandler(service budgetCategory.BudgetCategoryService, logger utils.Logger) BudgetCategoryApplictionService {
	return &BudgetCategoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h BudgetCategoryHandler) CreateBudgetCategory(ctx context.Context, req dto.CreateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error) {
	return h.service.CreateBudgetCategory(ctx, req)
}

func (h BudgetCategoryHandler) UpdateBudgetCategory(ctx context.Context, req dto.UpdateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error) {
	return h.service.UpdateBudgetCategory(ctx, req)
}

func (h BudgetCategoryHandler) DeleteBudgetCategory(ctx context.Context, req dto.DeleteBudgetCategoryRequest) error {
	return h.service.DeleteBudgetCategory(ctx, req)
}

func (h BudgetCategoryHandler) CreateBudgetCategoryAction(ctx context.Context, req dto.CreateBudgetCategoryRequest, maker action.User) (string, error) {
	return h.service.CreateBudgetCategoryAction(ctx, req, maker)
}

func (h BudgetCategoryHandler) UpdateBudgetCategoryAction(ctx context.Context, req dto.UpdateBudgetCategoryRequest, maker action.User) (string, error) {
	return h.service.UpdateBudgetCategoryAction(ctx, req, maker)
}

func (h BudgetCategoryHandler) DeleteBudgetCategoryAction(ctx context.Context, req dto.DeleteBudgetCategoryRequest, maker action.User) (string, error) {
	return h.service.DeleteBudgetCategoryAction(ctx, req, maker)
}

func (h BudgetCategoryHandler) ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action.User) error {
	return h.service.ApproveBudgetCategoryAction(ctx, actionId, approve, checker)
}

func (h BudgetCategoryHandler) GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*budgetCategory.BudgetCategory, error) {
	return h.service.GetBudgetCategory(ctx, budgetCategory)
}

func (h BudgetCategoryHandler) GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*budgetCategory.BudgetCategory, error) {
	return h.service.GetAllBudgetCategory(ctx, getAllBudgetCategory)
}

func (h BudgetCategoryHandler) ApproveBudgetCategoryActionHTTP(w http.ResponseWriter, r *http.Request) {
	var req dto.ApproveBudgetCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}
	if req.ActionID == "" {
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": "action_id is required"})
		return
	}

	checker := action.User{
		UserCode:    r.Context().Value(constant.ContextKey("user_id")).(string),
		FullName:    r.Context().Value(constant.ContextKey("full_name")).(string),
		PhoneNumber: r.Context().Value(constant.ContextKey("phone_number")).(string),
	}

	err := h.ApproveBudgetCategoryAction(r.Context(), req.ActionID, req.Approve, checker)
	if err != nil && err == mongo.ErrNoDocuments {
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": "action not found"})
		return
	}

	if err != nil && err == dto.ErrActionApproved {
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": "action already approved"})
		return
	}
	if err != nil && err == dto.ErrActionRejected {
		common.ErrorResponse(w, http.StatusBadRequest, "BAD_REQUEST", map[string]interface{}{"error": "action already rejected"})
		return
	}

	if err != nil {
		common.ErrorResponse(w, http.StatusInternalServerError, "BAD_REQUEST", map[string]interface{}{"error": err.Error()})
		return
	}

	common.SuccessResponse(w, http.StatusOK, "SUCCESS", map[string]string{"status": "success"})
}
