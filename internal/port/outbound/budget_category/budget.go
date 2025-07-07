package budget_category

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
)

type BudgetCategoryOutbound interface {
	CreateBudgetCategoryAction(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest, maker action_entity.User) (string, error)
	UpdateBudgetCategoryAction(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest, maker action_entity.User) (string, error)
	DeleteBudgetCategoryAction(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest, maker action_entity.User) (string, error)
	ApproveBudgetCategoryAction(ctx context.Context, actionId string, approve bool, checker action_entity.User) error
	CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (entity.BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (entity.BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error
	GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*entity.BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*entity.BudgetCategory, error)
}
