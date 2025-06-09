package account

import (
    "context"
    "cbe-super-app-member-users/internal/domain/account"
)

type ApplicationHandler interface {
    CreateAccount(ctx context.Context, req CreateAccountRequest) (*CreateAccountResponse, error)
}

type applicationHandler struct {
    domainService *account.AccountService
}

func NewApplicationHandler(domainService *account.AccountService) ApplicationHandler {
    return &applicationHandler{domainService: domainService}
}

func (h *applicationHandler) CreateAccount(ctx context.Context, req CreateAccountRequest) (*CreateAccountResponse, error) {
    result, err := h.domainService.CreateAccount(ctx, req.UserID)
    if err != nil {
        return nil, err
    }
    return &CreateAccountResponse{
        CustomerNumber: result.CustomerNumber,
        AccountNumber:  result.AccountNumber,
    }, nil
}