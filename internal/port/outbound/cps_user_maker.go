package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
)

type OutboundInfra interface {
	CreateUserRequest(ctx context.Context, user action.CPSUser, maker action.User) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, updated action.CPSUser, maker action.User) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context, actionCode string) ([]action.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionID string, approve bool, reason *string) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*action.CPSUser, error)
	UpdateOneServiceDetailRequest(ctx context.Context, id string, update service.Service) error
	//UpdateOneSeviceDeatil(ctx context.Context, pd service.Service) error
}
