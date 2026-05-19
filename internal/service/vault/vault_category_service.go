package vault

import (
	"cbe-super-app-cps-action/internal/constants"
	vault_category_dto "cbe-super-app-cps-action/internal/constants/dto/vault"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	helperr "cbe-super-app-cps-action/internal/service/vault/core"
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
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	vCategory, _ := s.repo.FindByName(ctx, req.Name)
	if vCategory != nil {
		span.AddEvent("Duplicate vault category", trace.WithAttributes(attribute.String("name", req.Name)))
		return "", errors.New(localization.ErrorDuplicateVaultCategory.Code)

	}

	coverImageUrl, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage, string(constants.VaultCategoryFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("File upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[VaultCatSvc][Create] upload err")
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
		Name:          req.Name,
		CoverImageURL: coverImageUrl,
		InterestType:  strings.ToUpper(req.InterestType),
		Deadlock:      *req.Deadlock,
		Tiers:         tiers,
		IsActive:      true,
	}

	cpsActionModel := lib.CpsModelBuilder("", makerData, nil, req_data, string(constants.RequestCreateVaultCategory), string(constants.CREATE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[VaultCatSvc][Create] cps action err: %v", err)
		return "", err
	}
	span.AddEvent("VaultCategory created", trace.WithAttributes(attribute.String("name", req.Name)))
	return "", nil
}

func (s *vaultCategoryService) FindAllVaultCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*imodel.VaultCategory], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "FindAllVaultCategories", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	if filterParams == nil {
		f := &types.Filter{}
		filterParams = f

	}
	entities, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[VaultCatSvc][FindAll] err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return &types.PaginatedResponse[[]*imodel.VaultCategory]{
		Data: entities.Data,
		Meta: entities.Meta,
	}, nil
}

func (s *vaultCategoryService) GetVaultCategory(ctx context.Context, id string) (*imodel.VaultCategory, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorVaultCategoryNotFound.Code {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][GetByID] err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return entity, nil
}

func (s *vaultCategoryService) UpdateVaultCategory(ctx context.Context, id string, req *vault_category_dto.UpdateCategoryRequest) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	if req.Name != nil {
		dup, err := s.repo.FindByName(ctx, *req.Name)
		if err != nil && err.Error() != localization.ErrorVaultCategoryNotFound.Code {
			log.Errorf("[VaultCatSvc][Update] find err: %v", err)
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}

		if dup != nil && dup.ID != id {
			if strings.EqualFold(dup.Name, *req.Name) {
				return "", errors.New(localization.ErrorDuplicateVaultCategory.Code)
			}
		}
	}

	prev, err := s.repo.FindByID(ctx, id)
	if err != nil || prev == nil {
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) || prev == nil {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return "", errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Update] find err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	updatedName := strings.ToUpper(prev.Name)
	updatedCover := prev.CoverImageURL

	if req.Name != nil {
		updatedName = strings.ToUpper(*req.Name)
	}

	var coverImageUrl string
	if req.CoverImage != nil {
		coverImageUrl, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage, string(constants.VaultCategoryFolderName), *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("File upload failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
			log.Errorf("[VaultCatSvc][Update] upload err")
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
		updatedCover = coverImageUrl
	}

	interestType := prev.InterestType
	if req.InterestType != nil {
		interestType = strings.ToUpper(*req.InterestType)
	}

	deadlock := prev.Deadlock
	if req.Deadlock != nil {
		deadlock = *req.Deadlock
	}

	tiers := prev.Tiers
	if req.Tiers != nil {
		mapped := make([]imodel.VaultTiers, 0, len(req.Tiers))
		for i, t := range req.Tiers {
			var name, tierInterest, min, max string
			var id string
			if i < len(prev.Tiers) {
				// preserve existing values as defaults
				id = prev.Tiers[i].ID
				name = prev.Tiers[i].Name
				tierInterest = prev.Tiers[i].TierInterest
				min = prev.Tiers[i].MinAmount
				max = prev.Tiers[i].MaxAmount
			}
			if t.Name != nil {
				name = *t.Name
			}
			if t.TierInterest != nil {
				tierInterest = *t.TierInterest
			}
			if t.Min != nil {
				min = *t.Min
			}
			if t.Max != nil {
				max = *t.Max
			}
			mapped = append(mapped, imodel.VaultTiers{
				ID:           id,
				Name:         name,
				TierInterest: tierInterest,
				MinAmount:    min,
				MaxAmount:    max,
			})
		}
		tiers = mapped
	}

	req_data := &imodel.VaultCategory{
		Name:          updatedName,
		CoverImageURL: updatedCover,
		InterestType:  interestType,
		Deadlock:      deadlock,
		Tiers:         tiers,
		UpdatedAt:     time.Now(),
		IsActive:      prev.IsActive,
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(id, makerData, prev, req_data, string(constants.RequestUpdateVaultCategory), string(constants.UPDATE))

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		if s.logger != nil {
			log.Errorf("[VaultCatSvc][Update] cps action err: %v", err)
		}
		return "", err
	}

	span.AddEvent("VaultCategory updated", trace.WithAttributes(attribute.String("id", id)))
	return id, nil
}

func (s *vaultCategoryService) DeleteVaultCategory(ctx context.Context, id string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	log.Infof("[VaultCatSvc][Delete] id: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Delete] id required")
		return "", errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return "", errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Delete] find err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	if exist.IsActive {
		span.AddEvent("Cannot delete active", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Delete] active id: %s", id)
		return "", errors.New(localization.ErrorCannotDeleteActiveVaultCategory.Code)
	}
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		if s.logger != nil {
			log.Errorf("[VaultCatSvc][Delete] incomplete user")
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, exist, string(constants.RequestDeleteVaultCategory), string(constants.DELETE))
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		if s.logger != nil {
			log.Errorf("[VaultCatSvc][Delete] cps action err: %v", err)
		}
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	span.AddEvent("VaultCategory deleted", trace.WithAttributes(attribute.String("id", id)))
	return id, nil
}

func (s *vaultCategoryService) EnableVaultCategory(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()

	log.Infof("[VaultCatSvc][Enable] id: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Enable] id required")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Enable] find err: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	if exist.IsActive {
		span.AddEvent("Already enabled", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Enable] already enabled id: %s", id)
		return errors.New(localization.ErrorVaultAlreadyEnabled.Code)
	}
	updated := *exist
	updated.IsActive = true
	updated.UpdatedAt = time.Now()
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
		if s.logger != nil {
			log.Errorf("[VaultCatSvc][Enable] incomplete user")
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, updated, string(constants.RequestEnableVaultCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultCategoryService) DisableVaultCategory(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DisableVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	log.Infof("[VaultCatSvc][Disable] id: %s", id)
	if id == "" {
		span.AddEvent("ID not set", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Disable] id required")
		return errors.New(localization.ErrorIdNotSetOnQueryParam.Code)
	}
	exist, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Disable] find err: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, sql.ErrNoRows) || err.Error() == localization.ErrorVaultCategoryNotFound.Code {
			span.AddEvent("Not found", trace.WithAttributes(attribute.String("id", id)))
			return errors.New(localization.ErrorVaultCategoryNotFound.Code)
		}
		span.AddEvent("Unexpected error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	if !exist.IsActive {
		span.AddEvent("Already disabled", trace.WithAttributes(attribute.String("id", id)))
		log.Errorf("[VaultCatSvc][Disable] already disabled id: %s", id)
		return errors.New(localization.ErrorVaultAlreadyDisabled.Code)
	}
	updated := *exist
	updated.IsActive = false
	updated.UpdatedAt = time.Now()
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		if s.logger != nil {
			span.AddEvent("Incomplete user context", trace.WithAttributes(attribute.String("id", id)))
			log.Errorf("[VaultCatSvc][Disable] incomplete user")
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	cpsActionModel := lib.CpsModelBuilder(id, maker, exist, updated, string(constants.RequestDisAbleVaultCategory), string(constants.UPDATE))
	return s.cpsService.CreateCPSAction(ctx, &cpsActionModel)
}

func (s *vaultCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "AuthorizeVaultCategory", "vaultCategoryService", "vaultCategoryService")
	defer span.End()
	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to marshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[VaultCatSvc][Authorize] marshal err: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateVaultCategory):
		span.AddEvent("RequestCreateVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))

		err = json.Unmarshal(marshaled, &actionMap)
		if err != nil {
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[VaultCatSvc][Authorize] unmarshal create err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		actionData := helperr.CategoryMapper(actionMap.(map[string]interface{}))

		_, err := s.repo.Create(ctx, &actionData)
		if err != nil {
			span.AddEvent("Failed to create vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[VaultCatSvc][Authorize] create err: %v", err)
			return nil, err
		}
	case string(constants.RequestUpdateVaultCategory):
		span.AddEvent("RequestUpdateVaultCategory", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))

		err = json.Unmarshal(marshaled, &actionMap)
		if err != nil {
			span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[VaultCatSvc][Authorize] unmarshal update err: %v", err)
			return nil, errors.New(localization.ErrorInvalidActionData.Code)
		}

		actionData := helperr.CategoryMapper(actionMap.(map[string]interface{}))

		if err := s.repo.Update(ctx, cpsAction.UniqueId, &actionData); err != nil {
			span.AddEvent("Failed to update vault category", trace.WithAttributes(attribute.String("error", err.Error())))
			if err.Error() == localization.ErrorDuplicateVaultCategory.Code {
				return nil, err
			}
			log.Errorf("[VaultCatSvc][Authorize] update err: %v", err)
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

	// case string(constants.RequestCreateWithdrawal):
	// 	span.AddEvent("RequestCreateWithdrawal", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))

	// 	err = json.Unmarshal(marshaled, &actionMap)
	// 	if err != nil {
	// 		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		log.Errorf("[VaultCatSvc][Authorize] unmarshal withdrawal err: %v", err)
	// 		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	// 	}

	// 	actionData := helperr.WithdrawalMapper(actionMap.(map[string]interface{}))

	// 	if err := s.AuthorizeWithdrawalCreate(ctx, &actionData); err != nil {
	// 		span.AddEvent("Failed to create withdrawal request", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		log.Errorf("[VaultCatSvc][Authorize] withdrawal create err: %v", err)
	// 		return nil, err
	// 	}
	// 	return cpsAction, nil

	case string(constants.RequestUnlockDeadlock):
		span.AddEvent("RequestUpdateDeadlock", trace.WithAttributes(attribute.String("id", cpsAction.UniqueId)))

		// err = json.Unmarshal(marshaled, &actionMap)
		// if err != nil {
		// 	span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(attribute.String("error", err.Error())))
		// 	log.Errorf("[VaultCatSvc][Authorize] unmarshal withdrawal update err: %v", err)
		// 	return nil, errors.New(localization.ErrorInvalidActionData.Code)
		// }

		// status := helperr.WithdrawalStatusMapper(actionMap.(map[string]interface{}))
		// if status == "" {
		// 	return cpsAction, fmt.Errorf("withdrawal status required")
		// }

		if err := s.AuthorizeDeadlockStatusUpdate(ctx, cpsAction.UniqueId, "DISABLED"); err != nil {
			span.AddEvent("Failed to update deadlock request", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[VaultCatSvc][Authorize] deadlock update err: %v", err)
			return nil, err
		}
		return cpsAction, nil
	}
	return cpsAction, nil
}
