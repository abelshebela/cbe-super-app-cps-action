package amount_based_auth_domain

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type Repository interface {
	UpdateAuthTier(ctx context.Context, request UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsAction, error)
	ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsAction, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectCPSAction) (*model.CpsAction, error)
}
