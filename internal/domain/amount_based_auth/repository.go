package amount_based_auth_domain

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository interface {
	UpdateAuthTier(request AmountBasedAuthRequest, ctx context.Context, logger utils.Logger) (AuthTier, error)
	ApproveAmountBasedAuth(ctx context.Context, id string) (string, error)
}
