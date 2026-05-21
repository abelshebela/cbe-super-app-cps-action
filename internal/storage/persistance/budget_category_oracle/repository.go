package budget_category_oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Repository struct {
	db            *sql.DB
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

var _ storage.BudgetCategoryOracleRepository = (*Repository)(nil)

func NewBudgetCategoryOracleRepository(db *sql.DB, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.BudgetCategoryOracleRepository {
	return &Repository{
		db:            db,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (r *Repository) Create(ctx context.Context, bc *imodel.BudgetCategoryOracle) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	// Updated to match new table structure with ACCOUNT_TYPE, CREATED_AT, LAST_MODIFIED_AT
	q := `INSERT INTO BUDGET_CATEGORIES (ID, NAME, ACCOUNT_TYPE, COLOR, ICON, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT)
		VALUES (SYS_GUID(), :1, :2, :3, :4, :5, :6, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := r.db.ExecContext(ctx, q,
		bc.Name,
		bc.AccountType,
		bc.Color,
		bc.Icon,
		bc.IsEnabled,
		bc.IsDeleted,
	)
	if err != nil {
		log.Errorf("[BudgetCategoryOracle][Create] insert failed: %v", err)
		return err
	}
	created, err := r.FindByName(ctx, bc.Name)
	if err != nil {
		return err
	}
	if created != nil {
		r.kafkaProducer.PublishMessage(ctx, created, string(constants.ClientOrchestrationBudgetCategoryTopic), string(constants.ClientOrchestrationBudgetCategoryTopic), "new budget category created")
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, id string, bc *imodel.BudgetCategoryOracle) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `UPDATE BUDGET_CATEGORIES SET NAME = :1, ACCOUNT_TYPE = :2, COLOR = :3, ICON = :4, LAST_MODIFIED_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:5)`
	_, err := r.db.ExecContext(ctx, q,
		bc.Name,
		bc.AccountType,
		bc.Color,
		bc.Icon,
		id,
	)
	if err != nil {
		log.Errorf("[BudgetCategoryOracle][Update] failed: %v", err)
		return err
	}
	updated, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if updated != nil {
		r.kafkaProducer.PublishMessage(ctx, updated, string(constants.ClientOrchestrationBudgetCategoryTopic), string(constants.ClientOrchestrationBudgetCategoryTopic), "budget category updated")
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	q := `UPDATE BUDGET_CATEGORIES SET IS_DELETED = 1, DELETED_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:1)`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		log.Errorf("[BudgetCategoryOracle][Delete] failed: %v", err)
		return err
	}
	return nil
}

func (r *Repository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	// Each :n is a distinct bind for godror; repeating :1 still expects one value per placeholder.
	q := `UPDATE BUDGET_CATEGORIES SET IS_ENABLED = :1, LAST_MODIFIED_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:2)`
	var v int
	if enable {
		v = 1
	} else {
		v = 0
	}
	_, err := r.db.ExecContext(ctx, q, v, id)
	if err != nil {
		log.Errorf("[BudgetCategoryOracle][EnableOrDisable] failed: %v", err)
		return err
	}
	updated, err := r.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if updated != nil {
		r.kafkaProducer.PublishMessage(ctx, updated, string(constants.ClientOrchestrationBudgetCategoryTopic), string(constants.ClientOrchestrationBudgetCategoryTopic), "budget category enable/disable updated")
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*imodel.BudgetCategoryOracle, error) {
	// Updated to match new table structure with ACCOUNT_TYPE, CREATED_AT, LAST_MODIFIED_AT
	q := `SELECT RAWTOHEX(ID) AS id, NAME, ACCOUNT_TYPE, COLOR, ICON, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
		FROM BUDGET_CATEGORIES WHERE ID = HEXTORAW(:1)`
	row := r.db.QueryRowContext(ctx, q, id)
	var bc imodel.BudgetCategoryOracle
	var createdAt, lastModifiedAt, deletedAt sql.NullTime
	err := row.Scan(
		&bc.ID,
		&bc.Name,
		&bc.AccountType,
		&bc.Color,
		&bc.Icon,
		&bc.IsEnabled,
		&bc.IsDeleted,
		&createdAt,
		&lastModifiedAt,
		&deletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, localization.ErrorResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	bc.CreatedAt = formatNullTime(createdAt)
	bc.LastModifiedAt = formatNullTime(lastModifiedAt)
	bc.DeletedAt = formatNullTime(deletedAt)
	return &bc, nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (*imodel.BudgetCategoryOracle, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Debugf("[FindByName] Querying for budget category with name: %s", name)
	q := `SELECT RAWTOHEX(ID) AS id, NAME, ACCOUNT_TYPE, COLOR, ICON, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
		FROM BUDGET_CATEGORIES WHERE UPPER(NAME) = UPPER(:1) AND IS_DELETED = 0`
	row := r.db.QueryRowContext(ctx, q, name)
	var bc imodel.BudgetCategoryOracle
	var createdAt, lastModifiedAt, deletedAt sql.NullTime
	err := row.Scan(
		&bc.ID,
		&bc.Name,
		&bc.AccountType,
		&bc.Color,
		&bc.Icon,
		&bc.IsEnabled,
		&bc.IsDeleted,
		&createdAt,
		&lastModifiedAt,
		&deletedAt,
	)
	if err == sql.ErrNoRows {
		log.Infof("[FindByName] No budget category found for name: %s", name)
		return nil, nil
	}
	if err != nil {
		log.Errorf("[FindByName] QueryRowContext/Scan error for name %s: %v", name, err)
		return nil, err
	}
	log.Debugf("[FindByName] Found budget category: ID=%s, Name=%s", bc.ID, bc.Name)
	bc.CreatedAt = formatNullTime(createdAt)
	bc.LastModifiedAt = formatNullTime(lastModifiedAt)
	bc.DeletedAt = formatNullTime(deletedAt)
	log.Debugf("[FindByName] Returning budget category: %+v", bc)
	return &bc, nil
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}

func (r *Repository) FindAllWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.BudgetCategoryOracle], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if filterParams == nil {
		fp := types.Filter{Page: 1, PerPage: 10}
		filterParams = &fp
	}
	var filters []string
	var args []interface{}
	idx := 1

	if filterParams.Search != "" {
		search := "%" + strings.ToUpper(filterParams.Search) + "%"
		// Each :n is a distinct bind for godror; repeating the same idx still expects one value per placeholder.
		filters = append(filters, fmt.Sprintf(`(UPPER(NAME) LIKE :%d OR UPPER(ACCOUNT_TYPE) LIKE :%d)`, idx, idx+1))
		args = append(args, search, search)
		idx += 2
	}

	if filterParams.Filters != nil {
		if v, ok := filterParams.Filters["enabled"]; ok {
			on := false
			switch b := v.(type) {
			case bool:
				on = b
			case string:
				on = strings.EqualFold(b, "true") || b == "1"
			}
			if on {
				filters = append(filters, "IS_ENABLED = 1")
			} else {
				filters = append(filters, "IS_ENABLED = 0")
			}
		}
		if v, ok := filterParams.Filters["name"]; ok {
			if s, ok := v.(string); ok && s != "" {
				filters = append(filters, fmt.Sprintf("UPPER(NAME) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
		if v, ok := filterParams.Filters["type"]; ok {
			if s, ok := v.(string); ok && s != "" {
				filters = append(filters, fmt.Sprintf("UPPER(ACCOUNT_TYPE) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
	}

	whereClause := "1=1"
	if len(filters) > 0 {
		whereClause = strings.Join(filters, " AND ")
	}
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM BUDGET_CATEGORIES WHERE IS_DELETED = 0 AND %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Errorf("[BudgetCategoryOracle][FindAllWithPagination] count failed: %v", err)
		return nil, err
	}

	page := filterParams.Page
	if page < 1 {
		page = 1
	}
	perPage := filterParams.PerPage
	if perPage < 1 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	limit := perPage

	if total == 0 {
		meta := local_util.BuildPaginationMeta(0, page, perPage)
		return &types.PaginatedResponse[[]imodel.BudgetCategoryOracle]{Data: []imodel.BudgetCategoryOracle{}, Meta: meta}, nil
	}

	if int64(offset) >= total && total > 0 {
		offset = 0
		limit = int(total)
	} else if int64(offset)+int64(limit) > total && total > 0 {
		limit = int(total) - offset
	}

	selectQuery := fmt.Sprintf(
		`SELECT RAWTOHEX(ID) AS id, NAME, ACCOUNT_TYPE, COLOR, ICON, IS_ENABLED, IS_DELETED, CREATED_AT, LAST_MODIFIED_AT, DELETED_AT
		 FROM BUDGET_CATEGORIES WHERE IS_DELETED = 0 AND %s ORDER BY CREATED_AT DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`,
		whereClause, idx, idx+1,
	)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		log.Errorf("[BudgetCategoryOracle][FindAllWithPagination] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []imodel.BudgetCategoryOracle
	for rows.Next() {
		var bc imodel.BudgetCategoryOracle
		var createdAt, lastModifiedAt, deletedAt sql.NullTime
		if err := rows.Scan(
			&bc.ID,
			&bc.Name,
			&bc.AccountType,
			&bc.Color,
			&bc.Icon,
			&bc.IsEnabled,
			&bc.IsDeleted,
			&createdAt,
			&lastModifiedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		bc.CreatedAt = formatNullTime(createdAt)
		bc.LastModifiedAt = formatNullTime(lastModifiedAt)
		bc.DeletedAt = formatNullTime(deletedAt)
		list = append(list, bc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	return &types.PaginatedResponse[[]imodel.BudgetCategoryOracle]{Data: list, Meta: meta}, nil
}

func (r *Repository) CheckBudgetCatagoryINUse(ctx context.Context, catagory_id string) (bool, error) {
	query := `SELECT COUNT(*) FROM BUDGET_CATAGORY_ALLOCATIONS WHERE BUDGET_ID = HEXTORAW(:1)`

	row := r.db.QueryRowContext(ctx, query, catagory_id)

	var count int
	if err := row.Scan(&count); err != nil {
		return false, err
	}

	return count > 0, nil

}
