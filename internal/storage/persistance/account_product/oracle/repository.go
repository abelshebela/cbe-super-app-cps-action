package account_product_oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	ap_core "cbe-super-app-cps-action/internal/storage/persistance/account_product/oracle/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type repository struct {
	db     DBTX
	logger utils.Logger
}

var _ storage.AccountProductRepository = (*repository)(nil)

func NewRepository(db DBTX, logger utils.Logger) storage.AccountProductRepository {
	return &repository{db: db, logger: logger}
}

func (r *repository) withTx(tx *sql.Tx) *repository {
	return &repository{db: tx, logger: r.logger}
}

func (r *repository) Create(ctx context.Context, ap *imodel.AccountProduct) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APOracle][Create] code=%s name=%s", ap.CBSProductCode, ap.ProductName)

	if dbHandle, ok := r.db.(*sql.DB); ok {
		tx, err := dbHandle.BeginTx(ctx, nil)
		if err != nil {
			log.Errorf("[APOracle][Create] BeginTx: %v", err)
			return local_util.HandleDBError(err)
		}
		defer func() { _ = tx.Rollback() }()
		if err := r.withTx(tx).createInner(ctx, ap); err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			log.Errorf("[APOracle][Create] Commit: %v", err)
			return local_util.HandleDBError(err)
		}
		return nil
	}

	return r.createInner(ctx, ap)
}

func (r *repository) createInner(ctx context.Context, ap *imodel.AccountProduct) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `INSERT INTO ACCOUNT_PRODUCTS
		(CBS_PRODUCT_CODE, PRODUCT_NAME, PRODUCT_TAG_LINE,
		 ACCOUNT_CATEGORY_ID, ACCOUNT_CURRENCY,
		 MINIMUM_OPENING_BALANCE, MINIMUM_MAINTENANCE_FEE, INTEREST_FEE,
		 FAQ_URL, PRODUCT_FEATURES, HAS_PHYSICAL_CARD, HAS_VIRTUAL_CARD,
		 PRODUCT_ICON, PRODUCT_COVER_IMAGE, IS_ENABLED, IS_DELETED)
		VALUES (:1, :2, :3, HEXTORAW(:4), :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16)`

	if _, err := r.db.ExecContext(ctx, q,
		ap.CBSProductCode,
		ap.ProductName,
		ap.ProductTagLine,
		ap.AccountCategoryID,
		ap.AccountCurrency,
		ap.MinimumOpeningBalance,
		ap.MinimumMaintenanceFee,
		ap.InterestFee,
		ap_core.NullStringFromString(ap.FaqURL),
		ap_core.NullStringFromString(ap.ProductFeatures),
		ap_core.BoolToInt(ap.HasPhysicalCard),
		ap_core.BoolToInt(ap.HasVirtualCard),
		ap_core.NullStringFromString(ap.ProductIcon),
		ap_core.NullStringFromString(ap.ProductCoverImage),
		ap_core.BoolToInt(ap.IsEnabled),
		ap_core.BoolToInt(ap.IsDeleted),
	); err != nil {
		log.Errorf("[APOracle][Create] insert: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, ap *imodel.AccountProduct) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APOracle][Update] id=%s", id)

	q := `UPDATE ACCOUNT_PRODUCTS
		SET CBS_PRODUCT_CODE = :1, PRODUCT_NAME = :2, PRODUCT_TAG_LINE = :3,
		    ACCOUNT_CATEGORY_ID = HEXTORAW(:4), ACCOUNT_CURRENCY = :5,
		    MINIMUM_OPENING_BALANCE = :6, MINIMUM_MAINTENANCE_FEE = :7, INTEREST_FEE = :8,
		    FAQ_URL = :9, PRODUCT_FEATURES = :10, HAS_PHYSICAL_CARD = :11,
		    HAS_VIRTUAL_CARD = :12, PRODUCT_ICON = :13, PRODUCT_COVER_IMAGE = :14,
		    LAST_MODIFIED_AT = :15
		WHERE ID = HEXTORAW(:16) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q,
		ap.CBSProductCode,
		ap.ProductName,
		ap.ProductTagLine,
		ap.AccountCategoryID,
		ap.AccountCurrency,
		ap.MinimumOpeningBalance,
		ap.MinimumMaintenanceFee,
		ap.InterestFee,
		ap_core.NullStringFromString(ap.FaqURL),
		ap_core.NullStringFromString(ap.ProductFeatures),
		ap_core.BoolToInt(ap.HasPhysicalCard),
		ap_core.BoolToInt(ap.HasVirtualCard),
		ap_core.NullStringFromString(ap.ProductIcon),
		ap_core.NullStringFromString(ap.ProductCoverImage),
		time.Now(),
		id,
	)
	if err != nil {
		log.Errorf("[APOracle][Update] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APOracle][Delete] id=%s", id)

	q := `UPDATE ACCOUNT_PRODUCTS
		SET IS_DELETED = 1, IS_ENABLED = 0, DELETED_AT = :1, LAST_MODIFIED_AT = :1
		WHERE ID = HEXTORAW(:2) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, time.Now(), id)
	if err != nil {
		log.Errorf("[APOracle][Delete] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APOracle][EnableOrDisable] id=%s enable=%v", id, enable)

	q := `UPDATE ACCOUNT_PRODUCTS
		SET IS_ENABLED = :1, LAST_MODIFIED_AT = :2
		WHERE ID = HEXTORAW(:3) AND IS_DELETED = 0`

	res, err := r.db.ExecContext(ctx, q, ap_core.BoolToInt(enable), time.Now(), id)
	if err != nil {
		log.Errorf("[APOracle][EnableOrDisable] exec: %v", err)
		return local_util.HandleDBError(err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	return nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*imodel.AccountProduct, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APOracle][FindByID] id=%s", id)

	q := `SELECT ` + apSelectCols + apFromTable + `WHERE ap.ID = HEXTORAW(:1) AND ap.IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, id)
	out, err := scanAPRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[APOracle][FindByID] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindByCBSCode(ctx context.Context, code string) (*imodel.AccountProduct, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)
	log.Infof("[APOracle][FindByCBSCode] code=%s", code)

	q := `SELECT ` + apSelectCols + apFromTable + `WHERE ap.CBS_PRODUCT_CODE = :1 AND ap.IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, code)
	out, err := scanAPRow(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		log.Errorf("[APOracle][FindByCBSCode] scan: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	return out, nil
}

func (r *repository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.AccountProduct], error) {
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

	conds := []string{"ap.IS_DELETED = 0"}
	args := []interface{}{}
	idx := 1

	if s := strings.TrimSpace(filterParam.Search); s != "" {
		conds = append(conds, fmt.Sprintf(
			"(UPPER(ap.PRODUCT_NAME) LIKE UPPER(:%d) OR UPPER(ap.CBS_PRODUCT_CODE) LIKE UPPER(:%d))",
			idx, idx+1,
		))
		args = append(args, "%"+s+"%", "%"+s+"%")
		idx += 2
	}

	if filterParam.Filters != nil {
		// if v, ok := filterParam.Filters["product_line"]; ok {
		// 	if s, _ := v.(string); strings.TrimSpace(s) != "" {
		// 		conds = append(conds, fmt.Sprintf("UPPER(ap.PRODUCT_LINE) = UPPER(:%d)", idx))
		// 		args = append(args, s)
		// 		idx++
		// 	}
		// }
		if v, ok := filterParam.Filters["account_category_id"]; ok {
			if s, _ := v.(string); strings.TrimSpace(s) != "" {
				conds = append(conds, fmt.Sprintf("ap.ACCOUNT_CATEGORY_ID = HEXTORAW(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParam.Filters["is_enabled"]; ok {
			if b, applied := lib.BoolFromInterface(v); applied {
				conds = append(conds, fmt.Sprintf("ap.IS_ENABLED = :%d", idx))
				args = append(args, ap_core.BoolToInt(b))
				idx++
			}
		}
	}

	where := strings.Join(conds, " AND ")

	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ACCOUNT_PRODUCTS ap WHERE `+where, args...).Scan(&total); err != nil {
		log.Errorf("[APOracle][FindAll] count: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	if total == 0 {
		return &types.PaginatedResponse[[]imodel.AccountProduct]{Data: []imodel.AccountProduct{}, Meta: meta}, nil
	}

	selectQ := `SELECT ` + apSelectCols + apFromTable + `WHERE ` + where +
		` ORDER BY ap.CREATED_AT DESC OFFSET :` + fmt.Sprintf("%d", idx) +
		` ROWS FETCH NEXT :` + fmt.Sprintf("%d", idx+1) + ` ROWS ONLY`
	args = append(args, offset, perPage)

	rows, err := r.db.QueryContext(ctx, selectQ, args...)
	if err != nil {
		log.Errorf("[APOracle][FindAll] query: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer rows.Close()

	out := []imodel.AccountProduct{}
	for rows.Next() {
		item, err := scanAPRow(rows)
		if err != nil {
			log.Errorf("[APOracle][FindAll] scan: %v", err)
			return nil, local_util.HandleDBError(err)
		}
		out = append(out, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return &types.PaginatedResponse[[]imodel.AccountProduct]{Data: out, Meta: meta}, nil
}
