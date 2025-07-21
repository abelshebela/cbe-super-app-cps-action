package amountbasedauth

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	amount_based_auth "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/amount_based_auth"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type AmountBasedAuthRepo interface {
	UpdateAmountBasedAuth(ctx context.Context, request amount_based_auth.UpdateAmountBasedAuth, cpsAction model.CreateCPSAction) (*model.CpsActionNormalized, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	RejectAmountBasedAuth(ctx context.Context, id string, cpsAction model.RejectAuthTierCPSAction) (*model.CpsActionNormalized, error)
	GetAllAmountBasedDetail(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*amount_based_auth.AuthTier], error)
}
