package access_items_relation_oracle

import (
	"context"
	"database/sql"

	local_model "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository struct {
	db     *sql.DB
	logger utils.Logger
}

func NewRepository(db *sql.DB, logger utils.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// FindParentChildRelationship implements the same contract as Mongo aggregation on access_items_relation.
func (r *Repository) FindParentChildRelationship(ctx context.Context) ([]local_model.AccessItemRelation, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `SELECT RAWTOHEX(PARENT_ID), RAWTOHEX(CHILD_ID) FROM ACCESS_LIST_RELATIONS ORDER BY PARENT_ID, CHILD_ID`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		log.Errorf("[AccessItemsRelationOracle][FindParentChildRelationship] query failed: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	var out []local_model.AccessItemRelation
	for rows.Next() {
		var rel local_model.AccessItemRelation
		if err := rows.Scan(&rel.ParentKey, &rel.ChildKey); err != nil {
			log.Errorf("[AccessItemsRelationOracle][FindParentChildRelationship] scan failed: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, rel)
	}
	if err := rows.Err(); err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}
