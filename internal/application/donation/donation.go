package donation

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	cps_const "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	cps_service "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/services"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DonationAbstract interface {
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest, maker cps_entities.User) error
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest, maker cps_entities.User) error
	DeleteDonationCategory(ctx context.Context, id string, maker cps_entities.User) error

	// Donation Company methods
	CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest, maker cps_entities.User) error
	FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error)
	FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error)
	UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest, maker cps_entities.User) error
	DeleteDonationCompany(ctx context.Context, id string, maker cps_entities.User) error

	// Donation methods
	CreateDonation(ctx context.Context, donation dto.DonationRequest, maker cps_entities.User) error
	FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error)
	FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error)
	UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest, maker cps_entities.User) error
	DeleteDonation(ctx context.Context, id string, maker cps_entities.User) error
}

type DonationService interface {
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) (*dto.DonationCategoryResponse, error)
	UploadIcon(ctx context.Context, icon *multipart.FileHeader) (string, error)
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error)

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
	UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) (*dto.DonationRequest, error)
	ValidateCompanyExists(ctx context.Context, companyID string) error
	ValidateCategoryExists(ctx context.Context, categoryID string) error
	DonationTitleExists(ctx context.Context, title string) (bool, error)
}

type DonationStore struct {
	service    DonationService
	cpsService cps_service.CPSActionService
	logger     utils.Logger
}

func NewDonationApplication(service DonationService, cpsService cps_service.CPSActionService, logger utils.Logger) DonationAbstract {
	return &DonationStore{
		service:    service,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (a *DonationStore) handleCPSAction(ctx context.Context, maker cps_entities.User, requestAction cps_const.RequestAction, curData, prevData interface{}, actionType cps_const.ActionType, uniqueID string) error {
	cpsAction := a.cpsService.BuildCPSAction(ctx, cps_entities.CreateCPSRequest{
		User:          maker,
		CurData:       curData,
		PrevData:      prevData,
		RequestAction: requestAction,
		ActionStatus:  cps_const.ActionPending,
		ActionType:    actionType,
	})

	// Set the unique ID if provided
	if uniqueID != "" {
		cpsAction.UniqueID = uniqueID
	}

	_, err := a.cpsService.CreateCPSAction(ctx, cpsAction)
	if err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest, maker cps_entities.User) error {
	iconURL, err := a.service.UploadIcon(ctx, donation.Icon)
	if err != nil {
		a.logger.Errorf("failed to upload icon: %v", err)
		return err
	}

	cpsRequest := dto.DonationCategoryCPSRequest{
		CategoryName: donation.CategoryName,
		Icon:         iconURL,
	}

	if err := a.handleCPSAction(ctx, maker, cps_const.RequestCreateDonationCategory, cpsRequest, nil, cps_const.ActionCreate, ""); err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error) {
	paginatedResponse, err := a.service.FetchDonationCategory(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("failed to fetch donation categories: %v", err)
		return nil, err
	}
	return paginatedResponse, nil
}

func (a *DonationStore) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	donationCategory, err := a.service.FetchDonationCategoryByID(ctx, id)
	if err != nil {
		a.logger.Errorf("failed to fetch donation category by ID: %v", err)
		return nil, err
	}
	return donationCategory, nil
}

func (a *DonationStore) UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest, maker cps_entities.User) error {
	var iconURL string
	var err error

	if donation.Icon != nil {
		iconURL, err = a.service.UploadIcon(ctx, donation.Icon)
		if err != nil {
			a.logger.Errorf("failed to upload icon: %v", err)
			return err
		}
	}

	cpsRequest := dto.DonationCategoryCPSRequest{
		CategoryName: donation.CategoryName,
		Icon:         iconURL,
	}

	prevData := map[string]interface{}{"id": id}
	if err := a.handleCPSAction(ctx, maker, cps_const.RequestUpdateDonationCategory, cpsRequest, prevData, cps_const.ActionUpdate, ""); err != nil {
		return err
	}
	return nil
}

func (a *DonationStore) DeleteDonationCategory(ctx context.Context, id string, maker cps_entities.User) error {
	return nil
}

func (a *DonationStore) CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest, maker cps_entities.User) error {
	// Validate company name uniqueness
	exist, err := a.service.DonationCompanyNameExists(ctx, company.CompanyName)
	if err != nil {
		a.logger.Errorf("failed to check company name: %v", err)
		return err
	}
	if exist {
		return fmt.Errorf("COMPANY_NAME_ALREADY_EXISTS")
	}

	// Validate account number uniqueness
	exist, err = a.service.DonationCompanyAccountExists(ctx, company.AccountNumber)
	if err != nil {
		a.logger.Errorf("failed to check account number: %v", err)
		return err
	}
	if exist {
		return fmt.Errorf("ACCOUNT_NUMBER_ALREADY_EXISTS")
	}

	// Validate account number with external API
	if err := a.service.ValidateAccountNumber(ctx, company.AccountNumber); err != nil {
		a.logger.Errorf("failed to validate account number: %v", err)
		return err
	}

	// Validate logo is required
	if company.CompanyLogo == nil {
		return fmt.Errorf("LOGO_IS_REQUIRED")
	}

	logoURL, err := a.service.UploadLogo(ctx, company.CompanyLogo)
	if err != nil {
		a.logger.Errorf("failed to upload logo: %v", err)
		return err
	}

	cpsRequest := dto.DonationCompanyCPSRequest{
		CompanyName:   company.CompanyName,
		CompanyLogo:   logoURL,
		AccountNumber: company.AccountNumber,
	}

	if err := a.handleCPSAction(ctx, maker, cps_const.RequestCreateDonationCompany, cpsRequest, nil, cps_const.ActionCreate, ""); err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error) {
	paginatedResponse, err := a.service.FetchDonationCompany(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("failed to fetch donation companies: %v", err)
		return nil, err
	}
	return paginatedResponse, nil
}

func (a *DonationStore) FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error) {
	donationCompany, err := a.service.FetchDonationCompanyByID(ctx, id)
	if err != nil {
		a.logger.Errorf("failed to fetch donation company by ID: %v", err)
		return nil, err
	}
	return donationCompany, nil
}

func (a *DonationStore) UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest, maker cps_entities.User) error {
	// First, fetch the existing donation company to get current data
	existingCompany, err := a.service.FetchDonationCompanyByID(ctx, id)
	if err != nil {
		a.logger.Errorf("failed to fetch existing company: %v", err)
		return err
	}

	// Validate company name uniqueness if company name is being updated
	if company.CompanyName != "" && company.CompanyName != existingCompany.CompanyName {
		exist, err := a.service.DonationCompanyNameExists(ctx, company.CompanyName)
		if err != nil {
			a.logger.Errorf("failed to check company name: %v", err)
			return err
		}
		if exist {
			return fmt.Errorf("COMPANY_NAME_ALREADY_EXISTS")
		}
	}

	// Validate account number uniqueness if account number is being updated
	if company.AccountNumber != "" && company.AccountNumber != existingCompany.AccountNumber {
		exist, err := a.service.DonationCompanyAccountExists(ctx, company.AccountNumber)
		if err != nil {
			a.logger.Errorf("failed to check account number: %v", err)
			return err
		}
		if exist {
			return fmt.Errorf("ACCOUNT_NUMBER_ALREADY_EXISTS")
		}

		// Validate account number with external API
		if err := a.service.ValidateAccountNumber(ctx, company.AccountNumber); err != nil {
			a.logger.Errorf("failed to validate account number: %v", err)
			return err
		}
	}

	var logoURL string
	var err2 error

	if company.CompanyLogo != nil {
		logoURL, err2 = a.service.UploadLogo(ctx, company.CompanyLogo)
		if err2 != nil {
			a.logger.Errorf("failed to upload logo: %v", err2)
			return err2
		}
	}

	cpsRequest := dto.DonationCompanyCPSRequest{
		CompanyName:   company.CompanyName,
		CompanyLogo:   logoURL,
		AccountNumber: company.AccountNumber,
	}

	prevData := map[string]interface{}{"id": id}
	if err := a.handleCPSAction(ctx, maker, cps_const.RequestUpdateDonationCompany, cpsRequest, prevData, cps_const.ActionUpdate, id); err != nil {
		return err
	}
	return nil
}

func (a *DonationStore) DeleteDonationCompany(ctx context.Context, id string, maker cps_entities.User) error {
	return nil
}

func (a *DonationStore) CreateDonation(ctx context.Context, donation dto.DonationRequest, maker cps_entities.User) error {
	// Validate donation title uniqueness
	exist, err := a.service.DonationTitleExists(ctx, donation.Title)
	if err != nil {
		a.logger.Errorf("failed to check donation title: %v", err)
		return err
	}
	if exist {
		return fmt.Errorf("DONATION_TITLE_ALREADY_EXISTS")
	}

	// Validate company exists
	if err := a.service.ValidateCompanyExists(ctx, donation.CompanyID); err != nil {
		a.logger.Errorf("failed to validate company: %v", err)
		return err
	}

	// Validate category exists
	if err := a.service.ValidateCategoryExists(ctx, donation.CategoryID); err != nil {
		a.logger.Errorf("failed to validate category: %v", err)
		return err
	}

	// Validate at least one image is required
	if len(donation.DonationImages) == 0 {
		return fmt.Errorf("AT_LEAST_ONE_IMAGE_REQUIRED")
	}

	imageURLs, err := a.service.UploadDonationImages(ctx, donation.DonationImages)
	if err != nil {
		a.logger.Errorf("failed to upload images: %v", err)
		return err
	}

	startDate := donation.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	cpsRequest := dto.DonationCPSRequest{
		CompanyID:           donation.CompanyID,
		CategoryID:          donation.CategoryID,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      imageURLs,
		EndDate:             donation.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		StartDate:           startDate.Format("2006-01-02T15:04:05Z07:00"),
	}

	if err := a.handleCPSAction(ctx, maker, cps_const.RequestCreateDonation, cpsRequest, nil, cps_const.ActionCreate, ""); err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error) {
	paginatedResponse, err := a.service.FetchDonation(ctx, filterParams)
	if err != nil {
		a.logger.Errorf("failed to fetch donations: %v", err)
		return nil, err
	}
	return paginatedResponse, nil
}

func (a *DonationStore) FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error) {
	donation, err := a.service.FetchDonationByID(ctx, id)
	if err != nil {
		a.logger.Errorf("failed to fetch donation by ID: %v", err)
		return nil, err
	}
	return donation, nil
}

func (a *DonationStore) UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest, maker cps_entities.User) error {
	// First, fetch the existing donation to get current data
	existingDonation, err := a.service.FetchDonationByID(ctx, id)
	if err != nil {
		a.logger.Errorf("failed to fetch existing donation: %v", err)
		return err
	}

	// Validate donation title uniqueness if title is being updated
	if donation.Title != "" && donation.Title != existingDonation.Title {
		exist, err := a.service.DonationTitleExists(ctx, donation.Title)
		if err != nil {
			a.logger.Errorf("failed to check donation title: %v", err)
			return err
		}
		if exist {
			return fmt.Errorf("DONATION_TITLE_ALREADY_EXISTS")
		}
	}

	// Validate company exists if company ID is being updated
	if donation.CompanyID != "" && donation.CompanyID != existingDonation.Company.ID {
		if err := a.service.ValidateCompanyExists(ctx, donation.CompanyID); err != nil {
			a.logger.Errorf("failed to validate company: %v", err)
			return err
		}
	}

	// Validate category exists if category ID is being updated
	if donation.CategoryID != "" && donation.CategoryID != existingDonation.Category.ID {
		if err := a.service.ValidateCategoryExists(ctx, donation.CategoryID); err != nil {
			a.logger.Errorf("failed to validate category: %v", err)
			return err
		}
	}

	var imageURLs []string
	var err2 error

	if len(donation.DonationImages) > 0 {
		imageURLs, err2 = a.service.UploadDonationImages(ctx, donation.DonationImages)
		if err2 != nil {
			a.logger.Errorf("failed to upload images: %v", err2)
			return err2
		}
	}

	startDate := donation.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	// Create CPS request with only the changed data
	cpsRequest := dto.DonationCPSRequest{
		CompanyID:           donation.CompanyID,
		CategoryID:          donation.CategoryID,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      imageURLs,
		EndDate:             donation.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		StartDate:           startDate.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Create previous data from existing donation
	prevData := map[string]interface{}{
		"id":                   id,
		"company_id":           existingDonation.Company.ID,
		"category_id":          existingDonation.Category.ID,
		"title":                existingDonation.Title,
		"is_featured":          existingDonation.IsFeatured,
		"target":               existingDonation.Target,
		"donation_description": existingDonation.DonationDescription,
		"donation_images":      existingDonation.DonationImages,
		"end_date":             existingDonation.EndDate,
		"start_date":           existingDonation.StartDate,
	}

	// Store donation ID in unique ID field for easy retrieval
	uniqueID := fmt.Sprintf("DONATION-%s", id)

	if err := a.handleCPSAction(ctx, maker, cps_const.RequestUpdateDonation, cpsRequest, prevData, cps_const.ActionUpdate, uniqueID); err != nil {
		a.logger.Errorf("failed to create CPS action: %v", err)
		return err
	}
	return nil
}

func (a *DonationStore) DeleteDonation(ctx context.Context, id string, maker cps_entities.User) error {
	return nil
}
