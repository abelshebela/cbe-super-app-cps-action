package Productcode

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/product_code"
	product_code_outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/product_code"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constan "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	shared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// NotificationPersistence implements the NotificationRepository interface
type RepositoryImpl struct {
	producCodeDal dal.MongoDal[model.Service, model.Service]
	logger        shared.Logger
}

// InitNotificationPersistence initializes the notification persistence layer
func InitProductCodePersistence(client *mongo.Client, dbName string, collection string, logger shared.Logger) product_code_outbound.Repository {
	return &RepositoryImpl{
		producCodeDal: dal.NewMongoDal[model.Service, model.Service](client, dbName, collection),
		logger:        logger,
	}
}

func (p *RepositoryImpl) FetchByID(ctx context.Context, id string) (*domain.ProductCode, error) {
	objID, err := common_util.ParsePrimitiveObjectID(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"_id":        objID,
		"is_deleted": false,
	}

	service, err := p.producCodeDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(common_util.NotFound)
		}

		return nil, fmt.Errorf(common_util.GeneralDBQueryFailed)
	}

	result := mappers.ToProducCode(*service)
	return result, nil
}
func (r *RepositoryImpl) FetchAll(ctx context.Context, filterParams *constan.Filter) (*utils.PaginatedResponse[[]*domain.ProductCode], error) {
	filter := bson.M{"is_deleted": false}

	if filterParams.Search != "" {
		filter["service_name"] = bson.M{"$regex": filterParams.Search, "$options": "i"}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	productCodeDocs, err := r.producCodeDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
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
		return nil, fmt.Errorf(utils.GeneralDBQueryFailed)
	}

	meta := utils.BuildPaginationMeta(total, filterParams.Page, limit)
	return &utils.PaginatedResponse[[]*domain.ProductCode]{
		Data: productCodes,
		Meta: meta,
	}, nil
}

func (p *RepositoryImpl) Update(ctx context.Context, productCode *domain.ProductCode) (*domain.ProductCode, error) {
	objID, err := common_util.ParsePrimitiveObjectID(productCode.ID)

	if err != nil {
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
			return nil, fmt.Errorf(common_util.NotFound)
		}
		return nil, fmt.Errorf(common_util.GeneralDBUpdateFailed)
	}

	return mappers.ToProducCode(res), nil
}
