package transaction_limit

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TransactionLimitStorage struct {
	dal    dal.MongoDal[imodel.TransactionLimit, imodel.TransactionLimit]
	logger utils.Logger
}

func NewTransactionLimitRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, collection string, logger utils.Logger) storage.TransactionLimitRepository {
	return &TransactionLimitStorage{
		dal:    dal.NewMongoDal[imodel.TransactionLimit, imodel.TransactionLimit](client, cfg, dbName, collection),
		logger: logger,
	}
}

func (r *TransactionLimitStorage) FindByCustomerNumber(ctx context.Context, customerNumber string) (*imodel.TransactionLimit, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	customerNumber = strings.TrimSpace(customerNumber)
	if customerNumber == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	filter := bson.M{"customer_number": customerNumber}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[TransactionLimit][FindByCustomerNumber] failed for customer_number=%s: %v", customerNumber, err)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	return result, nil
}

func (r *TransactionLimitStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.TransactionLimit], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	searchKeys := bson.M{}
	allowedKeys := []string{"customer_number"}

	if filterParam.Search != "" {
		searchKeys["$or"] = []bson.M{
			{"customer_number": bson.M{"$regex": filterParam.Search, "$options": "i"}},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	if skip < 0 {
		skip = 0
	}
	if limit <= 0 {
		limit = 10
	}

	results, err := r.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[TransactionLimit][FindAllWithPagination] failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	totalDocs, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		totalDocs = int64(skip) + int64(len(results))
	}

	page := filterParam.Page
	if page <= 0 {
		page = 1
	}

	totalPages := int((totalDocs + limit - 1) / limit)
	if totalPages == 0 {
		totalPages = 1
	}

	pagingCounter := (int64(page)-1)*limit + 1
	hasPrev := page > 1
	hasNext := page < totalPages

	var prevPage *int
	if hasPrev {
		p := page - 1
		prevPage = &p
	}
	var nextPage *int
	if hasNext {
		n := page + 1
		nextPage = &n
	}

	return &types.PaginatedResponse[[]imodel.TransactionLimit]{
		Data: results,
		Meta: types.PaginationMeta{
			TotalDocs:     totalDocs,
			Limit:         int(limit),
			TotalPages:    totalPages,
			Page:          page,
			PagingCounter: int(pagingCounter),
			HasPrevPage:   hasPrev,
			HasNextPage:   hasNext,
			PrevPage:      prevPage,
			NextPage:      nextPage,
		},
	}, nil
}
