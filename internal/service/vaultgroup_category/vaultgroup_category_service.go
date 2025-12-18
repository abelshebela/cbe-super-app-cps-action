package vaultgroupcategory

import (
	"cbe-super-app-cps-action/internal/constants"
	vaultgroup_category "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helperr "cbe-super-app-cps-action/internal/service/vaultgroup_category/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type vaultgroupCategoryService struct {
	repo        storage.VaultGroupCategoryRepository
	cpsService  service.CPSActionService
	logger      shared_utils.Logger
	minio       *s3.Client
	minioPubUrl string
	bucketName  string
	cfg         *config.VaultConfig
}

func NewVaultGroupCategoryService(re storage.VaultGroupCategoryRepository, cpsS service.CPSActionService, logger shared_utils.Logger, minio *s3.Client, minioPubUrl string, bucketName string, cfg *config.VaultConfig) *vaultgroupCategoryService {
	return &vaultgroupCategoryService{
		repo:        re,
		cpsService:  cpsS,
		logger:      logger,
		minio:       minio,
		minioPubUrl: minioPubUrl,
		bucketName:  bucketName,
		cfg:         cfg,
	}
}

func (s *vaultgroupCategoryService) CreateVaultGroupCategory(ctx context.Context, req *vaultgroup_category.CreateVaultGroupCategoryRequest) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()
	_, err := s.repo.GetGroupcategoryByName(ctx, req.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			makerData := local_util.ExtractUserFromContext(ctx)

			coverImageUrl, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage, "", *s.cfg, "", s.logger)
			if err != nil {
				span.AddEvent("File upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
				s.logger.Errorf(localization.ErrorFileUploadFailed.Code)
			}

			req_data := &model.VaultCategory{
				Name:         req.Name,
				CategoryType: req.CategoryType,
				CoverImage:   coverImageUrl,
				IsActive:     false,
			}

			cpsActionModel := lib.CpsModelBuilder("", makerData, nil, req_data, string(constants.RequestCreateVaultGroupCategory), string(constants.CREATE))
			if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
				span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error())))
				if s.logger != nil {
					s.logger.Errorf("failed to create CPS action for vault group category | action=%s | err=%v", constants.RequestCreateVaultGroupCategory, err)
				}
				return "", err
			}
			span.AddEvent("VaultGroupCategory created", trace.WithAttributes(attribute.String("name", req.Name)))
			return "", nil
		}
	}
	span.AddEvent("Duplicate group vault category", trace.WithAttributes(attribute.String("name", req.Name)))
	return "", errors.New(localization.ErrorDuplicateGroupVaultCategory.Code)
}

func (s *vaultgroupCategoryService) FindAllVaultGroupCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*vaultgroup_category.VaultGroupCategoryResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllVaultGroupCategories", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("failed to fetch vault group categories | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	resp := make([]*vaultgroup_category.VaultGroupCategoryResponse, 0, len(entities.Data))
	for _, en := range entities.Data {
		span.AddEvent("Mapping vault group category to response", trace.WithAttributes(attribute.String("id", en.ID)))
		resp = append(resp, helperr.MapVaultGroupCategoryToResponse(en))
	}
	return &types.PaginatedResponse[[]*vaultgroup_category.VaultGroupCategoryResponse]{
		Data: resp,
		Meta: entities.Meta,
	}, nil
}

func (s *vaultgroupCategoryService) GetVaultGroupCategory(ctx context.Context, id string) (*vaultgroup_category.VaultGroupCategoryResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return nil, errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return helperr.MapVaultGroupCategoryToResponse(entity), nil
}
func (s *vaultgroupCategoryService) UpdateVaultGroupCategory(ctx context.Context, id string, req *vaultgroup_category.UpdateVaultGroupCategoryRequest) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()
	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return "", errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	updatedName := strings.ToUpper(prev.Name)
	updatedCover := prev.CoverImage
	updatedCategoryType := prev.CategoryType

	if req.Name != nil {
		updatedName = *req.Name
	}

	if req.CategoryType != "" {
		updatedCategoryType = req.CategoryType
	}

	var coverImageUrl string
	if req.CoverImage != nil {
		coverImageUrl, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage, "", *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("File upload failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
			s.logger.Errorf(localization.ErrorFileUploadFailed.Code)
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		updatedCover = coverImageUrl
	}

	req_data := &model.VaultCategory{
		Name:         updatedName,
		CategoryType: updatedCategoryType,
		CoverImage:   updatedCover,
		UpdatedAt:    time.Now(),
		IsActive:     prev.IsActive,
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, makerData, prev, req_data, string(constants.RequestUpdateVaultGroupCategory), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault group category | action=%s | err=%v",
				constants.RequestUpdateVaultGroupCategory, err)
		}
		return "", err
	}

	span.AddEvent("VaultGroupCategory updated", trace.WithAttributes(attribute.String("id", id)))
	return id, nil
}

func (s *vaultgroupCategoryService) DeleteVaultGroupCategory(ctx context.Context, id string) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()
	s.logger.Infof("Deleting vault group category with ID: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category id is required")
		return "", errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return "", errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	if exist.IsDeleted {
		span.AddEvent("Already deleted", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category already deleted with id: %s", id)
		return "", errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
	}
	if exist.IsActive {
		span.AddEvent("Cannot delete active", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("cannot delete active vault group category with id: %s", id)
		return "", errors.New(localization.ErrorCannotDeleteActiveVaultGroupCategory.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault group category action | context = %v", maker)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	// mongosafeExist := helperr.ConvertVaultGroupCategoryToMongoSafe(exist)

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, exist, string(constants.RequestDeleteVaultGroupCategory), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("failed to create CPS action for vault group category | action=%s | err=%v", constants.RequestDeleteVaultGroupCategory, err)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	span.AddEvent("VaultGroupCategory deleted", trace.WithAttributes(attribute.String("id", id)))
	return id, nil
}
func (s *vaultgroupCategoryService) EnableVaultGroupCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()

	s.logger.Infof("Enabling vault group category with ID: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category id is required")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return localization.ErrorVaultGroupCategoryNotFound
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if exist.IsDeleted {
		span.AddEvent("Already deleted", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category already deleted with id: %s", id)
		return localization.ErrorVaultGroupCategooryAlreadyDeleted
	}
	if exist.IsActive {
		span.AddEvent("Already enabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category already enabled with id: %s", id)
		return errors.New(localization.ErrorVaultGroupAlreadyEnabled.Code)
	}
	updated := *exist
	updated.IsActive = true
	updated.UpdatedAt = time.Now()
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		if s.logger != nil {
			s.logger.Errorf("incomplete user context for vault group category action | context = %v", maker)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	// Convert to MongoDB-safe format
	// mongoSafeExist := helperr.ConvertVaultGroupCategoryToMongoSafe(exist)
	// mongoSafeUpdated := helperr.ConvertVaultGroupCategoryToMongoSafe(updated)
	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, updated, string(constants.RequestEnableVaultGroupCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultgroupCategoryService) DisableVaultGroupCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
	defer span.End()
	s.logger.Infof("Disabling vault group category with ID: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category id is required")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		s.logger.Errorf("failed to fetch vault group category by id | err=%v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return errors.New(localization.ErrorVaultGroupCategoryNotFound.Code)
		}
		span.AddEvent("Unexpected error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if !exist.IsActive {
		span.AddEvent("Already disabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("vault group category already disabled with id: %s", id)
		return errors.New(localization.ErrorVaultGroupAlreadyDisabled.Code)
	}
	updated := *exist
	updated.IsActive = false
	updated.UpdatedAt = time.Now()
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
			s.logger.Errorf("incomplete user context for vault group category action | context = %v", maker)
		}
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	// Convert to MongoDB-safe format
	// mongoSafeExist := helperr.ConvertVaultGroupCategoryToMongoSafe(exist)
	// mongoSafeUpdated := helperr.ConvertVaultGroupCategoryToMongoSafe(updated)
	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, updated, string(constants.RequestDisAbleVaultGroupCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultgroupCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeVaultGroupCategory", "vaultgroupCategoryService", "vaultgroupCategoryService")
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
	case string(constants.RequestCreateVaultGroupCategory):
		span.AddEvent("RequestCreateVaultGroupCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		fmt.Println("========ACTION DATA NAME========", actionData.CategoryType)
		_, err := s.repo.Create(ctx, &actionData)
		if err != nil {
			span.AddEvent("Failed to create vault group category", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("[VaultCategory Authorize] failed to authorize category creation %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateVaultGroupCategory):
		span.AddEvent("RequestUpdateVaultGroupCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		if err := s.repo.Update(ctx, cpsAction.UniqueId, &actionData); err != nil {
			span.AddEvent("Failed to update vault group category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
	case string(constants.RequestDeleteVaultGroupCategory):
		span.AddEvent("RequestDeleteVaultGroupCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		_, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind vault group category from CPSAction", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if _, err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			span.AddEvent("Failed to delete vault group category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestEnableVaultGroupCategory):
		span.AddEvent("RequestEnableVaultGroupCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		_, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind vault group category from CPSAction", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			span.AddEvent("Failed to enable vault group category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	case string(constants.RequestDisAbleVaultGroupCategory):
		span.AddEvent("RequestDisAbleVaultGroupCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))
		_, err := helperr.BindVaultGroupCategoryFromCPSAction(cpsAction.CurrentAction)
		if err != nil {
			span.AddEvent("Failed to bind vault group category from CPSAction", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			span.AddEvent("Failed to disable vault group category", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		return cpsAction, nil

	}
	// return nil, errors.New(localization.ErrorInvalidRequest.Code)
	return cpsAction, nil
}
