package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundPasswordRuleInfra interface {
	CreatePasswordRuleUpdateAction(ctx context.Context, rule *action.PasswordRule, maker action.User) (actionID string, err error)
	ApproveOrRejectPasswordRuleAction(ctx context.Context, actionID string, approve bool, checker action.User, rejectionReason *string) error
	GetPasswordRuleUpdateActionByID(ctx context.Context, actionID string) (*action.CPSAction, error)
	GetUpdateAction(ctx context.Context, maker action.User) (*action.CPSAction, error)
	GetCurrentPasswordRule(ctx context.Context) (*action.PasswordRule, error)
}
