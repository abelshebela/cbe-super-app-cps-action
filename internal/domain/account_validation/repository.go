package account_validation

import (
	"cbe-super-app-cps-action/internal/domain/action"
	"context"
)

type Repository interface {
	GetAccountValidationByID(ctx context.Context, id string) (ValidationRule, error)
	UpdateAccountValidation(ctx context.Context, id string, update ValidationRule) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
}
