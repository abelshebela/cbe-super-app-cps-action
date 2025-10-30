package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"time"
)

func DonationNameExists(ctx context.Context, categoryName string, donationCategoryRepo storage.DonationCategoryRepository) (bool, error) {
	// Use the repository method to find donation categories
	categories, err := donationCategoryRepo.FindAllWithPagination(ctx, types.Filter{
		Search:  categoryName,
		Page:    1,
		PerPage: 1,
	})
	if err != nil {
		return false, errors.New(localization.ErrorDonationCategoryLookupFailed.Code)
	}

	// Check if any categories were found
	return len(categories.Data) > 0, nil
}
func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func MapToDonationCategory(categoryName, iconURL string,enabled bool) *model.DonationCategory {
	return &model.DonationCategory{
		CategoryName:   categoryName,
		Icon:           iconURL,
		IsDeleted:      false,
		Enabled: enabled,
		LastModifiedAt: time.Now(),
	}
}

func IsDataSimilar(request donation_category.DonationCategoryRequest, existing *model.DonationCategory) bool {
	// Check if category name is the same (if provided in request)
	if request.CategoryName != "" && request.CategoryName != existing.CategoryName {
		return false
	}

	if request.Icon != nil {
		return false
	}

	if request.CategoryName == "" {
		return true
	}
	return true
}
