package updatedbulkservice

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	bulk_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/bulk_service"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type BulkServiceRepository interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*bulk_entity.APPAccessList], error)
	EnableOrDisableBulkService(ctx context.Context, keys []string, cpsAction model.CPSAction, requestActionType model.RequestAction) error
	AuthorizeBulkServiceEnable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
	AuthorizeBulkServiceDisable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
}
