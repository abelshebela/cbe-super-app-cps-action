package repository

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
)

type CPSUserRepo interface {
    CreateUserRequest(ctx context.Context, action action.CPSUser, maker action.User) error
    UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) error

    GetPendingUserActions(ctx context.Context) ([]action.CPSAction, error)
    ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error
}