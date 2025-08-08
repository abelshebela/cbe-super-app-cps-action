package service

import (
	"context"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/repository"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type CustomerDomain struct {
	customerRepo domain.CustomerRepository
	logger       utils.Logger
}

type CustomerService interface {
	GetCustomersDetail(ctx context.Context, kycLevel int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error)
	GetCustomerByID(ctx context.Context, id string) (*entity.User, error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error)
}

func IntiCustomerDomain(customerRepo domain.CustomerRepository, logger utils.Logger) CustomerService {
	return &CustomerDomain{
		customerRepo: customerRepo,
		logger:       logger,
	}
}

func (c *CustomerDomain) GetCustomersDetail(ctx context.Context, kyc_level int, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error) {
	customers, err := c.customerRepo.GetCustomersDetail(ctx, kyc_level, filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *CustomerDomain) GetCustomerByID(ctx context.Context, id string) (*entity.User, error) {
	customer, err := c.customerRepo.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *CustomerDomain) GetBlockedCustomer(ctx context.Context, filerParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error) {
	customers, err := c.customerRepo.GetBlockedCustomer(ctx, filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
