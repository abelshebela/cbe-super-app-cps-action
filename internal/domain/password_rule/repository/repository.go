package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type PasswordRuleRepository interface {
	GetAllPasswordRules(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*action.PasswordRule], error)
	GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error)
	GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error)
	GetCurrentPasswordRule(ctx context.Context) (*action.PasswordRule, error)

	UpdatePasswordRule(ctx context.Context, rule action.PasswordRule) error
	UpdateCpsAction(ctx context.Context, action action.CPSAction) error
	CreateCpsAction(ctx context.Context, action action.CPSAction) (action.CPSAction, error)
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
}
