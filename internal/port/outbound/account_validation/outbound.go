package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/account_validation"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundInfra interface {
	GetAccountValidationByID(ctx context.Context, id string) (account_validation.ValidationRule, error)
	UpdateAccountValidation(ctx context.Context, id string, update account_validation.ValidationRule) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
}
