package domain

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/hq/models"
)

type Repository interface {
	GetHQByID(ctx context.Context, id string) (models.HQ, error)
	UpdateHQ(ctx context.Context, id string, update models.HQ) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
}
