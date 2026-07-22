package transaction_limit

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
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
