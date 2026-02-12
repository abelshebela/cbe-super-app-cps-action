package sqlc

import (
	"context"
)

type Querier interface {
	SaveAmountTier(ctx context.Context, arg VaultAmountTier) (string, error)
	FindAmountTier(ctx context.Context, arg VaultAmountTierParams) ([]VaultAmountTier, error)
	FindAmountTierById(ctx context.Context, id string) (VaultAmountTier, error)
	EnableOrDisableAmountTier(ctx context.Context, id string, isEnabled bool) (string, error)
	DeleteAmountTier(ctx context.Context, id string) (string, error)
	UpdateAmountTier(ctx context.Context, arg VaultAmountTier) (string, error)
}

var _ Querier = (*Queries)(nil)
