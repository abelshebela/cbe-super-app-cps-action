package storage

import (
	"cbe-super-app-budget/internal/constants/dto"
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type SpendingRepository interface {
	Save(ctx context.Context, spending *dto.Spending) error
	Find(ctx context.Context) ([]dto.Spending, error)
	FindById(ctx context.Context, filter bson.M) (*dto.Spending, error)
	Update(ctx context.Context, update *dto.Spending) error
	Delete(ctx context.Context, id string) error
}

type BudgetRepository interface {
	Save(ctx context.Context, budget *dto.Budget) error
	Find(ctx context.Context) ([]dto.Budget, error)
	FindById(ctx context.Context, filter bson.M) (*dto.Budget, error)
	Update(ctx context.Context, update *dto.Budget) error
	Delete(ctx context.Context, id string) error
}

type BudgetCategoryRepository interface {
	Save(ctx context.Context, budget *dto.BudgetCategory) error
	Find(ctx context.Context) ([]dto.BudgetCategory, error)
	FindById(ctx context.Context, filter bson.M) (*dto.BudgetCategory, error)
	Update(ctx context.Context, update *dto.BudgetCategory) error
	Delete(ctx context.Context, id string) error
}
