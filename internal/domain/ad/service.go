package ad

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	util_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// AdvertService defines the interface for advert business logic
type AdvertService interface {
	CreateAdvert(ctx context.Context, request AdvertRequest) (*Advert, error)
	UpdateAdvert(ctx context.Context, id string, request AdvertRequest) (*Advert, *Advert, error)
	DeleteAdvert(ctx context.Context, id string) (*Advert, *Advert, error)
	EnableDisableAdvert(ctx context.Context, id string, enable bool) (*Advert, *Advert, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
	FetchAdvertByID(ctx context.Context, id string) (*Advert, error)
	FetchAdverts(ctx context.Context, filterParam *util_constant.Filter) (*utils.PaginatedResponse[[]*Advert], error)
}

// Service implements AdvertService
type Service struct {
	Repository  ADRepository
	logger      shared.Logger
	minioClient config.MinioClientInterface
	bucketName  string
	cfg         *config.VaultConfig
}

// NewAdvertService creates a new advert service instance
func NewAdvertService(repository ADRepository, minioClient config.MinioClientInterface, bucketName string, cfg *config.VaultConfig, logger shared.Logger) AdvertService {
	return &Service{
		Repository:  repository,
		logger:      logger,
		minioClient: minioClient,
		bucketName:  bucketName,
		cfg:         cfg,
	}
}

// CreateAdvert prepares a new advert without persisting
func (s *Service) CreateAdvert(ctx context.Context, request AdvertRequest) (*Advert, error) {
	s.logger.Infof("Creating advert, title: %s", request.Title)

	url, err := utils.UploadFileToMinio(ctx, s.minioClient, s.bucketName, request.BannerImage, "advert", s.cfg.MinioEndPoint, s.logger)
	if err != nil {
		s.logger.Errorf("Failed to upload banner image: %v", err)
		return nil, err
	}

	result := &Advert{
		Title:         request.Title,
		Description:   request.Description,
		BannerImage:   url,
		AdvertFor:     AdvertFor(request.AdvertFor),
		Date:          AdvertDate(request.Date),
		CreatedAt:     time.Now(),
		LastUpdatedAt: time.Now(),
		Enabled:       false,
		IsDeleted:     false,
	}

	s.logger.Infof("Advert constructed successfully, title: %s", result.Title)
	return result, nil
}

// UpdateAdvert updates an existing advert without persisting
func (s *Service) UpdateAdvert(ctx context.Context, id string, request AdvertRequest) (*Advert, *Advert, error) {
	s.logger.Infof("Updating advert, id: %s", id)

	prevAdvert, err := s.Repository.FetchAdvertByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch advert, id: %s, error: %v", id, err)
		return nil, nil, err
	}

	var url string
	if request.BannerImage != nil {
		url, err = utils.UploadFileToMinio(ctx, s.minioClient, s.bucketName, request.BannerImage, "advert", s.cfg.MinioEndPoint, s.logger)
		if err != nil {
			s.logger.Errorf("Failed to upload banner image: %v", err)
			return nil, nil, err
		}
	}

	curAdvert := &Advert{
		ID:            prevAdvert.ID,
		Title:         nonEmptyString(request.Title, prevAdvert.Title),
		Description:   nonEmptyString(request.Description, prevAdvert.Description),
		BannerImage:   nonEmptyString(url, prevAdvert.BannerImage),
		AdvertFor:     nonEmptyAdvertFor(AdvertFor(request.AdvertFor), prevAdvert.AdvertFor),
		Date:          nonEmptyAdvertDate(AdvertDate(request.Date), prevAdvert.Date),
		Enabled:       prevAdvert.Enabled,
		IsDeleted:     prevAdvert.IsDeleted,
		CreatedAt:     prevAdvert.CreatedAt,
		LastUpdatedAt: time.Now(),
	}

	s.logger.Infof("Advert updated successfully, id: %s", id)
	return curAdvert, prevAdvert, nil
}

// DeleteAdvert soft-deletes an advert without persisting
func (s *Service) DeleteAdvert(ctx context.Context, id string) (*Advert, *Advert, error) {
	s.logger.Infof("Deleting advert, id: %s", id)

	prevAdvert, err := s.Repository.FetchAdvertByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch advert, id: %s, error: %v", id, err)
		return nil, nil, err
	}

	curAdvert := generateAdvert(*prevAdvert)
	curAdvert.IsDeleted = true
	curAdvert.DeletedAt = time.Now()
	curAdvert.LastUpdatedAt = time.Now()

	s.logger.Infof("Advert marked for deletion, id: %s", id)
	return curAdvert, prevAdvert, nil
}

// EnableDisableAdvert enables or disables an advert without persisting
func (s *Service) EnableDisableAdvert(ctx context.Context, id string, enable bool) (*Advert, *Advert, error) {
	s.logger.Infof("EnableDisable advert, id: %s, enable: %v", id, enable)

	prevAdvert, err := s.Repository.FetchAdvertByID(ctx, id)
	if err != nil {
		s.logger.Errorf("Failed to fetch advert, id: %s, error: %v", id, err)
		return nil, nil, err
	}

	if enable && prevAdvert.Enabled {
		s.logger.Errorf("Advert already enabled, id: %s", id)
		return nil, nil, fmt.Errorf(utils.ErrAlreadyEnabled)
	}

	if !enable && !prevAdvert.Enabled {
		s.logger.Errorf("Advert already disabled, id: %s", id)
		return nil, nil, fmt.Errorf(utils.ErrAlreadyDisabled)
	}

	curAdvert := generateAdvert(*prevAdvert)
	curAdvert.Enabled = enable
	curAdvert.LastUpdatedAt = time.Now()

	s.logger.Infof("Advert enable/disable prepared, id: %s, enable: %v", id, enable)
	return curAdvert, prevAdvert, nil
}

// Authorize handles persistence for advert actions
func (s *Service) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	s.logger.Infof("Authorizing advert action, action: %s", action.RequestAction)

	var advert *Advert
	if err := utils.BindAction(action.CurrentAction, &advert); err != nil {
		s.logger.Errorf("Failed to bind current action to advert: %v", err)
		return nil, fmt.Errorf(utils.InvalidActionData)
	}

	var err error
	switch action.RequestAction {
	case constant.RequestCreateAdvert:
		advert, err = s.Repository.CreateAdvert(ctx, advert)
	case constant.RequestUpdateAdvert:
		advert, err = s.Repository.UpdateAdvert(ctx, advert)
	case constant.RequestDeleteAdvert:
		advert, err = s.Repository.DeleteAdvert(ctx, advert.ID)
	case constant.RequestEnableAdvert:
		advert, err = s.Repository.EnableDisableAdvert(ctx, advert.ID, true)
	case constant.RequestDisableAdvert:
		advert, err = s.Repository.EnableDisableAdvert(ctx, advert.ID, false)
	default:
		s.logger.Errorf("Unsupported action requested, action: %s", action.RequestAction)
		return nil, fmt.Errorf(utils.ErrUnsupported)
	}

	if err != nil {
		s.logger.Errorf("Failed to process advert action, action: %s, error: %v", action.RequestAction, err)
		return nil, err
	}

	action.CurrentAction = advert
	s.logger.Infof("Authorization completed for action, action: %s, id: %s", action.RequestAction, advert.ID)
	return action, nil
}

// FetchAdvertByID fetches an advert by ID
func (s *Service) FetchAdvertByID(ctx context.Context, id string) (*Advert, error) {
	s.logger.Infof("Fetching advert by ID, id: %s", id)
	return s.Repository.FetchAdvertByID(ctx, id)
}

// FetchAdverts fetches adverts with pagination and filtering
func (s *Service) FetchAdverts(ctx context.Context, filterParam *util_constant.Filter) (*utils.PaginatedResponse[[]*Advert], error) {
	s.logger.Infof("Fetching adverts with filter, filter: %v", filterParam)
	return s.Repository.FetchAdverts(ctx, filterParam)
}
