package outbound

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type OutboundInfra interface {
	CreateUserRequest(ctx context.Context, user action.CPSUser, maker action.User) error
	UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) error
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
}
