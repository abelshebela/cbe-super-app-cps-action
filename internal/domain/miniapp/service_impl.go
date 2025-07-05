package miniapp

import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	
	"context")

type MiniAppStore struct {
	Repository Repository
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp MiniApp, maker model.User) (string, error)
	CheckMiniApp(ctx context.Context, actionId string, action bool, checker model.User) error
}

func NewService(repository Repository) MiniAppService {
	return &MiniAppStore{
		Repository: repository,
	}
}
