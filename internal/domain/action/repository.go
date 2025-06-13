package action

import "context"

type Repository interface {
	GetAllHqServices(ctx context.Context) ([]Service, error)
	GetAllHqServicesPaginated(ctx context.Context, offset, limit int) ([]Service, error)
	GetHqServiceById(ctx context.Context, id string) (Service, error)
	UpdateHqService(ctx context.Context, service Service) error

	CreateCpsAction(ctx context.Context, Action CPSAction) (CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action CPSAction) error
	FetchCpsActionById(ctx context.Context, Action_Id string) (CPSAction, error)
}
