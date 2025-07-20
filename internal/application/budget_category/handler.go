package budget_category

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	budgetCategory "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/utils/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BudgetCategoryApplicationService interface {
	CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error)
	FindBudgetCategoryById(ctx context.Context, id string) (*budgetCategory.BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error

	ApproveAction(ctx context.Context, req dto.ApproveBudgetCategoryRequest, checker action.User) (action.CPSAction, error)
	CreateAction(ctx context.Context, data interface{}, maker action.User) (action.CPSAction, error)

	GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*budgetCategory.BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*budgetCategory.BudgetCategory, error)
	ApproveBudgetCategoryActionHTTP(w http.ResponseWriter, r *http.Request)
}

type BudgetCategoryHandler struct {
	service budgetCategory.BudgetCategoryService
	logger  utils.Logger
}

func InitBudgetCategoryHandler(service budgetCategory.BudgetCategoryService, logger utils.Logger) BudgetCategoryApplicationService {
	return &BudgetCategoryHandler{
		service: service,
		logger:  logger,
	}
}

func (h BudgetCategoryHandler) CreateBudgetCategory(ctx context.Context, req dto.CreateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error) {
	return h.service.CreateBudgetCategory(ctx, req)
}

func (h BudgetCategoryHandler) FindBudgetCategoryById(ctx context.Context, id string) (*budgetCategory.BudgetCategory, error) {
	return h.service.FindBudgetCategoryById(ctx, id)
}

func (h BudgetCategoryHandler) UpdateBudgetCategory(ctx context.Context, req dto.UpdateBudgetCategoryRequest) (budgetCategory.BudgetCategory, error) {
	return h.service.UpdateBudgetCategory(ctx, req)
}

func (h BudgetCategoryHandler) DeleteBudgetCategory(ctx context.Context, req dto.DeleteBudgetCategoryRequest) error {
	return h.service.DeleteBudgetCategory(ctx, req)
}

func (h BudgetCategoryHandler) CreateAction(ctx context.Context, data interface{}, maker action.User) (action.CPSAction, error) {
	return h.service.CreateAction(ctx, data, maker)
}

func (h BudgetCategoryHandler) ApproveAction(ctx context.Context, req dto.ApproveBudgetCategoryRequest, checker action.User) (action.CPSAction, error) {
	return h.service.ApproveAction(ctx, req, checker)
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
		h.logger.Errorf("failed to decode request body: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", map[string]interface{}{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		h.logger.Errorf("validation failed: %v", err)
		common.ErrorResponse(w, http.StatusBadRequest, "INVALID_INPUT", common.FormatValidationErrors(err))
		return
	}

	user := ctx_util.ExtractUserContext(r)
	if user.IsIncomplete() {
		h.logger.Errorf("user context is incomplete")
		common.ErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", map[string]interface{}{"error": "user context is incomplete"})
		return
	}

	checker := action.User{
		UserID:      user.UserID,
		UserCode:    user.UserCode,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}

	action, err := h.ApproveAction(r.Context(), req, checker)
	if err != nil && err == mongo.ErrNoDocuments {
		h.logger.Errorf("action not found")
		common.ErrorResponse(w, http.StatusNotFound, "ACTION_NOT_FOUND", map[string]interface{}{"error": "action not found"})
		return
	}

	if err != nil && err == dto.ErrActionApproved {
		h.logger.Errorf("action already approved")
		common.ErrorResponse(w, http.StatusBadRequest, "ACTION_ALREADY_APPROVED", map[string]interface{}{"error": "action already approved"})
		return
	}
	if err != nil && err == dto.ErrActionRejected {
		h.logger.Errorf("action already rejected")
		common.ErrorResponse(w, http.StatusBadRequest, "ACTION_ALREADY_REJECTED", map[string]interface{}{"error": "action already rejected"})
		return
	}

	if err != nil {
		h.logger.Errorf("approve action failed: %v", err)
		common.ErrorResponse(w, http.StatusInternalServerError, "UNHANDLED_SERVER_ERROR", map[string]interface{}{"error": err.Error()})
		return
	}

	common.SuccessResponse(w, http.StatusOK, "SUCCESS", map[string]interface{}{"action": action})
}
