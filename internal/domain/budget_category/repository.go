package budget_category

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type BudgetCategoryRepository interface {
	CreateAction(ctx context.Context, data interface{}, maker action.User) (action.CPSAction, error)
	UpdateAction(ctx context.Context, actionId string, checker action.User, status action.ActionStatus) (action.CPSAction, error)
	FindActionById(ctx context.Context, actionId string) (*action.CPSAction, error)
	ApproveAction(ctx context.Context, approveRequest dto.ApproveBudgetCategoryRequest, checker action.User) (action.CPSAction, error)

	CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (BudgetCategory, error)
	FindBudgetCategoryById(ctx context.Context, id string) (*BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error

	GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*BudgetCategory, error)
}
