package core

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"time"

	donation_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/donation"
)

func DonationNameExists(
	ctx context.Context,
	categoryName string,
	donationCategoryRepo storage.DonationCategoryRepository,
) (bool, error) {

	category, err := donationCategoryRepo.FindByName(ctx, categoryName)
	if err != nil {
		// Check if the error message is the "not found" error code
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return false, nil
		}
		return false, errors.New(localization.ErrorDonationCategoryLookupFailed.Code)
	}

	return category != nil && category.ID != "", nil
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func MapToDonationCategory(categoryName, iconURL string, enabled bool) *donation_model.DonationCategory {
	return &donation_model.DonationCategory{
		CategoryName:   categoryName,
		Icon:           iconURL,
		IsDeleted:      false,
		Enabled:        enabled,
		LastModifiedAt: time.Now(),
	}
}

func IsDataSimilar(request donation_category.DonationCategoryRequest, existing *donation_model.DonationCategory) bool {
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
