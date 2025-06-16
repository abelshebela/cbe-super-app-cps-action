package action

import "context"

type ServiceImpl interface {
	GetServicePaginated(ctx context.Context, limit, offset int) ([]Service, error)
	UpdateServiceFlagRequest(ctx context.Context, id string, action bool, maker_id string) (string, error)
	UpdateServiceFlag(ctx context.Context, action_id string, action bool, checker_id string) error
<<<<<<< HEAD

	GetAccountByAccount(ctx context.Context, Account string) ([]LinkedAccount, error)
	RemoveCifRequest(ctx context.Context, id []string, action bool, maker_id string) (string, error)
	RemoveCif(ctx context.Context, action_id string, action bool, checker_id string) error
=======
>>>>>>> b69ee66 (feature added)
}

type ServiceStore struct {
	repository Repository
}

func NewService(repo Repository) ServiceImpl {
	return &ServiceStore{
		repository: repo,
	}
}
