package productcode

import (
	"context"
	"time"

	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/pkgs/utils"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// productCodeService implements Service
type productCodeService struct {//todo: make the struct 
	repo   storage.ProductCodeRepository
	logger shared_utils.Logger
}

// NewService creates a new productCodeService
func NewproductCodeService(repo storage.ProductCodeRepository, logger shared_utils.Logger) service.ProductCodeService {
	return &productCodeService{
		repo:   repo,
		logger: logger,
	}
}


func (s *productCodeService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Authorization requested for action: %s", cpsAction.RequestAction)
	cpsAction.ActionStatus = "APPROVED"
	return cpsAction, nil
}

// FetchByID fetches a product code by ID
func (s *productCodeService) FetchProductCodeByID(ctx context.Context, id string) (*model.ProductCode, error) {
	productCode, err := s.repo.FetchByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[ProductCode.FetchByID] failed to fetch product code, id: %s, error: %v", id, err)
		return nil, err
	}
	return productCode, nil
}

// FetchAll fetches all product codes with pagination and filtering
func (s *productCodeService) FetchAllProductCodes(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error) {
	response, err := s.repo.FetchAll(ctx, filterParams)
	if err != nil {
		s.logger.Errorf("[ProductCode.FetchAll] failed to fetch product codes, error: %v", err)
		return nil, err
	}
	return response, nil
}

// Update updates a product code
func (s *productCodeService) UpdateProductCode(ctx context.Context, request productcode.UpdateProductCodeRequest) (*model.ProductCode, *model.ProductCode, error) {
	existing, err := s.repo.FetchByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("[ProductCode.Update] failed to fetch existing product code, id: %s, error: %v", request.ID, err)
		return nil, nil, err
	}

	updated := &model.ProductCode{
		ID:                 request.ID,
		ProductName:        utils.NonEmptyString(request.ProductName, existing.ProductName),
		CBEProductCodes:    model.NonEmptyProductCodes(request.CBEProductCodes, existing.CBEProductCodes),
		CBEIFBProductCodes: model.NonEmptyProductCodes(request.CBEIFBProductCodes, existing.CBEIFBProductCodes),
		CreatedAt:          existing.CreatedAt,
		LastUpdatedAt:      time.Now(),
	}

	return updated, existing, nil
}

// Helper functions


// TODO: Add CPS action handling when the service is available
// func (s *productCodeService) handleCPSAction(ctx context.Context, maker cps_entities.User, requestAction constant.RequestAction, curData, prevData interface{}, actionType constant.ActionType) error {
//     // This method needs to be implemented when CPS action service is available
//     // It should create CPS actions for approval workflow
//     return fmt.Errorf("CPS action handling not yet implemented")
// }
