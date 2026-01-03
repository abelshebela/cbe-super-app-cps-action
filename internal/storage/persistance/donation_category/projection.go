package donation_category

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DonationCategoryMapper maps DonationCategory model to BSON for database operations
func DonationCategoryMapper(data model.DonationCategory) bson.M {
	result := bson.M{}
	if data.CategoryName != "" {
		result["category_name"] = data.CategoryName
	}
	if data.Icon != "" {
		result["donation_icon"] = data.Icon
	}
	result["enabled"] = data.Enabled
	result["last_modified_at"] = data.LastModifiedAt

	return result
}

// MapToDonationCategoryListResponse maps model to DTO response
func MapToDonationCategoryListResponse(category *model.DonationCategory) *donation_category.DonationCategoryListResponse {
	return &donation_category.DonationCategoryListResponse{
		ID:             category.ID.Hex(),
		CategoryName:   category.CategoryName,
		Icon:           category.Icon,
		Enabled:        category.Enabled,
		IsDeleted:      category.IsDeleted,
		CreatedAt:      category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastModifiedAt: category.LastModifiedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// MapToDonationCategoryListResponses maps slice of models to slice of DTO responses
func MapToDonationCategoryListResponses(categories []model.DonationCategory) []donation_category.DonationCategoryListResponse {
	var responses []donation_category.DonationCategoryListResponse
	for i := range categories {
		responses = append(responses, *MapToDonationCategoryListResponse(&categories[i]))
	}
	return responses
}
