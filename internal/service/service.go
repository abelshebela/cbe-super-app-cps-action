package service

import (
	"cbe-super-app-budget/internal/constants/dto"
	"context"
)

type SpendingService interface {
	Get(ctx context.Context) ([]dto.Spending, error)
	GetOne(ctx context.Context, id string) (*dto.Spending, error)
	Add(ctx context.Context, req dto.SpendingRequest) (*dto.Spending, error)
	Modify(ctx context.Context, req dto.SpendingRequest) error
	Remove(ctx context.Context, id string) error
}
