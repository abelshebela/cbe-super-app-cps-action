package budget_category

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type Repository interface {
	CreateBudgetCategoryAction(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest, maker action.User) (string, error)
	UpdateBudgetCategoryAction(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest, maker action.User) (string, error)
	DeleteBudgetCategoryAction(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest, maker action.User) (string, error)
	ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action.User) error

	CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error

	GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*BudgetCategory, error)
}
