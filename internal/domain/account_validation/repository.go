package account_validation

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type Repository interface {
	GetAccountValidationByID(ctx context.Context, id string) (ValidationRule, error)
	UpdateAccountValidation(ctx context.Context, id string, update ValidationRule) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error)
}
