package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
)

type CPSUserRepo interface {
	CreateUserRequest(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, cpsAction model.CPSAction, userCode string) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionCode string, checker model.CPSAction) (*model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
}
