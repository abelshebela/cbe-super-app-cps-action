package outbound

import (
	"context"

	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type CPSOutboundInfra interface {
	CreateCPSAction(ctx context.Context, Action domain.CPSAction) (*domain.CPSAction, error)
	UpdateCPSAction(ctx context.Context, action domain.CPSAction) (*domain.CPSAction, error)
	CPSActionExists(ctx context.Context, uniqueID string) (bool, error)
}
