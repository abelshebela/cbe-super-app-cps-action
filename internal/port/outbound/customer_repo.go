package outbound

import (
	"context"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerDetailRepository interface {
	GetCustomersDetail(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
}
