package productcode

import (
	"context"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/event/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/pkgs/utils"

	"go.mongodb.org/mongo-driver/mongo"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type productCodeService struct {
	repo       storage.ProductCodeRepository
	cpsService service.CPSActionService
	logger     shared_utils.Logger
}

func NewProductCodeService(repo storage.ProductCodeRepository, cpsService service.CPSActionService, logger shared_utils.Logger) service.ProductCodeService {
	return &productCodeService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *productCodeService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Authorization requested for action: %s", cpsAction.RequestAction)
	cpsAction.ActionStatus = "APPROVED"
	var new model.ProductCode
	err := core.BindAction(cpsAction.CurrentAction, &new)
	if err != nil {
		s.logger.Errorf("failed to bind current action to product code: %v", err)
		return nil, fmt.Errorf("%v", localization.ErrorInvalidRequest.Code)
	}
	return cpsAction, s.repo.Update(ctx, &new)
}

func (s *productCodeService) FetchProductCodeByID(ctx context.Context, id string) (*model.ProductCode, error) {
	productCode, err := s.repo.FetchByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[ProductCode.FetchByID] failed to fetch product code, id: %s, error: %v", id, err)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("%v", localization.ErrorProductCodeNotFound.Code)
		}
		return nil, err
	}
	return productCode, nil
}

func (s *productCodeService) FetchAllProductCodes(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error) {
	response, err := s.repo.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		s.logger.Errorf("[ProductCode.FetchAll] failed to fetch product codes, error: %v", err)
		return nil, fmt.Errorf("%v", localization.ErrorProductCodeNotFound.Code)
	}

	return response, nil
}

func (s *productCodeService) UpdateProductCode(ctx context.Context, request productcode.UpdateProductCodeRequest) (*model.ProductCode, *model.ProductCode, error) {
	makerData := utils.ExtractUserFromContext(ctx)
	existing, err := s.repo.FetchByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("[ProductCode.Update] failed to fetch existing product code, id: %s, error: %v", request.ID, err)
		if err == mongo.ErrNoDocuments {
			return nil, nil, fmt.Errorf("%v", localization.ErrorProductCodeNotFound.Code)
		}
		return nil, nil, err
	}

	fmt.Println("this is the data from the request", request)
	updated := &model.ProductCode{
		ID:                 existing.ID,
		ProductName:        utils.NonEmptyString(request.ProductName, existing.ProductName),
		CBEProductCodes:    model.NonEmptyProductCodes(request.CBEProductCodes, existing.CBEProductCodes),
		CBEIFBProductCodes: model.NonEmptyProductCodes(request.CBEIFBProductCodes, existing.CBEIFBProductCodes),
		CreatedAt:          existing.CreatedAt,
		LastUpdatedAt:      time.Now(),
	}
	thereIsUpdate := true
	if existing.ProductName == updated.ProductName &&
		existing.CBEIFBProductCodes == updated.CBEIFBProductCodes &&
		existing.CBEProductCodes == updated.CBEProductCodes {
		thereIsUpdate = false
	}
	if thereIsUpdate {
		s.logger.Infof("ProductCode update detected, creating CPS action for approval. ProductCode ID: %s", updated.ID)
		cpsActionData := lib.CpsModelBuilder(updated.ID, makerData, *existing, *updated, string(constants.RequestUpdateProductCode), constants.UPDATE)
		err = s.cpsService.CreateCPSAction(ctx, &cpsActionData)
	} else {
		fmt.Printf("No changes detected for ProductCode ID: %s. Update request ignored.\n", updated.ID)
		fmt.Println("prev", existing, "new", updated)
		return nil, nil, fmt.Errorf("%s", localization.ErrorNoChangesDetected.Code)
	}
	return existing, updated, err
}
