package vault

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage/persistance/vault/gen/sqlc"
	"context"
	"database/sql"
)

func (r *VaultCategoryRepository) FindAllTransactionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.VaultTransaction], error) {
	// q := sqlc.New(r.db)

	params := sqlc.FindVaultCategoryParams{}

	if filterParam.PerPage > 0 {
		params.Limit = sql.NullInt64{Int64: int64(filterParam.PerPage), Valid: true}
	} else {
		params.Limit = sql.NullInt64{Int64: 50, Valid: true}
	}

	if filterParam.Page > 0 {
		params.Page = sql.NullInt64{Int64: int64(filterParam.Page), Valid: true}
	} else {
		params.Page = sql.NullInt64{Int64: 1, Valid: true}
	}

	// rows, err := q.FindVaultTransactions(ctx, params)
	// if err != nil {
	// 	r.logger.Errorf("failed to find vault transactions: %v", err)
	// 	if err == sql.ErrNoRows {
	// 		return nil, errors.New(localization.ErrorVaultCategoryNotFound.Code)
	// 	}
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }

	// categories := make([]*imodel.VaultTransaction, 0, len(rows))
	// var total int64

	// for _, row := range rows {
	// 	c := &imodel.VaultTransaction{}

	// 	// total = row.TotalCount
	// 	categories = append(categories, c)
	// }

	// resp := types.PaginatedResponse[[]*imodel.VaultCategory]{
	// 	Data: categories,
	// 	Meta: types.PaginationMeta{
	// 		TotalDocs:     total,
	// 		Limit:         limit,
	// 		TotalPages:    totalPages,
	// 		Page:          page,
	// 		PagingCounter: pagingCounter,
	// 		HasPrevPage:   hasPrevPage,
	// 		HasNextPage:   hasNextPage,
	// 		PrevPage:      prevPage,
	// 		NextPage:      nextPage,
	// 	},
	// }

	// return &resp, nil

	return nil, nil
}
