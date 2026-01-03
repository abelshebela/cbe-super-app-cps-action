package hq

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type HQStorage struct {
	dal        dal.MongoDal[model.HQ, model.HQ]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewHQRepository(client *mongo.Client,cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.HQRepository {
	return &HQStorage{
		dal:        dal.NewMongoDal[model.HQ, model.HQ](client,cfg, dbName, collection),
		client:     client,
		logger:     logger,
		collection: client.Database(dbName).Collection(collection),
	}
}

func (h *HQStorage) FindByID(ctx context.Context, id string) (*model.HQ, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID}

	result, err := h.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		return nil, err
	}
	return result, nil
}
func (p *HQStorage) Find(ctx context.Context, projections ...bson.M) (*model.HQ, error) {
	projection := bson.M{}
	if len(projections) > 0 {
		projection = projections[0]
	}
	result, err := p.dal.FindOne(ctx, bson.M{}, projection)
	if err != nil {
		p.logger.Errorf("failed to fetch HQ: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if result == nil {
		p.logger.Errorf("HQ not found (single)")
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}
	return result, nil
}

func (p *HQStorage) Update(ctx context.Context, field string, value interface{}, now time.Time) error {

	hqDoc, err := p.Find(ctx)
	if err != nil {
		return err
	}
	updateDoc := bson.M{field: value}
	switch field {
	case "block_time":
		updateDoc["updated_at_block"] = now
	case "archive_time":
		updateDoc["updated_at_archive"] = now
	case "password_expiry":
		updateDoc["updated_at_password_expiry"] = now
	case "total_cap":
		updateDoc["updated_at_total_cap"] = now
	}
	_, err = p.dal.UpdateOne(ctx, bson.M{"_id": hqDoc.ID}, updateDoc)
	if err != nil {
		p.logger.Errorf("failed to update HQ field %s: %v", field, err)
		return err
	}
	return nil
}

func (h *HQStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.HQ], error) {

	allowedKeys := []string{}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowedKeys)

	data, err := h.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := h.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.HQ]{
		Data: data,
		Meta: meta,
	}, nil
}
