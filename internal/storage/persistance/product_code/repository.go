package productcode

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	cps_errors "cbe-super-app-cps-action/internal/constants/errors"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductCodeStorage struct {
	producCodeDal dal.MongoDal[model.ServiceDetails, model.ServiceDetails]
	logger        utils.Logger
	collection    *mongo.Collection
}

func NewProductCodeRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ProductCodeRepository {
	return &ProductCodeStorage{
		producCodeDal: dal.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, collection),
		logger:        logger,
		collection:    client.Database(dbName).Collection(collection),
	}
}

func (p *ProductCodeStorage) FetchByID(ctx context.Context, id string) (*model.ProductCode, error) {
	p.logger.Infof("[productcode.FetchByID] Fetching product code with ID: %s", id)

	objID, ok := local_util.StringToObjectID(id)
	if !ok {
		return nil, cps_errors.ErrInvalidID
	}
	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}
	services, err := p.producCodeDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Errorf("[productcode.FetchByID] Product code not found for ID %s: %v", id, err)
			return nil, cps_errors.ErrProductCodeNotFound
		}
		p.logger.Errorf("[productcode.FetchByID] Database query failed for ID %s: %v", id, err)
		return nil, cps_errors.ErrGeneralDBQueryFailed
	}
	pc := productcode.ToProducCode(*services)

	p.logger.Infof("[productcode.FetchByID] Successfully fetched product code with ID: %s", id)
	return pc, nil
}

func (r *ProductCodeStorage) FetchAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error) {
	r.logger.Infof("[productcode.FetchAll] Fetching product codes with filter: %+v", filterParams)

	filter := local_util.BuildFilter(filterParams, r.logger)
	if filterParams.Search != "" {
		filter["service_name"] = bson.M{"$regex": filterParams.Search, "$options": "i"}
	}
	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$facet", Value: bson.M{
			"data": []bson.D{
				{{Key: "$skip", Value: int64(skip)}},
				{{Key: "$limit", Value: int64(limit)}},
			},
			"total": []bson.D{
				{{Key: "$count", Value: "count"}},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[productcode.FetchAll] Failed to execute aggregation: %v", err)
		return nil, cps_errors.ErrGeneralDBQueryFailed
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []*model.ServiceDetails `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("[productcode.FetchAll] Failed to decode aggregation results: %v", err)
		return nil, cps_errors.ErrGeneralDBQueryFailed
	}

	if len(results) == 0 {
		return &types.PaginatedResponse[[]*model.ProductCode]{
			Data: []*model.ProductCode{},
			Meta: local_util.BuildPaginationMeta(0, filterParams.Page, limit),
		}, nil
	}

	var productCodes []*model.ProductCode
	for _, doc := range results[0].Data {
		productCodes = append(productCodes, productcode.ToProducCode(*doc))
	}

	var total int64
	if len(results[0].Total) > 0 {
		total = results[0].Total[0].Count
	}

	meta := local_util.BuildPaginationMeta(total, filterParams.Page, limit)
	r.logger.Infof("[productcode.FetchAll] Successfully fetched %d product codes, total: %d", len(productCodes), total)
	return &types.PaginatedResponse[[]*model.ProductCode]{
		Data: productCodes,
		Meta: meta,
	}, nil
}

func (p *ProductCodeStorage) Update(ctx context.Context, productCode *model.ProductCode) error {
	p.logger.Infof("[productcode.Update] Updating product code with ID: %s", productCode.ID)

	objID, ok := local_util.StringToObjectID(productCode.ID)
	if !ok {
		return cps_errors.ErrInvalidID
	}
	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}
	update := ProductCodeMapper(*productCode)

	_, err := p.producCodeDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Errorf("[productcode.Update] Product code not found for ID %s: %v", productCode.ID, err)
			return cps_errors.ErrProductCodeNotFound
		}
		p.logger.Errorf("[productcode.Update] Database update failed for ID %s: %v", productCode.ID, err)
		return cps_errors.ErrGeneralDBQueryFailed
	}

	return nil
}
