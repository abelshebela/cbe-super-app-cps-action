package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundPasswordRuleInfra interface {
	GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error)
	GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error)
	GetCurrentPasswordRule(ctx context.Context) (*action.PasswordRule, error)
	UpdatePasswordRule(ctx context.Context, rule action.PasswordRule) error
	UpdateCpsAction(ctx context.Context, action action.CPSAction) error
	CreateCpsAction(ctx context.Context, action action.CPSAction) (action.CPSAction, error)
}
