package cps_action

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	local_utils "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CPSActionStorage struct {
	dal    dal.MongoDal[model.CPSAction, model.CPSAction]
	client *mongo.Client
	logger utils.Logger
}

func NewCPSActionRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.CPSActionRepository {
	return &CPSActionStorage{
		dal:    dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

// Ensure CPSActionRepository implements the storage.CPSActionRepository interface

func (r *CPSActionStorage) Save(ctx context.Context, cpsAction *model.CPSAction) error {

	_, err := r.dal.InsertOne(ctx, *cpsAction)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}
func (s *CPSActionStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	// 1. Base filter (only active records)
	filter := bson.M{
		"is_deleted": false,
		"department": department,
	}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"field1", "field2", "field3"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["field1"] = searchRegex // choose your searchable field(s)
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: data,
		Meta: meta,
	}, nil
}
func (r *CPSActionStorage) FindOne(ctx context.Context, filter model.CPSAction) (*model.CPSAction, error) {
	filterMap := BuildCPSActionFilter(filter)

	data, err := r.dal.FindOne(ctx, filterMap, nil)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, errors.New(localization.ErrorUnexpectedError.Message)
		}
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
	return data, nil
}

func (r *CPSActionStorage) Update(ctx context.Context, actionCode string, update model.CPSAction) error {
	filterMap := BuildCPSActionFilter(update)
	updateMap := BuildCPSActionUpdateMap(update)

	_, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}
func (r *CPSActionStorage) Delete(ctx context.Context, id string) error {
	idObj, ok := local_utils.StringToObjectID(id)
	if !ok {
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	filterMap := BuildCPSActionFilter(model.CPSAction{ID: idObj})

	err := r.dal.DeleteOne(ctx, filterMap)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	return nil
}
