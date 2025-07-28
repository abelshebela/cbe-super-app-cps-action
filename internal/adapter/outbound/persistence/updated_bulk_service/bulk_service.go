package updatedbulkservice

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/updated_bulk_service"
	outbound "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound/updated_bulk_service"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
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

func (b BulkServicePersistence) EnableOrDisableBulkService(ctx context.Context, serviceCodes []string, cpsAction model.CPSAction, requestActionType model.RequestAction) error {
	filter := bson.M{"service_code": bson.M{"$in": serviceCodes}, "is_deleted": false}

	// Check if the service exists
	_, err := b.mongoDalbulkService.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("SURVICE_NOT_FOUND")
		}
		return fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}

	// Check if there is a pending action
	makerData := ctx_util.ExtractContext(ctx)
	pendingFilter := bson.M{
		"maker_id":       makerData.UserID,
		"department":     makerData.Department,
		"action_status":  "PENDING",
		"request_action": requestActionType,
	}
	projection := bson.M{}
	pendingAction, err := b.mongoDalCpsAction.FindOne(ctx, pendingFilter, projection)
	if err != nil && err != mongo.ErrNoDocuments {
		return fmt.Errorf("DATABASE_ERROR_CHECKING_PENDING_ACTION")
	}
	if pendingAction != nil {
		return fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	// Create the CPS action
	_, err = b.mongoDalCpsAction.InsertOne(ctx, cpsAction)
	if err != nil {
		return fmt.Errorf("database error while creating user request action")
	}

	return nil
}

func (b BulkServicePersistence) AuthorizeBulkServiceEnable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	data, err := common_util.JsonUnmarshal[[]domain.ServiceDetails](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"service_code": service.ServiceCode}
		update := bson.M{"enabled": true}
		_, err = b.mongoDalbulkService.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}

func (b BulkServicePersistence) AuthorizeBulkServiceDisable(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	data, err := common_util.JsonUnmarshal[[]domain.ServiceDetails](action.CurrentAction)
	if err != nil {
		return nil, err
	}

	for _, service := range *data {
		filter := bson.M{"service_code": service.ServiceCode}
		update := bson.M{"enabled": false}
		_, err = b.mongoDalbulkService.UpdateOne(ctx, filter, update)
		if err != nil {
			return nil, err
		}
	}

	return action, nil
}
