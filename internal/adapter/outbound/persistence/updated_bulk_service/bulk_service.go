package updatedbulkservice

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/updated_bulk_service"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/updated_bulk_service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BulkServicePersistence struct {
	mongoDalCpsAction   dal.MongoDal[model.CPSAction, model.CPSAction]
	mongoDalbulkService dal.MongoDal[domain.ServiceDetails, domain.ServiceDetails]
	logger              utils.Logger
}

func InitBulkServicePersistence(client *mongo.Client, dbName string, collections []string, logger utils.Logger) outbound.BulkServiceRepository {
	mongoDalCpsAction := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collections[0])
	mongoDalbulkService := dal.NewMongoDal[domain.ServiceDetails, domain.ServiceDetails](client, dbName, collections[1])
	return &BulkServicePersistence{
		mongoDalCpsAction:   mongoDalCpsAction,
		mongoDalbulkService: mongoDalbulkService,
		logger:              logger,
	}
}

func (b BulkServicePersistence) GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*domain.ServiceDetails], error) {
	filter := bson.M{
		"is_deleted": false,
	}
	projection := bson.M{}

	// Add search functionality
	if filterParams.Search != "" {
		searchFilter := bson.M{
			"$or": []bson.M{
				{"key": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"serviceCode": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"serviceName": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"serviceType": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"paymentType": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"prefix": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"productCodes.PRD": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"productCodes.TRXN": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"IFBproductCodes.PRD": bson.M{"$regex": filterParams.Search, "$options": "i"}},
				{"IFBproductCodes.TRXN": bson.M{"$regex": filterParams.Search, "$options": "i"}},
			},
		}

		filter = bson.M{
			"$and": []bson.M{
				{"is_deleted": false},
				searchFilter,
			},
		}
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	bulkServices, err := b.mongoDalbulkService.FindAllWithPagination(ctx, filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	totalDocs, err := b.mongoDalbulkService.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := common_util.BuildPaginationMeta(totalDocs, page, limit)

	return &common_util.PaginatedResponse[[]*domain.ServiceDetails]{
		Data: bulkServices,
		Meta: meta,
	}, nil
}
