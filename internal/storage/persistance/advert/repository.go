package advert

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"
	"regexp"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AdvertStorage struct {
	dal    dal.MongoDal[model.Advert, model.Advert]
	client *mongo.Client
	logger utils.Logger
	collection   *mongo.Collection
}

func NewAdvertRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.AdvertRepository {
	return &AdvertStorage{
		dal:    dal.NewMongoDal[model.Advert, model.Advert](client, dbName, collection),
		client: client,
		logger: logger,
		collection:    client.Database(dbName).Collection(collection),
	}
}

func (a *AdvertStorage) Create(ctx context.Context, advert *model.Advert) error {
	_, err := a.dal.InsertOne(ctx, *advert)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *AdvertStorage) Update(ctx context.Context, id string, advert *model.Advert) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AdvertMapper(*advert)

	_, err = a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		a.logger.Errorf("UpdateOne failed: %v for filter:%v", err, filter)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *AdvertStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return a.dal.DeleteOne(ctx, filter)
}

func (a *AdvertStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable}
	_, err = a.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
	}
	return nil
}

func (a *AdvertStorage) FindByID(ctx context.Context, id string) (*model.Advert, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	result, err := a.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (a *AdvertStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Advert], error) {

	allowedKeys := []string{"title", "description", "enabled"}

	// Build search keys for $or search on title & description
	searchKeys := bson.M{}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"title": searchRegex},
			{"description": searchRegex},
			{"is_deleted": false},
		}
	}

	// Use FilterBuilder to construct filter + pagination
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// Fetch data with final filter
	data, err := a.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := a.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.Advert]{
		Data: data,
		Meta: meta,
	}, nil
}
func (a *AdvertStorage) FindByTitle(ctx context.Context, title string) (*model.Advert, error) {
a.logger.Infof("[Adver.FindByTitle] Searching for Advert by title: %s", title)

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
			a.logger.Infof(" No for Advert by title: %s", title)
			return nil, nil
		}
		a.logger.Errorf(" Database query failed for advert title  %s: %v", title, err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}

	return result,nil
}
