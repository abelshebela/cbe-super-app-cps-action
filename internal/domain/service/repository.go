package service

import "context"

type Repository interface {
	GetAllServiceDetails(ctx context.Context) ([]*Service, error)
	GetOneServiceDetail(ctx context.Context, id string) (Service, error)
	UpdateOneServiceDetail(ctx context.Context, id string, update Service) error

	UpdateOneSeviceDeatil(ctx context.Context, pd Service) error
}
