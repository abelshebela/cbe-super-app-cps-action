package account_application

import (
	"context"

	accountPort "cbe-super-app-member-users/internal/port/inbound/account"

	"cbe-super-app-member-users/internal/application/dto"
)

func (h AccountHandler) CreateAccount(ctx context.Context, userID string) (*accountPort.AccountCreationResult, error) {
	domainResult, err := h.accountService.CreateAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &accountPort.AccountCreationResult{
		CustomerNumber: domainResult.CustomerNumber,
		AccountNumber:  domainResult.AccountNumber,
	}, nil
}

func (h AccountHandler) HandleCreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	portResult, err := h.CreateAccount(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &dto.CreateAccountResponse{
		CustomerNumber: portResult.CustomerNumber,
		AccountNumber:  portResult.AccountNumber,
	}, nil
}
