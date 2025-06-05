package users

import (
	 domain "cbe-super-app-member-users/internal/domain/users"
)

type FetchLinkedAccountsRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type FullNameDto struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`
}

type LinkedAccountDetailDto struct {
	AccountNumber     string `json:"account_number"`
	AccountBranchCode string `json:"account_branch_code"`
	LinkedBranch      string `json:"linked_branch"`
	IsAccountActive   bool   `json:"is_account_active"`
	LinkedStatus      bool   `json:"linked_status"`
	CurrencyCode      string `json:"currency_code"`
}

type LinkedAccountResponseDto struct {
	UserID         string                  `json:"user_id"`
	FullName       string             `json:"full_name"`
	LinkedAccounts []LinkedAccountDetailDto `json:"linked_accounts"`
}

func ToLinkedAccountResponseDto(response *domain.LinkedAccountResponse) *LinkedAccountResponseDto {
	if response == nil {
		return nil
	}
	dtoAccounts := make([]LinkedAccountDetailDto, len(response.LinkedAccounts))
	for i, account := range response.LinkedAccounts {
		dtoAccounts[i] = LinkedAccountDetailDto{
			AccountNumber:     account.AccountNumber,
			AccountBranchCode: account.AccountBranchCode,
			LinkedBranch:      account.LinkedBranch,
			IsAccountActive:   account.IsAccountActive,
			LinkedStatus:      account.LinkedStatus,
			CurrencyCode:      account.CurrencyCode,
		}
	}
	return &LinkedAccountResponseDto{
        UserID: response.UserID,
        FullName: response.FullName,
        LinkedAccounts: dtoAccounts,
    }
}


type GenerateOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

type GenerateOTPResponse struct {
	OTP string `json:"otp"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	OTP    string `json:"otp" binding:"required"`
}

type VerifyOTPResponse struct {
	Success bool   `json:"success"`
	Email   string `json:"email"`
}