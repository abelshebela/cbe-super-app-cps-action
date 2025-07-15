package amountbasedauth

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AmountBasedAuthRepo interface {
	UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth.UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error)
	ApproveAmountBasedAuth(ctx context.Context, id string, cpsAction model.AuthorizeCPSAction) (*model.CpsActionNormalized, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error)
	GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*amount_based_auth.AmountBasedAuthRespose, error)
}
