package miniapp

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"context"
)

type MiniAppStore struct {
	Repository Repository
	logger     utils.Logger
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User, departmen string) (string, error)
	CheckMiniApp(ctx context.Context, actionId string, action bool, checker model.User, department string) (model.CPSAction, error)
	UpdateMiniAppAction(ctx context.Context, data MiniApp, maker model.User) (string, error)

	DeleteMiniAppAction(ctx context.Context, maker model.User, id string) (string, error)
	ListMiniApp(ctx context.Context) ([]*model.MiniApp, error)
	DetailMiniAppByID(ctx context.Context, id string) (model.MiniApp, error)
}

func NewService(repository Repository, logger utils.Logger) MiniAppService {
	return &MiniAppStore{
		Repository: repository,
		logger:     logger,
	}
}
