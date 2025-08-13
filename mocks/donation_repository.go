package mocks

import (
	"context"

	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"github.com/stretchr/testify/mock"
)

type DonationRepository struct {
	mock.Mock
}

// Donation Category methods
func (m *DonationRepository) DonationNameExists(ctx context.Context, CategoryName string) (bool, error) {
	args := m.Called(ctx, CategoryName)
	return args.Bool(0), args.Error(1)
}

func (m *DonationRepository) CreateDonationCategory(ctx context.Context, donation dto.DonationCategoryRequest) error {
	args := m.Called(ctx, donation)
	return args.Error(0)
}

func (m *DonationRepository) CreateDonationCategoryWithURL(ctx context.Context, categoryName, iconURL string) error {
	args := m.Called(ctx, categoryName, iconURL)
	return args.Error(0)
}

func (m *DonationRepository) FetchDonationCategory(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse], error) {
	args := m.Called(ctx, filterParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common_util.PaginatedResponse[[]*dto.DonationCategoryListResponse]), args.Error(1)
}

func (m *DonationRepository) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationCategoryListResponse), args.Error(1)
}

func (m *DonationRepository) UpdateDonationCategory(ctx context.Context, id string, donation dto.DonationCategoryRequest) (*dto.DonationCategoryRequest, error) {
	args := m.Called(ctx, id, donation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationCategoryRequest), args.Error(1)
}

// Donation Company methods
func (m *DonationRepository) DonationCompanyNameExists(ctx context.Context, companyName string) (bool, error) {
	args := m.Called(ctx, companyName)
	return args.Bool(0), args.Error(1)
}

func (m *DonationRepository) DonationCompanyAccountExists(ctx context.Context, accountNumber string) (bool, error) {
	args := m.Called(ctx, accountNumber)
	return args.Bool(0), args.Error(1)
}

func (m *DonationRepository) CreateDonationCompany(ctx context.Context, company dto.DonationCompanyRequest) error {
	args := m.Called(ctx, company)
	return args.Error(0)
}

func (m *DonationRepository) CreateDonationCompanyWithURL(ctx context.Context, companyName, logoURL, accountNumber string) error {
	args := m.Called(ctx, companyName, logoURL, accountNumber)
	return args.Error(0)
}

func (m *DonationRepository) FetchDonationCompany(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse], error) {
	args := m.Called(ctx, filterParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common_util.PaginatedResponse[[]*dto.DonationCompanyListResponse]), args.Error(1)
}

func (m *DonationRepository) FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationCompanyListResponse), args.Error(1)
}

func (m *DonationRepository) UpdateDonationCompany(ctx context.Context, id string, company dto.DonationCompanyRequest) (*dto.DonationCompanyRequest, error) {
	args := m.Called(ctx, id, company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationCompanyRequest), args.Error(1)
}

func (m *DonationRepository) UpdateDonationCompanyWithLogoURL(ctx context.Context, id string, company dto.DonationCompanyRequest, logoURL string) (*dto.DonationCompanyRequest, error) {
	args := m.Called(ctx, id, company, logoURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationCompanyRequest), args.Error(1)
}

// Donation methods
func (m *DonationRepository) DonationTitleExists(ctx context.Context, title string) (bool, error) {
	args := m.Called(ctx, title)
	return args.Bool(0), args.Error(1)
}

func (m *DonationRepository) DonationCompanyExists(ctx context.Context, companyID string) (bool, error) {
	args := m.Called(ctx, companyID)
	return args.Bool(0), args.Error(1)
}

func (m *DonationRepository) DonationCategoryExists(ctx context.Context, categoryID string) (bool, error) {
	args := m.Called(ctx, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *DonationRepository) CreateDonation(ctx context.Context, donation dto.DonationRequest) error {
	args := m.Called(ctx, donation)
	return args.Error(0)
}

func (m *DonationRepository) CreateDonationWithURLs(ctx context.Context, donation dto.DonationCPSRequest) error {
	args := m.Called(ctx, donation)
	return args.Error(0)
}

func (m *DonationRepository) FetchDonation(ctx context.Context, filterParams *constant.MongoFilter) (*common_util.PaginatedResponse[[]*dto.DonationListResponse], error) {
	args := m.Called(ctx, filterParams)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*common_util.PaginatedResponse[[]*dto.DonationListResponse]), args.Error(1)
}

func (m *DonationRepository) FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationListResponse), args.Error(1)
}

func (m *DonationRepository) FetchDonationByCode(ctx context.Context, donationCode string) (*dto.DonationListResponse, error) {
	args := m.Called(ctx, donationCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationListResponse), args.Error(1)
}

func (m *DonationRepository) UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) (*dto.DonationRequest, error) {
	args := m.Called(ctx, id, donation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationRequest), args.Error(1)
}

func (m *DonationRepository) UpdateDonationWithImageURLs(ctx context.Context, id string, donation dto.DonationRequest, imageURLs []string) (*dto.DonationRequest, error) {
	args := m.Called(ctx, id, donation, imageURLs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.DonationRequest), args.Error(1)
}
