package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	bankvault "cbe-super-app-cps-action/internal/constants/dto/bankvault"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helperr "cbe-super-app-cps-action/internal/service/bankvault/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type bankVaultService struct {
	repo       storage.BankVaultRepository
	cpsService service.CPSActionService
	logger     shared_utils.Logger
}

func NewBankVaultService(re storage.BankVaultRepository, cpsS service.CPSActionService, logger shared_utils.Logger) *bankVaultService {
	return &bankVaultService{
		repo:       re,
		cpsService: cpsS,
		logger:     logger,
	}
}

func (s *bankVaultService) CreateBankVault(ctx context.Context, req *model.BankVaultProduct) (string, error) {
	err := s.repo.FindBankVaultByName(ctx, req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			makerData := local_util.ExtractUserFromContext(ctx)
			current := helperr.ConvertBankVaultToMongoSafe(req)
			cpsActionModel := lib.CpsModelBuilder("", makerData, nil, current, string(constants.RequestCreateBankVault), string(constants.CREATE))
			if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
				if s.logger != nil {
					s.logger.Errorf("failed to create CPS action for bank vault | action=%s | err=%v", constants.RequestCreateBankVault, err)
				}
				return "", err
			}
			return "", nil
		}
	}
	return "", errors.New(localization.ErrorDuplicateBankProduct.Code)
}

func (s *bankVaultService) FindAllBankVaults(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*bankvault.BankVaultProductResponse], error) {
	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("failed to fetch bank vault products: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	resp := make([]*bankvault.BankVaultProductResponse, 0, len(entities.Data))
	for _, e := range entities.Data {
		resp = append(resp, helperr.MapBankVaultToResponse(e))
	}
	return &types.PaginatedResponse[[]*bankvault.BankVaultProductResponse]{
		Data: resp,
		Meta: entities.Meta,
	}, nil

}
func (s *bankVaultService) GetBankVault(ctx context.Context, id string) (*bankvault.BankVaultProductResponse, error) {
	enitity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return helperr.MapBankVaultToResponse(enitity), nil
}

func (s *bankVaultService) UpdateBankVault(ctx context.Context, id string, req *model.UpdateBankVault) (string, error) {
	req.UpdatedAt = time.Now().UTC()
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	current := helperr.BuildUpdateBankVault(prev, req)
	makerData := local_util.ExtractUserFromContext(ctx)

	// Convert to MongoDB-safe format
	mongoSafePrev := helperr.ConvertBankVaultToMongoSafe(prev)
	mongoSafeCurrent := helperr.ConvertBankVaultToMongoSafe(&current)

	cpsActionModel := lib.CpsModelBuilder(id, makerData, mongoSafePrev, mongoSafeCurrent, string(constants.RequestUpdateBankVault), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for bank vault | action=%s | err=%v", constants.RequestUpdateBankVault, err)
		}
		return "", err
	}
	return id, nil
}
func (s *bankVaultService) DeleteBankVault(ctx context.Context, id string) (string, error) {
	s.logger.Infof("Deleting bank vault request: %s", id)
	if id == "" {
		s.logger.Errorf("Empty id")
		return "", errors.New(localization.ErrorUnexpectedError.Message)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Error fetching bank vault product: %v", err)
		return "", nil
	}

	if exist.IsActive {
		s.logger.Errorf("Cannot delete active bank vault product: %s", id)
		return "", errors.New(localization.ErrorCannotDeletedBankProduct.Code)
	}

	if exist.IsDeleted {
		s.logger.Errorf("Bank vault product already deleted: %s", id)
		return "", errors.New(localization.ErrorBankVaultProductAlreadyDeleted.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for bank vault action | context = %v", maker)
		}
		return "", fmt.Errorf("INCOMPLETE_USER_INFO")
	}
	// Convert to MongoDB-safe format
	mongoSafeExist := helperr.ConvertBankVaultToMongoSafe(exist)
	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafeExist, nil, string(constants.RequestDeleteBankVault), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for bank vault | action=%s | err=%v", constants.RequestDeleteBankVault, err)
		}
		return "", err
	}
	return id, nil
}

func (s *bankVaultService) EnableBankVault(ctx context.Context, id string) error {
	s.logger.Infof("Enabling bank vault: %s", id)
	if id == "" {
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}
	if prev.IsActive {
		s.logger.Infof("Bank vault already enabled: %s", id)
		return errors.New(localization.ErrorBankVaultAlreadyEnabled.Code)
	}
	updated := prev
	updated.IsActive = true
	maker := local_util.ExtractUserFromContext(ctx)

	// Convert to MongoDB-safe format
	mongoSafePrev := helperr.ConvertBankVaultToMongoSafe(prev)
	mongoSafeUpdated := helperr.ConvertBankVaultToMongoSafe(updated)
	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafePrev, mongoSafeUpdated, string(constants.RequestEnableBankVault), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *bankVaultService) DisableBankVault(ctx context.Context, id string) error {
	s.logger.Infof("Disabling bank vault: %s", id)
	if id == "" {
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}
	if !prev.IsActive {
		return errors.New(localization.ErrorBankVaultAlreadyDisabled.Code)
	}
	updated := prev
	updated.IsActive = false
	maker := local_util.ExtractUserFromContext(ctx)

	// Convert to MongoDB-safe format
	mongoSafePrev := helperr.ConvertBankVaultToMongoSafe(prev)
	mongoSafeUpdated := helperr.ConvertBankVaultToMongoSafe(updated)
	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafePrev, mongoSafeUpdated, string(constants.RequestDisAbleBankVault), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}
func (s *bankVaultService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateBankVault):
		bankvault, err := helperr.BindBankVaultFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if _, err := s.repo.Create(ctx, bankvault); err != nil {
			return nil, err
		}
		return cpsAction, nil

	case string(constants.RequestUpdateBankVault):
		bankvault, err := helperr.BindBankVaultUpdateFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Update(ctx, cpsAction.UniqueId, bankvault); err != nil {
			return nil, err
		}
		return cpsAction, nil

	case string(constants.RequestDeleteBankVault):
		if _, err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			return nil, err
		}
		return cpsAction, nil
	case string(constants.RequestEnableBankVault):
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			return nil, err
		}
		return cpsAction, nil
	case string(constants.RequestDisAbleBankVault):
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			return nil, err
		}
		return cpsAction, nil
	}

	return nil, errors.New(localization.ErrorInvalidRequest.Code)
}
