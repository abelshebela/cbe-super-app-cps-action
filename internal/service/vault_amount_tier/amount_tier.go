package vaultamounttier

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/pkgs/utils"
	"context"

	vault_amount_dto "cbe-super-app-cps-action/internal/constants/dto/vault_amount_tier"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	helperr "cbe-super-app-cps-action/internal/service/vault_amount_tier/core"
	"encoding/json"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type VaultAmountTierService struct {
	repo       storage.VaultAmountTierRepository
	cpsService service.CPSActionService
	logger     shared_utils.Logger
}

func NewVaultAmountTierService(repo storage.VaultAmountTierRepository, cpsService service.CPSActionService, logger shared_utils.Logger) *VaultAmountTierService {
	return &VaultAmountTierService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *VaultAmountTierService) CreateAmountTier(ctx context.Context, req *vault_amount_dto.VaultAmountTierRequest) (string, error) {
	current := &model.VaultAmountTier{
		VaultCategoryID: req.VaultCategoryID,
		MinAmount:       req.MinAmount,
		MaxAmount:       req.MaxAmount,
		Interest:        req.Interest,
		IsActive:        true,
	}

	makerData := utils.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder("", makerData, nil, current, string(constants.RequestCreateVaultAmountTier), string(constants.CREATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[VaultTierSvc][Create] cps action err: %v", err)
		return "", err
	}

	return "", nil
}

func (s *VaultAmountTierService) FindAllAmountTiers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.VaultAmountTier], error) {
	return s.repo.FindAllWithPagination(ctx, *filterParams)
}

func (s *VaultAmountTierService) GetAmountTier(ctx context.Context, id string) (*model.VaultAmountTier, error) {
	if id == "" {
		return nil, nil
	}
	return s.repo.FindByID(ctx, id)
}

func (s *VaultAmountTierService) UpdateAmountTier(ctx context.Context, id string, req *vault_amount_dto.UpdateVaultAmountTierRequest) (string, error) {
	current := &model.VaultAmountTier{
		ID:        id,
		MinAmount: req.MinAmount,
		MaxAmount: req.MaxAmount,
		Interest:  req.Interest,
	}

	makerData := utils.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(id, makerData, nil, current, string(constants.RequestUpdateVaultAmountTier), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[VaultTierSvc][Update] cps action err: %v", err)
		return "", err
	}

	return "", nil
}

func (s *VaultAmountTierService) DeleteAmountTier(ctx context.Context, id string) (string, error) {
	current := &model.VaultAmountTier{
		IsDeleted: true,
	}

	makerData := utils.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(id, makerData, nil, current, string(constants.RequestDeleteVaultAmountTier), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[VaultTierSvc][Delete] cps action err: %v", err)
		return "", err
	}

	return "", nil
}

func (s *VaultAmountTierService) EnableOrDisableAmountTier(ctx context.Context, id string, enable bool) (string, error) {
	current := &model.VaultAmountTier{
		IsActive: enable,
	}

	makerData := utils.ExtractUserFromContext(ctx)
	var requestType string
	if enable {
		requestType = string(constants.RequestEnableVaultAmountTier)
	} else {
		requestType = string(constants.RequestDisAbleVaultAmountTier)
	}
	cpsAction := lib.CpsModelBuilder(id, makerData, nil, current, requestType, string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[VaultTierSvc][EnableDisable] cps action err: %v", err)
		return "", err
	}

	return "", nil
}

func (s *VaultAmountTierService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := utils.TraceLogger(ctx, "service", "Authorize", "Bank Vault", "Authorize")
	defer span.End()

	s.logger.Infof("[VaultTierSvc][Authorize] action: %s", cpsAction.RequestAction)

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to marshal CurrentAction")
		s.logger.Errorf("[VaultTierSvc][Authorize] marshal err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction")
		s.logger.Errorf("[VaultTierSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	actionData := helperr.AmountTierMapper(actionMap.(map[string]interface{}))

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateVaultAmountTier):
		_, err := s.repo.Create(ctx, &actionData)
		if err != nil {
			span.AddEvent("Failed to create vault amount tier")
			s.logger.Errorf("[VaultTierSvc][Authorize] create err: %v", err)
			return nil, err
		}
		return cpsAction, nil

	case string(constants.RequestUpdateVaultAmountTier):
		span.AddEvent("RequestUpdateVaultAmountTier")
		if err := s.repo.Update(ctx, cpsAction.UniqueId, &actionData); err != nil {
			span.AddEvent("Failed to update vault amount tier")
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestDeleteVaultAmountTier):
		span.AddEvent("RequestDeleteVaultAmountTier")
		_, err := helperr.BindVaultAmountTierFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind vault amount tier from CPSAction")
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		_, err = s.repo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to delete vault amount tier")
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestEnableVaultAmountTier):
		span.AddEvent("RequestEnableVaultAmountTier")
		_, err := helperr.BindVaultAmountTierFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind vault amount tier from CPSAction")
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			span.AddEvent("Failed to enable vault amount tier")
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestDisAbleVaultAmountTier):

		span.AddEvent("RequestDisAbleVaultAmountTier")
		_, err := helperr.BindVaultAmountTierFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind vault amount tier from CPSAction")
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			span.AddEvent("Failed to disable vault amount tier")
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil
	}

	span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
	s.logger.Errorf("[VaultTierSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
	return nil, errors.New(localization.ErrorInvalidRequest.Code)
}
