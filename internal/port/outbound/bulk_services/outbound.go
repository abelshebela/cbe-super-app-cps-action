package outbound

import (
	"context"

	domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundInfra interface {
	GetAllHqServices(ctx context.Context) ([]domain.Service, error)
	GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]domain.Service, error)
	GetHqServiceById(ctx context.Context, id string) (domain.Service, error)
	UpdateHqService(ctx context.Context, service domain.Service) error

	CreateCpsAction(ctx context.Context, Action domain.CPSAction) (domain.CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action domain.CPSAction) error
	FetchCpsActionById(ctx context.Context, Action_Id string) (domain.CPSAction, error)
	FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (domain.CPSAction, error)

	FetchAccountsByCif(ctx context.Context, cif string) ([]domain.LinkedAccount, error)
	UpdateAccounts(ctx context.Context, linkedAccounts []domain.LinkedAccount) ([]domain.LinkedAccount, error)
	UpdateAccount(ctx context.Context, linkedAccount domain.LinkedAccount) (domain.LinkedAccount, error)
}
