package updatedbulkservice

import (
	"context"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	bulk_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/updated_bulk_service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type BulkServiceRepository interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*bulk_entity.ServiceDetails], error)
	// EnableDisableUser(ctx context.Context, userCode string, cpsAction model.CPSAction, requestActionType model.RequestAction) error
	// AuthorizeUserEnable(ctx context.Context, action *action_entity.CPSAction) (*action_entity.CPSAction, error)
	// AuthorizeUserDisable(ctx context.Context, action *action_entity.CPSAction) (*action_entity.CPSAction, error)
}
