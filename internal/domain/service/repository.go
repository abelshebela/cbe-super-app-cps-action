package service

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type Repository interface {
	GetAllService() ([]*ServiceResponse, error)
	GetOneService(id string) (Service, error)
	UpdateOneService(id string, update any) error
	CreateAction(ctx context.Context, req action.CPSAction) error
}
