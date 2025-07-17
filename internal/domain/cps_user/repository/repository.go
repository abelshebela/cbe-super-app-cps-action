package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
)

type CPSUserRepo interface {
	CreateUserRequest(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, cpsAction model.CPSAction, userCode string) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionCode string, checker model.CPSAction) (*model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error)
	GetAllCPSUsers(ctx context.Context) ([]model.CPSUser, error)
}
