package donation

import (
	"cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DonationMapper maps Donation model to BSON for database operations
func DonationMapper(data model.Donation) bson.M {
	result := bson.M{}

	if data.DonationCode != "" {
		result["donation_code"] = data.DonationCode
	}
	if !data.CompanyID.IsZero() {
		result["company_id"] = data.CompanyID
	}
	if !data.CategoryID.IsZero() {
		result["category_id"] = data.CategoryID
	}
	if data.Title != "" {
		result["title"] = data.Title
	}
	result["is_featured"] = data.IsFeatured
	if data.Target != 0 {
		result["target"] = data.Target
	}
	if data.DonationDescription != "" {
		result["donation_description"] = data.DonationDescription
	}
	if len(data.DonationImages) > 0 {
		result["donation_images"] = data.DonationImages
	}
	if data.CoverImage != "" {
		result["cover_image"] = data.CoverImage
	}
	if !data.EndDate.IsZero() {
		result["end_date"] = data.EndDate
	}
	if !data.StartDate.IsZero() {
		result["start_date"] = data.StartDate
	}
	result["is_deleted"] = data.IsDeleted
	if !data.CreatedAt.IsZero() {
		result["created_at"] = data.CreatedAt
	}

	return result
}
