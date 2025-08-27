package customer

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerService struct {
	repo   storage.CustomerRepository
	logger utils.Logger
}

func NewCustomerService(repo storage.CustomerRepository, logger utils.Logger) service.CustomerService {
	return &customerService{
		repo:   repo,
		logger: logger,
	}
}

func (c *customerService) GetCustomersDetail(ctx context.Context, kyc_level int, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error) {
	filter := &types.Filter{
		Page:    1,
		PerPage: 10,
		Search:  "searchTerm",
		Filters: map[string]interface{}{
			"kyc_level": kyc_level,
		},
	}
	customers, err := c.repo.FindAllWithPagination(ctx, *filter)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (s *customerService) GetCustomerByID(ctx context.Context, id string) (*model.User, error) {

	customer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *customerService) GetBlockedCustomer(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error) {
	customers, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
