package productcode

import (
	"context"
	"fmt"
	"time"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	constan "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

)

// Service defines the interface for product code business logic
type Service interface {
	FetchByID(ctx context.Context, id string) (*ProductCode, error)
	FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*ProductCode], error)
	Update(ctx context.Context, request UpdateProductCodeRequest) (*ProductCode, *ProductCode, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

// ServiceImpl implements Service
type ServiceImpl struct {
	repo   Repository
	logger shared_utils.Logger
}

// NewService creates a new ServiceImpl
func NewService(repo Repository, logger shared_utils.Logger) Service {
	return &ServiceImpl{repo: repo, logger: logger}
}

// FetchByID fetches a product code by ID
func (s *ServiceImpl) FetchByID(ctx context.Context, id string) (*ProductCode, error) {
	productCode, err := s.repo.FetchByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return productCode, nil
}

// FetchAll fetches all product codes with pagination and filtering
func (s *ServiceImpl) FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*ProductCode], error) {
	response, err := s.repo.FetchAll(ctx, filterParams)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Update updates a product code
func (s *ServiceImpl) Update(ctx context.Context, request UpdateProductCodeRequest) (*ProductCode, *ProductCode, error) {
	existing, err := s.repo.FetchByID(ctx, request.ID)
	if err != nil {
		return nil, nil, err
	}
	updated := &ProductCode{
		ID:                 request.ID,
		ProductName:        nonEmptyString(request.ProductName, existing.ProductName),
		CBEProductCodes:    nonEmptyProductCodes(request.CBEProductCodes, existing.CBEProductCodes),
		CBEIFBProductCodes: nonEmptyProductCodes(request.CBEIFBProductCodes, existing.CBEIFBProductCodes),
		CreatedAt:          existing.CreatedAt,
		LastUpdatedAt:      time.Now(),
	}

	return updated, existing, nil
}

func (s *ServiceImpl) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction

	var productCode *ProductCode
	var err error

	bindErr := common_util.BindAction(action.CurrentAction, &productCode)
	if bindErr != nil {
		s.logger.Errorf("Failed to bind current action to productCode: %v", bindErr)
		return nil, fmt.Errorf(common_util.InvalidActionData)
	}

	switch requestedAction {
	case cps_const.RequestUpdateProductCode:
		productCode, err = s.repo.Update(ctx, productCode)
		if err != nil {
			s.logger.Errorf("Failed to create productCode in repository", "error", err)
			return nil, err
		}

	default:
		s.logger.Errorf("Unsupported action requested", "action", requestedAction)
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	action.CurrentAction = productCode
	s.logger.Infof("Authorization completed for action", "action", requestedAction, "id", productCode.ID)
	return action, nil
}
