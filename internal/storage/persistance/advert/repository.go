package advert

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"regexp"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AdvertStorage struct {
	dal        dal.MongoDal[model.Advert, model.Advert]
	client     *mongo.Client
	logger     utils.Logger
	collection *mongo.Collection
}

func NewAdvertRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.AdvertRepository {
	return &AdvertStorage{
		dal:        dal.NewMongoDal[model.Advert, model.Advert](client, cfg, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

func (a *AdvertStorage) Create(ctx context.Context, advert *model.Advert) error {
	a.logger.Infof("[Create] creating advert")
	_, err := a.dal.InsertOne(ctx, *advert)
	if err != nil {
		a.logger.Errorf("[Create] failed to create advert: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[Create] advert created successfully")
	return nil
}

func (a *AdvertStorage) Update(ctx context.Context, id string, advert *model.Advert) error {
	a.logger.Infof("[Update] updating advert for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AdvertMapper(*advert)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		a.logger.Errorf("[Update] failed to update advert: %v", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[Update] advert updated successfully")
	return nil
}

func (a *AdvertStorage) Delete(ctx context.Context, id string) error {
	a.logger.Infof("[Delete] deleting advert for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	err = a.dal.DeleteOne(ctx, filter)
	if err != nil {
		a.logger.Errorf("[Delete] failed to delete advert: %v", err)
		return err
	}
	a.logger.Infof("[Delete] advert deleted successfully")
	return nil
}

func (a *AdvertStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	a.logger.Infof("[EnableOrDisable] processing advert enable/disable for id: %s, enabled: %v", id, enable)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = a.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[EnableOrDisable] advert not found")
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[EnableOrDisable] failed to enable/disable advert: %v", err)
		return err
	}
	a.logger.Infof("[EnableOrDisable] advert enable/disable completed successfully")
	return nil
}

func (a *AdvertStorage) FindByID(ctx context.Context, id string) (*model.Advert, error) {
	a.logger.Infof("[FindByID] fetching advert by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			a.logger.Errorf("[FindByID] advert not found")
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		a.logger.Errorf("[FindByID] failed to find advert: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	a.logger.Infof("[FindByID] advert retrieved successfully")
	return result, nil
}

func (a *AdvertStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Advert], error) {

	allowedKeys := []string{"title", "description", "enabled","advert_for"}

	// Build search keys for $or search on title & description
	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"title": searchRegex},
			{"description": searchRegex},
			// {"is_deleted": false},
		}
	}

	// Use FilterBuilder to construct filter + pagination
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// Fetch data with final filter
	data, err := a.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		a.logger.Errorf("[FindAllWithPagination] failed to fetch adverts: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.dal.TotalCount(ctx, filter)
	if err != nil {
		a.logger.Errorf("[FindAllWithPagination] failed to count adverts: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	a.logger.Infof("[FindAllWithPagination] retrieved %d adverts", len(data))

	return &types.PaginatedResponse[[]model.Advert]{
		Data: data,
		Meta: meta,
	}, nil
}

func (a *AdvertStorage) FindByTitle(ctx context.Context, title string) (*model.Advert, error) {
	a.logger.Infof("[FindByTitle] searching for advert by title")

	filter := bson.M{
		"title": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(title) + "$", // exact match, case-insensitive
			"$options": "i",
		},
		"is_deleted": false,
	}

	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			a.logger.Infof("[FindByTitle] advert not found")
			return nil, nil
		}
		a.logger.Errorf("[FindByTitle] failed to find advert: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}
	a.logger.Infof("[FindByTitle] advert retrieved successfully")
	return result, nil
}
