package productcode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/productcode/core"
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
	s.logger.Infof("[Authorize] authorizing product code action: %s", cpsAction.RequestAction)
	cpsAction.ActionStatus = "APPROVED"
	var new model.ProductCode
	err := core.BindAction(cpsAction.CurrentAction, &new)
	if err != nil {
		s.logger.Errorf("[Authorize] failed to bind current action to product code: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}
	new.ID = cpsAction.UniqueId
	err = s.repo.Update(ctx, &new)
	if err != nil {
		s.logger.Errorf("[Authorize] failed to update product code: %v", err)
		return nil, err
	}
	s.logger.Infof("[Authorize] product code authorized successfully for id: %s", cpsAction.UniqueId)
	return cpsAction, nil
}

func (s *productCodeService) FetchProductCodeByID(ctx context.Context, id string) (*model.ProductCode, error) {
	productCode, err := s.repo.FetchByID(ctx, id)
	if err != nil {
		s.logger.Errorf("[ProductCode.FetchByID] failed to fetch product code, id: %s, error: %v", id, err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorProductCodeNotFound.Code)
		}
		return nil, err
	}
	return productCode, nil
}

func (s *productCodeService) FetchAllProductCodes(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error) {
	response, err := s.repo.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		s.logger.Errorf("[FetchAllProductCodes] failed to fetch product codes: %v", err)
		return nil, errors.New(localization.ErrorProductCodeNotFound.Code)
	}
	s.logger.Infof("[FetchAllProductCodes] retrieved %d product codes", len(response.Data))
	return response, nil
}

func (s *productCodeService) UpdateProductCode(ctx context.Context, request productcode.UpdateProductCodeRequest) (*model.ProductCode, *model.ProductCode, error) {
	makerData := utils.ExtractUserFromContext(ctx)
	existing, err := s.repo.FetchByID(ctx, request.ID)
	if err != nil {
		s.logger.Errorf("[ProductCode.Update] failed to fetch existing product code, id: %s, error: %v", request.ID, err)
		if err == mongo.ErrNoDocuments {
			return nil, nil, errors.New(localization.ErrorProductCodeNotFound.Code)
		}
		return nil, nil, err
	}

	updated := &model.ProductCode{
		ID:                 existing.ID,
		ProductName:        utils.NonEmptyString(request.ProductName, existing.ProductName),
		CBEProductCodes:    model.NonEmptyProductCodes(request.CBEProductCodes, existing.CBEProductCodes),
		CBEIFBProductCodes: model.NonEmptyProductCodes(request.CBEIFBProductCodes, existing.CBEIFBProductCodes),
		CreatedAt:          existing.CreatedAt,
		LastUpdatedAt:      time.Now(),
	}

	// Check for duplicate product name if name is being changed
	if existing.ProductName != updated.ProductName {
		s.logger.Infof("[ProductCode.Update] Product name changed from '%s' to '%s', checking for duplicates", existing.ProductName, updated.ProductName)
		duplicate, err := s.repo.FindByName(ctx, updated.ProductName)
		if err != nil {
			s.logger.Errorf("[ProductCode.Update] failed to check for duplicate product name: %v", err)
			return nil, nil, err
		}
		if duplicate != nil && duplicate.ID != existing.ID {
			s.logger.Warnf("[ProductCode.Update] duplicate product name detected: '%s' already exists with ID: %s", updated.ProductName, duplicate.ID)
			return nil, nil, fmt.Errorf("%s", localization.ErrorDuplicateProductName.Code)
		}
	}

	// Check for duplicate PRD codes if either CBE or CBE IFB PRD has changed
	cbePRDChanged := existing.CBEProductCodes.PRD != updated.CBEProductCodes.PRD
	cbeIFBPRDChanged := existing.CBEIFBProductCodes.PRD != updated.CBEIFBProductCodes.PRD

	if cbePRDChanged || cbeIFBPRDChanged {
		s.logger.Infof("[ProductCode.Update] PRD values changed, checking for duplicates. CBE PRD: %s, CBE IFB PRD: %s",
			updated.CBEProductCodes.PRD, updated.CBEIFBProductCodes.PRD)

		duplicates, err := s.repo.FindByPRD(ctx, updated.CBEProductCodes.PRD, updated.CBEIFBProductCodes.PRD)
		if err != nil {
			s.logger.Errorf("[ProductCode.Update] failed to check for duplicate PRD codes: %v", err)
			return nil, nil, err
		}

		// Check if any duplicate belongs to a different product code
		for _, dup := range duplicates {
			if dup.ID != existing.ID {
				// Determine which PRD field is duplicated
				if updated.CBEProductCodes.PRD != "" && dup.CBEProductCodes.PRD == updated.CBEProductCodes.PRD {
					s.logger.Warnf("[ProductCode.Update] duplicate CBE PRD detected: '%s' already exists in product code ID: %s",
						updated.CBEProductCodes.PRD, dup.ID)
					return nil, nil, fmt.Errorf("%s", localization.ErrorDuplicateCBEProductCode.Code)
				}
				if updated.CBEIFBProductCodes.PRD != "" && dup.CBEIFBProductCodes.PRD == updated.CBEIFBProductCodes.PRD {
					s.logger.Warnf("[ProductCode.Update] duplicate CBE IFB PRD detected: '%s' already exists in product code ID: %s",
						updated.CBEIFBProductCodes.PRD, dup.ID)
					return nil, nil, fmt.Errorf("%s", localization.ErrorDuplicateCBEIFBProductCode.Code)
				}
			}
		}
		s.logger.Infof("[ProductCode.Update] No duplicate PRD codes found")
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
		if err != nil {
			s.logger.Errorf("[UpdateProductCode] failed to create CPS action: %v", err)
			return nil, nil, err
		}
		s.logger.Infof("[UpdateProductCode] product code update request created successfully for id: %s", updated.ID)
	} else {
		s.logger.Warnf("[UpdateProductCode] no changes detected for product code id: %s", updated.ID)
		return nil, nil, fmt.Errorf("%s", localization.ErrorNoChangesDetected.Code)
	}
	return existing, updated, err
}
