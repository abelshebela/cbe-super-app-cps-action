package miniapp_oracle

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
	"regexp"
	"time"

	querypkg "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/mini_app/oracle/query"

	// sharedModel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	constants "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"


)
type miniAppCategoryOraclePersistence struct {
	logger shared_utils.Logger
	db     *sql.DB
}

func NewCategoryOracleRepository(logger shared_utils.Logger, db *sql.DB) storage.MiniAppCategoryRepository {
	return &miniAppCategoryOraclePersistence{logger: logger, db: db}
}

func (m *miniAppCategoryOraclePersistence) Create(ctx context.Context, category *local_model.MiniAppCategory) error {
	q := querypkg.MiniAppCategoryInsert
	_, err := m.db.ExecContext(
		ctx, q,
		 category.Name, category.Icon, boolToInt64(category.IsEnabled), boolToInt64(category.IsDeleted),
		category.CreatedAt, category.UpdatedAt,
	)
	if err != nil {
		return constants.ErrDatabaseError
	}
	return nil
}

func (m *miniAppCategoryOraclePersistence) Update(ctx context.Context, category *local_model.MiniAppCategory, id string) error {
	if !isHexID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppCategoryUpdate
	res, err := m.db.ExecContext(ctx, q, category.Name, category.Icon, category.UpdatedAt, id)
	if err != nil {
		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrCategoryNotFound
	}
	return nil
}

func (m *miniAppCategoryOraclePersistence) Delete(ctx context.Context, id string) error {
	return m.EnableOrDisableDelete(ctx, id, true)
}

func (m *miniAppCategoryOraclePersistence) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	if !isHexID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppCategoryToggleEnabled
	res, err := m.db.ExecContext(ctx, q, boolToInt64(enable), time.Now(), id)
	if err != nil {
		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrCategoryNotFound
	}
	return nil
}

func (m *miniAppCategoryOraclePersistence) EnableOrDisableDelete(ctx context.Context, id string, isDelete bool) error {
	if !isHexID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppCategorySoftDelete
	res, err := m.db.ExecContext(ctx, q, time.Now(), time.Now(), id)
	if err != nil {
		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 && isDelete {
		return constants.ErrCategoryNotFound
	}
	return nil
}

func boolToInt64(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

func isHexID(v string) bool {
	if len(v) != 32 {
		return false
	}
	ok, _ := regexp.MatchString("^[0-9a-fA-F]{32}$", v)
	return ok
}