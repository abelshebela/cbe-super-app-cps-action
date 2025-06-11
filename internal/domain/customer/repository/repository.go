package repository

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/customer/entity"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/utils"
)

type CustomerRepository interface {
	GetCustomersDetail(ctx context.Context,filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetCustomerByID(ctx context.Context, id string) (*member.User,error)
}
