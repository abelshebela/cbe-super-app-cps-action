package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/customer/entity"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerRepository interface {
	GetCustomersDetail(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
}
