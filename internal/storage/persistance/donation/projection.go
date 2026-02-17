package donation

import (
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/types"
	"time"

	donation_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/donation"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func MapToDonationListResponse(donation imodel.Donation, company *donation_model.DonationCompany, category *donation_model.DonationCategory) *donation_dto.DonationListResponse {
	donationImages := make([]types.DonationImage, len(donation.DonationImages))
	for i, img := range donation.DonationImages {
		donationImages[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt,
		}
	}

	companyResponse := donation_dto.Company{
		ID:            company.ID.Hex(),
		CompanyName:   company.CompanyName,
		CompanyLogo:   company.CompanyLogo,
		AccountNumber: company.AccountNumber,
		Enabled:       company.Enabled,
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
		Target:              string(donation.Target),
		DonationDescription: donation.DonationDescription,
		DonationImages:      donationImages,
		CurrentAmount:       string(donation.CurrentAmount),
		CoverImage:          donation.CoverImage,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		IsDeleted:           donation.IsDeleted,
		CreatedAt:           donation.CreatedAt.Format(time.RFC3339),
		LastModifiedAt:      donation.LastModifiedAt.Format(time.RFC3339),
		Enabled:             donation.Enabled,
	}
}

func DonationMapper(donation imodel.Donation) bson.M {
	updateData := bson.M{

		"donation_code":        donation.DonationCode,
		"company_id":           donation.CompanyID,
		"category_id":          donation.CategoryID,
		"title":                donation.Title,
		"is_featured":          donation.IsFeatured,
		"target":               donation.Target,
		"current_amount":       donation.CurrentAmount,
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
