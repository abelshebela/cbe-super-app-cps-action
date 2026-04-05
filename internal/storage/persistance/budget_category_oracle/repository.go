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
	// Oracle column ENABLED (legacy); not IS_ENABLED. See db/migrations/000010_budget_categories_normalize_enabled_column.up.sql.
	// No IS_DELETED; removals use DELETE FROM (see Delete).
	q := `INSERT INTO BUDGET_CATEGORIES (ID, NAME, COLOR, ICON, TYPE, ENABLED, CREATE_AT, UPDATE_AT)
		VALUES (SYS_GUID(), :1, :2, :3, :4, :5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	_, err := r.db.ExecContext(ctx, q,
		bc.Name,
		bc.Color,
		bc.Icon,
		bc.Type,
		bc.IsEnabled,
	)
	if err != nil {
		r.logger.Errorf("[BudgetCategoryOracle][Create] insert failed: %v", err)
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
	q := `UPDATE BUDGET_CATEGORIES SET NAME = :1, COLOR = :2, ICON = :3, TYPE = :4, ENABLED = :5, UPDATE_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:6)`
	_, err := r.db.ExecContext(ctx, q,
		bc.Name,
		bc.Color,
		bc.Icon,
		bc.Type,
		bc.IsEnabled,
		id,
	)
	if err != nil {
		r.logger.Errorf("[BudgetCategoryOracle][Update] failed: %v", err)
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
	q := `DELETE FROM BUDGET_CATEGORIES WHERE ID = HEXTORAW(:1)`
	_, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		r.logger.Errorf("[BudgetCategoryOracle][Delete] failed: %v", err)
		return err
	}
	return nil
}

func (r *Repository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	q := `UPDATE BUDGET_CATEGORIES SET ENABLED = :1, UPDATE_AT = CURRENT_TIMESTAMP WHERE ID = HEXTORAW(:2)`
	var v int
	if enable {
		v = 1
	} else {
		v = 0
	}
	_, err := r.db.ExecContext(ctx, q, v, id)
	if err != nil {
		r.logger.Errorf("[BudgetCategoryOracle][EnableOrDisable] failed: %v", err)
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
	q := `SELECT RAWTOHEX(ID) AS id, NAME, COLOR, ICON, TYPE, ENABLED, CREATE_AT, UPDATE_AT
		FROM BUDGET_CATEGORIES WHERE ID = HEXTORAW(:1)`
	row := r.db.QueryRowContext(ctx, q, id)
	var bc imodel.BudgetCategoryOracle
	var createAt, updateAt sql.NullTime
	err := row.Scan(
		&bc.ID,
		&bc.Name,
		&bc.Color,
		&bc.Icon,
		&bc.Type,
		&bc.IsEnabled,
		&createAt,
		&updateAt,
	)
	if err == sql.ErrNoRows {
		return nil, localization.ErrorResourceNotFound
	}
	if err != nil {
		return nil, err
	}
	bc.IsDeleted = 0
	bc.CreateAt = formatNullTime(createAt)
	bc.UpdateAt = formatNullTime(updateAt)
	return &bc, nil
}

func (r *Repository) FindByName(ctx context.Context, name string) (*imodel.BudgetCategoryOracle, error) {
	q := `SELECT RAWTOHEX(ID) AS id, NAME, COLOR, ICON, TYPE, ENABLED, CREATE_AT, UPDATE_AT
		FROM BUDGET_CATEGORIES WHERE UPPER(NAME) = UPPER(:1)`
	row := r.db.QueryRowContext(ctx, q, name)
	var bc imodel.BudgetCategoryOracle
	var createAt, updateAt sql.NullTime
	err := row.Scan(
		&bc.ID,
		&bc.Name,
		&bc.Color,
		&bc.Icon,
		&bc.Type,
		&bc.IsEnabled,
		&createAt,
		&updateAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	bc.IsDeleted = 0
	bc.CreateAt = formatNullTime(createAt)
	bc.UpdateAt = formatNullTime(updateAt)
	return &bc, nil
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.UTC().Format(time.RFC3339)
}

func (r *Repository) FindAllWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.BudgetCategoryOracle], error) {
	if filterParams == nil {
		fp := types.Filter{Page: 1, PerPage: 10}
		filterParams = &fp
	}
	var filters []string
	var args []interface{}
	idx := 1

	if filterParams.Search != "" {
		search := "%" + strings.ToUpper(filterParams.Search) + "%"
		filters = append(filters, fmt.Sprintf("(UPPER(NAME) LIKE :%d OR UPPER(TYPE) LIKE :%d)", idx, idx))
		args = append(args, search)
		idx++
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
				filters = append(filters, "ENABLED = 1")
			} else {
				filters = append(filters, "ENABLED = 0")
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
				filters = append(filters, fmt.Sprintf("UPPER(TYPE) = UPPER(:%d)", idx))
				args = append(args, s)
				idx++
			}
		}
	}

	whereClause := "1=1"
	if len(filters) > 0 {
		whereClause = strings.Join(filters, " AND ")
	}
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM BUDGET_CATEGORIES WHERE %s", whereClause)
	var total int64
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		r.logger.Errorf("[BudgetCategoryOracle][FindAllWithPagination] count failed: %v", err)
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
		`SELECT RAWTOHEX(ID) AS id, NAME, COLOR, ICON, TYPE, ENABLED, CREATE_AT, UPDATE_AT
		 FROM BUDGET_CATEGORIES WHERE %s ORDER BY CREATE_AT DESC OFFSET :%d ROWS FETCH NEXT :%d ROWS ONLY`,
		whereClause, idx, idx+1,
	)
	args = append(args, offset, limit)

	rows, err := r.db.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		r.logger.Errorf("[BudgetCategoryOracle][FindAllWithPagination] query failed: %v", err)
		return nil, err
	}
	defer rows.Close()

	var list []imodel.BudgetCategoryOracle
	for rows.Next() {
		var bc imodel.BudgetCategoryOracle
		var createAt, updateAt sql.NullTime
		if err := rows.Scan(
			&bc.ID,
			&bc.Name,
			&bc.Color,
			&bc.Icon,
			&bc.Type,
			&bc.IsEnabled,
			&createAt,
			&updateAt,
		); err != nil {
			return nil, err
		}
		bc.IsDeleted = 0
		bc.CreateAt = formatNullTime(createAt)
		bc.UpdateAt = formatNullTime(updateAt)
		list = append(list, bc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, page, perPage)
	return &types.PaginatedResponse[[]imodel.BudgetCategoryOracle]{Data: list, Meta: meta}, nil
}
