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
	FullName       FullNameDto             `json:"full_name"`
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
		FullName: FullNameDto{
			FirstName:  response.FullName.FirstName,
			MiddleName: response.FullName.MiddleName,
			LastName:   response.FullName.LastName,
		},
		LinkedAccounts: dtoAccounts,
	}
}