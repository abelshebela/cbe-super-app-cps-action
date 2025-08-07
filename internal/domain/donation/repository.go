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
}
