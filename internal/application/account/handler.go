package account

import (
	"context"
	accountPort "cbe-super-app-member-users/internal/port/inbound/account"
	domainAccount "cbe-super-app-member-users/internal/domain/account"
)

type ApplicationHandler struct {
	accountService *domainAccount.AccountService
}

func NewApplicationHandler(accountService *domainAccount.AccountService) *ApplicationHandler {
	return &ApplicationHandler{
		accountService: accountService,
	}
}

type CreateAccountRequest struct {
	UserID string `json:"user_id"`
}

type CreateAccountResponse struct {
	CustomerNumber string `json:"customer_number"`
	AccountNumber  string `json:"account_number"`
}

func (h *ApplicationHandler) CreateAccount(ctx context.Context, userID string) (*accountPort.AccountCreationResult, error) {
	domainResult, err := h.accountService.CreateAccount(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return &accountPort.AccountCreationResult{
		CustomerNumber: domainResult.CustomerNumber,
		AccountNumber:  domainResult.AccountNumber,
	}, nil
}

func (h *ApplicationHandler) HandleCreateAccount(ctx context.Context, req CreateAccountRequest) (*CreateAccountResponse, error) {
	portResult, err := h.CreateAccount(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	
	return &CreateAccountResponse{
		CustomerNumber: portResult.CustomerNumber,
		AccountNumber:  portResult.AccountNumber,
	}, nil
}