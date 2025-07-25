package customer

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/service"

	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/customer/entity"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	GetCustomersDeatil(ctx context.Context, kyc_level int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error)
	GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error)
	GetCustomerByID(ctx context.Context, id string) (*entity.User, error)
}

type CustomerHandler struct {
	domain service.CustomerService
	logger utils.Logger
}

func InitCustomerHandler(customerDomain service.CustomerService, logger utils.Logger) ApplicationService {
	return CustomerHandler{
		domain: customerDomain,
		logger: logger,
	}
}

func (c CustomerHandler) GetCustomersDeatil(ctx context.Context, kyc_level int, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error) {
	customers, err := c.domain.GetCustomersDetail(ctx, kyc_level, filterParams)
	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (c CustomerHandler) GetCustomerByID(ctx context.Context, id string) (*entity.User, error) {
	customer, err := c.domain.GetCustomerByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c CustomerHandler) GetBlockedCustomer(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entity.User], error) {
	customers, err := c.domain.GetBlockedCustomer(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return customers, nil
}
