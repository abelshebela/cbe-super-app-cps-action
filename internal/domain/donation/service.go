package donation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	"go.mongodb.org/mongo-driver/v2/bson"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"mime/multipart"

	accountlookup "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/core_banking_calls"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/google/uuid"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DonationService interface {
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) (*dto.DonationCategoryResponse, error)
	UploadIcon(ctx context.Context, icon *multipart.FileHeader) (string, error)
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error)
	Authorize(ctx context.Context, action *entities.CPSAction) (*entities.CPSAction, error)

	// Donation Company methods
	CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest) (*dto.DonationCompanyResponse, error)
	UploadLogo(ctx context.Context, logo *multipart.FileHeader) (string, error)
	FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error)
	FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error)
	UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest) (*dto.DonationCompanyRequest, error)
	ValidateAccountNumber(ctx context.Context, accountNumber string) error
	DonationCompanyNameExists(ctx context.Context, companyName string) (bool, error)
	DonationCompanyAccountExists(ctx context.Context, accountNumber string) (bool, error)

	// Donation methods
	CreateDonation(ctx context.Context, donation dto.DonationRequest) (*dto.DonationResponse, error)
	UploadDonationImages(ctx context.Context, images []*multipart.FileHeader) ([]string, error)
	FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error)
	FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error)
	FetchDonationByCode(ctx context.Context, donationCode string) (*dto.DonationListResponse, error)
	UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) (*dto.DonationRequest, error)
	ValidateCompanyExists(ctx context.Context, companyID string) error
	ValidateCategoryExists(ctx context.Context, categoryID string) error
	DonationTitleExists(ctx context.Context, title string) (bool, error)
	GenerateImageID() string
}

type Service struct {
	Repository    DonationRepository
	logger        shared_utils.Logger
	minio         config.MinioClientInterface
	bucketName    string
	cfg           *config.VaultConfig
	accountLookup accountlookup.Account
}

func NewDonationService(repository DonationRepository,
	minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	logger shared_utils.Logger,
	accountLookup accountlookup.Account,
) DonationService {
	return &Service{
		Repository:    repository,
		logger:        logger,
		minio:         minio,
		bucketName:    bucketName,
		cfg:           cfg,
		accountLookup: accountLookup,
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
		// Extract donation category ID from unique ID (format: "{id}")
		categoryID := ""
		if action.UniqueID != "" {
			categoryID = action.UniqueID
		} else {
			// Fallback to previous action ID if unique ID is not set
			if prevData, ok := action.PreviousAction.(map[string]interface{}); ok {
				if id, exists := prevData["id"].(string); exists {
					categoryID = id
				}
			}
		}

		if categoryID == "" {
			e.logger.Errorf("donation category ID not found in unique ID or previous action")
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		var donationCPS *dto.DonationCategoryCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &donationCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation Category: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		// Build patch request with only non-empty fields to avoid overwriting existing data
		updateRequest := dto.DonationCategoryRequest{}

		// Only include fields that are actually being updated (not empty)
		if donationCPS.CategoryName != "" {
			updateRequest.CategoryName = donationCPS.CategoryName
		}
		// Icon is handled separately and will be nil in the update request

		_, err := e.Repository.UpdateDonationCategory(ctx, categoryID, updateRequest)
		if err != nil {
			e.logger.Errorf("failed to update donation category: %v", err)
			return nil, err
		}

	case cps_const.RequestCreateDonationCompany:
		var companyCPS *dto.DonationCompanyCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &companyCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation company: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		if companyCPS.CompanyLogo == "" {
			e.logger.Errorf("logo URL is empty in CPS request")
			return nil, fmt.Errorf("LOGO_URL_MISSING")
		}

		if err := e.ValidateAccountNumber(ctx, companyCPS.AccountNumber); err != nil {
			return nil, err
		}

		err = e.Repository.CreateDonationCompanyWithURL(ctx, companyCPS.CompanyName, companyCPS.CompanyLogo, companyCPS.AccountNumber)
		if err != nil {
			e.logger.Errorf("failed to create donation company: %v", err)
			return nil, err
		}

	case cps_const.RequestUpdateDonationCompany:
		// Extract donation company ID from unique ID (format: "{id}")
		companyID := ""
		if action.UniqueID != "" {
			companyID = action.UniqueID
		} else {
			// Fallback to previous action ID if unique ID is not set
			if prevData, ok := action.PreviousAction.(map[string]interface{}); ok {
				if id, exists := prevData["id"].(string); exists {
					companyID = id
				}
			}
		}

		if companyID == "" {
			e.logger.Errorf("donation company ID not found in unique ID or previous action")
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		var companyCPS *dto.DonationCompanyCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &companyCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation company: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		if companyCPS.AccountNumber != "" {
			if err := e.ValidateAccountNumber(ctx, companyCPS.AccountNumber); err != nil {
				return nil, err
			}
		}

		// Build patch request with only non-empty fields to avoid overwriting existing data
		updateRequest := dto.DonationCompanyRequest{}

		// Only include fields that are actually being updated (not empty)
		if companyCPS.CompanyName != "" {
			updateRequest.CompanyName = companyCPS.CompanyName
		}
		if companyCPS.AccountNumber != "" {
			updateRequest.AccountNumber = companyCPS.AccountNumber
		}
		// Logo is handled separately and will be nil in the update request

		_, err := e.Repository.UpdateDonationCompany(ctx, companyID, updateRequest)
		if err != nil {
			e.logger.Errorf("failed to update donation company: %v", err)
			return nil, err
		}

	case cps_const.RequestCreateDonation:
		var donationCPS *dto.DonationCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &donationCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		if len(donationCPS.DonationImages) == 0 {
			e.logger.Errorf("donation images are empty in CPS request")
			return nil, fmt.Errorf("DONATION_IMAGES_MISSING")
		}

		if err := e.ValidateCompanyExists(ctx, donationCPS.CompanyID); err != nil {
			return nil, err
		}

		if err := e.ValidateCategoryExists(ctx, donationCPS.CategoryID); err != nil {
			return nil, err
		}

		err = e.Repository.CreateDonationWithURLs(ctx, *donationCPS)
		if err != nil {
			e.logger.Errorf("failed to create donation: %v", err)
			return nil, err
		}

	case cps_const.RequestUpdateDonation:
		var donationCPS dto.DonationCPSRequest
		bindErr := common_util.BindAction(action.CurrentAction, &donationCPS)
		if bindErr != nil {
			e.logger.Errorf("failed to bind current action to donation: %v", bindErr)
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		// Extract donation ID from unique ID (format: "{id}")
		donationID := ""
		if action.UniqueID != "" {
			donationID = action.UniqueID
		} else {
			// Fallback to previous action ID if unique ID is not set
			if prevData, ok := action.PreviousAction.(map[string]interface{}); ok {
				if id, exists := prevData["id"].(string); exists {
					donationID = id
				}
			}
		}

		if donationID == "" {
			e.logger.Errorf("donation ID not found in unique ID or previous action")
			return nil, fmt.Errorf(common_util.InvalidActionData)
		}

		if donationCPS.CompanyID != "" {
			if err := e.ValidateCompanyExists(ctx, donationCPS.CompanyID); err != nil {
				return nil, err
			}
		}

		if donationCPS.CategoryID != "" {
			if err := e.ValidateCategoryExists(ctx, donationCPS.CategoryID); err != nil {
				return nil, err
			}
		}

		// Create update request with only changed fields
		// Build patch request with only non-empty fields to avoid overwriting existing data
		updateRequest := dto.DonationRequest{}

		// Only include fields that are actually being updated (not empty)
		if donationCPS.CompanyID != "" {
			updateRequest.CompanyID = donationCPS.CompanyID
		}
		if donationCPS.CategoryID != "" {
			updateRequest.CategoryID = donationCPS.CategoryID
		}
		if donationCPS.Title != "" {
			updateRequest.Title = donationCPS.Title
		}
		// Always include boolean fields as they can be false
		updateRequest.IsFeatured = donationCPS.IsFeatured
		if donationCPS.Target > 0 {
			updateRequest.Target = donationCPS.Target
		}
		if donationCPS.DonationDescription != "" {
			updateRequest.DonationDescription = donationCPS.DonationDescription
		}
		// Images are handled separately
		updateRequest.DonationImages = nil

		// Parse dates only if provided
		if donationCPS.StartDate != "" {
			if startDate, err := time.Parse("2006-01-02T15:04:05Z07:00", donationCPS.StartDate); err == nil {
				updateRequest.StartDate = startDate
			}
		}
		if donationCPS.EndDate != "" {
			if endDate, err := time.Parse("2006-01-02T15:04:05Z07:00", donationCPS.EndDate); err == nil {
				updateRequest.EndDate = endDate
			}
		}

		_, err := e.Repository.UpdateDonation(ctx, donationID, updateRequest)
		if err != nil {
			e.logger.Errorf("failed to update donation: %v", err)
			return nil, fmt.Errorf("failed to update donation: %v", err)
		}

		// Handle image operations if present
		if len(donationCPS.DonationImages) > 0 {
			var previousImages []dto.DonationImage

			// fmt.Printf("PreviousAction type: %T\n", action.PreviousAction)
			// fmt.Printf("PreviousAction value: %+v\n", action.PreviousAction)

			if prevData, ok := action.PreviousAction.(map[string]interface{}); ok {
				// fmt.Println("Previous action data:", prevData)
				if imagesData, exists := prevData["donation_images"]; exists {
					if imagesBytes, err := json.Marshal(imagesData); err == nil {
						if err := json.Unmarshal(imagesBytes, &previousImages); err == nil {
							fmt.Println("Previous images parsed:", previousImages)
						}
					}
					// fmt.Println("imaaaaaaaaaaaggggeeeeeesssss:", imagesData)
				}
			} else if bsonDoc, ok := action.PreviousAction.(bson.D); ok {
				// fmt.Println("Previous action is bson.D, converting...")
				prevData := make(map[string]interface{})
				for _, elem := range bsonDoc {
					prevData[elem.Key] = elem.Value
				}
				// fmt.Println("Converted previous action data:", prevData)
				if imagesData, exists := prevData["donation_images"]; exists {
					if imagesBytes, err := json.Marshal(imagesData); err == nil {
						if err := json.Unmarshal(imagesBytes, &previousImages); err == nil {
							fmt.Println("Previous images parsed:", previousImages)
						}
					}
					// fmt.Println("imaaaaaaaaaaaggggeeeeeesssss:", imagesData)
				}
			} else {
				fmt.Println("Failed to cast PreviousAction to expected types")
			}

			for _, image := range donationCPS.DonationImages {
				if image.ID != "" && image.PhotoURL != "" {
					createdAt, err := time.Parse("2006-01-02T15:04:05Z07:00", image.CreatedAt)
					if err != nil {
						createdAt = time.Now()
					}

					imageExists := false
					for _, prevImage := range previousImages {
						if prevImage.ID == image.ID {
							imageExists = true
							break
						}
					}

					if imageExists {
						if err := e.Repository.UpdateDonationImage(ctx, donationID, image.ID, image.PhotoURL); err != nil {
							e.logger.Errorf("failed to update donation image: %v", err)
							return nil, fmt.Errorf("failed to update donation image: %v", err)
						}
					} else {
						donationImage := dto.DonationImage{
							ID:        image.ID,
							PhotoURL:  image.PhotoURL,
							CreatedAt: createdAt.Format("2006-01-02T15:04:05Z07:00"),
						}
						if err := e.Repository.AddDonationImage(ctx, donationID, donationImage); err != nil {
							e.logger.Errorf("failed to add donation image: %v", err)
							return nil, fmt.Errorf("failed to add donation image: %v", err)
						}
					}
				}
			}
		}

		// Handle image deletion if present
		if donationCPS.ImageIDToDelete != "" {
			if err := e.Repository.DeleteDonationImage(ctx, donationID, donationCPS.ImageIDToDelete); err != nil {
				e.logger.Errorf("failed to delete donation image: %v", err)
				return nil, fmt.Errorf("failed to delete donation image: %v", err)
			}
		}

		e.logger.Infof("Successfully approved and authorized CPS action: %s", action.ID)
		return action, nil

	default:
		return nil, fmt.Errorf(common_util.ErrUnsupported)
	}

	return action, nil
}

func (e *Service) UploadLogo(ctx context.Context, logo *multipart.FileHeader) (string, error) {
	if logo == nil {
		return "", fmt.Errorf("LOGO_IS_REQUIRED")
	}

	if e.minio == nil {
		e.logger.Errorf("MinIO client is nil")
		return "", fmt.Errorf("MINIO_CLIENT_NOT_CONFIGURED")
	}

	if e.cfg == nil {
		e.logger.Errorf("Configuration is nil")
		return "", fmt.Errorf("CONFIGURATION_NOT_LOADED")
	}

	logoURL, err := common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, logo, "donation_company_logo", e.cfg.MinioEndPoint, e.logger)
	if err != nil {
		e.logger.Errorf("Failed to upload logo to MinIO bucket '%s': %v", e.bucketName, err)
		if strings.Contains(err.Error(), "time") || strings.Contains(err.Error(), "server") || strings.Contains(err.Error(), "difference") {
			return "", fmt.Errorf("MINIO_TIME_SYNC_ERROR")
		}
		if strings.Contains(err.Error(), "bucket") || strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "not valid") {
			return "", fmt.Errorf("MINIO_BUCKET_ERROR")
		}
		if strings.Contains(err.Error(), "UNHANDLED_SERVER_ERROR") {
			return "", fmt.Errorf("MINIO_TIME_SYNC_ERROR")
		}
		return "", fmt.Errorf("FAILED_TO_UPLOAD_LOGO")
	}

	e.logger.Infof("Successfully uploaded donation company logo with URL: %s", logoURL)
	return logoURL, nil
}

func (e *Service) CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest) (*dto.DonationCompanyResponse, error) {
	exist, err := e.Repository.DonationCompanyNameExists(ctx, company.CompanyName)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("COMPANY_NAME_ALREADY_EXISTS")
	}

	exist, err = e.Repository.DonationCompanyAccountExists(ctx, company.AccountNumber)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("ACCOUNT_NUMBER_ALREADY_EXISTS")
	}

	if err := e.ValidateAccountNumber(ctx, company.AccountNumber); err != nil {
		return nil, err
	}

	if company.CompanyLogo == nil {
		return nil, fmt.Errorf("LOGO_IS_REQUIRED")
	}

	logoURL, err := e.UploadLogo(ctx, company.CompanyLogo)
	if err != nil {
		return nil, err
	}

	result := dto.DonationCompanyResponse{
		CompanyName:   company.CompanyName,
		CompanyLogo:   logoURL,
		AccountNumber: company.AccountNumber,
	}
	return &result, nil
}

func (e *Service) UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest) (*dto.DonationCompanyRequest, error) {
	e.logger.Infof("Updating company ", "id", id)
	_, err := e.Repository.FetchDonationCompanyByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to fetch company", "id", id, "error", err)
		return nil, err
	}

	var logoURL string
	if company.CompanyLogo != nil {
		if e.minio == nil {
			e.logger.Errorf("MinIO client is nil")
			return nil, fmt.Errorf("MINIO_CLIENT_NOT_CONFIGURED")
		}

		if e.cfg == nil {
			e.logger.Errorf("Configuration is nil")
			return nil, fmt.Errorf("CONFIGURATION_NOT_LOADED")
		}

		e.logger.Infof("Starting updated donation company logo upload to MinIO bucket: %s", e.bucketName)
		e.logger.Infof("MinIO endpoint: %s", e.cfg.MinioEndPoint)

		logoURL, err = common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, company.CompanyLogo, "donation_company_logo", e.cfg.MinioEndPoint, e.logger)
		if err != nil {
			e.logger.Errorf("Failed to upload logo to MinIO bucket '%s': %v", e.bucketName, err)
			if strings.Contains(err.Error(), "time") || strings.Contains(err.Error(), "server") || strings.Contains(err.Error(), "difference") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			if strings.Contains(err.Error(), "bucket") || strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "not valid") {
				return nil, fmt.Errorf("MINIO_BUCKET_ERROR")
			}
			if strings.Contains(err.Error(), "UNHANDLED_SERVER_ERROR") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			return nil, fmt.Errorf("FAILED_TO_UPLOAD_LOGO")
		}
		e.logger.Infof("Successfully uploaded updated donation company logo with URL: %s", logoURL)
	}

	if company.CompanyName != "" {
		exist, err := e.Repository.DonationCompanyNameExists(ctx, company.CompanyName)
		if err != nil {
			return nil, err
		}

		if exist {
			return nil, fmt.Errorf("COMPANY_NAME_ALREADY_EXISTS")
		}
	}

	if company.AccountNumber != "" {
		exist, err := e.Repository.DonationCompanyAccountExists(ctx, company.AccountNumber)
		if err != nil {
			return nil, err
		}

		if exist {
			return nil, fmt.Errorf("ACCOUNT_NUMBER_ALREADY_EXISTS")
		}

		if err := e.ValidateAccountNumber(ctx, company.AccountNumber); err != nil {
			return nil, err
		}
	}

	// Create a new request with the logo URL if uploaded
	updateRequest := company
	if logoURL != "" {
		// We need to pass the logo URL to the repository
		// Since the repository expects DonationCompanyRequest, we'll clear the file and pass URL separately
		updateRequest.CompanyLogo = nil
	}

	updatedCompany, err := e.Repository.UpdateDonationCompanyWithLogoURL(ctx, id, updateRequest, logoURL)
	if err != nil {
		e.logger.Errorf("Failed to update donation company in database: %v", err)
		return nil, err
	}

	return updatedCompany, nil
}

func (e *Service) FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error) {
	return e.Repository.FetchDonationCompany(ctx, filterParams)
}

func (e *Service) FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error) {
	return e.Repository.FetchDonationCompanyByID(ctx, id)
}

func (e *Service) ValidateAccountNumber(ctx context.Context, accountNumber string) error {
	if accountNumber == "" {
		return fmt.Errorf("ACCOUNT_NUMBER_IS_REQUIRED")
	}

	accountReq := model.AccountLookUpRequest{
		AccountNumber: accountNumber,
	}

	_, err := e.accountLookup.LookupAccountByAccountNumber(ctx, accountReq)
	if err != nil {
		e.logger.Errorf("account validation failed for account %s: %v", accountNumber, err)
		return fmt.Errorf("ACCOUNT_NOT_FOUND")
	}

	return nil
}

func (e *Service) UploadDonationImages(ctx context.Context, images []*multipart.FileHeader) ([]string, error) {
	if len(images) == 0 {
		return nil, fmt.Errorf("AT_LEAST_ONE_IMAGE_REQUIRED")
	}

	if e.minio == nil {
		e.logger.Errorf("MinIO client is nil")
		return nil, fmt.Errorf("MINIO_CLIENT_NOT_CONFIGURED")
	}

	if e.cfg == nil {
		e.logger.Errorf("Configuration is nil")
		return nil, fmt.Errorf("CONFIGURATION_NOT_LOADED")
	}

	var imageURLs []string
	for i, image := range images {
		if image == nil {
			continue
		}

		imageURL, err := common_util.UploadFileToMinio(ctx, e.minio, e.bucketName, image, "donation_images", e.cfg.MinioEndPoint, e.logger)
		if err != nil {
			e.logger.Errorf("Failed to upload donation image %d to MinIO bucket '%s': %v", i+1, e.bucketName, err)
			if strings.Contains(err.Error(), "time") || strings.Contains(err.Error(), "server") || strings.Contains(err.Error(), "difference") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			if strings.Contains(err.Error(), "bucket") || strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "not valid") {
				return nil, fmt.Errorf("MINIO_BUCKET_ERROR")
			}
			if strings.Contains(err.Error(), "UNHANDLED_SERVER_ERROR") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			return nil, fmt.Errorf("FAILED_TO_UPLOAD_IMAGE_%d", i+1)
		}

		imageURLs = append(imageURLs, imageURL)
		e.logger.Infof("Successfully uploaded donation image %d with URL: %s", i+1, imageURL)
	}

	return imageURLs, nil
}

func (e *Service) CreateDonation(ctx context.Context, donation dto.DonationRequest) (*dto.DonationResponse, error) {
	exist, err := e.Repository.DonationTitleExists(ctx, donation.Title)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, fmt.Errorf("DONATION_TITLE_ALREADY_EXISTS")
	}

	if err := e.ValidateCompanyExists(ctx, donation.CompanyID); err != nil {
		return nil, err
	}

	if err := e.ValidateCategoryExists(ctx, donation.CategoryID); err != nil {
		return nil, err
	}

	if len(donation.DonationImages) == 0 {
		return nil, fmt.Errorf("AT_LEAST_ONE_IMAGE_REQUIRED")
	}

	imageURLs, err := e.UploadDonationImages(ctx, donation.DonationImages)
	if err != nil {
		return nil, err
	}

	startDate := donation.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	donationImages := make([]dto.DonationImage, len(imageURLs))
	for i, url := range imageURLs {
		donationImages[i] = dto.DonationImage{
			ID:        e.GenerateImageID(),
			PhotoURL:  url,
			CreatedAt: time.Now().Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	result := dto.DonationResponse{
		CompanyID:           donation.CompanyID,
		CategoryID:          donation.CategoryID,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      donationImages,
		EndDate:             donation.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		StartDate:           startDate.Format("2006-01-02T15:04:05Z07:00"),
	}
	return &result, nil
}

func (e *Service) UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) (*dto.DonationRequest, error) {
	e.logger.Infof("Updating donation ", "id", id)
	_, err := e.Repository.FetchDonationByID(ctx, id)
	if err != nil {
		e.logger.Errorf("Failed to fetch donation", "id", id, "error", err)
		return nil, err
	}

	var imageURLs []string
	if len(donation.DonationImages) > 0 {
		if e.minio == nil {
			e.logger.Errorf("MinIO client is nil")
			return nil, fmt.Errorf("MINIO_CLIENT_NOT_CONFIGURED")
		}

		if e.cfg == nil {
			e.logger.Errorf("Configuration is nil")
			return nil, fmt.Errorf("CONFIGURATION_NOT_LOADED")
		}

		e.logger.Infof("Starting updated donation images upload to MinIO bucket: %s", e.bucketName)
		e.logger.Infof("MinIO endpoint: %s", e.cfg.MinioEndPoint)

		imageURLs, err = e.UploadDonationImages(ctx, donation.DonationImages)
		if err != nil {
			e.logger.Errorf("Failed to upload images to MinIO bucket '%s': %v", e.bucketName, err)
			if strings.Contains(err.Error(), "time") || strings.Contains(err.Error(), "server") || strings.Contains(err.Error(), "difference") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			if strings.Contains(err.Error(), "bucket") || strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "not valid") {
				return nil, fmt.Errorf("MINIO_BUCKET_ERROR")
			}
			if strings.Contains(err.Error(), "UNHANDLED_SERVER_ERROR") {
				return nil, fmt.Errorf("MINIO_TIME_SYNC_ERROR")
			}
			return nil, fmt.Errorf("FAILED_TO_UPLOAD_IMAGES")
		}
		e.logger.Infof("Successfully uploaded updated donation images with URLs: %v", imageURLs)

		donation.DonationImages = nil
	}

	if donation.Title != "" {
		exist, err := e.Repository.DonationTitleExists(ctx, donation.Title)
		if err != nil {
			return nil, err
		}

		if exist {
			return nil, fmt.Errorf("DONATION_TITLE_ALREADY_EXISTS")
		}
	}

	if donation.CompanyID != "" {
		if err := e.ValidateCompanyExists(ctx, donation.CompanyID); err != nil {
			return nil, err
		}
	}

	if donation.CategoryID != "" {
		if err := e.ValidateCategoryExists(ctx, donation.CategoryID); err != nil {
			return nil, err
		}
	}

	updatedDonation, err := e.Repository.UpdateDonationWithImageURLs(ctx, id, donation, imageURLs)
	if err != nil {
		e.logger.Errorf("Failed to update donation in database: %v", err)
		return nil, err
	}

	return updatedDonation, nil
}

func (e *Service) FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error) {
	return e.Repository.FetchDonation(ctx, filterParams)
}

func (e *Service) FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error) {
	return e.Repository.FetchDonationByID(ctx, id)
}

func (e *Service) FetchDonationByCode(ctx context.Context, donationCode string) (*dto.DonationListResponse, error) {
	return e.Repository.FetchDonationByCode(ctx, donationCode)
}

func (e *Service) ValidateCompanyExists(ctx context.Context, companyID string) error {
	exist, err := e.Repository.DonationCompanyExists(ctx, companyID)
	if err != nil {
		// Check if it's an invalid ID format error
		if err.Error() == "INVALID_ID_FORMAT" {
			return fmt.Errorf("INVALID_ID_FORMAT")
		}
		return fmt.Errorf("COMPANY_LOOKUP_FAILED")
	}
	if !exist {
		return fmt.Errorf("COMPANY_NOT_FOUND")
	}
	return nil
}

func (e *Service) ValidateCategoryExists(ctx context.Context, categoryID string) error {
	exist, err := e.Repository.DonationCategoryExists(ctx, categoryID)
	if err != nil {
		// Check if it's an invalid ID format error
		if err.Error() == "INVALID_ID_FORMAT" {
			return fmt.Errorf("INVALID_ID_FORMAT")
		}
		return fmt.Errorf("CATEGORY_LOOKUP_FAILED")
	}
	if !exist {
		return fmt.Errorf("CATEGORY_NOT_FOUND")
	}
	return nil
}

func (e *Service) DonationTitleExists(ctx context.Context, title string) (bool, error) {
	exist, err := e.Repository.DonationTitleExists(ctx, title)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (e *Service) DonationCompanyNameExists(ctx context.Context, companyName string) (bool, error) {
	exist, err := e.Repository.DonationCompanyNameExists(ctx, companyName)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (e *Service) DonationCompanyAccountExists(ctx context.Context, accountNumber string) (bool, error) {
	exist, err := e.Repository.DonationCompanyAccountExists(ctx, accountNumber)
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (e *Service) GenerateImageID() string {
	return uuid.New().String()
}
