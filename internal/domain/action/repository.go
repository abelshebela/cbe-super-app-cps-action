package action

import "context"

type ActionRepository interface {
	GetAllHqServices(ctx context.Context) ([]ServiceDetails, error)
	GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]ServiceDetails, error)
	GetHqServiceById(ctx context.Context, id string) (ServiceDetails, error)
	UpdateHqService(ctx context.Context, service ServiceDetails) error

	CreateCpsAction(ctx context.Context, Action CPSAction) (CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action CPSAction) error
	FetchCpsActionById(ctx context.Context, ActionID string) (CPSAction, error)
	FetchLastCpsActionByMakerID(ctx context.Context, makerID string) (CPSAction, error)

	FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]LinkedAccount, error)
	FetchLinkedAccountById(ctx context.Context, id []string) ([]LinkedAccount, error)
	UpdateAccounts(ctx context.Context, linkedAccounts []LinkedAccount) ([]LinkedAccount, error)
	UpdateAccount(ctx context.Context, linkedAccount LinkedAccount) (LinkedAccount, error)
}
