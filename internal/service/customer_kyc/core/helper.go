package core

import (
	// "cbe-super-app-cps-action/internal/constants"
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

	// Workaround: coreio maps Detail.Customer to T24 EMPLOYERSNAME which is always empty.
	// The real customer number is in the response's transactionId but not exposed by the library.
	// We recover it via PhoneLookup which queries T24 by MNEMONIC and returns CustomerID.

	// if response.AccountCreationDetail != nil && response.AccountCreationDetail.Detail.Customer == "" && response.AccountCreationDetail.Detail.CoCode != "" {
	// 	lookup, lookupErr := coreAPI.PhoneLookup(ctx, core.PhoneLookupParam{PhoneNumber: response.AccountCreationDetail.Detail.})
	// 	if lookupErr != nil {
	// 		logger.Warnf("[CreateAccountToCore] customer number lookup by mnemonic=%s failed: %v", response.AccountCreationDetail.Menmonic, lookupErr)
	// 	} else if lookup != nil && lookup.Success && lookup.Detail != nil && lookup.Detail.CustomerID != "" {
	// 		response.AccountCreationDetail.Detail.Customer = lookup.Detail.CustomerID
	// 		logger.Infof("[CreateAccountToCore] resolved customer number %s via mnemonic lookup", response.AccountCreationDetail.Detail.Customer)
	// 	} else {
	// 		logger.Warnf("[CreateAccountToCore] mnemonic lookup returned no CustomerID for mnemonic=%s", response.AccountCreationDetail.Menmonic)
	// 	}
	// }

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

		CustomerStatus:      boolToCustomerStatus(c.Enabled),
		KYCStatus:           string(c.KYCStatus),
		MoneyLaunderingFree: nil,
		TermsAndConditions:  "",

		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.LastModifiedAt.Format(time.RFC3339),
	}

	return &response
}

func MapSelfActivationUserToResponsePaginated(c *types.PaginatedResponse[[]imodel.SelfActivationUser]) *types.PaginatedResponse[[]dto.SelfActivationUser] {
	if c == nil {
		return nil
	}

	responses := make([]dto.SelfActivationUser, 0, len(c.Data))
	for _, item := range c.Data {
		if mapped := MapSelfActivationUserToResponse(&item); mapped != nil {
			responses = append(responses, *mapped)
		}
	}

	return &types.PaginatedResponse[[]dto.SelfActivationUser]{
		Data: responses,
		Meta: c.Meta,
	}
}

func MapSelfActivationUserToResponse(u *imodel.SelfActivationUser) *dto.SelfActivationUser {
	if u == nil {
		return nil
	}

	return &dto.SelfActivationUser{
		ID:             u.ID,
		Sub:            u.Sub,
		Name:           u.Name,
		Email:          u.Email,
		PhoneNumber:    u.PhoneNumber,
		Gender:         u.Gender,
		Picture:        u.Picture,
		Nationality:    u.Nationality,
		BirthDate:      u.BirthDate,
		Address: dto.CustomerAddress{
			Zone:   u.Address.Zone,
			Kebele: u.Address.Kebele,
			Woreda: u.Address.Woreda,
			Region: u.Address.Region,
		},
		Enabled:        u.Enabled,
		AccountNumbers: u.AccountNumbers,
		CustomerID:     u.CustomerID,
		ExpiryDate:     u.ExpiryDate,
		IssueDate:      u.IssueDate,
		KYC: dto.KycInfo{
			PhoneMismatch:  u.KYC.PhoneMismatch,
			BelowThreshold: u.KYC.BelowThreshold,
			ContainsANDOR:  u.KYC.ContainsANDOR,
			KYCStatus:      u.KYC.KYCStatus,
		},
		ComplyCube:      mapSelfActivationComplyCube(u.ComplyCube),
		KYCRejectReason: u.KYCRejectReason,
		CreatedAt:       u.CreatedAt,
		LastModifiedAt:  u.LastModifiedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func mapSelfActivationComplyCube(c imodel.ComplyCube) dto.ComplyCube {
	return dto.ComplyCube{
		DocumentID:      c.DocumentID,
		LiveVideoID:     c.LiveVideoID,
		DocumentType:    c.DocumentType,
		IdentityCheckID: c.IdentityCheckID,
		IdentityCheck: dto.IdentityCheck{
			ID:             c.IdentityCheck.ID,
			ClientID:       c.IdentityCheck.ClientID,
			LiveVideoID:    c.IdentityCheck.LiveVideoID,
			DocumentID:     c.IdentityCheck.DocumentID,
			EntityName:     c.IdentityCheck.EntityName,
			Type:           c.IdentityCheck.Type,
			Status:         c.IdentityCheck.Status,
			InitialOutcome: c.IdentityCheck.InitialOutcome,
			Result: dto.IdentityResult{
				Outcome: c.IdentityCheck.Result.Outcome,
				Breakdown: dto.IdentityBreakdown{
					IntegrityAnalysis: dto.IntegrityAnalysis{
						FaceDetection: c.IdentityCheck.Result.Breakdown.IntegrityAnalysis.FaceDetection,
					},
					FaceAnalysis: dto.FaceAnalysis{
						FacialSimilarity:       c.IdentityCheck.Result.Breakdown.FaceAnalysis.FacialSimilarity,
						PreviouslyEnrolledFace: c.IdentityCheck.Result.Breakdown.FaceAnalysis.PreviouslyEnrolledFace,
						Breakdown: dto.FaceAnalysisBreakdown{
							FacialSimilarityScore: c.IdentityCheck.Result.Breakdown.FaceAnalysis.Breakdown.FacialSimilarityScore,
						},
					},
					AuthenticityAnalysis: dto.AuthenticityAnalysis{
						SpoofedImageAnalysis:            c.IdentityCheck.Result.Breakdown.AuthenticityAnalysis.SpoofedImageAnalysis,
						LivenessCheck:                   c.IdentityCheck.Result.Breakdown.AuthenticityAnalysis.LivenessCheck,
						LivenessVoiceChallengeAnalysis:  c.IdentityCheck.Result.Breakdown.AuthenticityAnalysis.LivenessVoiceChallengeAnalysis,
						LivenessActionChallengeAnalysis: c.IdentityCheck.Result.Breakdown.AuthenticityAnalysis.LivenessActionChallengeAnalysis,
						Breakdown: dto.AuthenticityAnalysisBreakdown{
							LivenessCheckScore: c.IdentityCheck.Result.Breakdown.AuthenticityAnalysis.Breakdown.LivenessCheckScore,
						},
					},
				},
			},
			Metadata: dto.IdentityMetadata{
				LiveVideo: dto.LiveVideoMetadata{
					Language: c.IdentityCheck.Metadata.LiveVideo.Language,
				},
			},
			CreatedAt: c.IdentityCheck.CreatedAt,
			UpdatedAt: c.IdentityCheck.UpdatedAt,
		},
		IdentityOutcome: c.IdentityOutcome,
		IdentityStatus:  c.IdentityStatus,
		UpdatedAt:       c.UpdatedAt,
	}
}

func MapSelfActivationUserInfo(u *imodel.UserInfo) *dto.UserInfo {
	if u == nil {
		return nil
	}

	return &dto.UserInfo{
		ID:          u.ID,
		UserCode:    u.UserCode,
		FullName:    u.FullName,
		Email:       u.Email,
		Department:  u.Department,
		PhoneNumber: u.PhoneNumber,
	}
}
