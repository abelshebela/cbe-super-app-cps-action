package customer

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	GetCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetFaydaCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error)
	GetFaydaCustomerByID(ctx context.Context, id string) (*member.User, error)
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

func (c CustomerHandler) GetFaydaCustomersDeatil(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error) {
	customers, err := c.domain.GetFaydaCustomersDeatil(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c CustomerHandler) GetFaydaCustomerByID(ctx context.Context, id string) (*member.User, error) {
	customer, err := c.domain.GetFaydaCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c CustomerHandler) GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*entity.CustomerRespose, error) {
	customers, err := c.domain.GetBlockedCustomer(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return customers, nil
}
