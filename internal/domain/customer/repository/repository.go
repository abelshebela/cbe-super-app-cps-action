package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
)

type CustomerRepository interface {
	GetCustomersDetail(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetFaydaCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetFaydaCustomerByID(ctx context.Context, id string) (*member.User, error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
}
