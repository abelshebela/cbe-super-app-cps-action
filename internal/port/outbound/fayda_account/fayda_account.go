package faydaaccount

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
)

type FaydaRepository interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
	RejectFaydaAccountDisable(ctx context.Context, req entity.CPSAction) (*entity.CPSAction, error)
}
