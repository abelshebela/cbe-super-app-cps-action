package core

import (
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"context"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func CompanyNameExists(ctx context.Context, companyName string, donationCompanyRepo storage.DonationCompanyRepository) (bool, error) {
	// Use the repository method to find donation companies
	companies, err := donationCompanyRepo.FindAllWithPagination(ctx, types.Filter{
		Search:  companyName,
		Page:    1,
		PerPage: 1,
	})
	if err != nil {
		return false, errors.New(localization.ErrorDonationCompanyLookupFailed.Code)
	}

	// Check if any companies were found
	return len(companies.Data) > 0, nil
}

func CheckIfAccountExists(ctx context.Context, accountNumber string, repo storage.DonationCompanyRepository) error {
	donCompany, err := repo.FindByAccountNumber(ctx, accountNumber)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil
		}
		return err
	}
	if donCompany.AccountNumber != "" {
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}
	return nil
}

func ValidateAccountNumberWithExternalAPI(ctx context.Context, accountNumber string, accountLookupService account_lookup.Account) (*model.AccountDetail, error) {
	accountRequest := model.AccountLookUpRequest{
		AccountNumber: accountNumber,
	}
	accountDetail, err := accountLookupService.LookupAccountByAccountNumber(ctx, accountRequest)
	if err != nil {
		return nil, err
	}

	if accountDetail == nil {
		return nil, err
	}

	return accountDetail, nil
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

// MapToDonationCompanyResponse creates a response DTO from request DTO and logo URL
func MapToDonationCompanyResponse(donationCompany dto.DonationCompanyRequest, logoURL string) dto.DonationCompanyResponse {
	return dto.DonationCompanyResponse{
		CompanyName:    donationCompany.CompanyName,
		CompanyCode:    donationCompany.CompanyCode,
		CompanyLogo:    logoURL,
		AccountNumber:  donationCompany.AccountNumber,
		PhoneNumber:    donationCompany.PhoneNumber,
		Email:          donationCompany.Email,
		Address:        donationCompany.Address,
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      time.Now().Format(time.RFC3339),
		LastModifiedAt: time.Now().Format(time.RFC3339),
	}
}

// MapToDonationCompanyCPSRequest creates a CPS request DTO from request DTO and logo URL
func MapToDonationCompanyCPSRequest(id string, donationCompany dto.DonationCompanyRequest, logoURL string) dto.DonationCompanyCPSRequest {
	return dto.DonationCompanyCPSRequest{
		ID:            id,
		CompanyName:   donationCompany.CompanyName,
		CompanyCode:   donationCompany.CompanyCode,
		CompanyLogo:   logoURL,
		AccountNumber: donationCompany.AccountNumber,
		PhoneNumber:   donationCompany.PhoneNumber,
		Email:         donationCompany.Email,
		Address:       donationCompany.Address,
	}
}

func MapToDonationCompanyonUpdateCPSRequest(id string, existing dto.DonationCompanyListResponse, donationCompany dto.DonationCompanyRequest, logoURL string) *model.DonationCompany {
	result := &model.DonationCompany{}

	if donationCompany.CompanyName != "" && donationCompany.CompanyName != existing.CompanyName {
		result.CompanyName = donationCompany.CompanyName
	} else {
		result.CompanyName = existing.CompanyName
	}

	if donationCompany.CompanyCode != "" && donationCompany.CompanyCode != existing.CompanyCode {
		result.CompanyCode = "DON-COMPANY-" + donationCompany.CompanyCode
	} else {
		result.CompanyCode = existing.CompanyCode
	}

	if donationCompany.AccountNumber != existing.AccountNumber {
		result.AccountNumber = donationCompany.AccountNumber
	} else {
		result.AccountNumber = existing.AccountNumber
	}

	if donationCompany.PhoneNumber != "" && donationCompany.PhoneNumber != existing.PhoneNumber {
		result.PhoneNumber = donationCompany.PhoneNumber
	} else {
		result.PhoneNumber = existing.PhoneNumber
	}

	if donationCompany.Email != "" && donationCompany.Email != existing.Email {
		result.Email = donationCompany.Email
	} else {
		result.Email = existing.Email
	}

	if donationCompany.Address != "" && donationCompany.Address != existing.Address {
		result.Address = donationCompany.Address
	} else {
		result.Address = existing.Address
	}

	if logoURL != existing.CompanyLogo {
		result.CompanyLogo = logoURL
	} else {
		result.CompanyLogo = existing.CompanyLogo
	}
	createdAt, _ := time.Parse(time.RFC3339, existing.CreatedAt)
	result.CreatedAt = createdAt
	result.LastModifiedAt = time.Now()

	result.Enabled = existing.Enabled

	return result
}

func IsDataSimilar(request dto.DonationCompanyRequest, existing *model.DonationCompany) bool {
	if request.CompanyName != "" && request.CompanyName != existing.CompanyName {
		return false
	}

	if request.CompanyCode != "" && request.CompanyCode != existing.CompanyCode {
		return false
	}

	if request.CompanyLogo != nil {
		return false
	}

	if request.AccountNumber != "" && request.AccountNumber != existing.AccountNumber {
		return false
	}

	if request.Address != "" && request.Address != existing.Address {
		return false
	}

	if request.PhoneNumber != "" && request.PhoneNumber != existing.PhoneNumber {
		return false
	}

	if request.Email != "" && request.Email != existing.Email {
		return false
	}

	return true
}

// CheckDataSimilarityAndValidation checks if data is similar and validates uniqueness
func CheckDataSimilarityAndValidation(ctx context.Context, request dto.DonationCompanyRequest, existing *dto.DonationCompanyListResponse, donationCompanyRepo storage.DonationCompanyRepository, accountLookupService account_lookup.Account) error {
	// Convert existing DTO to model for similarity check
	existingModel := &model.DonationCompany{
		CompanyName:   existing.CompanyName,
		CompanyCode:   existing.CompanyCode,
		CompanyLogo:   existing.CompanyLogo,
		AccountNumber: existing.AccountNumber,
		IsDeleted:     existing.IsDeleted,
	}

	// Check if data is similar to existing data
	if IsDataSimilar(request, existingModel) {
		return errors.New(localization.ErrorNoChangesToUpdate.Code)
	}

	// Check if company name is being updated and if it already exists
	if request.CompanyName != "" && request.CompanyName != existing.CompanyName {
		ok, err := CompanyNameExists(ctx, request.CompanyName, donationCompanyRepo)
		if err != nil {
			return err
		}
		if ok {
			return errors.New(localization.ErrorCompanyNameAlreadyExists.Code)
		}
	}

	if request.CompanyCode != "" && request.CompanyCode != existing.CompanyCode {
		ok, err := CompanyNameExists(ctx, request.CompanyCode, donationCompanyRepo)
		if err != nil {
			return err
		}
		if ok {
			return errors.New(localization.ErrorCompanyNameAlreadyExists.Code)
		}
	}

	// Check if account number is being updated and if it already exists
	if request.AccountNumber != "" && request.AccountNumber != existing.AccountNumber {
		if err := CheckIfAccountExists(ctx, request.AccountNumber, donationCompanyRepo); err != nil {
			return err
		}
	}

	return nil
}

// ConvertDonationListResponseToModel converts a DonationListResponse DTO to a model.Donation
// This is used when we need to update donations and preserve all existing fields
func ConvertDonationListResponseToModel(donationResponse *donation_dto.DonationListResponse) *model.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(donationResponse.Company.ID)
	categoryObjID, _ := bson.ObjectIDFromHex(donationResponse.Category.ID)

	donationImages := make([]types.DonationImage, len(donationResponse.DonationImages))
	for i, img := range donationResponse.DonationImages {
		createdAt, _ := time.Parse(time.RFC3339, img.CreatedAt)
		donationImages[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: createdAt,
		}
	}

	startTime, _ := time.Parse(time.RFC3339, donationResponse.StartDate)
	endTime, _ := time.Parse(time.RFC3339, donationResponse.EndDate)
	createdAt, _ := time.Parse(time.RFC3339, donationResponse.CreatedAt)
	lastModifiedAt, _ := time.Parse(time.RFC3339, donationResponse.LastModifiedAt)

	donationID, _ := bson.ObjectIDFromHex(donationResponse.ID)

	return &model.Donation{
		ID:                  donationID,
		DonationCode:        donationResponse.DonationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               donationResponse.Title,
		IsFeatured:          donationResponse.IsFeatured,
		Target:              donationResponse.Target,
		CurrentAmount:       donationResponse.CurrentAmount,
		DonationDescription: donationResponse.DonationDescription,
		DonationImages:      donationImages,
		CoverImage:          donationResponse.CoverImage,
		StartDate:           startTime,
		EndDate:             endTime,
		Enabled:             donationResponse.Enabled,
		IsDeleted:           donationResponse.IsDeleted,
		CreatedAt:           createdAt,
		LastModifiedAt:      lastModifiedAt,
	}
}
