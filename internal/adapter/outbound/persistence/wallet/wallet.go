package wallet

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/wallet"
	wallet_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/wallet"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type WalletPersistence struct {
	walletDal dal.MongoDal[model.WalletDocument, model.WalletDocument]
	logger    utils.Logger
}

func InitWalletPersistence(client *mongo.Client, dbName string, collection string, logger utils.Logger) wallet_outbound.WalletRepository {
	return &WalletPersistence{
		walletDal: dal.NewMongoDal[model.WalletDocument, model.WalletDocument](client, dbName, collection),
		logger:    logger,
	}
}

func (w *WalletPersistence) CreateWallet(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error) {
	walletDoc, err := mappers.ToWalletDocument(&wallet)
	if err != nil {
		return nil, err
	}

	res, err := w.walletDal.InsertOne(ctx, *walletDoc)
	if err != nil {
		return nil, fmt.Errorf(common_util.GeneralDBInsertFailed)
	}
	result := mappers.ToDomainWallet(res)
	return result, nil
}

func (w *WalletPersistence) FetchWalletByID(ctx context.Context, id string) (*entity.Wallet, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	walletDoc, err := w.walletDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	result := mappers.ToDomainWallet(*walletDoc)
	return result, nil
}

func (w *WalletPersistence) FetchWallet(ctx context.Context, filterParam *constant.MongoFilter) (*common_util.PaginatedResponse[[]*entity.Wallet], error) {
	filter := bson.M{"is_deleted": false}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"code": searchRegex},
			{"name": searchRegex},
			{"avatar": searchRegex},
		}
	}

	if filterParam.Filters != nil {
		allowedKeys := []string{"enabled"}
		handlers := map[string]func(interface{}) interface{}{
			"is_deleted": func(value interface{}) interface{} {
				if str, ok := value.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						return parsed
					}
				}
				return value
			},
			"enabled": func(value interface{}) interface{} {
				if str, ok := value.(string); ok {
					if parsed, err := strconv.ParseBool(str); err == nil {
						return parsed
					}
				}
				return value
			},
		}

		enhancedFilter := common_util.BuildMongoFilterWithHandlers(filterParam.Filters, allowedKeys, handlers)
		for key, value := range enhancedFilter {
			filter[key] = value
		}
	}

	skip := (filterParam.Page - 1) * filterParam.PerPage
	limit := filterParam.PerPage

	walletDocs, err := w.walletDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return &common_util.PaginatedResponse[[]*entity.Wallet]{
				Data: []*entity.Wallet{},
				Meta: common_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage),
			}, nil
		}
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	var wallets []*entity.Wallet
	for _, doc := range walletDocs {
		converted := mappers.ToDomainWallet(*doc)
		wallets = append(wallets, converted)
	}

	total, err := w.walletDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := common_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &common_util.PaginatedResponse[[]*entity.Wallet]{
		Data: wallets,
		Meta: meta,
	}, nil
}

func (w *WalletPersistence) UpdateWallet(ctx context.Context, wallet entity.Wallet) (*entity.Wallet, error) {
	objID, err := common_util.ParsePrimitiveObjectID(wallet.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{}
	if wallet.Name != "" {
		update["name"] = wallet.Name
	}
	if wallet.Code != "" {
		update["code"] = wallet.Code
	}
	if wallet.Avatar != "" {
		update["avatar"] = wallet.Avatar
	}
	update["last_modified_at"] = time.Now()

	if len(update) == 1 {
		return nil, fmt.Errorf(common_util.NoDataProvidedForUpdate)
	}

	walletDoc, err := w.walletDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		w.logger.Warnf(err.Error(), "while updating wallet")
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToDomainWallet(walletDoc)
	return result, nil
}

func (w *WalletPersistence) DeleteWallet(ctx context.Context, id string) (*entity.Wallet, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	walletDoc, err := w.walletDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToDomainWallet(walletDoc)
	return result, nil
}

func (w *WalletPersistence) EnableDisableWallet(ctx context.Context, id string, enable bool) (*entity.Wallet, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"enabled":          enable,
		"last_modified_at": time.Now(),
	}

	walletDoc, err := w.walletDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Errorf("Wallet with ID %s not found for enable/disable", id)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		w.logger.Errorf("Failed to update enabled state for Wallet ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToDomainWallet(walletDoc)
	return result, nil
}

func (w *WalletPersistence) WalletNameExists(ctx context.Context, name string, code string, id *string) (bool, error) {
	if name == "" {
		return false, fmt.Errorf(common_util.InvalidInput)
	}

	filter := bson.M{
		"$and": []bson.M{
			{"is_deleted": false},
			{
				"$or": []bson.M{
					{"name": bson.M{"$regex": fmt.Sprintf("^%s$", name), "$options": "i"}},
					{"code": bson.M{"$regex": fmt.Sprintf("^%s$", code), "$options": "i"}},
				},
			},
		},
	}

	if id != nil {
		objID, err := common_util.ParsePrimitiveObjectID(*id)
		if err != nil {
			return false, fmt.Errorf(common_util.InvalidID)
		}
		filter["_id"] = bson.M{"$ne": objID}
	}

	count, err := w.walletDal.TotalCount(ctx, filter)
	if err != nil {
		w.logger.Warnf("Failed to check wallet name existence for name %s: %v", name, err)
		return false, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	return count > 0, nil
}
