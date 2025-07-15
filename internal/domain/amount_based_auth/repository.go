package amount_based_auth_domain

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Repository interface {
	UpdateAmountBasedAuth(ctx context.Context, request UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error)
	ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error)
	GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*AmountBasedAuthRespose, error)
}
