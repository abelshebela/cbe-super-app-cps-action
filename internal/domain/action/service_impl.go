package action

import "context"

type ServiceImpl interface {
	GetServicePaginated(ctx context.Context, limit, offset int) ([]Service, error)
	UpdateServiceFlagRequest(ctx context.Context, id string, action bool, maker_id string) (string, error)
	UpdateServiceFlag(ctx context.Context, action_id string, action bool, checker_id string) error
}

type ServiceStore struct {
	repository Repository
}

func NewService(repo Repository) ServiceImpl {
	return &ServiceStore{
		repository: repo,
	}
}
