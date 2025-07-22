package service

import (
	"context"

	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/customer"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type CustomerDomain struct {
	customerService outbound.CustomerRepository
	logger          utils.Logger
}

type CustomerService interface {
	GetCustomersDetail(ctx context.Context, kycLevel int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
}

func IntiCustomerDomain(customerRepo outbound.CustomerRepository, logger utils.Logger) CustomerService {
	return &CustomerDomain{
		customerService: customerRepo,
		logger:          logger,
	}
}

func (c *CustomerDomain) GetCustomersDetail(ctx context.Context, kyc_level int, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	customers, err := c.customerService.GetCustomersDetail(ctx, kyc_level, filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *CustomerDomain) GetCustomerByID(ctx context.Context, id string) (*member.User, error) {
	customer, err := c.customerService.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *CustomerDomain) GetBlockedCustomer(ctx context.Context, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	customers, err := c.customerService.GetBlockedCustomer(ctx, filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
