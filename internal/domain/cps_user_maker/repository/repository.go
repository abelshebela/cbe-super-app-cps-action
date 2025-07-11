package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type CPSUserRepo interface {
	CreateUserRequest(ctx context.Context, action action.CPSUser, maker action.User) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) (*model.CPSAction, error)

	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
}
