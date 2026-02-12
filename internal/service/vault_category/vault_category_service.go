package vaultcategory

import (
	"cbe-super-app-cps-action/internal/constants"
	vault_category_dto "cbe-super-app-cps-action/internal/constants/dto/vault_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helperr "cbe-super-app-cps-action/internal/service/vault_category/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"strings"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type vaultCategoryService struct {
	repo        storage.VaultCategoryRepository
	cpsService  service.CPSActionService
	logger      shared_utils.Logger
	minio       *s3.Client
	minioPubUrl string
	bucketName  string
	cfg         *config.VaultConfig
}

func NewVaultCategoryService(re storage.VaultCategoryRepository, cpsS service.CPSActionService, logger shared_utils.Logger, minio *s3.Client, minioPubUrl string, bucketName string, cfg *config.VaultConfig) *vaultCategoryService {
	return &vaultCategoryService{
		repo:        re,
		cpsService:  cpsS,
		logger:      logger,
		minio:       minio,
		minioPubUrl: minioPubUrl,
		bucketName:  bucketName,
		cfg:         cfg,
	}
}

func (s *vaultCategoryService) CreateVaultCategory(ctx context.Context, req *vault_category_dto.CreateCategoryRequest) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	vCategory, _ := s.repo.FindByName(ctx, req.Name)
	if vCategory != nil {
		span.AddEvent("Duplicate vault category", trace.WithAttributes(attribute.String("name", req.Name)))
		return "", errors.New(localization.ErrorDuplicateGroupVaultCategory.Code)

	}

	coverImageUrl, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage, string(constants.VaultGroupCategoryFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("File upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf(localization.ErrorFileUploadFailed.Code)
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	var tiers []imodel.VaultTiers
	for _, t := range req.Tiers {
		tiers = append(tiers, imodel.VaultTiers{
			Name:         t.Name,
			TierInterest: t.TierInterest,
			MinAmount:    t.Min,
			MaxAmount:    t.Max,
		})
	}

	req_data := &imodel.VaultCategory{
		Name:             req.Name,
		CoverImageURL:    coverImageUrl,
		InterestType:     req.InterestType,
		CategoryInterest: req.CategoryInterest,
		Deadlock:         *req.Deadlock,
		Tiers:            tiers,
		IsActive:         true,
	}

	cpsActionModel := lib.CpsModelBuilder("", makerData, nil, req_data, string(constants.RequestCreateVaultCategory), string(constants.CREATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to create CPS action for vault category | action=%s | err=%v", constants.RequestCreateVaultCategory, err)
		return "", err
	}
	span.AddEvent("VaultCategory created", trace.WithAttributes(attribute.String("name", req.Name)))
	return "", nil
}

func (s *vaultCategoryService) FindAllVaultCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.VaultCategory], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllVaultCategories", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to fetch vault categories | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	resp := make([]*imodel.VaultCategory, 0, len(entities.Data))
	for _, en := range entities.Data {
		span.AddEvent("Mapping vault category to response", trace.WithAttributes(attribute.String("id", en.ID)))
		resp = append(resp, helperr.MapVaultCategoryToResponse(en))
	}
	return &types.PaginatedResponse[[]*imodel.VaultCategory]{
		Data: resp,
		Meta: entities.Meta,
	}, nil
}

func (s *vaultCategoryService) GetVaultCategory(ctx context.Context, id string) (*imodel.VaultCategory, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault category by id | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return helperr.MapVaultCategoryToResponse(entity), nil
}

func (s *vaultCategoryService) UpdateVaultCategory(ctx context.Context, id string, req *vault_category_dto.UpdateCategoryRequest) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return "", errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault category by id | err=%v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	updatedName := strings.ToUpper(prev.Name)
	updatedCover := prev.CoverImageURL

	if req.Name != nil {
		updatedName = strings.ToUpper(*req.Name)
	}

	var coverImageUrl string
	if req.CoverImage != nil {
		coverImageUrl, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage, string(constants.VaultGroupCategoryFolderName), *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("File upload failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
			s.logger.Errorf(localization.ErrorFileUploadFailed.Code)
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		updatedCover = coverImageUrl
	}

	interestType := prev.InterestType
	if req.InterestType != nil {
		interestType = *req.InterestType
	}

	categoryInterest := prev.CategoryInterest
	if req.CategoryInterest != nil {
		categoryInterest = *req.CategoryInterest
	}

	deadlock := prev.Deadlock
	if req.Deadlock != nil {
		deadlock = *req.Deadlock
	}

	req_data := &imodel.VaultCategory{
		Name:             updatedName,
		CoverImageURL:    updatedCover,
		InterestType:     interestType,
		CategoryInterest: categoryInterest,
		Deadlock:         deadlock,
		Tiers:            prev.Tiers,
		UpdatedAt:        time.Now(),
		IsActive:         prev.IsActive,
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, makerData, prev, req_data, string(constants.RequestUpdateVaultCategory), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault category | action=%s | err=%v",
				constants.RequestUpdateVaultCategory, err)
		}
		return "", err
	}

	span.AddEvent("VaultCategory updated", trace.WithAttributes(attribute.String("id", id)))
	return id, nil
}

func (s *vaultCategoryService) DeleteVaultCategory(ctx context.Context, id string) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	s.logger.Infof("Deleting vault category with ID: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault category id is required")
		return "", errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return "", errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault category by id | err=%v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	if exist.IsActive {
		span.AddEvent("Cannot delete active", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("cannot delete active vault category with id: %s", id)
		return "", errors.New(localization.ErrorCannotDeleteActiveVaultGroupCategory.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault category action | context = %v", maker)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, exist, string(constants.RequestDeleteVaultCategory), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault category | action=%s | err=%v", constants.RequestDeleteVaultCategory, err)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	span.AddEvent("VaultCategory deleted", trace.WithAttributes(attribute.String("id", id)))
	return id, nil
}

func (s *vaultCategoryService) EnableVaultCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	s.logger.Infof("Enabling vault category with ID: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault category id is required")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if exist.IsActive {
		span.AddEvent("Already enabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault category already enabled with id: %s", id)
		return errors.New(localization.ErrorVaultGroupAlreadyEnabled.Code)
	}
	updated := *exist
	updated.IsActive = true
	updated.UpdatedAt = time.Now()
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault category action | context = %v", maker)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, updated, string(constants.RequestEnableVaultCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultCategoryService) DisableVaultCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	s.logger.Infof("Disabling vault category with ID: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault category id is required")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("Unexpected error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if !exist.IsActive {
		span.AddEvent("Already disabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault category already disabled with id: %s", id)
		return errors.New(localization.ErrorVaultGroupAlreadyDisabled.Code)
	}
	updated := *exist
	updated.IsActive = false
	updated.UpdatedAt = time.Now()
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
			s.logger.Errorf("incomplete user context for vault category action | context = %v", maker)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, updated, string(constants.RequestDisAbleVaultCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to marshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to marshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	actionData := helperr.CategoryMapper(actionMap.(map[string]interface{}))

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateVaultCategory):
		span.AddEvent("RequestCreateVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		_, err := s.repo.Create(ctx, &actionData)
		if err != nil {
			span.AddEvent("Failed to create vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("[VaultCategory Authorize] failed to authorize category creation %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateVaultCategory):
		span.AddEvent("RequestUpdateVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		if err := s.repo.Update(ctx, cpsAction.UniqueId, &actionData); err != nil {
			span.AddEvent("Failed to update vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
	case string(constants.RequestDeleteVaultCategory):
		span.AddEvent("RequestDeleteVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		if _, err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			span.AddEvent("Failed to delete vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestEnableVaultCategory):
		span.AddEvent("RequestEnableVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			span.AddEvent("Failed to enable vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestDisAbleVaultCategory):
		span.AddEvent("RequestDisAbleVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			span.AddEvent("Failed to disable vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	}
	return cpsAction, nil
}
