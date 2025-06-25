package users_inbound

import (
	// "context"
	"net/http"
)

type InBound interface {
	FetchLinkedAccounts(w http.ResponseWriter, r *http.Request)
	GenerateEmailOTP(w http.ResponseWriter, r *http.Request)
	VerifyEmailOTP(w http.ResponseWriter, r *http.Request)
	UpdateProfilePicture(w http.ResponseWriter, r *http.Request)
	UnlinkDevice(w http.ResponseWriter, r *http.Request)
}

// type UserPort interface {
// 	FetchLinkedAccounts(ctx context.Context, userID string) (*LinkedAccountResponse, error)
// 	GenerateEmailOTP(ctx context.Context, req OTPRequest) (string, error)
// 	VerifyEmailOTP(ctx context.Context, verification OTPVerification) error
// }

type LinkedAccountDetail struct {
	AccountNumber     string `json:"account_number"`
	AccountBranchCode string `json:"account_branch_code"`
	LinkedBranch      string `json:"linked_branch"`
	IsAccountActive   bool   `json:"is_account_active"`
	LinkedStatus      bool   `json:"linked_status"`
	CurrencyCode      string `json:"currency_code"`
}

type LinkedAccountResponse struct {
	UserID         string                `json:"user_id"`
	FullName       string                `json:"full_name"`
	LinkedAccounts []LinkedAccountDetail `json:"linked_accounts"`
}

type OTPRequest struct {
	UserID string
	Email  string
}

type OTPVerification struct {
	UserID string
	Email  string
	OTP    string
}

