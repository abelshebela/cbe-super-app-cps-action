package miniapp

import "context"

type MiniAppStore struct {
	repository Repository
}
type MiniAppService interface {
	CreateMiniAppAction(ctx context.Context, miniApp MiniApp, makerId string) (string, error)
	CheckMiniApp(ctx context.Context, actionId string, action bool, checkerId string) error
}

func NewService(repository Repository) MiniAppService {
	return &MiniAppStore{
		repository: repository,
	}
}
