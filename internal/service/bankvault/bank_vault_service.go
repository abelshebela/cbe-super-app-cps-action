package bankvault

import (
	"cbe-super-app-cps-action/internal/constants"
	bankvault "cbe-super-app-cps-action/internal/constants/dto/bankvault"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helper "cbe-super-app-cps-action/internal/service/bankvault/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBankVault", "Bank Vault", "CreateBankVault")
	defer span.End()

	err := s.repo.FindBankVaultByName(ctx, req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			makerData := local_util.ExtractUserFromContext(ctx)
			current := helper.ConvertBankVaultToMongoSafe(req)
			cpsActionModel := lib.CpsModelBuilder("", makerData, nil, current, string(constants.RequestCreateBankVault), string(constants.CREATE))
			if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
				span.AddEvent("failed to create CPS action", trace.WithAttributes(
					attribute.String("error", err.Error()),
					attribute.String("name", req.Name),
				))
				if s.logger != nil {
					s.logger.Errorf("failed to create CPS action for bank vault | action=%s | err=%v", constants.RequestCreateBankVault, err)
				}
				return "", err
			}
			return "", nil
		}
	}
	span.AddEvent("duplicate bank product", trace.WithAttributes(attribute.String("name", req.Name)))
	return "", errors.New(localization.ErrorDuplicateBankProduct.Code)
}

func (s *bankVaultService) FindAllBankVaults(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*bankvault.BankVaultProductResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllBankVaults", "Bank Vault", "FindAllBankVaults")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[FindAllBankVaults] failed to fetch bank vault products", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[FindAllBankVaults] failed to fetch bank vault products: %v", err)
		return nil, err
	}
	resp := make([]*bankvault.BankVaultProductResponse, 0, len(entities.Data))
	for _, e := range entities.Data {
		resp = append(resp, helper.MapBankVaultToResponse(e))
	}
	s.logger.Infof("[FindAllBankVaults] retrieved %d bank vault products", len(resp))
	return &types.PaginatedResponse[[]*bankvault.BankVaultProductResponse]{
		Data: resp,
		Meta: entities.Meta,
	}, nil

}
func (s *bankVaultService) GetBankVault(ctx context.Context, id string) (*bankvault.BankVaultProductResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetBankVault", "Bank Vault", "GetBankVault")
	defer span.End()

	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[GetBankVault] failed to fetch bank vault", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[GetBankVault] failed to fetch bank vault: %v", err)
		return nil, err
	}
	s.logger.Infof("[GetBankVault] bank vault retrieved successfully for id: %s", id)
	return helper.MapBankVaultToResponse(entity), nil
}

func (s *bankVaultService) UpdateBankVault(ctx context.Context, id string, req *model.UpdateBankVault) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBankVault", "Bank Vault", "UpdateBankVault")
	defer span.End()

	s.logger.Infof("[UpdateBankVault] updating bank vault for id: %s", id)
	req.UpdatedAt = time.Now().UTC()
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[UpdateBankVault] failed to find bank vault", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[UpdateBankVault] failed to find bank vault: %v", err)
		return "", err
	}
	current := helper.BuildUpdateBankVault(prev, req)
	makerData := local_util.ExtractUserFromContext(ctx)

	// Convert to MongoDB-safe format
	mongoSafePrev := helper.ConvertBankVaultToMongoSafe(prev)
	mongoSafeCurrent := helper.ConvertBankVaultToMongoSafe(&current)

	cpsActionModel := lib.CpsModelBuilder(id, makerData, mongoSafePrev, mongoSafeCurrent, string(constants.RequestUpdateBankVault), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[UpdateBankVault] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[UpdateBankVault] failed to create CPS action: %v", err)
		return "", err
	}
	s.logger.Infof("[UpdateBankVault] bank vault update request created successfully for id: %s", id)
	return id, nil
}
func (s *bankVaultService) DeleteBankVault(ctx context.Context, id string) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteBankVault", "Bank Vault", "DeleteBankVault")
	defer span.End()

	s.logger.Infof("[DeleteBankVault] deleting bank vault for id: %s", id)
	if id == "" {
		span.AddEvent("[DeleteBankVault] empty id provided")
		s.logger.Errorf("[DeleteBankVault] empty id provided")
		return "", errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[DeleteBankVault] failed to find bank vault", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[DeleteBankVault] failed to find bank vault: %v", err)
		return "", err
	}

	if exist.IsActive {
		span.AddEvent("[DeleteBankVault] cannot delete active bank vault product", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[DeleteBankVault] cannot delete active bank vault product")
		return "", errors.New(localization.ErrorCannotDeletedBankProduct.Code)
	}

	if exist.IsDeleted {
		span.AddEvent("[DeleteBankVault] bank vault product already deleted", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[DeleteBankVault] bank vault product already deleted")
		return "", errors.New(localization.ErrorBankVaultProductAlreadyDeleted.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("[DeleteBankVault] incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[DeleteBankVault] incomplete user context")
		return "", fmt.Errorf(localization.ErrorIncompleteUserInfo.Code)
	}
	// Convert to MongoDB-safe format
	mongoSafeExist := helper.ConvertBankVaultToMongoSafe(exist)
	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafeExist, nil, string(constants.RequestDeleteBankVault), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[DeleteBankVault] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[DeleteBankVault] failed to create CPS action: %v", err)
		return "", err
	}
	s.logger.Infof("[DeleteBankVault] bank vault deletion request created successfully for id: %s", id)
	return id, nil
}

func (s *bankVaultService) EnableBankVault(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableBankVault", "Bank Vault", "EnableBankVault")
	defer span.End()

	s.logger.Infof("Enabling bank vault: %s", id)
	if id == "" {
		span.AddEvent("[EnableBankVault] empty id provided")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("[EnableBankVault] bank vault not found", trace.WithAttributes(attribute.String("id", id)))
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("[EnableBankVault] failed to find bank vault", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if prev.IsActive {
		span.AddEvent("[EnableBankVault] bank vault already enabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Infof("Bank vault already enabled: %s", id)
		return errors.New(localization.ErrorBankVaultAlreadyEnabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)

	mongoSafePrev := helper.ConvertBankVaultToMongoSafe(prev)
	updated := *prev
	updated.IsActive = true
	updated.UpdatedAt = time.Now()
	mongoSafeUpdated := helper.ConvertBankVaultToMongoSafe(&updated)

	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafePrev, mongoSafeUpdated, string(constants.RequestEnableBankVault), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *bankVaultService) DisableBankVault(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableBankVault", "Bank Vault", "DisableBankVault")
	defer span.End()

	s.logger.Infof("[DisableBankVault] disabling bank vault for id: %s", id)
	if id == "" {
		span.AddEvent("[DisableBankVault] empty id provided")
		s.logger.Errorf("[DisableBankVault] empty id provided")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("[DisableBankVault] bank vault not found", trace.WithAttributes(attribute.String("id", id)))
			s.logger.Errorf("[DisableBankVault] bank vault not found: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("[DisableBankVault] failed to find bank vault", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[DisableBankVault] failed to find bank vault: %v", err)
		return err
	}
	if !prev.IsActive {
		span.AddEvent("[DisableBankVault] bank vault already disabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[DisableBankVault] bank vault already disabled")
		return errors.New(localization.ErrorBankVaultAlreadyDisabled.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)

	mongoSafePrev := helper.ConvertBankVaultToMongoSafe(prev)
	updated := *prev
	updated.IsActive = false
	updated.UpdatedAt = time.Now()
	mongoSafeUpdated := helper.ConvertBankVaultToMongoSafe(&updated)

	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafePrev, mongoSafeUpdated, string(constants.RequestDisAbleBankVault), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[DisableBankVault] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[DisableBankVault] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[DisableBankVault] bank vault disable request created successfully for id: %s", id)
	return nil
}

func (s *bankVaultService) FindAllBankLockedVaultsWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.LockedVault], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllBankLockedVaultsWithPagination", "Bank Vault", "FindAllBankLockedVaultsWithPagination")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	results, err := s.repo.FindAllBankLockedVaultsWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[FindAllBankLockedVaultsWithPagination] failed to fetch bank locked vaults", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[FindAllBankLockedVaultsWithPagination] failed to fetch bank locked vaults: %v", err)
		return nil, err
	}
	s.logger.Infof("[FindAllBankLockedVaultsWithPagination] retrieved %d bank locked vaults", len(results.Data))
	return &types.PaginatedResponse[[]*model.LockedVault]{
		Data: results.Data,
		Meta: results.Meta,
	}, nil
}

// func (s *bankVaultService) GetBankLockedVault(ctx context.Context, id string) (*bankvault.LockedVaultResponse, error) {
// 	return nil, nil
// }

func (s *bankVaultService) FindAllGroupVaultsWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.GroupVault], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllGroupVaultsWithPagination", "Bank Vault", "FindAllGroupVaultsWithPagination")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f
	}

	results, err := s.repo.FindAllGroupVaultWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[FindAllGroupVaultsWithPagination] failed to fetch group vaults", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[FindAllGroupVaultsWithPagination] failed to fetch group vaults: %v", err)
		return nil, err
	}
	s.logger.Infof("[FindAllGroupVaultsWithPagination] retrieved %d group vaults", len(results.Data))
	return &types.PaginatedResponse[[]*model.GroupVault]{
		Data: results.Data,
		Meta: results.Meta,
	}, nil

}

// func (s *bankVaultService) GetGroupVault(ctx context.Context, id string) (*bankvault.GroupVaultResponse, error) {
// 	return nil, nil
// }

func (s *bankVaultService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Bank Vault", "Authorize")
	defer span.End()

	s.logger.Infof("[Authorize] authorizing bank vault action: %s", cpsAction.RequestAction)
	switch cpsAction.RequestAction {
	case string(constants.RequestCreateBankVault):
		bankvault, err := helper.BindBankVaultFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("[Authorize] failed to bind bank vault from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to bind bank vault from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if _, err := s.repo.Create(ctx, bankvault); err != nil {
			span.AddEvent("[Authorize] failed to create bank vault", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to create bank vault: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] bank vault created successfully")
		return cpsAction, nil

	case string(constants.RequestUpdateBankVault):
		bankvault, err := helper.BindBankVaultUpdateFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("[Authorize] failed to bind bank vault update from action", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to bind bank vault update from action: %v", err)
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Update(ctx, cpsAction.UniqueId, bankvault); err != nil {
			span.AddEvent("[Authorize] failed to update bank vault", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to update bank vault: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] bank vault updated successfully for id: %s", cpsAction.UniqueId)
		return cpsAction, nil

	case string(constants.RequestDeleteBankVault):
		if _, err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			span.AddEvent("[Authorize] failed to delete bank vault", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to delete bank vault: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] bank vault deleted successfully for id: %s", cpsAction.UniqueId)
		return cpsAction, nil
	case string(constants.RequestEnableBankVault):
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			span.AddEvent("[Authorize] failed to enable bank vault", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to enable bank vault: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] bank vault enabled successfully for id: %s", cpsAction.UniqueId)
		return cpsAction, nil
	case string(constants.RequestDisAbleBankVault):
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			span.AddEvent("[Authorize] failed to disable bank vault", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to disable bank vault: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] bank vault disabled successfully for id: %s", cpsAction.UniqueId)
		return cpsAction, nil
	}

	span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
	s.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
	return nil, errors.New(localization.ErrorInvalidRequest.Code)
}
