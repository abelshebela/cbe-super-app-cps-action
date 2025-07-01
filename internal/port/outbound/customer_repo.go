package outbound

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/customer/entity"

	constant "cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerDetailRepository interface {
	GetCustomersDetail(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
}
