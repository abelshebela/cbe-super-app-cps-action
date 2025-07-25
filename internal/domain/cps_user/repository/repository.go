package repository

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	userDTO "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_user/dto"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type CPSUserRepo interface {
	CreateUserRequest(ctx context.Context, cpsAction model.CPSAction) (*model.CPSAction, error)
	UpdateUserRequest(ctx context.Context, cpsAction model.CPSAction, userCode string) (*model.CPSAction, error)
	GetPendingUserActions(ctx context.Context) ([]model.CPSAction, error)
	ApproveUserAction(ctx context.Context, actionCode string, checker model.CPSAction) (*model.CPSAction, error)
	FetchUserByUserCode(ctx context.Context, userCode string) (*userDTO.CPSUserDTO, error)
	GetAllCPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*userDTO.CPSUserDTO], error)
	FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.ActionResponse, error)
	GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error)
	AuthorizeUserCreate(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeUserUpdate(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeUserDelete(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeUserEnable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeUserDisable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	DeleteUserRequest(ctx context.Context, userCode string, cpsAction model.CPSAction) (*model.CPSAction, error)
	EnableDisableUser(ctx context.Context, userCode string, cpsAction model.CPSAction, requestActionType model.RequestAction) error
}
