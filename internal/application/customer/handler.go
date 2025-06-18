package customer

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/customer/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/customer/service"

	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	GetCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
}

type CustomerHandler struct {
	domain *service.CustomerDomain
	logger utils.Logger
}

func InitCustomerHandler(customerDomain *service.CustomerDomain, logger utils.Logger) ApplicationService {
	return CustomerHandler{
		domain: customerDomain,
		logger: logger,
	}
}

func (c CustomerHandler) GetCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error) {
	customers, err := c.domain.GetCustomersDetail(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c CustomerHandler) GetCustomerByID(ctx context.Context, id string) (*member.User, error) {
	customer, err := c.domain.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
