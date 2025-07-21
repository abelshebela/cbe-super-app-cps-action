package miniapp

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type MiniRepository interface {
	CreateMiniAppAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	CreateMiniApp(context.Context, *entities.CPSAction) (*model.MiniApp, error)
	GetMiniAppActionId(ctx context.Context, action_id string) (model.CPSAction, error)
	UpdateCpsAction(ctx context.Context, action model.CPSAction) (model.CPSAction, error)
	ListMiniApp(ctx context.Context) ([]*model.MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
	UpdateMinApp(ctx context.Context, action *entities.CPSAction, id string) (*model.MiniApp, error)
	DeleteMiniAppAction(ctx context.Context, id string) (*model.MiniApp, error)
	ExtractActionData(currentAction interface{}) (*model.MiniApp, error)
}
