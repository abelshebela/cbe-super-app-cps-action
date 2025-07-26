package customer

import (
	"context"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type CustomerRepository interface {
	GetCustomersDetail(ctx context.Context, kyc_level int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error)
	GetCustomerByID(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	CheckUserExist(ctx context.Context, user *entity.User) (bool, error)
}
