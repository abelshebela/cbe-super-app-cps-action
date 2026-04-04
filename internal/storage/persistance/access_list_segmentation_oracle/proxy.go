package access_list_segmentation_oracle

import (
	"context"
	"database/sql"

	local_model "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	access_items_relation_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_items_relation_oracle"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// accessListSegmentationProxy embeds the Mongo implementation and overrides only
// FindParentChildRelationship to read edges from Oracle (ACCESS_ITEMS_RELATION).
type accessListSegmentationProxy struct {
	storage.AccessListSegmentationRepository
	relations *access_items_relation_oracle.Repository
}

func NewAccessListSegmentationProxy(delegate storage.AccessListSegmentationRepository, db *sql.DB, log utils.Logger) storage.AccessListSegmentationRepository {
	if db == nil {
		return delegate
	}
	return &accessListSegmentationProxy{
		AccessListSegmentationRepository: delegate,
		relations:                        access_items_relation_oracle.NewRepository(db, log),
	}
}

func (p *accessListSegmentationProxy) FindParentChildRelationship(ctx context.Context) ([]local_model.AccessItemRelation, error) {
	return p.relations.FindParentChildRelationship(ctx)
}
