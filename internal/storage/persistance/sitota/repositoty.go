package sitota

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"cbe-super-app-cps-action/internal/storage/persistance/sitota/sqlc"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type SitotaRepository struct {
	db     *sql.DB
	logger shared_utils.Logger
}

func NewSitotaRepository(db *sql.DB, logger shared_utils.Logger) storage.SitotaRepository {
	return &SitotaRepository{
		db:     db,
		logger: logger,
	}
}

func (r *SitotaRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.SitotaTransaction], error) {
	q := sqlc.New(r.db)

	params := sqlc.FindSitotaTransactionsParams{}
	params.Search = filterParam.Search
	if filterParam.Page > 0 && filterParam.PerPage > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	} else if filterParam.PerPage > 0 {
		params.Page = sql.NullInt64{Int64: 1, Valid: true}
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	}

	rows, err := q.FindSitotaTransactions(ctx, params)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	var sitotas []*model.SitotaTransaction
	var total int64
	for _, t := range rows {
		var createdAt time.Time
		var updatedAt time.Time

		if t.CreatedAt.Valid {
			createdAt = t.CreatedAt.Time
		}
		if t.UpdatedAt.Valid {
			updatedAt = t.UpdatedAt.Time
		}

		sitotas = append(sitotas, &model.SitotaTransaction{
			ID:                     t.ID,
			SenderName:             t.SenderName,
			SenderAccountNumber:    t.SenderAccountNumber,
			RecipientName:          t.RecipientName,
			RecipientAccountNumber: t.RecipientAccountNumber,
			SitotaAmount:           t.SitotaAmount,
			Status:                 t.Status,
			CreatedAt:              createdAt,
			UpdatedAt:              updatedAt,
		})
		total = t.TotalCount
	}

	limit := int(params.Limit.Int64)
	if limit == 0 {
		limit = 50
	}
	page := int(params.Page.Int64)
	if page == 0 {
		page = 1
	}

	return &types.PaginatedResponse[[]*model.SitotaTransaction]{
		Data: sitotas,
		Meta: types.PaginationMeta{
			TotalDocs:  total,
			Limit:      limit,
			Page:       page,
			TotalPages: 0,
		},
	}, nil
}

func (r *SitotaRepository) Get(ctx context.Context, id string) (*model.SitotaTransaction, error) {
	q := sqlc.New(r.db)

	t, err := q.FindSitotaTransactionByID(ctx, id)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}

	var createdAt time.Time
	var updatedAt time.Time

	if t.CreatedAt.Valid {
		createdAt = t.CreatedAt.Time
	}
	if t.UpdatedAt.Valid {
		updatedAt = t.UpdatedAt.Time
	}

	return &model.SitotaTransaction{
		ID:                     t.ID,
		SenderName:             t.SenderName,
		SenderAccountNumber:    t.SenderAccountNumber,
		RecipientName:          t.RecipientName,
		RecipientAccountNumber: t.RecipientAccountNumber,
		SitotaAmount:           t.SitotaAmount,
		Status:                 t.Status,
		CreatedAt:              createdAt,
		UpdatedAt:              updatedAt,
	}, nil
}
