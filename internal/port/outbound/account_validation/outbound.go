package outbound

import (
	"cbe-super-app-cps-action/internal/domain/account_validation"
	"cbe-super-app-cps-action/internal/domain/action"
	"context"
)

type OutboundInfra interface {
	GetAccountValidationByID(ctx context.Context, id string) (account_validation.ValidationRule, error)
	UpdateAccountValidation(ctx context.Context, id string, update account_validation.ValidationRule) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
}
