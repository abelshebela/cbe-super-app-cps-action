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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func extractDuplicateContract(message string) string {
	parts := strings.Fields(message)

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

func CreateAccountToCore(ctx context.Context, data core.CreateCustomerParam, accountLookupService accountLookup.Account, coreAPI core.CBECoreAPIInterface, cfg *config.VaultConfig, logger utils.Logger) (*core.CusteomerAccountCreationResponse, error) {
	response, err := coreAPI.AccountCreate(ctx, data, cfg.AccountOpeningURL, "6501")
	if err != nil {
		logger.Errorf("failed to create customer: %v", err)
		return nil, err
	}

	if response == nil {
		return nil, fmt.Errorf("customer creation returned nil response")
	}

	if !response.AccountCreationDetail.Success {
		logger.Errorf("customer creation rejected by core: %v", response.AccountCreationDetail.Messages)
		return nil, fmt.Errorf("customer creation failed: %s", strings.Join(response.AccountCreationDetail.Messages, ", "))
	}

	return response, nil
}

func MapCustomerKYCToResponsePaginated(c *types.PaginatedResponse[[]imodel.CustomerKYC]) *types.PaginatedResponse[[]dto.CustomerKYCResponse] {

	if c == nil {
		return nil
	}

	responses := make([]dto.CustomerKYCResponse, 0, len(c.Data))

	for _, item := range c.Data {
		var livenessVideo string
		if item.ComplyCube != nil {
			livenessVideo = item.ComplyCube.LiveVideoID
		}

		response := dto.CustomerKYCResponse{
			ID: item.ID.Hex(),

			PersonalInformation: dto.PersonalInformation{
				FullName:      item.KYCData.FullName,
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
				LivenessVideo: livenessVideo,
				IDCardFront:   item.KYCData.DocumentFront,
				IDCardBack:    item.KYCData.DocumentBack,
			},

			Review: item.Review,

			CustomerStatus:      boolToCustomerStatus(item.Enabled),
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

	var livenessVideo string
	if c.ComplyCube != nil {
		livenessVideo = c.ComplyCube.LiveVideoID
	}

	response := dto.CustomerKYCResponse{
		ID: c.ID.Hex(),

		PersonalInformation: dto.PersonalInformation{
			FullName:      c.KYCData.FullName,
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
			LivenessVideo: livenessVideo,
			IDCardFront:   c.KYCData.DocumentFront,
			IDCardBack:    c.KYCData.DocumentBack,
		},

		Review: c.Review,

		CustomerStatus:      boolToCustomerStatus(c.Enabled),
		MoneyLaunderingFree: nil,
		TermsAndConditions:  "",

		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.LastModifiedAt.Format(time.RFC3339),
	}

	return &response
}

func MapSelfActivationUserToResponsePaginated(c *types.PaginatedResponse[[]imodel.SelfActivationUser]) *types.PaginatedResponse[[]dto.CustomerKYCResponse] {
	if c == nil {
		return nil
	}

	responses := make([]dto.CustomerKYCResponse, 0, len(c.Data))
	for _, item := range c.Data {
		if mapped := MapSelfActivationUserToResponse(&item); mapped != nil {
			responses = append(responses, *mapped)
		}
	}

	return &types.PaginatedResponse[[]dto.CustomerKYCResponse]{
		Data: responses,
		Meta: c.Meta,
	}
}

func MapSelfActivationUserToResponse(u *imodel.SelfActivationUser) *dto.CustomerKYCResponse {
	if u == nil {
		return nil
	}

	response := dto.CustomerKYCResponse{
		ID:             u.ID.Hex(),
		Sub:            u.Sub,
		CustomerNumber: u.CustomerNumber,
		PersonalInformation: dto.PersonalInformation{
			FullName:       u.Name,
			Email:          u.Email,
			PhoneNumber:    u.PhoneNumber,
			Gender:         u.Gender,
			Nationality:    u.Nationality,
			DateOfBirth:    u.BirthDate,
			MaritalStatus:  "",
			AccountNumbers: u.ChosenAccounts,
		},
		ResidentialAddress: dto.ResidentialAddress{
			Zone:   u.Address.Zone,
			Kebele: u.Address.Kebele,
			Wereda: u.Address.Woreda,
			Region: u.Address.Region,
		},
		FinancialInformation: dto.FinancialInformation{},
		CapturedDocuments: dto.CapturedDocuments{
			Photo:         u.Picture,
			LivenessVideo: u.SelfiePhoto,
		},
		CustomerStatus: "NEW",
		KYC: dto.KycInfo{
			PhoneMismatch:  u.KYC.PhoneMismatch,
			BelowThreshold: u.KYC.BelowThreshold,
			ContainsANDOR:  u.KYC.ContainsANDOR,
		},
		ComplyCube: u.ComplyCube,

		Review: u.Review,

		MoneyLaunderingFree: nil,
		TermsAndConditions:  "",
		CreatedAt:           u.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           u.UpdatedAt.Format(time.RFC3339),
	}

	return &response
}

func MapSelfActivationUserInfo(u *imodel.UserInfo) *imodel.UserInfo {
	if u == nil {
		return nil
	}

	return &imodel.UserInfo{
		ID:          u.ID,
		UserCode:    u.UserCode,
		FullName:    u.FullName,
		Email:       u.Email,
		Department:  u.Department,
		PhoneNumber: u.PhoneNumber,
	}
}

func BuildRow(request imodel.ExportSelfActivationRequest) []string {
	var registrationDate string
	if !request.RegistrationDate.IsZero() {
		registrationDate = request.RegistrationDate.Format(time.RFC3339)
	}

	return []string{
		request.CustomerName,
		request.PhoneNumber,
		request.Gender,
		request.DateOfBirth,
		request.Region,
		registrationDate,
		request.RejectionReason,
		request.CustomerStatus,
		request.KYCStatus,
	}
}
