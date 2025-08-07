package donation

import (
	"context"
	"fmt"
	"strings"

	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"mime/multipart"

	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DonationService interface {
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) (*dto.DonationCategoryResponse, error)
	UploadIcon(ctx context.Context, icon *multipart.FileHeader) (string, error)
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)
}

type Service struct {
	Repository DonationRepository
	logger     shared_utils.Logger
	minio      config.MinioClientInterface
	bucketName string
	cfg        *config.VaultConfig
}

func NewDonationService(repository DonationRepository,
	minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	logger shared_utils.Logger,
) DonationService {
	return &Service{
		Repository: repository,
		logger:     logger,
		minio:      minio,
		cfg:        cfg,
		bucketName: bucketName,
	}
}

func (e *Service) UploadIcon(ctx context.Context, icon *multipart.FileHeader) (string, error) {
	if icon == nil {
		return "", fmt.Errorf("ICON_IS_REQUIRED")
	}

	if e.minio == nil {
		e.logger.Errorf("MinIO client is nil")
		return "", fmt.Errorf("MINIO_CLIENT_NOT_CONFIGURED")
	}

	if e.cfg == nil {
		e.logger.Errorf("Configuration is nil")
		return "", fmt.Errorf("CONFIGURATION_NOT_LOADED")
	}

	iconURL, err := common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, icon, "donation_category_icon", e.cfg.MinioEndPoint, e.logger)
	if err != nil {
		e.logger.Errorf("Failed to upload icon to MinIO bucket '%s': %v", e.bucketName, err)
		if strings.Contains(err.Error(), "time") || strings.Contains(err.Error(), "server") || strings.Contains(err.Error(), "difference") {
			return "", fmt.Errorf("MINIO_TIME_SYNC_ERROR")
		}
		if strings.Contains(err.Error(), "bucket") || strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "not valid") {
			return "", fmt.Errorf("MINIO_BUCKET_ERROR")
		}
		if strings.Contains(err.Error(), "UNHANDLED_SERVER_ERROR") {
			return "", fmt.Errorf("MINIO_TIME_SYNC_ERROR")
		}
		return "", fmt.Errorf("FAILED_TO_UPLOAD_ICON")
	}

	e.logger.Infof("Successfully uploaded donation category icon with URL: %s", iconURL)
	return iconURL, nil
}

func (e *Service) CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) (*dto.DonationCategoryResponse, error) {
	exist, err := e.Repository.DonationNameExists(ctx, donation.CategoryName)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("CATEGORY_NAME_ALREADY_EXISTS")
	}

	if donation.Icon == nil {
		return nil, fmt.Errorf("ICON_IS_REQUIRED")
	}

	iconURL, err := e.UploadIcon(ctx, donation.Icon)
	if err != nil {
		return nil, err
	}

	result := dto.DonationCategoryResponse{
		CategoryName: donation.CategoryName,
		Icon:         iconURL,
	}
	return &result, nil
}

func (e *Service) UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error) {
	e.logger.Infof("Updating category ", "id", id)
	_, err := e.Repository.FetchDonationCategoryByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to fetch category name", "id", id, "error", err)
		return nil, err
	}

	if donation.Icon != nil {
		// Check if MinIO client is properly configured
		if e.minio == nil {
			e.logger.Errorf("MinIO client is nil")
			return nil, fmt.Errorf("MINIO_CLIENT_NOT_CONFIGURED")
		}

		if e.cfg == nil {
			e.logger.Errorf("Configuration is nil")
			return nil, fmt.Errorf("CONFIGURATION_NOT_LOADED")
		}

		e.logger.Infof("Starting updated donation category icon upload to MinIO bucket: %s", e.bucketName)
		e.logger.Infof("MinIO endpoint: %s", e.cfg.MinioEndPoint)

		iconURL, err := common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, donation.Icon, "donation_category_icon", e.cfg.MinioEndPoint, e.logger)
		if err != nil {
			e.logger.Errorf("Failed to upload icon to MinIO bucket '%s': %v", e.bucketName, err)
			// Check if it's a time synchronization issue
			if strings.Contains(err.Error(), "time") || strings.Contains(err.Error(), "server") || strings.Contains(err.Error(), "difference") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			// Check if it's a bucket creation issue
			if strings.Contains(err.Error(), "bucket") || strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "not valid") {
				return nil, fmt.Errorf("MINIO_BUCKET_ERROR")
			}
			// Check if it's a generic server error (which could be time sync)
			if strings.Contains(err.Error(), "UNHANDLED_SERVER_ERROR") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			return nil, fmt.Errorf("FAILED_TO_UPLOAD_ICON")
		}
		e.logger.Infof("Successfully uploaded updated donation category icon with URL: %s", iconURL)

		// Update the donation with the new URL
		donation.Icon = nil // Clear the file since we're storing URL
	}

	if donation.CategoryName != "" {
		exist, err := e.Repository.DonationNameExists(ctx, donation.CategoryName)
		if err != nil {
			return nil, err
		}

		if exist {
			return nil, fmt.Errorf("CATEGORY_NAME_ALREADY_EXISTS")
		}
	}

	// Update in database
	updatedDonation, err := e.Repository.UpdateDonationCategory(ctx, id, donation)
	if err != nil {
		e.logger.Errorf("Failed to update donation category in database: %v", err)
		return nil, err
	}

	return updatedDonation, nil
}

func (e *Service) FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error) {
	return e.Repository.FetchDonationCategory(ctx, filterParams)
}

func (e *Service) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	return e.Repository.FetchDonationCategoryByID(ctx, id)
}

func (e *Service) Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error) {
	requestedAction := action.RequestAction
	var err error

	castToBsonM := func(input interface{}) (bson.M, error) {
		raw, err := bson.Marshal(input)
		if err != nil {
			return nil, err
		}
		var out bson.M
		err = bson.Unmarshal(raw, &out)
		return out, err
	}

	switch requestedAction {
	case cps_const.RequestCreateDonationCategory:
		var donationCPS *dto.DonationCategoryCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &donationCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation Category: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		if donationCPS.Icon == "" {
			e.logger.Errorf("icon URL is empty in CPS request")
			return nil, fmt.Errorf("ICON_URL_MISSING")
		}

		err = e.Repository.CreateDonationCategoryWithURL(ctx, donationCPS.CategoryName, donationCPS.Icon)
		if err != nil {
			e.logger.Errorf("failed to create donation category: %v", err)
			return nil, err
		}

	case cps_const.RequestUpdateDonationCategory:
		if action.PreviousAction == nil {
			e.logger.Errorf("previous action is required for update")
			return nil, fmt.Errorf("PREVIOUS_ACTION_REQUIRED")
		}

		prevAction, err := castToBsonM(action.PreviousAction)
		if err != nil {
			e.logger.Errorf("invalid previous action format: %v", err)
			return nil, fmt.Errorf("INVALID_PREVIOUS_ACTION_FORMAT")
		}

		id, exists := prevAction["id"]
		if !exists {
			e.logger.Errorf("id not found in previous action")
			return nil, fmt.Errorf("ID_NOT_FOUND")
		}

		idStr, ok := id.(string)
		if !ok {
			e.logger.Errorf("invalid id format")
			return nil, fmt.Errorf("INVALID_ID_FORMAT")
		}

		var donationCPS *dto.DonationCategoryCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &donationCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation Category: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		updateRequest := dto.DonationCategoryRequest{
			CategoryName: donationCPS.CategoryName,
			Icon:         nil,
		}

		_, err = e.Repository.UpdateDonationCategory(ctx, idStr, updateRequest)
		if err != nil {
			e.logger.Errorf("failed to update donation category: %v", err)
			return nil, err
		}

	default:
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	return action, nil
}
