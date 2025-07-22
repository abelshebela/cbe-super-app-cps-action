package customer

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerRepository interface {
	GetCustomersDetail(ctx context.Context, kyc_level int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
}
