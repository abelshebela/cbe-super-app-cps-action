package miniapp

import (
	"context"

	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type Repository interface {
	CreateMiniAppAction(ctx context.Context, action domain.CPSAction) (domain.CPSAction, error)
	CreateMiniApp(ctx context.Context, miniapp MiniApp) error

	GetMiniAppActionId(ctx context.Context, action_id string) (domain.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action domain.CPSAction) error
}
