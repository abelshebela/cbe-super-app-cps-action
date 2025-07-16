package budget_category

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget_category"
)

type BudgetCategoryOutbound interface {
	CreateAction(ctx context.Context, data interface{}, maker action_entity.User) (action_entity.CPSAction, error)
	ApproveAction(ctx context.Context, approveRequest dto.ApproveBudgetCategoryRequest, checker action_entity.User) (action_entity.CPSAction, error)

	CreateBudgetCategory(ctx context.Context, budgetCategory dto.CreateBudgetCategoryRequest) (entity.BudgetCategory, error)
	UpdateBudgetCategory(ctx context.Context, budgetCategory dto.UpdateBudgetCategoryRequest) (entity.BudgetCategory, error)
	DeleteBudgetCategory(ctx context.Context, budgetCategory dto.DeleteBudgetCategoryRequest) error

	GetBudgetCategory(ctx context.Context, budgetCategory dto.GetBudgetCategoryRequest) (*entity.BudgetCategory, error)
	GetAllBudgetCategory(ctx context.Context, getAllBudgetCategory dto.GetAllBudgetCategoryRequest) ([]*entity.BudgetCategory, error)

	FindActionById(ctx context.Context, actionId string) (*action_entity.CPSAction, error)

	FindBudgetCategoryById(ctx context.Context, id string) (*entity.BudgetCategory, error)

	UpdateAction(ctx context.Context, actionId string, checker action_entity.User, status action_entity.ActionStatus) (action_entity.CPSAction, error)
}
