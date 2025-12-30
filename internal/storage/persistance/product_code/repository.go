package productcode

import (
	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"regexp"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

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

func (r *ProductCodeStorage) FetchAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error) {
	r.logger.Infof("[FetchAll] fetching product codes")

	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	if filterParams.Search != "" {
		searchRegex := bson.M{"$regex": filterParams.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"service_name": searchRegex},
			{"cbe_product_codes.prd": searchRegex},
			{"cbe_ifb_product_codes.prd": searchRegex},
		}
	}
	allowedKeys := []string{"_id", "service_name", "created_at", "last_modified_at", "cbe_product_codes.prd", "cbe_ifb_product_codes.prd"}

	filter, skip, limit := lib.FilterBuilder(*filterParams, searchKeys, allowedKeys)
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
		r.logger.Errorf("[FetchAll] failed to execute aggregation: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []*model.ServiceDetails `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		r.logger.Errorf("[FetchAll] failed to decode aggregation results: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}

	if len(results) == 0 {
		return &types.PaginatedResponse[[]*model.ProductCode]{
			Data: []*model.ProductCode{},
			Meta: local_util.BuildPaginationMeta(0, filterParams.Page, int(limit)),
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

	meta := local_util.BuildPaginationMeta(total, filterParams.Page, int(limit))
	r.logger.Infof("[FetchAll] retrieved %d product codes", len(productCodes))
	return &types.PaginatedResponse[[]*model.ProductCode]{
		Data: productCodes,
		Meta: meta,
	}, nil
}

func (s *ProductCodeStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error) {
	s.logger.Infof("[FindAllWithPagination] fetching product codes with pagination")
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}
	allowedKeys := []string{"_id", "service_name", "created_at", "last_modified_at", "cbe_product_codes.prd", "cbe_ifb_product_codes.prd"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"service_name": searchRegex},
			{"cbe_product_codes.prd": searchRegex},
			{"cbe_ifb_product_codes.prd": searchRegex},
		}
	}
	filter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)
	data, err := s.producCodeDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to fetch product codes: %v", err)
		return nil, fmt.Errorf("%s", localization.ErrorUnexpectedError.Code)
	}
	pdata := []*model.ProductCode{}
	for i := range data {
		pdata = append(pdata, productcode.ToProducCode(data[i]))
	}
	total, err := s.producCodeDal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count product codes: %v", err)
		return nil, fmt.Errorf("%s", localization.ErrorUnexpectedError.Code)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d product codes", len(pdata))
	return &types.PaginatedResponse[[]*model.ProductCode]{
		Data: pdata,
		Meta: meta,
	}, nil
}

func (p *ProductCodeStorage) FetchByID(ctx context.Context, id string) (*model.ProductCode, error) {
	p.logger.Infof("[FetchByID] fetching product code by id: %s", id)

	objID, ok := local_util.StringToObjectID(id)
	if !ok {
		return nil, fmt.Errorf("%s", "ERROR_INVALID_ID")
	}
	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}
	services, err := p.producCodeDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Errorf("[FetchByID] product code not found")
			return nil, fmt.Errorf("%s", "ERROR_PRODCUT_CODE_NOT_FOUND")
		}
		p.logger.Errorf("[FetchByID] failed to fetch product code: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}
	pc := productcode.ToProducCode(*services)

	p.logger.Infof("[FetchByID] product code retrieved successfully")
	return pc, nil
}

func (p *ProductCodeStorage) Update(ctx context.Context, productCode *model.ProductCode) error {
	p.logger.Infof("[Update] updating product code for id: %s", productCode.ID)

	objID, ok := local_util.StringToObjectID(productCode.ID)
	if !ok {
		return fmt.Errorf("%s", "ERROR_INVALID_ID")
	}
	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}
	update := ProductCodeMapper(*productCode)

	_, err := p.producCodeDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Errorf("[Update] product code not found")
			return fmt.Errorf("%s", "ERROR_PRODCUT_CODE_NOT_FOUND")
		}
		p.logger.Errorf("[Update] failed to update product code: %v", err)
		return fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}
	p.logger.Infof("[Update] product code updated successfully")
	return nil
}
func (p *ProductCodeStorage) FindByName(ctx context.Context, name string) (*model.ProductCode, error) {
	p.logger.Infof("[FindByName] searching for product code by name")

	filter := bson.M{
		"service_name": bson.M{
			"$regex":   "^" + regexp.QuoteMeta(name) + "$", // exact match, case-insensitive
			"$options": "i",
		},
		"is_deleted": false,
	}

	result, err := p.producCodeDal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Infof("[FindByName] product code not found")
			return nil, nil
		}
		p.logger.Errorf("[FindByName] failed to find product code: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}

	pc := productcode.ToProducCode(*result)
	p.logger.Infof("[FindByName] product code retrieved successfully")
	return pc, nil
}

func (p *ProductCodeStorage) FindByPRD(ctx context.Context, cbePRD, cbeIFBPRD string) ([]*model.ProductCode, error) {
	p.logger.Infof("[FindByPRD] searching for product codes by PRD")

	// Build the $or query to check both PRD fields efficiently
	orConditions := []bson.M{}

	if cbePRD != "" {
		orConditions = append(orConditions, bson.M{"cbe_product_codes.prd": cbePRD})
	}

	if cbeIFBPRD != "" {
		orConditions = append(orConditions, bson.M{"cbe_ifb_product_codes.prd": cbeIFBPRD})
	}

	// If both are empty, return empty result
	if len(orConditions) == 0 {
		p.logger.Infof("[FindByPRD] both PRD values are empty, skipping search")
		return []*model.ProductCode{}, nil
	}

	filter := bson.M{
		"is_deleted": false,
		"$or":        orConditions,
	}

	results, err := p.producCodeDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		p.logger.Errorf("[FindByPRD] failed to find product codes: %v", err)
		return nil, fmt.Errorf("%s", "ERROR_PRODUCT_CODE_DATABASE_QUERY_FAILED")
	}

	// Convert ServiceDetails to ProductCode
	productCodes := make([]*model.ProductCode, 0, len(results))
	for _, result := range results {
		productCodes = append(productCodes, productcode.ToProducCode(result))
	}

	p.logger.Infof("[FindByPRD] retrieved %d product codes", len(productCodes))
	return productCodes, nil
}
