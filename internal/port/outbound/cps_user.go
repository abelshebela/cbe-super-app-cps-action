package outbound

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type OutboundInfra interface {
	CreateUserRequest(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, cpsAction model.CPSAction, userCode string) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionCode string, checker model.CPSAction) (*model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*model.CPSUser, error)
	GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*model.CPSUser], error)
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
	AuthorizeCreate(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeUpdate(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)

	// UpdateOneServiceDetailRequest(ctx context.Context, id string, update service.Service) error
	// UpdateOneSeviceDeatil(ctx context.Context, pd service.Service) error
}
