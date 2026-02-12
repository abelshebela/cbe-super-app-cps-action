package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/ad/core"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"mime/multipart"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_constant "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Service implements AdvertService
type advertService struct {
	Repository  storage.AdvertRepository
	cpsService  service.CPSActionService
	logger      utils.Logger
	minioClient *s3.Client
	bucketName  string
	cfg         *config.VaultConfig
}

// NewAdvertService creates a new advert service instance
func NewAdvertService(repository storage.AdvertRepository, cpsService service.CPSActionService, minioClient *s3.Client, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.AdvertService {
	return &advertService{
		Repository:  repository,
		cpsService:  cpsService,
		logger:      logger,
		minioClient: minioClient,
		bucketName:  bucketName,
		cfg:         cfg,
	}
}

// handleCPSAction encapsulates the common CPS action logic
func (s *advertService) handleCPSAction(ctx context.Context, uniqueID string, requestAction cpsaction.RequestAction, curData, prevData interface{}, actionType cpsaction.ActionType) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		s.logger.Errorf("[handleCPSAction] incomplete user context")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, maker, prevData, curData, string(requestAction), string(actionType))

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		s.logger.Errorf("[handleCPSAction] failed to create CPS action, action: %s, error: %v", requestAction, err)
		return err
	}
	return nil
}

// CreateAdvert prepares a new advert without persisting
func (s *advertService) CreateAdvert(ctx context.Context, ad *model.Advert, bannerImage *multipart.FileHeader) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateAdvert", "Ad", "CreateAdvert")
	defer span.End()

	s.logger.Infof("[CreateAdvert] creating advert")
	isDuplicate, err := s.Repository.FindByTitle(ctx, ad.Title)
	if err != nil {
		span.AddEvent("[CreateAdvert] failed to check for duplicate advert", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("title", ad.Title),
		))
		s.logger.Errorf("[CreateAdvert] failed to check for duplicate advert: %v", err)
		return err
	}
	if isDuplicate != nil {
		span.AddEvent("[CreateAdvert] duplicate advert title found", trace.WithAttributes(attribute.String("title", ad.Title)))
		s.logger.Errorf("[CreateAdvert] duplicate advert title found")
		return errors.New(localization.ErrorAdvertTitleAlreadyExists.Code)
	}

	url, err := lib.UploadFileToMinio(ctx, s.minioClient, s.bucketName, bannerImage, string(constants.AdFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("[CreateAdvert] failed to upload banner image", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("title", ad.Title),
		))
		s.logger.Errorf("[CreateAdvert] failed to upload banner image: %v", err)
		return errors.New(localization.ErrorFileUploadFailed.Code)
	}

	// update banner image url after uploading
	ad.BannerImage = url
	ad.Enabled = true

	err = s.handleCPSAction(ctx, "", cpsaction.RequestCreateAdvert, ad, nil, cpsaction.ActionCreate)
	if err != nil {
		span.AddEvent("[CreateAdvert] failed to handle CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("title", ad.Title),
		))
		s.logger.Errorf("[CreateAdvert] failed to handle CPS action: %v", err)
		return err
	}

	s.logger.Infof("[CreateAdvert] advert creation request created successfully")
	return nil
}

// FetchAdverts fetches adverts with pagination and filtering
func (s *advertService) FetchAdverts(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Advert], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchAdverts", "Ad", "FetchAdverts")
	defer span.End()

	result, err := s.Repository.FindAllWithPagination(ctx, filterParam)
	if err != nil {
		span.AddEvent("[FetchAdverts] failed to fetch adverts", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[FetchAdverts] failed to fetch adverts: %v", err)
		return nil, err
	}
	s.logger.Infof("[FetchAdverts] retrieved %d adverts", len(result.Data))
	return result, nil
}

// FetchAdvertByID fetches an advert by ID
func (s *advertService) FetchAdvertByID(ctx context.Context, id string) (*model.Advert, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchAdvertByID", "Ad", "FetchAdvertByID")
	defer span.End()

	result, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[FetchAdvertByID] failed to fetch advert", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[FetchAdvertByID] failed to fetch advert: %v", err)
		return nil, err
	}
	s.logger.Infof("[FetchAdvertByID] advert retrieved successfully for id: %s", id)
	return result, nil
}

// // UpdateAdvert updates an existing advert without persisting
func (s *advertService) UpdateAdvert(ctx context.Context, id string, ad *model.Advert, bannerImage *multipart.FileHeader) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateAdvert", "Ad", "UpdateAdvert")
	defer span.End()

	s.logger.Infof("[UpdateAdvert] updating advert for id: %s", id)
	prevAdvert, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[UpdateAdvert] failed to fetch advert", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[UpdateAdvert] failed to fetch advert: %v", err)
		return err
	}
	if ad.Title != "" {
		isDuplicate, err := s.Repository.FindByTitle(ctx, ad.Title)
		if err != nil {
			span.AddEvent("[UpdateAdvert] failed to check for duplicate advert", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			s.logger.Errorf("[UpdateAdvert] failed to check for duplicate advert: %v", err)
			return err
		}
		if isDuplicate != nil && isDuplicate.ID.Hex() != id {
			span.AddEvent("[UpdateAdvert] duplicate advert title found", trace.WithAttributes(attribute.String("title", ad.Title)))
			s.logger.Errorf("[UpdateAdvert] duplicate advert title found")
			return errors.New(localization.ErrorAdvertTitleAlreadyExists.Code)
		}
	}

	var url string
	if bannerImage != nil {

		url, err = lib.UploadFileToMinio(ctx, s.minioClient, s.bucketName, bannerImage, string(constants.AdFolderName), *s.cfg, "", s.logger)
		if err != nil {
			span.AddEvent("[UpdateAdvert] failed to upload banner image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			s.logger.Errorf("Failed to upload banner image: %v", err)
			return err
		}
	}

	curAdvert := &model.Advert{
		ID:            prevAdvert.ID,
		Title:         local_util.NonEmptyString(ad.Title, prevAdvert.Title),
		Description:   local_util.NonEmptyString(ad.Description, prevAdvert.Description),
		BannerImage:   local_util.NonEmptyString(url, prevAdvert.BannerImage),
		AdvertFor:     local_util.NonEmptyAdvertFor(shared_constant.AdvertFor(ad.AdvertFor), prevAdvert.AdvertFor),
		Enabled:       prevAdvert.Enabled,
		IsDeleted:     prevAdvert.IsDeleted,
		CreatedAt:     prevAdvert.CreatedAt,
		LastUpdatedAt: time.Now(),
	}
	err = s.handleCPSAction(ctx, id, cpsaction.RequestUpdateAdvert, curAdvert, prevAdvert, cpsaction.ActionUpdate)
	if err != nil {
		span.AddEvent("[UpdateAdvert] failed to handle CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[UpdateAdvert] failed to handle CPS action: %v", err)
		return err
	}
	s.logger.Infof("[UpdateAdvert] advert update request created successfully for id: %s", id)
	return nil
}

// DeleteAdvert soft-deletes an advert without persisting
func (s *advertService) DeleteAdvert(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteAdvert", "Ad", "DeleteAdvert")
	defer span.End()

	s.logger.Infof("[DeleteAdvert] deleting advert for id: %s", id)

	prevAdvert, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[DeleteAdvert] failed to fetch advert", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[DeleteAdvert] failed to fetch advert: %v", err)
		return err
	}

	curAdvert := core.GenerateAdvert(*prevAdvert)
	curAdvert.IsDeleted = true
	curAdvert.DeletedAt = time.Now()
	curAdvert.LastUpdatedAt = time.Now()

	err = s.handleCPSAction(ctx, curAdvert.ID.Hex(), cpsaction.RequestDeleteAdvert, curAdvert, prevAdvert, cpsaction.ActionDelete)
	if err != nil {
		span.AddEvent("[DeleteAdvert] failed to handle CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[DeleteAdvert] failed to handle CPS action: %v", err)
		return err
	}

	s.logger.Infof("[DeleteAdvert] advert deletion request created successfully for id: %s", id)
	return nil
}

// EnableDisableAdvert enables or disables an advert without persisting
func (s *advertService) EnableDisableAdvert(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisableAdvert", "Ad", "EnableDisableAdvert")
	defer span.End()

	s.logger.Infof("[EnableDisableAdvert] processing advert enable/disable for id: %s, enabled: %v", id, enable)

	prevAdvert, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[EnableDisableAdvert] failed to fetch advert", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[EnableDisableAdvert] failed to fetch advert: %v", err)
		return err
	}

	if enable && prevAdvert.Enabled {
		span.AddEvent("[EnableDisableAdvert] advert already enabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[EnableDisableAdvert] advert already enabled")
		return errors.New(localization.ErrorAdvertAlreadyEnabled.Code)
	}

	if !enable && !prevAdvert.Enabled {
		span.AddEvent("[EnableDisableAdvert] advert already disabled", trace.WithAttributes(attribute.String("id", id)))
		s.logger.Errorf("[EnableDisableAdvert] advert already disabled")
		return errors.New(localization.ErrorAdvertAlreadyDisabled.Code)
	}

	curAdvert := core.GenerateAdvert(*prevAdvert)
	curAdvert.Enabled = enable
	curAdvert.LastUpdatedAt = time.Now()

	var action cpsaction.RequestAction
	if enable {
		action = cpsaction.RequestEnableAdvert
	} else {
		action = cpsaction.RequestDisableAdvert
	}

	err = s.handleCPSAction(ctx, curAdvert.ID.Hex(), action, curAdvert, prevAdvert, cpsaction.ActionUpdate)
	if err != nil {
		span.AddEvent("[EnableDisableAdvert] failed to handle CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		s.logger.Errorf("[EnableDisableAdvert] failed to handle CPS action: %v", err)
		return err
	}

	s.logger.Infof("[EnableDisableAdvert] advert enable/disable request created successfully for id: %s", id)
	return nil
}

// Authorize handles persistence for advert actions
func (s *advertService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Ad", "Authorize")
	defer span.End()

	s.logger.Infof("[Authorize] authorizing advert action: %s", action.RequestAction)

	advert, err := local_util.JsonUnmarshal[model.Advert](action.CurrentAction)
	if err != nil {
		span.AddEvent("[Authorize] failed to unmarshal current action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		s.logger.Errorf("[Authorize] failed to unmarshal current action into advert: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch action.RequestAction {
	case string(cpsaction.RequestCreateAdvert):
		err = s.Repository.Create(ctx, advert)
		if err != nil {
			span.AddEvent("[Authorize] failed to create advert", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to create advert: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] advert created successfully")
	case string(cpsaction.RequestUpdateAdvert):
		err = s.Repository.Update(ctx, action.UniqueId, advert)
		if err != nil {
			span.AddEvent("[Authorize] failed to update advert", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to update advert: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] advert updated successfully for id: %s", action.UniqueId)
	case string(cpsaction.RequestDeleteAdvert):
		err = s.Repository.Delete(ctx, action.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] failed to delete advert", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to delete advert: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] advert deleted successfully for id: %s", action.UniqueId)
	case string(cpsaction.RequestEnableAdvert):
		err = s.Repository.EnableOrDisable(ctx, action.UniqueId, true)
		if err != nil {
			span.AddEvent("[Authorize] failed to enable advert", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to enable advert: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] advert enabled successfully for id: %s", action.UniqueId)
	case string(cpsaction.RequestDisableAdvert):
		err = s.Repository.EnableOrDisable(ctx, action.UniqueId, false)
		if err != nil {
			span.AddEvent("[Authorize] failed to disable advert", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			s.logger.Errorf("[Authorize] failed to disable advert: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] advert disabled successfully for id: %s", action.UniqueId)
	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", action.RequestAction)))
		s.logger.Errorf("[Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	s.logger.Infof("[Authorize] advert action authorized successfully: %s", action.RequestAction)

	action.CurrentAction = advert
	s.logger.Infof("Authorization completed for action, action: %s, id: %s", action.RequestAction, advert.ID)
	return action, nil
}
