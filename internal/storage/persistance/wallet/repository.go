package wallet

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type WalletStorage struct {
	dal    dal.MongoDal[model.Wallet, model.Wallet]
	logger utils.Logger
}

func NewWalletRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.WalletRepository {
	return &WalletStorage{
		dal:    dal.NewMongoDal[model.Wallet, model.Wallet](client, dbName, collection),
		logger: logger,
	}
}

func (w *WalletStorage) Create(ctx context.Context, wallet *model.Wallet) error {
	walletDoc, err := ToWalletDocument(*wallet)
	if err != nil {
		w.logger.Errorf("Failed to convert wallet to document: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_, err = w.dal.InsertOne(ctx, *walletDoc)
	if err != nil {
		w.logger.Errorf("Failed to insert wallet: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *WalletStorage) Update(ctx context.Context, id string, wallet *model.Wallet) error {
	var update bson.M
	objID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update = UpdateMapper(*wallet)

	if len(update) == 1 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorWalletNotFound.Code)
		}
		w.logger.Errorf("Failed to update wallet: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *WalletStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorWalletNotFound.Code)
		}
		w.logger.Errorf("Failed to delete wallet: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *WalletStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Warnf("Wallet ID %s not found for enable/disable", id)
			return errors.New(localization.ErrorWalletNotFound.Code)
		}
		w.logger.Errorf("Failed to enable/disable wallet ID %s: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *WalletStorage) FindByID(ctx context.Context, id string) (*model.Wallet, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	doc, err := w.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorWalletNotFound.Code)
		}
		w.logger.Errorf("FindByID wallet failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return doc, nil
}

func (w *WalletStorage) Find(ctx context.Context, code, name string) (*model.Wallet, error) {
	filter := bson.M{
		"is_deleted": false,
	}

	var orFilters []bson.M

	if code != "" {
		orFilters = append(orFilters, bson.M{
			"code": code,
		})
	}

	if name != "" {
		orFilters = append(orFilters, bson.M{
			"name": bson.M{
				"$regex":   name,
				"$options": "i",
			},
		})
	}

	if len(orFilters) > 0 {
		filter["$or"] = orFilters
	}

	doc, err := w.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Warnf("No wallet found with %s: %s", code, name)
			return nil, nil
		}
		w.logger.Errorf("FindBy code:%s name:%s wallet failed: %v", code, name, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return doc, nil
}

func (e *WalletStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error) {
	allowedKeys := []string{"name", "code", "enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowedKeys)

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
			{"code": searchRegex},
		}
	}
	docs, err := e.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		e.logger.Errorf("FindAllWithPagination Wallet failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := e.dal.TotalCount(ctx, filter)
	if err != nil {
		e.logger.Errorf("Count Wallet failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	e.logger.Infof("FindAllWithPagination returning %d wallets, total: %d", len(docs), total)

	return &types.PaginatedResponse[[]*model.Wallet]{
		Data: docs,
		Meta: meta,
	}, nil
}
