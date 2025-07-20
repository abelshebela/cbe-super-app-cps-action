package service

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/repository"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type CustomerDomain struct {
	customerService repository.CustomerRepository
	logger          utils.Logger
}

func IntiCustomerDomain(customerRepo repository.CustomerRepository, logger utils.Logger) *CustomerDomain {
	return &CustomerDomain{
		customerService: customerRepo,
		logger:          logger,
	}
}

func (c *CustomerDomain) GetCustomersDetail(ctx context.Context, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	customers, err := c.customerService.GetCustomersDetail(ctx, filerParams)
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

func (c *CustomerDomain) GetFaydaCustomersDeatil(ctx context.Context, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {

	customers, err := c.customerService.GetFaydaCustomersDeatil(ctx, filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
func (c *CustomerDomain) GetFaydaCustomerByID(ctx context.Context, id string) (*member.User, error) {
	customer, err := c.customerService.GetFaydaCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *CustomerDomain) GetBlockedCustomer(ctx context.Context, filerParams *constant.Filter) (*entity.CustomerRespose, error) {
	customers, err := c.customerService.GetBlockedCustomer(ctx, filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
