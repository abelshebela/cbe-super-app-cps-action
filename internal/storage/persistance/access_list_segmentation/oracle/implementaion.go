package access_list_segmentation_oracle

import (
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	shared_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type accessListSegmentationOracle struct {
	db     DBTX
	logger utils.Logger
}

func NewAccessListSegmentationOracle(db DBTX, logger utils.Logger) storage.AccessListSegmentationRepository {
	return &accessListSegmentationOracle{db: db, logger: logger}
}

// BulkDisable implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) BulkDisable(ctx context.Context, req access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest) error {
	panic("unimplemented")
}

// CreateAccountSegment implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) CreateAccountSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	panic("unimplemented")
}

// CreateBlockSegment implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) CreateBlockSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error {
	panic("unimplemented")
}

// EnableOrDisable implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	panic("unimplemented")
}

// FindAllBySegmentIDorSegmentCode implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllBySegmentIDorSegmentCode(ctx context.Context, segmentIDorCode string) ([]shared_model.APPAccessList, error) {
	panic("unimplemented")
}

// FindAllBySegmentIDorSegmentCodeAndKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllBySegmentIDorSegmentCodeAndKeys(ctx context.Context, segmentIDorCode string, keys []string) ([]model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AccessListSegmentation], error) {
	panic("unimplemented")
}

// FindByAccountSegmentationAndAccessListKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByAccountSegmentationAndAccessListKeys(ctx context.Context, customerSegments string, segmentKeys []string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindByID implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByID(ctx context.Context, id string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindByIDAndType implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByIDAndType(ctx context.Context, ids string, t string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindByIDS implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindByIDS(ctx context.Context, ids []string, t string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindBySegmentIDAndAccessListKeys implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindBySegmentIDAndAccessListKeys(ctx context.Context, id string, keys []string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindBySegmentationAndServiceID implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindBySegmentationAndServiceID(ctx context.Context, segmentationID string, serviceID string) (*model.AccessListSegmentation, error) {
	panic("unimplemented")
}

// FindParentChildRelationship implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) FindParentChildRelationship(ctx context.Context) ([]model.AccessItemRelation, error) {
	panic("unimplemented")
}

// Update implements [storage.AccessListSegmentationRepository].
func (q *accessListSegmentationOracle) Update(ctx context.Context, id string, accessListSegmentation model.AccessListSegmentation) error {
	panic("unimplemented")
}
