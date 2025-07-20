package account_validation

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
)

type AccountValidationRepository interface {
	GetAccountValidationByID(ctx context.Context, id string) (ValidationRule, error)
	Authorize(ctx context.Context, cpsAction *entities.CPSAction) (*entities.CPSAction, error)
	UpdateAccountValidation(ctx context.Context, id string, update ValidationRule) error
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
}
