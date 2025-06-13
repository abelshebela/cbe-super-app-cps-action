package service

import (
	"context"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/customer/entity"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/internal/domain/customer/repository"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-ms/utils"
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

func (c *CustomerDomain) GetCustomersDetail(ctx context.Context,filerParams *constant.Filter) (*entity.CustomerRespose, error) {
	customers, err := c.customerService.GetCustomersDetail(ctx,filerParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *CustomerDomain) GetCustomerByID(ctx context.Context, id string) (*member.User,error) {
	customer,err := c.customerService.GetCustomerByID(ctx,id)
	if err != nil {
		return nil,err
	}

	return customer,nil
}
