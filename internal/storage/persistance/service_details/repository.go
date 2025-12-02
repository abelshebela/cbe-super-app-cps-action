package service_details

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ServiceDetailsStorage struct {
	dal    dal.MongoDal[model.ServiceDetails, model.ServiceDetails]
	client *mongo.Client
	logger utils.Logger
}

func NewServiceDetailsRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ServiceDetailsRepository {
	return &ServiceDetailsStorage{
		dal:    dal.NewMongoDal[model.ServiceDetails, model.ServiceDetails](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (s *ServiceDetailsStorage) Create(ctx context.Context, details *model.ServiceDetails) error {
	_, err := s.dal.InsertOne(ctx, *details)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServiceDetailsStorage) Update(ctx context.Context, id string, details *model.ServiceDetails) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := ServiceDetailsMapper(*details)

	_, err = s.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (s *ServiceDetailsStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	return s.dal.DeleteOne(ctx, filter)
}

func (s *ServiceDetailsStorage) FindByID(ctx context.Context, projection bson.M, id string) (*model.ServiceDetails, error) {
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		s.logger.Errorf("Invalid ObjectID for fetch by id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": idObj, "is_deleted": false}

	result, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		s.logger.Errorf("Error finding CPSAction: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	s.logger.Infof("Successfully found ServiceDetails: %+v", result)
	return result, nil

}

func (s *ServiceDetailsStorage) FindAllWithPagination(ctx context.Context, projection bson.M, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ServiceDetails], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"payment_type", "min_amount", "enabled", "service_type", "service_code", "service_name"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{{"service_name": searchRegex}, {"service_code": searchRegex}, {"service_type": searchRegex}, {"product_codes": searchRegex}}
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, projection, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.ServiceDetails]{
		Data: data,
		Meta: meta,
	}, nil
}
