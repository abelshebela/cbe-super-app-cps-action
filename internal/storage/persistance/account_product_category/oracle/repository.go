package account_product_category_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	apc_core "github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/account_product_category/oracle/core"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type repository struct {
	db     DBTX
	logger utils.Logger
}

var _ storage.AccountProductCategoryRepository = (*repository)(nil)

func NewRepository(db DBTX, logger utils.Logger) storage.AccountProductCategoryRepository {
	return &repository{db: db, logger: logger}
}

func (r *repository) withTx(tx *sql.Tx) *repository {
	return &repository{db: tx, logger: r.logger}
}

func (r *repository) Create(ctx context.Context, apc *imodel.AccountProductCategory) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][Create] code=%s name=%s", apc.CBSCategoryCode, apc.CategoryName)

	if dbHandle, ok := r.db.(*sql.DB); ok {
		tx, err := dbHandle.BeginTx(ctx, nil)
		if err != nil {
			log.Errorf("[APCOracle][Create] BeginTx: %v", err)
			return local_util.HandleDBError(err)
		}
		defer func() { _ = tx.Rollback() }()
		if err := r.withTx(tx).createInner(ctx, apc); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			log.Errorf("[APCOracle][Create] Commit: %v", err)
			return local_util.HandleDBError(err)
		}
		return nil
	}
	return r.createInner(ctx, apc)
}

func (r *repository) createInner(ctx context.Context, apc *imodel.AccountProductCategory) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `INSERT INTO ACCOUNT_CATEGORIES
		(ACCOUNT_TYPE, CATEGORY_NAME, CBS_CATEGORY_CODE, DESCRIPTION, IS_ENABLED, IS_DELETED)
		VALUES (:1, :2, :3, :4, :5, :6)`

	if _, err := r.db.ExecContext(ctx, q,
		apc.AccountType,
		apc.CategoryName,
		apc.CBSCategoryCode,
		apc_core.NullStringFromString(apc.Description),
		apc_core.BoolToInt(apc.IsEnabled),
		apc_core.BoolToInt(apc.IsDeleted),
	); err != nil {
		log.Errorf("[APCOracle][Create] insert: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, apc *imodel.AccountProductCategory) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][Update] id=%s", id)

	q := `UPDATE ACCOUNT_CATEGORIES
		SET ACCOUNT_TYPE = :1, CATEGORY_NAME = :2, CBS_CATEGORY_CODE = :3,
		    DESCRIPTION = :4, LAST_MODIFIED_AT = :5
		WHERE ID = HEXTORAW(:6) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q,
		apc.AccountType,
		apc.CategoryName,
		apc.CBSCategoryCode,
		apc_core.NullStringFromString(apc.Description),
		time.Now(),
		id,
	)
	if err != nil {
		log.Errorf("[APCOracle][Update] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][Delete] id=%s", id)

	var linkedCount int64
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM ACCOUNT_PRODUCTS WHERE ACCOUNT_CATEGORY_ID = HEXTORAW(:1) AND IS_DELETED = 0`,
		id,
	).Scan(&linkedCount); err != nil {
		log.Errorf("[APCOracle][Delete] linked products check: %v", err)
		return local_util.HandleDBError(err)
	}
	if linkedCount > 0 {
		log.Errorf("[APCOracle][Delete] category id=%s has %d active product(s)", id, linkedCount)
		return errors.New(localization.ErrorAPCHasActiveProducts.Code)
	}

	res, err := r.db.ExecContext(ctx,
		`DELETE FROM ACCOUNT_CATEGORIES WHERE ID = HEXTORAW(:1)`,
		id,
	)
	if err != nil {
		log.Errorf("[APCOracle][Delete] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][EnableOrDisable] id=%s enable=%v", id, enable)

	q := `UPDATE ACCOUNT_CATEGORIES
		SET IS_ENABLED = :1, LAST_MODIFIED_AT = :2
		WHERE ID = HEXTORAW(:3) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, apc_core.BoolToInt(enable), time.Now(), id)
	if err != nil {
		log.Errorf("[APCOracle][EnableOrDisable] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*imodel.AccountProductCategory, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][FindByID] id=%s", id)

	q := `SELECT ` + apcSelectCols + apcFromTable + `WHERE ID = HEXTORAW(:1) AND IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, id)
	out, err := scanAPCRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[APCOracle][FindByID] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindByCBSCode(ctx context.Context, code string) (*imodel.AccountProductCategory, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][FindByCBSCode] code=%s", code)

	q := `SELECT ` + apcSelectCols + apcFromTable + `WHERE CBS_CATEGORY_CODE = :1 AND IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, code)
	out, err := scanAPCRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[APCOracle][FindByCBSCode] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindByCategoryName(ctx context.Context, name string) (*imodel.AccountProductCategory, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APCOracle][FindByCategoryName] name=%s", name)

	q := `SELECT ` + apcSelectCols + apcFromTable + `WHERE UPPER(CATEGORY_NAME) = UPPER(:1) AND IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, name)
	out, err := scanAPCRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[APCOracle][FindByCategoryName] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.AccountProductCategory], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	perPage := filterParam.PerPage
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	conds := []string{"IS_DELETED = 0"}
	args := []interface{}{}
	idx := 1

	if s := strings.TrimSpace(filterParam.Search); s != "" {
		conds = append(conds, fmt.Sprintf(
			"(UPPER(CATEGORY_NAME) LIKE UPPER(:%d) OR UPPER(CBS_CATEGORY_CODE) LIKE UPPER(:%d))",
			idx, idx+1,
		))
		args = append(args, "%"+s+"%", "%"+s+"%")
		idx += 2
	}

	if filterParam.Filters != nil {
		if v, ok := filterParam.Filters["account_type"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("UPPER(ACCOUNT_TYPE) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if b, applied := lib.BoolFromInterface(v); applied {
				conds = append(conds, fmt.Sprintf("IS_ENABLED = :%d", idx))
				args = append(args, apc_core.BoolToInt(b))
				idx++
			}
		}
	}

	where := strings.Join(conds, " AND ")

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ACCOUNT_CATEGORIES WHERE `+where, args...).Scan(&total); err != nil {
		log.Errorf("[APCOracle][FindAll] count: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]imodel.AccountProductCategory]{Data: []imodel.AccountProductCategory{}, Meta: meta}, nil
	}

	selectQ := `SELECT ` + apcSelectCols + apcFromTable + `WHERE ` + where +
		` ORDER BY CREATED_AT DESC OFFSET :` + fmt.Sprintf("%d", idx) +
		` ROWS FETCH NEXT :` + fmt.Sprintf("%d", idx+1) + ` ROWS ONLY`
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQ, args...)
	if err != nil {
		log.Errorf("[APCOracle][FindAll] query: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	out := []imodel.AccountProductCategory{}
	for rows.Next() {
		item, err := scanAPCRow(rows)
		if err != nil {
			log.Errorf("[APCOracle][FindAll] scan: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return &types.PaginatedResponse[[]imodel.AccountProductCategory]{Data: out, Meta: meta}, nil
}
