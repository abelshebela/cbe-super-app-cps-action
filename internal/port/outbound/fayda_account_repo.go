package outbound

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/fayda_account/entity"
)

type FaydaAccountRepository interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
	RejectFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
}
