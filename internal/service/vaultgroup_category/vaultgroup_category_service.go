package vaultgroupcategory

import (
	"cbe-super-app-cps-action/internal/constants"
	vaultgroup_category "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helperr "cbe-super-app-cps-action/internal/service/vaultgroup_category/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"time"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type vaultgroupCategoryService struct {
	repo       storage.VaultGroupCategoryRepository
	cpsService service.CPSActionService
	logger     shared_utils.Logger
}

func NewVaultGroupCategoryService(re storage.VaultGroupCategoryRepository, cpsS service.CPSActionService, logger shared_utils.Logger) *vaultgroupCategoryService {
	return &vaultgroupCategoryService{
		repo:       re,
		cpsService: cpsS,
		logger:     logger,
	}
}

func (s *vaultgroupCategoryService) CreateVaultGroupCategory(ctx context.Context, req *model.VaultGroupCategory) (string, error) {
	makerData := local_util.ExtractUserFromContext(ctx)
	mongoSafeReq := helperr.ConvertVaultGroupCategoryToMongoSafe(req)

	cpsActionModel := lib.CpsModelBuilder("", makerData, mongoSafeReq, mongoSafeReq, string(constants.RequestCreateVaultGroupCategory), string(constants.CREATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {

		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault group category | action=%s | err=%v", constants.RequestCreateVaultGroupCategory, err)
		}
		return "", errors.New(localization.ErrorCPSActionFailed.Code)
	}
	return req.ID, nil
}

func (s *vaultgroupCategoryService) FindAllVaultGroupCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*vaultgroup_category.VaultGroupCategoryResponse], error) {
	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("failed to fetch vault group categories | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	resp := make([]*vaultgroup_category.VaultGroupCategoryResponse, 0, len(entities.Data))
	for _, en := range entities.Data {
		resp = append(resp, helperr.MapVaultGroupCategoryToResponse(en))
	}
	return &types.PaginatedResponse[[]*vaultgroup_category.VaultGroupCategoryResponse]{
		Data: resp,
		Meta: entities.Meta,
	}, nil
}

func (s *vaultgroupCategoryService) GetVaultGroupCategory(ctx context.Context, id string) (*vaultgroup_category.VaultGroupCategoryResponse, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return helperr.MapVaultGroupCategoryToResponse(entity), nil
}

func (s *vaultgroupCategoryService) UpdateVaultGroupCategory(ctx context.Context, id string, req *model.VaultGroupCategory) (string, error) {
	req.UpdatedAt = time.Now().UTC()
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", localization.ErrorVaultGroupCategoryNotFound
		}
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	current := helperr.BuildUpdateVaultGroupCategory(prev, req)
	makerData := local_util.ExtractUserFromContext(ctx)

	mongosafePrev := helperr.ConvertVaultGroupCategoryToMongoSafe(prev)
	mongosafeCurrent := helperr.ConvertVaultGroupCategoryToMongoSafe(current)

	cpsActionModel := lib.CpsModelBuilder(id, makerData, mongosafePrev, mongosafeCurrent, string(constants.RequestUpdateVaultGroupCategory), string(constants.UPDATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {

		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault group category | action=%s | err=%v", constants.RequestUpdateVaultGroupCategory, err)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	return id, nil
}

func (s *vaultgroupCategoryService) DeleteVaultGroupCategory(ctx context.Context, id string) (string, error) {
	s.logger.Infof("Deleting vault group category with ID: %s", id)
	if id == "" {
		s.logger.Errorf("vault group category id is required")
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", localization.ErrorVaultGroupCategoryNotFound
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	if exist.IsDeleted {
		s.logger.Errorf("vault group category already deleted with id: %s", id)
		return "", localization.ErrorVaultGroupCategoryNotFound
	}
	if exist.IsActive {
		s.logger.Errorf("cannot delete active vault group category with id: %s", id)
		return "", errors.New(localization.ErrorCannotDeleteActiveVaultGroupCategory.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault group category action | context = %v", maker)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	mongosafeExist := helperr.ConvertVaultGroupCategoryToMongoSafe(exist)

	cpsActionModel := lib.CpsModelBuilder(id, maker, mongosafeExist, mongosafeExist, string(constants.RequestDeleteVaultGroupCategory), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault group category | action=%s | err=%v", constants.RequestDeleteVaultGroupCategory, err)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	return id, nil
}
func (s *vaultgroupCategoryService) EnableVaultGroupCategory(ctx context.Context, id string) error {
	s.logger.Infof("Enabling vault group category with ID: %s", id)
	if id == "" {
		s.logger.Errorf("vault group category id is required")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return localization.ErrorVaultGroupCategoryNotFound
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if exist.IsDeleted {
		s.logger.Errorf("vault group category already deleted with id: %s", id)
		return localization.ErrorVaultGroupCategooryAlreadyDeleted
	}
	if exist.IsActive {
		s.logger.Errorf("vault group category already enabled with id: %s", id)
		return errors.New(localization.ErrorVaultGroupAlreadyEnabled.Code)
	}
	updated := exist
	updated.IsActive = true
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault group category action | context = %v", maker)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	// Convert to MongoDB-safe format
	mongoSafeExist := helperr.ConvertVaultGroupCategoryToMongoSafe(exist)
	mongoSafeUpdated := helperr.ConvertVaultGroupCategoryToMongoSafe(updated)
	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafeExist, mongoSafeUpdated, string(constants.RequestEnableVaultGroupCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultgroupCategoryService) DisableVaultGroupCategory(ctx context.Context, id string) error {
	s.logger.Infof("Disabling vault group category with ID: %s", id)
	if id == "" {
		s.logger.Errorf("vault group category id is required")
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return localization.ErrorVaultGroupCategoryNotFound
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if !exist.IsActive {
		s.logger.Errorf("vault group category already disabled with id: %s", id)
		return errors.New(localization.ErrorVaultGroupAlreadyDisabled.Code)
	}
	updated := exist
	updated.IsActive = false
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault group category action | context = %v", maker)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	// Convert to MongoDB-safe format
	mongoSafeExist := helperr.ConvertVaultGroupCategoryToMongoSafe(exist)
	mongoSafeUpdated := helperr.ConvertVaultGroupCategoryToMongoSafe(updated)
	cpsActionModel := lib.CpsModelBuilder(id, maker, mongoSafeExist, mongoSafeUpdated, string(constants.RequestDisAbleVaultGroupCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultgroupCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	cpsAction.MakerActionTime = time.Now()
	cpsAction.LastModifiedAt = cpsAction.MakerActionTime

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateVaultGroupCategory):
		create, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if _, err := s.repo.Create(ctx, &create); err != nil {
			return nil, err
		}
		return cpsAction, nil

	case string(constants.RequestUpdateVaultGroupCategory):
		vaultgroup, err := helperr.BindVaultGroupCategoryUpdateFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		update := helperr.VaultGroupCategoryUpdate(&vaultgroup)
		if err := s.repo.Update(ctx, cpsAction.UniqueId, &update); err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestDeleteVaultGroupCategory):
		_, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestEnableVaultGroupCategory):
		_, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestDisAbleVaultGroupCategory):
		_, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	}
	return nil, errors.New(localization.ErrorInvalidRequest.Code)

}
