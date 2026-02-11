package donation_company

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func MapToDonationCompanyListResponse(company *model.DonationCompany) *dto.DonationCompanyListResponse {
	return &dto.DonationCompanyListResponse{
		ID:                company.ID.Hex(),
		CompanyName:       company.CompanyName,
		CompanyCode:       company.CompanyCode,
		CompanyLogo:       company.CompanyLogo,
		AccountNumber:     company.AccountNumber,
		AccountHolderName: company.AccountHolderName,
		PhoneNumber:       company.PhoneNumber,
		Email:             company.Email,
		Address:           company.Address,
		IsDeleted:         company.IsDeleted,
		Enabled:           company.Enabled,
		CreatedAt:         company.CreatedAt.Format(time.RFC3339),
		LastModifiedAt:    company.LastModifiedAt.Format(time.RFC3339),
	}
}

func MapToDonationCompanyListResponses(companies []model.DonationCompany) []dto.DonationCompanyListResponse {
	responses := make([]dto.DonationCompanyListResponse, len(companies))
	for i := range companies {
		responses[i] = *MapToDonationCompanyListResponse(&companies[i])
	}
	return responses
}

func DonationCompanyMapper(company model.DonationCompany) bson.M {
	return bson.M{
		"company_name":        company.CompanyName,
		"company_logo":        company.CompanyLogo,
		"company_code":        company.CompanyCode,
		"account_number":      company.AccountNumber,
		"phone_number":        company.PhoneNumber,
		"account_holder_name": company.AccountHolderName,
		"email":               company.Email,
		"address":             company.Address,
		"enabled":             company.Enabled,
		"is_deleted":          company.IsDeleted,
		"last_modified_at":    company.LastModifiedAt,
	}
}
