package core

import (
	dto "cbe-super-app-cps-action/internal/constants/dto/customer_kyc"
	"cbe-super-app-cps-action/internal/constants/model"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	accountLookup "cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hugokessem/coreio/core"

	// coreCustomer "github.com/hugokessem/coreio/lib/core/cusotmer/customer_creation"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func extractDuplicateContract(message string) string {
	parts := strings.Fields(message)

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

func CreateAccountToCore(ctx context.Context, data core.CreateCustomerParam, accountLookupService accountLookup.Account, coreAPI core.CBECoreAPIInterface, logger utils.Logger) (*core.CreateCustomerResult, error) {
	response, err := coreAPI.CreateCustomer(ctx, data)
	if err != nil {
		logger.Errorf("failed to create customer: %v", err)
		return nil, err
	}

	if response == nil {
		return nil, fmt.Errorf("customer creation returned nil response")
	}

	if !response.Success {
		logger.Errorf("customer creation rejected by core: %v", response.Messages)
		return nil, fmt.Errorf("customer creation failed: %s", strings.Join(response.Messages, ", "))
	}

	// Workaround: coreio maps Detail.Customer to T24 EMPLOYERSNAME which is always empty.
	// The real customer number is in the response's transactionId but not exposed by the library.
	// We recover it via PhoneLookup which queries T24 by MNEMONIC and returns CustomerID.
	if response.Detail != nil && response.Detail.Customer == "" && response.Detail.Menmonic != "" {
		lookup, lookupErr := coreAPI.PhoneLookup(ctx, core.PhoneLookupParam{PhoneNumber: response.Detail.Menmonic})
		if lookupErr != nil {
			logger.Warnf("[CreateAccountToCore] customer number lookup by mnemonic=%s failed: %v", response.Detail.Menmonic, lookupErr)
		} else if lookup != nil && lookup.Success && lookup.Detail != nil && lookup.Detail.CustomerID != "" {
			response.Detail.Customer = lookup.Detail.CustomerID
			logger.Infof("[CreateAccountToCore] resolved customer number %s via mnemonic lookup", response.Detail.Customer)
		} else {
			logger.Warnf("[CreateAccountToCore] mnemonic lookup returned no CustomerID for mnemonic=%s", response.Detail.Menmonic)
		}
	}

	return response, nil
}

func MapCustomerKYCToResponsePaginated(c *types.PaginatedResponse[[]imodel.CustomerKYC]) *types.PaginatedResponse[[]dto.CustomerKYCResponse] {

	if c == nil {
		return nil
	}

	responses := make([]dto.CustomerKYCResponse, 0, len(c.Data))

	for _, item := range c.Data {

		// Split full name into first/middle/last
		var firstName, middleName, lastName string

		nameParts := strings.Fields(item.KYCData.FullName)

		if len(nameParts) > 0 {
			firstName = nameParts[0]
		}

		if len(nameParts) == 2 {
			lastName = nameParts[1]
		}

		if len(nameParts) >= 3 {
			middleName = strings.Join(nameParts[1:len(nameParts)-1], " ")
			lastName = nameParts[len(nameParts)-1]
		}

		response := dto.CustomerKYCResponse{
			ID: item.ID.Hex(),

			PersonalInformation: dto.PersonalInformation{
				FirstName:     firstName,
				MiddleName:    middleName,
				LastName:      lastName,
				MotherName:    item.KYCData.MothersName,
				PhoneNumber:   item.KYCData.PhoneNumber,
				Gender:        item.KYCData.Gender,
				Nationality:   item.KYCData.Nationality,
				DateOfBirth:   item.KYCData.BirthDate.Format(time.RFC3339),
				MaritalStatus: item.KYCData.MaritalStatus,
			},

			ResidentialAddress: dto.ResidentialAddress{
				Country: item.KYCData.Country,
				Zone:    item.KYCData.Address.Zone,
				Region:  item.KYCData.Address.Region,
				Wereda:  item.KYCData.Address.Woreda,
				Kebele:  item.KYCData.Address.Kebele,
			},

			FinancialInformation: dto.FinancialInformation{
				EmploymentStatus:     item.KYCData.EmployementStatus,
				SourceOfIncome:       item.KYCData.SourceOfIncome,
				Occupation:           item.KYCData.Occupation,
				AverageMonthlyIncome: item.KYCData.MonthlyIncome,
			},

			CapturedDocuments: dto.CapturedDocuments{
				Photo:         item.KYCData.Picture,
				LivenessVideo: "",
				IDCardFront:   item.KYCData.DocumentFront,
				IDCardBack:    item.KYCData.DocumentBack,
			},

			CustomerStatus:      boolToCustomerStatus(item.Enabled),
			KYCStatus:           string(item.KYCStatus),
			MoneyLaunderingFree: nil,
			TermsAndConditions:  "",

			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			UpdatedAt: item.LastModifiedAt.Format(time.RFC3339),
		}

		responses = append(responses, response)
	}

	return &types.PaginatedResponse[[]dto.CustomerKYCResponse]{
		Data: responses,
		Meta: c.Meta,
	}
}

func boolToCustomerStatus(enabled bool) string {
	if enabled {
		return "ACTIVE"
	}

	return "INACTIVE"
}

func MapCustomerKYCToResponse(c *model.CustomerKYC) *dto.CustomerKYCResponse {

	if c == nil {
		return nil
	}

	// Split full name into first/middle/last
	var firstName, middleName, lastName string

	nameParts := strings.Fields(c.KYCData.FullName)

	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}

	if len(nameParts) == 2 {
		lastName = nameParts[1]
	}

	if len(nameParts) >= 3 {
		middleName = strings.Join(nameParts[1:len(nameParts)-1], " ")
		lastName = nameParts[len(nameParts)-1]
	}

	response := dto.CustomerKYCResponse{
		ID: c.ID.Hex(),

		PersonalInformation: dto.PersonalInformation{
			FirstName:     firstName,
			MiddleName:    middleName,
			LastName:      lastName,
			MotherName:    c.KYCData.MothersName,
			PhoneNumber:   c.KYCData.PhoneNumber,
			Gender:        c.KYCData.Gender,
			Nationality:   c.KYCData.Nationality,
			DateOfBirth:   c.KYCData.BirthDate.Format(time.RFC3339),
			MaritalStatus: c.KYCData.MaritalStatus,
		},

		ResidentialAddress: dto.ResidentialAddress{
			Country: c.KYCData.Country,
			Zone:    c.KYCData.Address.Zone,
			Region:  c.KYCData.Address.Region,
			Wereda:  c.KYCData.Address.Woreda,
			Kebele:  c.KYCData.Address.Kebele,
		},

		FinancialInformation: dto.FinancialInformation{
			EmploymentStatus:     c.KYCData.EmployementStatus,
			SourceOfIncome:       c.KYCData.SourceOfIncome,
			Occupation:           c.KYCData.Occupation,
			AverageMonthlyIncome: c.KYCData.MonthlyIncome,
		},

		CapturedDocuments: dto.CapturedDocuments{
			Photo:         c.KYCData.Picture,
			LivenessVideo: "",
			IDCardFront:   c.KYCData.DocumentFront,
			IDCardBack:    c.KYCData.DocumentBack,
		},

		CustomerStatus:      boolToCustomerStatus(c.Enabled),
		KYCStatus:           string(c.KYCStatus),
		MoneyLaunderingFree: nil,
		TermsAndConditions:  "",

		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.LastModifiedAt.Format(time.RFC3339),
	}

	return &response
}
