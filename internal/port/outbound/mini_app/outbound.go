package miniapp

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	// miniApp_domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
)

type Outbound interface {
	CreateMiniAppAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)
	DeleteMiniAppAction(ctx context.Context, action model.CPSAction, id string) (model.CPSAction, error)

	CreateMiniApp(ctx context.Context, miniapp model.MiniApp) error
	GetMiniAppActionId(ctx context.Context, action_id string) (model.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)
	ListMiniApp(ctx context.Context) ([]*model.MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}
