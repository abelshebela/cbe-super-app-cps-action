package dto
type FetchLinkedAccountsRequest struct {
	UserID string `json:"user_id"`
}

type GenerateOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type VerifyOTPRequest struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	OTP    string `json:"otp"`
}

type GenerateOTPResponse struct {
	OTP string `json:"otp"`
}

type VerifyOTPResponse struct {
	Success bool `json:"success"`
}