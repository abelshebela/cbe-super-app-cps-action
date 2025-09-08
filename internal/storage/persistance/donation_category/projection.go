package donation_category

import (
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/model"

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

	return result
}

// MapToDonationCategoryListResponse maps model to DTO response
func MapToDonationCategoryListResponse(category *model.DonationCategory) *donation_category.DonationCategoryListResponse {
	return &donation_category.DonationCategoryListResponse{
		ID:             category.ID.Hex(),
		CategoryName:   category.CategoryName,
		Icon:           category.Icon,
		IsDeleted:      category.IsDeleted,
		CreatedAt:      category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastModifiedAt: category.LastModifiedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// MapToDonationCategoryListResponses maps slice of models to slice of DTO responses
func MapToDonationCategoryListResponses(categories []*model.DonationCategory) []donation_category.DonationCategoryListResponse {
	var responses []donation_category.DonationCategoryListResponse
	for _, category := range categories {
		responses = append(responses, *MapToDonationCategoryListResponse(category))
	}
	return responses
}
