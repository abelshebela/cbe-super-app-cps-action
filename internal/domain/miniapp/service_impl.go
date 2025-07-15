package miniapp

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"context"
)

type MiniAppStore struct {
	Repository Repository
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User) (string, error)
	CheckMiniApp(ctx context.Context, actionId string, action bool, checker model.User) error
	UpdateMiniAppAction(ctx context.Context,
		data MiniApp, maker model.User) (string, error)

	DeleteMiniAppAction(ctx context.Context, maker model.User) (string, error)
	ListMiniApp(ctx context.Context) ([]*MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (MiniApp, error)
}

func NewService(repository Repository, logger utils.Logger) MiniAppService {
	return &MiniAppStore{
		Repository: repository,
	}
}
