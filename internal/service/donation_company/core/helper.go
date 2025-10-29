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

func AccountNumberExists(ctx context.Context, accountNumber string, donationCompanyRepo storage.DonationCompanyRepository) (bool, error) {
	// Use the repository method to find donation companies by account number
	companies, err := donationCompanyRepo.FindAllWithPagination(ctx, types.Filter{
		Search:  accountNumber,
		Page:    1,
		PerPage: 1,
	})
	if err != nil {
		return false, errors.New(localization.ErrorDonationCompanyLookupFailed.Code)
	}

	// Check if any companies were found with this account number
	for _, company := range companies.Data {
		if company.AccountNumber == accountNumber {
			return true, nil
		}
	}
	return false, nil
}

func ValidateAccountNumberWithExternalAPI(ctx context.Context, accountNumber string, accountLookupService account_lookup.Account) (*model.AccountInfo, error) {
	// Create account lookup request
	accountRequest := model.AccountLookUpRequest{
		AccountNumber: accountNumber,
	}

	// Call external account lookup service
	accountInfo, err := accountLookupService.LookupAccountByAccountNumber(ctx, accountRequest)
	if err != nil {
		return nil, err
	}

	// Check if account is active and valid
	if accountInfo == nil {
		return nil, errors.New(localization.ErrorAccountNumberNotFound.Code)
	}

	// Additional validation checks
	if !accountInfo.ActiveAccount {
		return nil, errors.New(localization.ErrorAccountNumberNotActive.Code)
	}

	if accountInfo.AccountFrozen {
		return nil, errors.New(localization.ErrorAccountNumberNotActive.Code)
	}

	if accountInfo.AccountDormant {
		return nil, errors.New(localization.ErrorAccountNumberNotActive.Code)
	}

	return accountInfo, nil
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func MapToDonationCompany(donationCompany *dto.DonationCompanyListResponse,Enabled bool) model.DonationCompany {
	return model.DonationCompany{
		CompanyName:    donationCompany.CompanyName,
		CompanyLogo:    donationCompany.CompanyLogo,
		AccountNumber:  donationCompany.AccountNumber,
		IsDeleted:      false,
		Enabled: Enabled,
		LastModifiedAt: time.Now(),
	}
}

// MapToDonationCompanyResponse creates a response DTO from request DTO and logo URL
func MapToDonationCompanyResponse(donationCompany dto.DonationCompanyRequest, logoURL string) dto.DonationCompanyResponse {
	return dto.DonationCompanyResponse{
		CompanyName:   donationCompany.CompanyName,
		CompanyLogo:   logoURL,
		AccountNumber: donationCompany.AccountNumber,
		PhoneNumber:   donationCompany.PhoneNumber,
		Email:         donationCompany.Email,
		Address:       donationCompany.Address,
		Enabled:       true,
	}
}

// MapToDonationCompanyCPSRequest creates a CPS request DTO from request DTO and logo URL
func MapToDonationCompanyCPSRequest(id string, donationCompany dto.DonationCompanyRequest, logoURL string) dto.DonationCompanyCPSRequest {
	return dto.DonationCompanyCPSRequest{
		ID:            id,
		CompanyName:   donationCompany.CompanyName,
		CompanyLogo:   logoURL,
		AccountNumber: donationCompany.AccountNumber,
	}
}

// MapToDonationCompanyRequest creates a request DTO from CPS request DTO
func MapToDonationCompanyRequest(cpsRequest dto.DonationCompanyCPSRequest) dto.DonationCompanyRequest {
	return dto.DonationCompanyRequest{
		CompanyName:   cpsRequest.CompanyName,
		AccountNumber: cpsRequest.AccountNumber,
	}
}

func IsDataSimilar(request dto.DonationCompanyRequest, existing *model.DonationCompany) bool {
	// Check if company name is the same (if provided in request)
	if request.CompanyName != "" && request.CompanyName != existing.CompanyName {
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

	// Check if account number is being updated and if it already exists
	if request.AccountNumber != "" && request.AccountNumber != existing.AccountNumber {
		ok, err := AccountNumberExists(ctx, request.AccountNumber, donationCompanyRepo)
		if err != nil {
			return err
		}
		if ok {
			return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
		}

		// Validate account number with external API
		if _, err := ValidateAccountNumberWithExternalAPI(ctx, request.AccountNumber, accountLookupService); err != nil {
			return err
		}
	}

	return nil
}
