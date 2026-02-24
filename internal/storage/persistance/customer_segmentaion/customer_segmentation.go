package customersegmentaion

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerStorage struct {
	cfg           *config.VaultConfig
	dal           dal.MongoDal[imodel.CustomerSegmentation, imodel.CustomerSegmentation]
	client        *mongo.Client
	dbName        string
	collection    string
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

func NewCustomerSegmentationRepository(client *mongo.Client, cfg *config.VaultConfig, dbName, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.CustomerSegmentationRepository {
	return &customerStorage{
		cfg:           cfg,
		dal:           dal.NewMongoDal[imodel.CustomerSegmentation, imodel.CustomerSegmentation](client, cfg, dbName, collection),
		client:        client,
		dbName:        dbName,
		collection:    collection,
		kafkaProducer: kafkaProducer,
		logger:        logger,
	}
}

func (r *customerStorage) Create(ctx context.Context, seg *imodel.CustomerSegmentation) error {
	_, err := r.dal.InsertOne(ctx, *seg)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Create] failed to create customer segmentation: %s", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *customerStorage) Update(ctx context.Context, id string, seg *imodel.CustomerSegmentation) error {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj}
	update := MapToCustomerSegUpdate(seg)
	updatedCustomerSegmentation, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Update] failed to update customer segmentation: %s", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	r.kafkaProducer.PublishMessage(
		ctx,
		updatedCustomerSegmentation,
		string(constants.ClientOrchestrationServicesTopic),
		string(constants.CustomerSegmentationUpdatedTopic),
		"customer segmentation updated",
	)

	return nil
}

func (r *customerStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][EnableOrDisable] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj, "is_deleted": false}
	update := bson.M{"is_enabled": enable, "last_modified_at": time.Now()}

	_, err = r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][EnableOrDisable] failed to enable/disable customer segmentation: %v", err)
		return local_util.HandleDBError(err)
	}
	return nil
}

func (r *customerStorage) Delete(ctx context.Context, id string) error {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj}
	_, err = r.dal.UpdateOne(ctx, filter, bson.M{"is_deleted": true})
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][Delete] failed to delete customer segmentation: %s", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (r *customerStorage) FindByID(ctx context.Context, id string) (*imodel.CustomerSegmentation, error) {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj, "is_deleted": false}
	seg, err := r.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindByID] failed to find customer segmentation: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return seg, nil
}

func (r *customerStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerSegmentation], error) {
	searchKeys := bson.M{}

	allowedKeys := []string{"search"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"customer_role.name": searchRegex},
			{"t24_customer_sub_segments.customer_segment": searchRegex},
			{"t24_customer_sub_segments.customer_sub_segment": searchRegex},
			{"t24_customer_sub_segments.customer_group": searchRegex},
		}
	}

	projection := bson.M{}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	if val, ok := filterParam.Filters["is_enabled"]; ok {
		if val != nil && val != "" {
			filter["is_enabled"] = val
		} else {
			delete(filter, "is_enabled")
		}
	}

	data, err := r.dal.FindAllWithPaginationE(ctx, filter, projection, skip, limit)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] failed to fetch customer segmentations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindAllWithPagination] failed to count customer segmentations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]imodel.CustomerSegmentation]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *customerStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CustomerSegmentation, error) {
	// New model: customer_segment is now inside t24_customer_sub_segments array
	filter := bson.M{
		"t24_customer_sub_segments.cust_segment": customerSegment,
		"is_deleted":                             false,
	}
	// Project only relevant fields (optional, can be nil for all fields)
	projection := bson.M{
		"customer_role":             1,
		"t24_customer_sub_segments": 1,
		"is_enabled":                1,
		"is_deleted":                1,
		"created_at":                1,
		"updated_at":                1,
	}
	result, err := r.dal.FindOne(ctx, filter, projection)
	if err != nil {
		r.logger.Errorf("[CustomerSegmentation][FindByCustomerSegmentation] failed to fetch customer segmentation: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}
