package mappers

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Helper function to convert string ID to ObjectID
func stringToObjectID(id string) (bson.ObjectID, error) {
	return common_util.ParsePrimitiveObjectID(id)
}

// Helper function to convert ObjectID to string
func objectIDToString(objID bson.ObjectID) string {
	return objID.Hex()
}

func ToDonationCategoryModel(d *dto.DonationCategoryRequest) *model.DonationCategory {
	return &model.DonationCategory{
		CategoryName: d.CategoryName,
		Icon:         "", // Will be set by service
		IsDeleted:    false,
		CreatedAt:    time.Now(),
	}
}

func ToDonationCategoryModelWithURL(d *dto.DonationCategoryRequest, iconURL string) *model.DonationCategory {
	return &model.DonationCategory{
		CategoryName: d.CategoryName,
		Icon:         iconURL,
		IsDeleted:    false,
		CreatedAt:    time.Now(),
	}
}

func ToDonationCategoryDTO(m *model.DonationCategory) *dto.DonationCategoryRequest {
	return &dto.DonationCategoryRequest{
		CategoryName: m.CategoryName,
		Icon:         nil, // URL is stored in Icon field
	}
}

func ToDonationCategoryListResponse(m *model.DonationCategory) *dto.DonationCategoryListResponse {
	return &dto.DonationCategoryListResponse{
		CategoryName:   m.CategoryName,
		Icon:           m.Icon,
		IsDeleted:      m.IsDeleted,
		CreatedAt:      m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastModifiedAt: m.LastModifiedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToDonationCompanyDTO(m *model.DonationCompany) *dto.DonationCompanyRequest {
	return &dto.DonationCompanyRequest{
		CompanyName:   m.CompanyName,
		CompanyLogo:   nil, // URL is stored in CompanyLogo field
		AccountNumber: m.AccountNumber,
	}
}

func ToDonationCompanyModel(d *dto.DonationCompanyRequest) *model.DonationCompany {
	return &model.DonationCompany{
		CompanyName:   d.CompanyName,
		CompanyLogo:   "", // Will be set by service
		AccountNumber: d.AccountNumber,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
	}
}

func ToDonationCompanyModelWithURL(d *dto.DonationCompanyRequest, logoURL string) *model.DonationCompany {
	return &model.DonationCompany{
		CompanyName:   d.CompanyName,
		CompanyLogo:   logoURL,
		AccountNumber: d.AccountNumber,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
	}
}

func ToDonationCompanyListResponse(m *model.DonationCompany) *dto.DonationCompanyListResponse {
	return &dto.DonationCompanyListResponse{
		CompanyName:    m.CompanyName,
		CompanyLogo:    m.CompanyLogo,
		AccountNumber:  m.AccountNumber,
		IsDeleted:      m.IsDeleted,
		CreatedAt:      m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastModifiedAt: m.LastModifiedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToDonationModel(d *dto.DonationRequest) *model.Donation {
	// Convert string IDs to ObjectIDs
	companyObjID, _ := stringToObjectID(d.CompanyID)
	categoryObjID, _ := stringToObjectID(d.CategoryID)

	return &model.Donation{
		DonationCode:        d.DonationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               d.Title,
		IsFeatured:          d.IsFeatured,
		Target:              d.Target,
		DonationDescription: d.DonationDescription,
		DonationImages:      []string{}, // Will be set by service
		EndDate:             d.EndDate,
		StartDate:           d.StartDate,
		IsDeleted:           false,
		CreatedAt:           time.Now(),
	}
}

func ToDonationModelWithURLs(d *dto.DonationRequest, imageURLs []string) *model.Donation {
	// Convert string IDs to ObjectIDs
	companyObjID, _ := stringToObjectID(d.CompanyID)
	categoryObjID, _ := stringToObjectID(d.CategoryID)

	return &model.Donation{
		DonationCode:        d.DonationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               d.Title,
		IsFeatured:          d.IsFeatured,
		Target:              d.Target,
		DonationDescription: d.DonationDescription,
		DonationImages:      imageURLs,
		EndDate:             d.EndDate,
		StartDate:           d.StartDate,
		IsDeleted:           false,
		CreatedAt:           time.Now(),
	}
}

func ToDonationDTO(m *model.Donation) *dto.DonationRequest {
	return &dto.DonationRequest{
		DonationCode:        m.DonationCode,
		CompanyID:           objectIDToString(m.CompanyID),
		CategoryID:          objectIDToString(m.CategoryID),
		Title:               m.Title,
		IsFeatured:          m.IsFeatured,
		Target:              m.Target,
		DonationDescription: m.DonationDescription,
		DonationImages:      nil, // URLs are stored in DonationImages field
		EndDate:             m.EndDate,
		StartDate:           m.StartDate,
	}
}

func ToDonationListResponse(
	m *model.Donation,
	company *model.DonationCompany,
	category *model.DonationCategory,
) *dto.DonationListResponse {
	return &dto.DonationListResponse{
		DonationCode: m.DonationCode,
		Company: dto.DonationCompanyCPSRequest{
			ID:            company.ID.Hex(),
			CompanyName:   company.CompanyName,
			CompanyLogo:   company.CompanyLogo,
			AccountNumber: company.AccountNumber,
		},
		Category: dto.DonationCategoryCPSRequest{
			ID:           category.ID.Hex(),
			CategoryName: category.CategoryName,
			Icon:         category.Icon,
		},
		Title:               m.Title,
		IsFeatured:          m.IsFeatured,
		Target:              m.Target,
		DonationDescription: m.DonationDescription,
		DonationImages:      m.DonationImages,
		EndDate:             m.EndDate.Format("2006-01-02T15:04:05Z07:00"),
		StartDate:           m.StartDate.Format("2006-01-02T15:04:05Z07:00"),
		IsDeleted:           m.IsDeleted,
		CreatedAt:           m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		LastModifiedAt:      m.LastModifiedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
