package miniapp

import (
	"context"

	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	miniApp_domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/miniapp"
)

type Outbound interface {
	CreateMiniAppAction(ctx context.Context, action domain.CPSAction) (domain.CPSAction, error)
	CreateMiniApp(ctx context.Context, miniapp miniApp_domain.MiniApp) error

	GetMiniAppActionId(ctx context.Context, action_id string) (domain.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action domain.CPSAction) error
}
