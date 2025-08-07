package mappers

import (
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	dto "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
)

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

func ToDonationCompanyModel(d *dto.DonationCompanyRequest) *model.DonationCompany {
	return &model.DonationCompany{
		CompanyName:   d.CompanyName,
		CompanyLogo:   d.CompanyLogo,
		AccountNumber: d.AccountNumber,
		CreatedAt:     time.Now(),
	}
}

func ToDonationCompanyDTO(m *model.DonationCompany) *dto.DonationCompanyRequest {
	return &dto.DonationCompanyRequest{
		CompanyName:   m.CompanyName,
		CompanyLogo:   m.CompanyLogo,
		AccountNumber: m.AccountNumber,
	}
}

func ToDonationModel(d *dto.DonationRequest) *model.Donation {
	return &model.Donation{
		CompanyID:           d.CompanyID,
		CategoryID:          d.CategoryID,
		Title:               d.Title,
		IsFeatured:          d.IsFeatured,
		Target:              d.Target,
		DonationDescription: d.DonationDescription,
		DonationImages:      d.DonationImages,
		EndDate:             d.EndDate,
		StartDate:           d.StartDate,
		CreatedAt:           time.Now(),
	}
}

func ToDonationDTO(m *model.Donation) *dto.DonationRequest {
	return &dto.DonationRequest{
		CompanyID:           m.CompanyID,
		CategoryID:          m.CategoryID,
		Title:               m.Title,
		IsFeatured:          m.IsFeatured,
		Target:              m.Target,
		DonationDescription: m.DonationDescription,
		DonationImages:      m.DonationImages,
		EndDate:             m.EndDate,
		StartDate:           m.StartDate,
	}
}
