package outbound

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type FaydaAccountRepository interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error)
	AuthorizeFaydaAccountDisable(ctx context.Context, req *entities.CPSAction) (*entities.CPSAction, error)
	GetAllFaydaAccounts(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
}
