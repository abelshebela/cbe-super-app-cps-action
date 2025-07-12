package action

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ServiceImpl interface {
	GetServicePaginated(ctx context.Context, limit, offset int) ([]ServiceDetails, error)
	UpdateServiceFlagRequest(ctx context.Context, id string, action bool, maker_id string) (string, error)
	UpdateServiceFlag(ctx context.Context, actionID string, action bool, checker_id string) error
	GetAccountByAccount(ctx context.Context, Account string) ([]LinkedAccount, error)
	RemoveCifRequest(ctx context.Context, id []string, action bool, maker_id string) (string, error)
	RemoveCif(ctx context.Context, actionID string, action bool, checker_id string) error
	CreateCpsAction(ctx context.Context, Action CPSAction) (CPSAction, error)
}

type ServiceStore struct {
	Repository ActionRepository
	Logger     utils.Logger
}

func NewService(repo ActionRepository, logger utils.Logger) ServiceImpl {
	return &ServiceStore{
		Repository: repo,
		Logger:     logger,
	}
}
