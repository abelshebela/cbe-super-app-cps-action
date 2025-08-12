package bps_user

import (
	"context"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type Repository interface {
	// Basic CRUD operations
	GetBPSUserByUserCode(ctx context.Context, userCode string) (*BPSUser, error)
	GetAllBPSUsers(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*BPSUser], error)

	// Enable/Disable operations with maker-checker flow
	EnableBPSUser(ctx context.Context, userCode string) error
	DisableBPSUser(ctx context.Context, userCode string) error
	UpdateBPSUserStatus(ctx context.Context, userCode string, enabled bool, updatedAt time.Time) error

	GetCPSActionByID(ctx context.Context, actionID string) (*model.CPSAction, error)

	CheckPendingAction(ctx context.Context, requestAction string, department string) (bool, error)

	// Authorization methods
	AuthorizeBPSUserEnable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
	AuthorizeBPSUserDisable(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}
