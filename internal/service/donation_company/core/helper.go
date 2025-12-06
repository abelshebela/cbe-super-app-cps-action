package core

import (
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

// func MapToDonationCompany(donationCompany *dto.DonationCompanyListResponse, Enabled bool) model.DonationCompany {
// 	return model.DonationCompany{
// 		CompanyName:    donationCompany.CompanyName,
// 		CompanyLogo:    donationCompany.CompanyLogo,
// 		AccountNumber:  donationCompany.AccountNumber,
// 		IsDeleted:      false,
// 		Enabled:        Enabled,
// 		LastModifiedAt: time.Now(),
// 	}
// }

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

	if donationCompany.PhoneNumber != "" && donationCompany.Address != existing.Address {
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

// MapToDonationCompanyRequest creates a request DTO from CPS request DTO
// func MapToDonationCompanyRequest(cpsRequest model.DonationCompany) *model.DonationCompany {
// 	return &model.DonationCompany{
// 		CompanyName:   cpsRequest.CompanyName,
// 		CompanyCode:   cpsRequest.CompanyCode,
// 		AccountNumber: cpsRequest.AccountNumber,
// 		PhoneNumber:   cpsRequest.PhoneNumber,
// 		Email:         cpsRequest.Email,
// 		Address:       cpsRequest.Address,
// 	}
// }

func IsDataSimilar(request dto.DonationCompanyRequest, existing *model.DonationCompany) bool {
	// Check if company name is the same (if provided in request)
	if request.CompanyName != "" && request.CompanyName != existing.CompanyName {
		return false
	}

	// Check if the company code is the same
	if request.CompanyCode != "" && request.CompanyCode != existing.CompanyCode {
		return false
	}

	// Check if account number is the same (if provided in request)
	if request.AccountNumber != "" && request.AccountNumber != existing.AccountNumber {
		return false
	}

	// Check if logo is being updated
	if request.CompanyLogo != nil {
		return false
	}

	// If no fields are provided, consider it similar
	if request.CompanyName == "" && request.AccountNumber == "" {
		return true
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
