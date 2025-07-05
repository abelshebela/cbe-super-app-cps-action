package amount_based_auth_domain

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type Repository interface {
	UpdateAuthTier(ctx context.Context, request UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CPSAction, error)
	ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CPSAction, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectCPSAction) (*model.CPSAction, error)
}
