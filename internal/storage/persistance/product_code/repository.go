package productcode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cbe-super-app-cps-action/internal/constants/dto/productcode"
	cps_errors "cbe-super-app-cps-action/internal/constants/errors"

	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductCodeStorage struct {
	producCodeDal dal.MongoDal[model.ServiceDetails, model.ServiceDetails]
	logger        utils.Logger
}

// NewProductCodeRepository creates a new repository instance
func NewProductCodeRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ProductCodeRepository {
	return &ProductCodeStorage{
		producCodeDal: dal.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, collection),
		logger:        logger,
	}
}

func (p *ProductCodeStorage) FetchByID(ctx context.Context, id string) (*model.ProductCode, error) {
	p.logger.Infof("[productcode.FetchByID] Fetching product code with ID: %s", id)

	objID, ok := local_util.StringToObjectID(id)
	if !ok {
		err := errors.New("invalid id")
		p.logger.Errorf("[productcode.FetchByID] Failed to parse ID %s: %v", id, err)
		return nil, err
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

	filter := bson.M{"is_deleted": false}
	if filterParams.Search != "" {
		filter["service_name"] = bson.M{"$regex": filterParams.Search, "$options": "i"}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage
	for k, v := range filterParams.Filters {
		if k == "_id" {
			val, ok := local_util.StringToObjectID(v.(string))
			if !ok {
				continue
			}
			v = val
		}
		filter[k] = v
	}
	fmt.Println("this is the final filter", filter)

	services, err := r.producCodeDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		r.logger.Errorf("[productcode.FetchAll] Failed to fetch product codes: %v", err)
		return nil, cps_errors.ErrGeneralDBQueryFailed
	}

	var productCodes []*model.ProductCode
	for _, doc := range services {
		pc_response := productcode.ToProducCode(*doc)
		productCodes = append(productCodes, &model.ProductCode{
			ID:                 pc_response.ID,
			ProductName:        pc_response.ProductName,
			CBEProductCodes:    model.ProductCodes(pc_response.CBEProductCodes),
			CBEIFBProductCodes: model.ProductCodes(pc_response.CBEIFBProductCodes),
			CreatedAt:          pc_response.CreatedAt,
			LastUpdatedAt:      pc_response.LastUpdatedAt,
		})
	}

	total, err := r.producCodeDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[productcode.FetchAll] Failed to count product codes: %v", err)
		return nil, cps_errors.ErrGeneralDBQueryFailed
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
		p.logger.Errorf("[productcode.Update] Failed to parse ID %s: %v", productCode.ID, "invalid product id")
		return fmt.Errorf("invalid product id")
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	update := bson.M{
		"service_name":          productCode.ProductName,
		"cbe_product_codes":     productCode.CBEProductCodes,
		"cbe_ifb_product_codes": productCode.CBEIFBProductCodes,
		"last_modified_at":      time.Now(),
	}

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
