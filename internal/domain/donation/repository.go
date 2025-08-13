package donation

import (
	"context"

	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type DonationRepository interface {
	DonationNameExists(ctx context.Context, CategoryName string) (bool, error)
	CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) error
	CreateDonationCategoryWithURL(ctx context.Context, categoryName, iconURL string) error
	FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error)

	// Donation Company methods
	DonationCompanyNameExists(ctx context.Context, companyName string) (bool, error)
	DonationCompanyAccountExists(ctx context.Context, accountNumber string) (bool, error)
	CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest) error
	CreateDonationCompanyWithURL(ctx context.Context, companyName, logoURL, accountNumber string) error
	FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error)
	FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error)
	UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest) (*dto.DonationCompanyRequest, error)
	UpdateDonationCompanyWithLogoURL(ctx context.Context, id string, company dto.DonationCompanyRequest, logoURL string) (*dto.DonationCompanyRequest, error)

	// Donation methods
	DonationTitleExists(ctx context.Context, title string) (bool, error)
	DonationCompanyExists(ctx context.Context, companyID string) (bool, error)
	DonationCategoryExists(ctx context.Context, categoryID string) (bool, error)
	CreateDonation(ctx context.Context, donation dto.DonationRequest) error
	CreateDonationWithURLs(ctx context.Context, donation dto.DonationCPSRequest) error
	FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error)
	FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error)
	FetchDonationByCode(ctx context.Context, donationCode string) (*dto.DonationListResponse, error)
	UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) (*dto.DonationRequest, error)
	UpdateDonationWithImageURLs(ctx context.Context, id string, donation dto.DonationRequest, imageURLs []string) (*dto.DonationRequest, error)

	UpdateDonationImage(ctx context.Context, donationID, imageID, photoURL string) error
	DeleteDonationImage(ctx context.Context, donationID, imageID string) error
	AddDonationImage(ctx context.Context, donationID string, image dto.DonationImage) error
}
