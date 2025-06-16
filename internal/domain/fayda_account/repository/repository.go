package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
)

type Repository interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
	RejectFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
}
