package donation

import (
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MapToDonationListResponse(donation *model.Donation, company *model.DonationCompany, category *model.DonationCategory) *donation_dto.DonationListResponse {
	donationImages := make([]donation_dto.DonationImage, len(donation.DonationImages))
	for i, img := range donation.DonationImages {
		donationImages[i] = donation_dto.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt.Format(time.RFC3339),
		}
	}

	companyResponse := donation_dto.Company{
		ID:            company.ID.Hex(),
		CompanyName:   company.CompanyName,
		CompanyLogo:   company.CompanyLogo,
		AccountNumber: company.AccountNumber,
	}

	categoryResponse := donation_dto.Category{
		ID:           category.ID.Hex(),
		CategoryName: category.CategoryName,
		Icon:         category.Icon,
	}

	return &donation_dto.DonationListResponse{
		ID:                  donation.ID.Hex(),
		DonationCode:        donation.DonationCode,
		Company:             companyResponse,
		Category:            categoryResponse,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      donationImages,
		CoverImage:          donation.CoverImage,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		IsDeleted:           donation.IsDeleted,
		CreatedAt:           donation.CreatedAt.Format(time.RFC3339),
		LastModifiedAt:      donation.LastModifiedAt.Format(time.RFC3339),
		Enabled:             donation.Enabled,
	}
}

func DonationMapper(donation model.Donation) bson.M {
	updateData := bson.M{
		
			"donation_code":        donation.DonationCode,
			"company_id":           donation.CompanyID,
			"category_id":          donation.CategoryID,
			"title":                donation.Title,
			"is_featured":          donation.IsFeatured,
			"target":               donation.Target,
			"donation_description": donation.DonationDescription,
			"donation_images":      donation.DonationImages,
			"cover_image":          donation.CoverImage,
			"end_date":             donation.EndDate,
			"start_date":           donation.StartDate,
			"enabled":              donation.Enabled,
			"last_modified_at":     donation.LastModifiedAt,
		
	}
	

	return updateData
}