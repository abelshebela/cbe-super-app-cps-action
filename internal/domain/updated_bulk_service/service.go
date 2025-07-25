package updatedbulkservice

import (
	"context"

	// "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	// action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type BulkService interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*ServiceDetails], error)
	// EnableDisableUser(ctx context.Context, userCode string, cpsAction model.CPSAction, requestActionType model.RequestAction) error
	// AuthorizeUserEnable(ctx context.Context, action *action_entity.CPSAction) (*action_entity.CPSAction, error)
	// AuthorizeUserDisable(ctx context.Context, action *action_entity.CPSAction) (*action_entity.CPSAction, error)
}

type bulkServiceImpl struct {
	repo   BulkServiceRespository
	logger utils.Logger
}

func NewBulkService(repo BulkServiceRespository, logger utils.Logger) BulkService {
	return &bulkServiceImpl{repo: repo, logger: logger}
}

func (b *bulkServiceImpl) GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*ServiceDetails], error) {
	return b.repo.GetAllBulkServices(ctx, filterParams)
}

// func (b *bulkServiceImpl) EnableDisableUser(ctx context.Context, userCode string, cpsAction model.CPSAction, requestActionType model.RequestAction) error {
// }

// func (b *bulkServiceImpl) AuthorizeUserEnable(ctx context.Context, action *action_entity.CPSAction) (*action_entity.CPSAction, error) {
// }

// func (b *bulkServiceImpl) AuthorizeUserDisable(ctx context.Context, action *action_entity.CPSAction) (*action_entity.CPSAction, error) {
// 	return nil, nil
// }
