package updatedbulkservice

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	bulk_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/updated_bulk_service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type BulkServiceRepository interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*bulk_entity.ServiceDetails], error)
	EnableOrDisableBulkService(ctx context.Context, userCode string, cpsAction model.CPSAction, requestActionType model.RequestAction) error
	AuthorizeBulkServiceEnable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeBulkServiceDisable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
}
