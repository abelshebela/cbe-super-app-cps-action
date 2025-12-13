package Topup

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TopupStorage struct {
	dal    dal.MongoDal[model.Topup, model.Topup]
	logger utils.Logger
}

func NewTopupRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.TopupRepository {
	return &TopupStorage{
		dal:    dal.NewMongoDal[model.Topup, model.Topup](client, dbName, collection),
		logger: logger,
	}
}

func (w *TopupStorage) Create(ctx context.Context, Topup *model.Topup) error {
	TopupDoc, err := ToTopupDocument(*Topup)
	if err != nil {
		w.logger.Errorf("Failed to convert Topup to document: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	_, err = w.dal.InsertOne(ctx, *TopupDoc)
	if err != nil {
		w.logger.Errorf("Failed to insert Topup: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *TopupStorage) Update(ctx context.Context, id string, Topup *model.Topup) error {
	var update bson.M
	objID, err := bson.ObjectIDFromHex(id)

	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update = UpdateMapper(*Topup)

	if len(update) == 1 {
		return errors.New(localization.ErrorNoDataProvided.Code)
	}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorTopupNotFound.Code)
		}
		w.logger.Errorf("Failed to update Topup: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *TopupStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"is_deleted": true, "deleted_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.New(localization.ErrorTopupNotFound.Code)
		}
		w.logger.Errorf("Failed to delete Topup: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *TopupStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = w.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Warnf("Topup ID %s not found for enable/disable", id)
			return errors.New(localization.ErrorTopupNotFound.Code)
		}
		w.logger.Errorf("Failed to enable/disable Topup ID %s: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (w *TopupStorage) FindByID(ctx context.Context, id string) (*model.Topup, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	doc, err := w.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New(localization.ErrorTopupNotFound.Code)
		}
		w.logger.Errorf("FindByID Topup failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return doc, nil
}

func (w *TopupStorage) Find(ctx context.Context, code, name string) (*model.Topup, error) {
	filter := bson.M{
		"is_deleted": false,
	}

	var orFilters []bson.M

	if code != "" {
		orFilters = append(orFilters, bson.M{"code": code})
	}

	if name != "" {
		orFilters = append(orFilters, bson.M{"name": bson.M{"$regex": name, "$options": "i"}})
	}

	if len(orFilters) > 0 {
		filter["$or"] = orFilters
	}

	doc, err := w.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			w.logger.Warnf("No Topup found with %s: %s", code, name)
			return nil, nil
		}
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return doc, nil
}

func (e *TopupStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Topup], error) {
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
		e.logger.Errorf("FindAllWithPagination Topup failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := e.dal.TotalCount(ctx, filter)
	if err != nil {
		e.logger.Errorf("Count Topup failed", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	e.logger.Infof("FindAllWithPagination returning %d Topups, total: %d", len(docs), total)

	return &types.PaginatedResponse[[]*model.Topup]{
		Data: docs,
		Meta: meta,
	}, nil
}
