package bps_user

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bps_user"
)

type BPSUserPersistence interface {
	// Basic CRUD operations
	GetBPSUserByUserCode(ctx context.Context, userCode string) (*bps_user.BPSUser, error)
	GetAllBPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*bps_user.BPSUser], error)
	EnableBPSUser(ctx context.Context, userCode string) error
	DisableBPSUser(ctx context.Context, userCode string) error
	UpdateBPSUserStatus(ctx context.Context, userCode string, enabled bool, updatedAt time.Time) error
	GetCPSActionByID(ctx context.Context, actionID string) (*model.CPSAction, error)
	// UpdateCPSAction(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
	CheckPendingAction(ctx context.Context, requestAction string, department string) (bool, error)
	AuthorizeBPSUserEnable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
	AuthorizeBPSUserDisable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}
