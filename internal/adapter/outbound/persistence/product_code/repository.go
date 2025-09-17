package productcode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	product_code_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/product_code"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constan "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// RepositoryImpl implements the ProductCodeRepository interface
type RepositoryImpl struct {
	producCodeDal dal.MongoDal[model.Service, model.Service]
	logger        shared.Logger
}

// InitProductCodePersistence initializes the product code persistence layer
func InitProductCodePersistence(client *mongo.Client, dbName string, collection string, logger shared.Logger) product_code_outbound.Repository {
	return &RepositoryImpl{
		producCodeDal: dal.NewMongoDal[model.Service, model.Service](client, dbName, collection),
		logger:        logger,
	}
}

func (p *RepositoryImpl) FetchByID(ctx context.Context, id string) (*domain.ProductCode, error) {
	p.logger.Infof("[productcode.FetchByID] Fetching product code with ID: %s", id)

	objID, err := common_util.ParsePrimitiveObjectID(id)
	if err != nil {
		p.logger.Errorf("[productcode.FetchByID] Failed to parse ID %s: %v", id, err)
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	service, err := p.producCodeDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Errorf("[productcode.FetchByID] Product code not found for ID %s: %v", id, err)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		p.logger.Errorf("[productcode.FetchByID] Database query failed for ID %s: %v", id, err)
		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	result := mappers.ToProducCode(*service)
	p.logger.Infof("[productcode.FetchByID] Successfully fetched product code with ID: %s", id)
	return result, nil
}

func (r *RepositoryImpl) FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*domain.ProductCode], error) {
	r.logger.Infof("[productcode.FetchAll] Fetching product codes with filter: %+v", filterParams)

	filter := bson.M{"is_deleted": false}
	if filterParams.Search != "" {
		filter["service_name"] = bson.M{"$regex": filterParams.Search, "$options": "i"}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	productCodeDocs, err := r.producCodeDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		r.logger.Errorf("[productcode.FetchAll] Failed to fetch product codes: %v", err)
		return nil, fmt.Errorf(utils.GeneralDBQueryFailed)
	}

	var productCodes []*domain.ProductCode
	for _, doc := range productCodeDocs {
		response := mappers.ToProducCode(*doc)
		productCodes = append(productCodes, &domain.ProductCode{
			ID:                 response.ID,
			ProductName:        response.ProductName,
			CBEProductCodes:    domain.ProductCodes(response.CBEProductCodes),
			CBEIFBProductCodes: domain.ProductCodes(response.CBEIFBProductCodes),
			CreatedAt:          response.CreatedAt,
			LastUpdatedAt:      response.LastUpdatedAt,
		})
	}

	total, err := r.producCodeDal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[productcode.FetchAll] Failed to count product codes: %v", err)
		return nil, fmt.Errorf(utils.GeneralDBQueryFailed)
	}

	meta := utils.BuildPaginationMeta(total, filterParams.Page, limit)
	r.logger.Infof("[productcode.FetchAll] Successfully fetched %d product codes, total: %d", len(productCodes), total)
	return &utils.PaginatedResponse[[]*domain.ProductCode]{
		Data: productCodes,
		Meta: meta,
	}, nil
}

func (p *RepositoryImpl) Update(ctx context.Context, productCode *domain.ProductCode) (*domain.ProductCode, error) {
	p.logger.Infof("[productcode.Update] Updating product code with ID: %s", productCode.ID)

	objID, err := common_util.ParsePrimitiveObjectID(productCode.ID)
	if err != nil {
		p.logger.Errorf("[productcode.Update] Failed to parse ID %s: %v", productCode.ID, err)
		return nil, err
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

	res, err := p.producCodeDal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			p.logger.Errorf("[productcode.Update] Product code not found for ID %s: %v", productCode.ID, err)
			return nil, fmt.Errorf(common_util.NotFound)
		}
		p.logger.Errorf("[productcode.Update] Database update failed for ID %s: %v", productCode.ID, err)
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	result := mappers.ToProducCode(res)
	p.logger.Infof("[productcode.Update] Successfully updated product code with ID: %s", productCode.ID)
	return result, nil
}
