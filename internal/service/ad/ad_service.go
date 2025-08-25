package ad

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"mime/multipart"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// Service implements AdvertService
type advertService struct {
	Repository  storage.AdvertRepository
	cpsService  service.CPSActionService
	logger      utils.Logger
	minioClient config.MinioClientInterface
	bucketName  string
	cfg         *config.VaultConfig
}

// NewAdvertService creates a new advert service instance
func NewAdvertService(repository storage.AdvertRepository, cpsService service.CPSActionService, minioClient config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.AdvertService {
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
		s.logger.Errorf("Incomplete user context for advert creation | context = %v", maker)
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, maker,prevData,curData, string(requestAction),string(actionType))

	err := s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		s.logger.Errorf("[event.handleCPSAction] failed to create CPS action, action: %s, error: %v", requestAction, err)
		return err
	}
	return nil
}

// CreateAdvert prepares a new advert without persisting
func (s *advertService) CreateAdvert(ctx context.Context, ad *model.Advert, bannerImage *multipart.FileHeader) error {
	s.logger.Infof("Creating advert, title: %s", ad.Title)

	url, err := lib.UploadFileToMinio(ctx, s.minioClient, s.bucketName, bannerImage, "advert", s.cfg.MinioEndPoint, s.logger)
	if err != nil {
		s.logger.Errorf("Failed to upload banner image: %v", err)
		return errors.New(localization.MsgFileUploadFailed)
	}

	// update banner image url after uploading
	ad.BannerImage = url

	err = s.handleCPSAction(ctx, "", cpsaction.RequestCreateAdvert, ad, nil, cpsaction.ActionCreate)
	if err != nil {
		s.logger.Errorf("Failed to handle CPS action for advert creation, title: %s, error: %v", ad.Title, err)
		return err
	}

	s.logger.Infof("Advert constructed successfully, title: %s", ad.Title)
	return nil
}

// FetchAdverts fetches adverts with pagination and filtering
func (s *advertService) FetchAdverts(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Advert], error) {
	s.logger.Infof("Fetching adverts with filter, filter: %v", filterParam)
	return s.Repository.FindAllWithPagination(ctx, filterParam)
}

// FetchAdvertByID fetches an advert by ID
func (s *advertService) FetchAdvertByID(ctx context.Context, id string) (*model.Advert, error) {
	s.logger.Infof("Fetching advert by ID, id: %s", id)
	return s.Repository.FindByID(ctx, id)
}

// // UpdateAdvert updates an existing advert without persisting
func (s *advertService) UpdateAdvert(ctx context.Context, id string, ad *model.Advert, bannerImage *multipart.FileHeader) error {
	s.logger.Infof("Updating advert, id: %s", id)

	prevAdvert, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch advert, id: %s, error: %v", id, err)
		return err
	}

	var url string
	if bannerImage != nil {
		url, err = local_util.UploadFileToMinio(ctx, s.minioClient, s.bucketName, bannerImage, "advert", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			s.logger.Errorf("Failed to upload banner image: %v", err)
			return err
		}
	}

	curAdvert := &model.Advert{
		ID:            prevAdvert.ID,
		Title:         local_util.NonEmptyString(ad.Title, prevAdvert.Title),
		Description:   local_util.NonEmptyString(ad.Description, prevAdvert.Description),
		BannerImage:   local_util.NonEmptyString(url, prevAdvert.BannerImage),
		AdvertFor:     local_util.NonEmptyAdvertFor(constants.AdvertFor(ad.AdvertFor), prevAdvert.AdvertFor),
		Date:          local_util.NonEmptyAdvertDate(types.AdvertDate(ad.Date), prevAdvert.Date),
		Enabled:       prevAdvert.Enabled,
		IsDeleted:     prevAdvert.IsDeleted,
		CreatedAt:     prevAdvert.CreatedAt,
		LastUpdatedAt: time.Now(),
	}

	err = s.handleCPSAction(ctx, curAdvert.ID.Hex(), cpsaction.RequestUpdateAdvert, curAdvert, prevAdvert, cpsaction.ActionUpdate)
	if err != nil {
		s.logger.Errorf("Failed to handle CPS action for advert update, id: %s, error: %v", id, err)
		return err
	}
	s.logger.Infof("Advert updated successfully, id: %s", id)
	return nil
}

// DeleteAdvert soft-deletes an advert without persisting
func (s *advertService) DeleteAdvert(ctx context.Context, id string) error {
	s.logger.Infof("Deleting advert, id: %s", id)

	prevAdvert, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch advert, id: %s, error: %v", id, err)
		return err
	}

	curAdvert := GenerateAdvert(*prevAdvert)
	curAdvert.IsDeleted = true
	curAdvert.DeletedAt = time.Now()
	curAdvert.LastUpdatedAt = time.Now()

	err = s.handleCPSAction(ctx, curAdvert.ID.Hex(), cpsaction.RequestDeleteAdvert, curAdvert, prevAdvert, cpsaction.ActionDelete)
	if err != nil {
		s.logger.Errorf("Failed to handle CPS action for advert deletion, id: %s, error: %v", id, err)
		return err
	}

	s.logger.Infof("Advert marked for deletion, id: %s", id)
	return nil
}

// EnableDisableAdvert enables or disables an advert without persisting
func (s *advertService) EnableDisableAdvert(ctx context.Context, id string, enable bool) error {
	s.logger.Infof("EnableDisable advert, id: %s, enable: %v", id, enable)

	prevAdvert, err := s.Repository.FindByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch advert, id: %s, error: %v", id, err)
		return err
	}

	if enable && prevAdvert.Enabled {
		s.logger.Errorf("Advert already enabled, id: %s", id)
		return errors.New(localization.ErrorAdvertAlreadyEnabled.Code)
	}

	if !enable && !prevAdvert.Enabled {
		s.logger.Errorf("Advert already disabled, id: %s", id)
		return errors.New(localization.ErrorAdvertAlreadyDisabled.Code)
	}

	curAdvert := GenerateAdvert(*prevAdvert)
	curAdvert.Enabled = enable
	curAdvert.LastUpdatedAt = time.Now()

	s.logger.Infof("Advert enable/disable prepared, id: %s, enable: %v", id, enable)

	var action cpsaction.RequestAction
	if enable {
		action = cpsaction.RequestEnableAdvert
	} else {
		action = cpsaction.RequestDisableAdvert
	}

	err = s.handleCPSAction(ctx, curAdvert.ID.Hex(), action, curAdvert, prevAdvert, cpsaction.ActionUpdate)
	if err != nil {
		s.logger.Errorf("Failed to handle CPS action for advert enable/disable, id: %s, error: %v", id, err)
		return err
	}

	s.logger.Infof("Advert enable/disable action created, id: %s, enable: %v", id, enable)
	return nil
}

// Authorize handles persistence for advert actions
func (s *advertService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("Authorizing advert action, action: %s", action.RequestAction)

	var advert *model.Advert
	if err := local_util.BindAction(action.CurrentAction, &advert); err != nil {
		s.logger.Errorf("Failed to bind current action to advert: %v", err)
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	var err error
	switch action.RequestAction {
	case string(cpsaction.RequestCreateAdvert):
		err = s.Repository.Create(ctx, advert)
	case string(cpsaction.RequestUpdateAdvert):
		err = s.Repository.Update(ctx, advert.ID.Hex(), advert)
	case string(cpsaction.RequestDeleteAdvert):
		err = s.Repository.Delete(ctx, advert.ID.Hex())
	case string(cpsaction.RequestEnableAdvert):
		err = s.Repository.EnableOrDisable(ctx, advert.ID.Hex(), true)
	case string(cpsaction.RequestDisableAdvert):
		err = s.Repository.EnableOrDisable(ctx, advert.ID.Hex(), false)
	default:
		s.logger.Errorf("Unsupported action requested, action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	if err != nil {
		s.logger.Errorf("Failed to process advert action, action: %s, error: %v", action.RequestAction, err)
		return nil, err
	}

	action.CurrentAction = advert
	s.logger.Infof("Authorization completed for action, action: %s, id: %s", action.RequestAction, advert.ID)
	return action, nil
}
