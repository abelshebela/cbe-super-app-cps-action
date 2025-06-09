package account

type CreateAccountRequest struct {
    UserID string `json:"user_id" binding:"required"`
}

type CreateAccountResponse struct {
    CustomerNumber string `json:"customer_number"`
    AccountNumber  string `json:"account_number"`
}