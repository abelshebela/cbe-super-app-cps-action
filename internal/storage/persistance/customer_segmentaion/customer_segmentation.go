package customersegmentaion

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerStorage struct {
	dal        dal.MongoDal[model.CustomerSegmentation, model.CustomerSegmentation]
	client     *mongo.Client
	dbName     string
	collection string
	logger     utils.Logger
}

func NewCustomerSegmentationRepository(client *mongo.Client, dbName, collection string, logger utils.Logger) storage.CustomerSegmentationRepository {
	return &customerStorage{
		dal:        dal.NewMongoDal[model.CustomerSegmentation, model.CustomerSegmentation](client, dbName, collection),
		client:     client,
		dbName:     dbName,
		collection: collection,
		logger:     logger,
	}
}

func (r *customerStorage) Create(ctx context.Context, seg *model.CustomerSegmentation) error {
	_, err := r.dal.InsertOne(ctx, *seg)
	if err != nil {
		r.logger.Errorf("Unable to create customer segmentation with error: %s", err)
		return err
	}
	return nil
}

func (r *customerStorage) Update(ctx context.Context, id string, seg *model.CustomerSegmentation) error {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj}
	update := MapToCustomerSegUpdate(seg)
	_, err = r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("Unable to update customer segmentation with error: %s", err)
		return err
	}
	return nil
}

func (r *customerStorage) Delete(ctx context.Context, id string) error {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj}
	_, err = r.dal.UpdateOne(ctx, filter, bson.M{"is_deleted": true})
	if err != nil {
		r.logger.Errorf("Unable to delete customer segmentation with error: %s", err)
		return err
	}
	return nil
}

func (r *customerStorage) FindByID(ctx context.Context, id string) (*model.CustomerSegmentation, error) {
	obj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": obj, "is_deleted": false}
	seg, err := r.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("[FindByID] customer segmentation not found")
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		r.logger.Errorf("Unable to find customer segmentation by ID with error: %s", err)
		return nil, err
	}

	return seg, nil
}

func (r *customerStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CustomerSegmentation], error) {
	searchKeys := bson.M{}

	allowedKeys := []string{"search"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"customer_role": searchRegex},
			{"customer_segment": searchRegex},
			{"customer_sub_segment": searchRegex},
			{"customer_group": searchRegex},
		}
	}

	projection := bson.M{}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := r.dal.FindAllWithPagination(ctx, filter, projection, skip, limit)
	if err != nil {
		r.logger.Errorf("[FindAllWithPagination] failed to fetch customer segmentations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		r.logger.Errorf("[FindAllWithPagination] failed to count customer segmentations: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.CustomerSegmentation]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *customerStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*model.CustomerSegmentation, error) {
	filter := bson.M{"customer_segment": customerSegment, "is_deleted": false}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		r.logger.Errorf("[FindByCustomerSegmentation] failed to fetch customer segmentation: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	return result, nil
}
