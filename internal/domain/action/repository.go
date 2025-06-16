package action

import "context"

type Repository interface {
	//GetMemberById(ctx context.Context, id string) (Member, error)
	//UpdateMember(ctx context.Context, member Member) (Member, error)

	GetAllHqServices(ctx context.Context) ([]Service, error)
	GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]Service, error)
	GetHqServiceById(ctx context.Context, id string) (Service, error)
	UpdateHqService(ctx context.Context, service Service) error

	CreateCpsAction(ctx context.Context, Action CPSAction) (CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action CPSAction) error
	FetchCpsActionById(ctx context.Context, Action_Id string) (CPSAction, error)
	FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (CPSAction, error)

	FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]LinkedAccount, error)
	FetchLinkedAccountById(ctx context.Context, id []string) ([]LinkedAccount, error)
	UpdateAccounts(ctx context.Context, linkedAccounts []LinkedAccount) ([]LinkedAccount, error)
	UpdateAccount(ctx context.Context, linkedAccount LinkedAccount) (LinkedAccount, error)
}
