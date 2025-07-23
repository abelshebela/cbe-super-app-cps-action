package service

import (
	"cbe-super-app-budget/internal/constants/dto"
	"context"
)

type BudgetService interface {
	Get(ctx context.Context) ([]dto.Budget, error)
	GetOne(ctx context.Context, id string) (*dto.Budget, error)
	Add(ctx context.Context, req dto.BudgetRequest) (*dto.Budget, error)
	Modify(ctx context.Context, req dto.BudgetRequest) error
	Remove(ctx context.Context, id string) error
}

type BudgetCategoryService interface {
	Get(ctx context.Context) ([]dto.BudgetCategory, error)
	GetOne(ctx context.Context, id string) (*dto.BudgetCategory, error)
	Add(ctx context.Context, req dto.BudgetCategoryRequest) (*dto.BudgetCategory, error)
	Modify(ctx context.Context, req dto.BudgetCategoryRequest) error
	Remove(ctx context.Context, id string) error
}

type SpendingService interface {
	Get(ctx context.Context) ([]dto.Spending, error)
	GetOne(ctx context.Context, id string) (*dto.Spending, error)
	Add(ctx context.Context, req dto.SpendingRequest) (*dto.Spending, error)
	Modify(ctx context.Context, req dto.SpendingRequest) error
	Remove(ctx context.Context, id string) error
}
