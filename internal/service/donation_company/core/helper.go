package core

import (
	dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"context"
	"encoding/json"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

func ValidateAccountNumberWithExternalAPI(ctx context.Context, accountNumber string, accountLookupService account_lookup.Account, logger utils.Logger) (*imodel.AccountDetail, error) {
	accountRequest := model.AccountLookUpRequest{
		AccountNumber: accountNumber,
	}
	accountDetail, err := accountLookupService.LookupAccountByAccountNumber(ctx, accountRequest)
	if err != nil {
		logger.Errorf("Error occurred while validating account number: %v", err)
		return nil, err
	}
	if accountDetail == nil {
		logger.Errorf("Account not found for number: %s", accountNumber)
		return nil, errors.New(localization.ErrorAccountNotFound.Code)
	}
	if accountDetail.Currency != "ETB" {
		logger.Errorf("Account currency not supported for number: %s", accountNumber)
		return nil, errors.New(localization.ErrorAccountCurrencyNotSupported.Code)
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
		CompanyName:        donationCompany.CompanyName,
		CompanyCode:        donationCompany.CompanyCode,
		CompanyDescription: donationCompany.CompanyDescription,
		CompanyLogo:        logoURL,
		PhoneNumber:        donationCompany.PhoneNumber,
		Email:              donationCompany.Email,
		Address:            donationCompany.Address,
		Enabled:            true,
		IsDeleted:          false,
		CreatedAt:          time.Now().Format(time.RFC3339),
		LastModifiedAt:     time.Now().Format(time.RFC3339),
	}
}

// MapToDonationCompanyCPSRequest creates a CPS request DTO from request DTO and logo URL
func MapToDonationCompanyCPSRequest(id string, donationCompany dto.DonationCompanyRequest, logoURL string) dto.DonationCompanyCPSRequest {
	return dto.DonationCompanyCPSRequest{
		ID:                 id,
		CompanyName:        donationCompany.CompanyName,
		CompanyCode:        donationCompany.CompanyCode,
		CompanyDescription: donationCompany.CompanyDescription,
		CompanyLogo:        logoURL,
		PhoneNumber:        donationCompany.PhoneNumber,
		Email:              donationCompany.Email,
		Address:            donationCompany.Address,
	}
}

func MapToDonationCompanyonUpdateCPSRequest(id string, existing dto.DonationCompanyListResponse, donationCompany dto.DonationCompanyRequest, logoURL string) *imodel.DonationCompanyOracle {
	result := &imodel.DonationCompanyOracle{}

	if donationCompany.CompanyName != "" && donationCompany.CompanyName != existing.CompanyName {
		result.CompanyName = donationCompany.CompanyName
	} else {
		result.CompanyName = existing.CompanyName
	}

	result.CompanyCode = existing.CompanyCode
	if donationCompany.CompanyDescription != "" {
		result.CompanyDescription = donationCompany.CompanyDescription
	} else {
		result.CompanyDescription = existing.CompanyDescription
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

func IsDataSimilar(request dto.DonationCompanyRequest, existing *imodel.DonationCompanyOracle) bool {
	if request.CompanyName != "" && request.CompanyName != existing.CompanyName {
		return false
	}
	if request.CompanyDescription != "" && request.CompanyDescription != existing.CompanyDescription {
		return false
	}
	if request.CompanyCode != "" && request.CompanyCode != existing.CompanyCode {
		return false
	}

	if request.CompanyLogo != nil {
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
func CheckDataSimilarityAndValidation(ctx context.Context, request dto.DonationCompanyRequest, existing *dto.DonationCompanyListResponse, donationCompanyRepo storage.DonationCompanyRepository) error {
	// Convert existing DTO to model for similarity check
	existingModel := &imodel.DonationCompanyOracle{
		CompanyName:        existing.CompanyName,
		CompanyCode:        existing.CompanyCode,
		CompanyDescription: existing.CompanyDescription,
		CompanyLogo:        existing.CompanyLogo,
		PhoneNumber:        existing.PhoneNumber,
		Email:              existing.Email,
		Address:            existing.Address,
		IsDeleted:          existing.IsDeleted,
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

	return nil
}
