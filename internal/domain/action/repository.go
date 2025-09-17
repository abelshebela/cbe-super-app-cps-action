package action

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type ActionRepository interface {
	GetAllHqServices(ctx context.Context) ([]ServiceDetails, error)
	GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]ServiceDetails, error)
	GetHqServiceById(ctx context.Context, id string) (ServiceDetails, error)
	UpdateHqService(ctx context.Context, service ServiceDetails) error

	CreateCpsAction(ctx context.Context, Action CPSAction) (CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action CPSAction) error
	FetchCpsActionById(ctx context.Context, actionCode string) (CPSAction, error)
	FetchLastCpsActionByMakerID(ctx context.Context, makerID string) (CPSAction, error)

	FetchAccountsByAccountNumber(ctx context.Context, accountNumber string) ([]LinkedAccount, error)
	FetchLinkedAccountById(ctx context.Context, id []string) ([]LinkedAccount, error)
	UpdateAccounts(ctx context.Context, linkedAccounts []LinkedAccount) ([]LinkedAccount, error)
	UpdateAccount(ctx context.Context, linkedAccount LinkedAccount) (LinkedAccount, error)
}

type IActionRepository interface {
	CreateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	GetCPSActionByID(ctx context.Context, id string) (*entities.CPSAction, error)
	UpdateCPSAction(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}
