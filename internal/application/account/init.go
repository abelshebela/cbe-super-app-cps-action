package account_application

import (
	"context"

	domainAccount "cbe-super-app-member-users/internal/domain/account"
	accountPort "cbe-super-app-member-users/internal/port/inbound/account"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	CreateAccount(ctx context.Context, userID string) (*accountPort.AccountCreationResult, error)
}

type AccountHandler struct {
	accountService *domainAccount.AccountService
	logger         utils.Logger
}

func InitAccountHandler(accountService *domainAccount.AccountService, logger utils.Logger) ApplicationService {
	return AccountHandler{
		accountService: accountService,
		logger:         logger,
	}
}
