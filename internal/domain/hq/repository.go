package hq

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
)

type Repository interface {
	GetHQByID(ctx context.Context, id string) (HQ, error)
	UpdateHQ(ctx context.Context, id string, update HQ) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
}
