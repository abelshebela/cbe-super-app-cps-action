package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type WalletStorage struct {
	dal        dal.MongoDal[local_model.Wallet, local_model.Wallet]
	serviceDal dal.MongoDal[model.Service, model.Service]
	collection *mongo.Collection
	logger     utils.Logger
}

func NewWalletRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection, ServicesCollection string, logger utils.Logger) storage.WalletRepository {
	collectionRef := client.Database(dbName).Collection(collection)

	return &WalletStorage{
		dal:        dal.NewMongoDal[local_model.Wallet, local_model.Wallet](client, cfg, dbName, collection),
		serviceDal: dal.NewMongoDal[model.Service, model.Service](client, cfg, dbName, ServicesCollection),
		collection: collectionRef,
		logger:     logger,
	}
}

func (w *WalletStorage) Create(ctx context.Context, wallet *local_model.Wallet) error {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	walletDoc, err := ToWalletDocument(*wallet)
	if err != nil {
		log.Errorf("[WalletStorage][Create] failed to convert wallet to document: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_, err = w.dal.InsertOne(ctx, *walletDoc)
	if err != nil {
		log.Errorf("[WalletStorage][Create] failed to insert wallet: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *WalletStorage) Update(ctx context.Context, id string, wallet *local_model.Wallet) error {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	var update bson.M
	objID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		log.Errorf("[WalletStorage][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update = UpdateMapper(*wallet)

	if len(update) == 1 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[WalletStorage][Update] failed to update wallet: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (w *WalletStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[WalletStorage][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[WalletStorage][Delete] failed to delete wallet: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (w *WalletStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[WalletStorage][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[WalletStorage][EnableOrDisable] failed to enable/disable wallet: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (w *WalletStorage) FindByID(ctx context.Context, id string) (*local_model.Wallet, error) {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[WalletStorage][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	doc, err := w.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[WalletStorage][FindByID] failed to find wallet: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return doc, nil
}

func (w *WalletStorage) Find(ctx context.Context, code, name string) (*local_model.Wallet, error) {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	filter := bson.M{
		"is_deleted": false,
	}

	var orFilters []bson.M

	if code != "" {
		orFilters = append(orFilters, bson.M{
			"unique_code": bson.M{
				"$regex":   code,
				"$options": "i",
			},
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
			log.Warnf("[WalletStorage][Find] no wallet found with code=%s name=%s", code, name)
			return nil, nil
		}
		log.Errorf("[WalletStorage][Find] failed to find wallet code=%s name=%s: %v", code, name, err)
		return nil, local_util.HandleDBError(err)
	}

	return doc, nil
}

func (e *WalletStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.Wallet], error) {
	log := local_util.LoggerFromCtx(ctx, e.logger)

	allowedKeys := []string{"name", "code", "enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowedKeys)
	if filterParam.Search != "" && filterParam.Search != "enabled" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
			{"unique_code": searchRegex},
		}
	}
	// Determine sort order by wallet name if requested
	sort := bson.D{}
	if val, ok := filterParam.Filters["sort_name"]; ok {
		if s, ok := val.(string); ok {
			switch s := s; s {
			case "ASC", "asc":
				sort = append(sort, bson.E{Key: "name", Value: 1})
			case "DESC", "desc":
				sort = append(sort, bson.E{Key: "name", Value: -1})
			}
		}
	}

	param := dal.FilterParam{
		Filter: filter,
		Sort:   sort,
		Skip:   skip,
		Limit:  limit,
	}

	docs, err := e.dal.FindAllWithPaginationD(ctx, param)
	if err != nil {
		log.Errorf("[WalletStorage][FindAllWithPagination] failed to fetch wallets: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := e.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[WalletStorage][FindAllWithPagination] failed to count wallets: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	log.Infof("[WalletStorage][FindAllWithPagination] returning %d wallets, total: %d", len(docs), total)

	return &types.PaginatedResponse[[]local_model.Wallet]{
		Data: docs,
		Meta: meta,
	}, nil
}

func (w *WalletStorage) FindAllWithPaginationForGRPC(
	ctx context.Context,
	filterParam types.Filter,
) (*types.PaginatedResponse[[]local_model.Wallet], error) {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	allowedKeys := []string{"name", "code", "enabled"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowedKeys)
	if filterParam.Search != "" && filterParam.Search != "enabled" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": searchRegex},
			{"unique_code": searchRegex},
		}
	}
	if filterParam.Search == "enabled" {
		filter["enabled"] = true
	}
	sort := bson.D{}
	if val, ok := filterParam.Filters["sort_name"]; ok {
		if s, ok := val.(string); ok {
			switch s := s; s {
			case "ASC", "asc":
				sort = append(sort, bson.E{Key: "name", Value: 1})
			case "DESC", "desc":
				sort = append(sort, bson.E{Key: "name", Value: -1})
			}
		}
	}
	if len(sort) == 0 {
		sort = append(sort, bson.E{Key: "created_at", Value: -1})
	}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},

		// Convert service_id string to ObjectId for lookup
		{{Key: "$addFields", Value: bson.M{
			"service_obj_id": bson.M{
				"$cond": bson.A{
					bson.M{
						"$and": bson.A{
							bson.M{"$ne": bson.A{"$service_id", ""}},
							bson.M{"$eq": bson.A{
								bson.M{"$strLenCP": "$service_id"},
								24,
							}},
						},
					},
					bson.M{"$toObjectId": "$service_id"},
					nil,
				},
			},
		}}},

		// Lookup services using converted ObjectId
		{{Key: "$lookup", Value: bson.M{
			"from":         "services",
			"localField":   "service_obj_id",
			"foreignField": "_id",
			"as":           "temp_service",
		}}},

		// Unwind safely
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$temp_service",
			"preserveNullAndEmptyArrays": true,
		}}},

		{{

			Key: "$addFields", Value: bson.M{
				"service_code": "$temp_service.service_code",
				"service_key":  "$temp_service.service_key",
				"cap":          "$temp_service.cap",
				"service_id": bson.M{
					"$cond": bson.A{
						bson.M{"$ifNull": bson.A{"$temp_service._id", false}},
						bson.M{"$toString": "$temp_service._id"},
						"$service_id_safe",
					},
				},
				// all below related to time is to handle potential date fields that are stored as strings, converting them to dates if needed
				"created_at": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{bson.M{"$type": "$created_at"}, "string"}},
						bson.M{
							"$dateFromString": bson.M{
								"dateString": "$created_at",
								"format":     "%Y-%m-%dT%H:%M:%S.%L",
								"onError":    nil,
								"onNull":     nil,
							},
						},
						"$created_at",
					},
				},
				"last_modified_at": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{bson.M{"$type": "$last_modified_at"}, "string"}},
						bson.M{
							"$dateFromString": bson.M{
								"dateString": "$last_modified_at",
								"format":     "%Y-%m-%dT%H:%M:%S.%L",
								"onError":    nil,
								"onNull":     nil,
							},
						},
						"$last_modified_at",
					},
				},
				"deleted_at": bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{bson.M{"$type": "$deleted_at"}, "string"}},
						bson.M{
							"$dateFromString": bson.M{
								"dateString": "$deleted_at",
								"format":     "%Y-%m-%dT%H:%M:%S.%L",
								"onError":    nil,
								"onNull":     nil,
							},
						},
						"$last_modified_at",
					},
				},
			},
		}},
		{{Key: "$sort", Value: sort}},
		// Cleanup
		{{Key: "$project", Value: bson.M{
			"temp_service":   0,
			"service_obj_id": 0,
		}}},

		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: limit}},
	}

	opts := options.Aggregate().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2,
	})

	cursor, err := w.collection.Aggregate(ctx, pipeline, opts)
	if err != nil {
		log.Errorf("[WalletStorage][FindAllWithPaginationForGRPC] aggregation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var wallets []local_model.Wallet
	if err := cursor.All(ctx, &wallets); err != nil {
		log.Errorf("[WalletStorage][FindAllWithPaginationForGRPC] decode failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := w.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[WalletStorage][FindAllWithPaginationForGRPC] failed to count wallets: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]local_model.Wallet]{
		Data: wallets,
		Meta: meta,
	}, nil
}

func (w *WalletStorage) FindByIDForGRPC(ctx context.Context, id string) (*local_model.GRPCWallet, error) {
	log := local_util.LoggerFromCtx(ctx, w.logger)

	doc, err := w.FindByID(ctx, id)
	if err != nil {
		return nil, err // Wallet itself not found, this should still error
	}

	grpcWallet := ToGRPCWallet(*doc)

	// Attempt to parse ServiceID
	objID, err := bson.ObjectIDFromHex(grpcWallet.ServiceID)
	if err != nil {
		// Log it, but don't fail the request. Return the wallet as is.
		log.Debugf("[WalletStorage][FindByIDForGRPC] wallet %s has invalid ServiceID: %v", id, grpcWallet.ServiceID)
		return &grpcWallet, nil
	}

	// Attempt to fetch Service
	service, err := w.serviceDal.FindOne(ctx, bson.M{"_id": objID}, nil)
	if err != nil {
		log.Warnf("[WalletStorage][FindByIDForGRPC] service not found for wallet %s: %v", id, err)
		return &grpcWallet, nil // Return wallet even if service lookup fails
	}

	// Enrich if found
	if service != nil {
		grpcWallet.ServiceCode = service.ServiceCode
		// grpcWallet.ServiceKey = service.ServiceName
	}

	return &grpcWallet, nil
}
