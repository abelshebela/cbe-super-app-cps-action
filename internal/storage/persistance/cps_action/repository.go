package cps_action

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"fmt"

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
	r.logger.Infof("Attempting to save CPSAction: %+v", cpsAction)
	fmt.Printf("Cps file of id is fff:%s\n\n\n", cpsAction.ID)
	cps, err := r.dal.InsertOne(ctx, *cpsAction)
	if err != nil {
		r.logger.Errorf("Failed to save CPSAction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	fmt.Printf("Cps file of id is:%s", cps.ID)
	r.logger.Infof("Successfully saved CPSAction with ActionCode: %s", cpsAction.ActionCode)
	return nil
}

func (s *CPSActionStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	s.logger.Infof("Finding all CPSActions with pagination. Department: %s, Filter: %+v", department, filterParam)
	// 1. Base filter (only active records)
	filter := bson.M{
		"is_deleted": false,
		"department": department,
	}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"action_status", "action_type", "request_action", "maker_phone_number", "checker_phone_number", "maker_name", "maker_id", "checker_name", "checker_phone_number", "checker_id"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["field1"] = searchRegex // choose your searchable field(s)
		searchKeys["$or"] = []bson.M{
			{"maker_name": searchRegex},
			{"maker_phone_number": searchRegex},
			{"checker_name": searchRegex},
			{"checker_phone_number": searchRegex},
		}
		s.logger.Infof("Search applied with regex: %v", searchRegex)
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	s.logger.Debugf("Mongo filter: %+v, skip: %d, limit: %d", filter, skip, limit)

	// 5. Fetch data
	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("Error fetching paginated CPSActions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("Error counting total CPSActions: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_utils.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
	s.logger.Infof("Successfully fetched paginated CPSActions. Total: %d", total)
	return &types.PaginatedResponse[[]*model.CPSAction]{
		Data: data,
		Meta: meta,
	}, nil
}

func (r *CPSActionStorage) FindOne(ctx context.Context, filter model.CPSAction) (*model.CPSAction, error) {
	r.logger.Infof("Finding one CPSAction with filter: %+v", filter)
	filterMap := BuildCPSActionFilter(filter)

	data, err := r.dal.FindOne(ctx, filterMap, nil)
	if err != nil {
		r.logger.Errorf("Error finding CPSAction: %v", err)
		code, _ := local_utils.HandleMongoError(err)
		return nil, errors.New(code)
	}
	r.logger.Infof("Successfully found CPSAction: %+v", data)
	return data, nil
}

func (r *CPSActionStorage) Update(ctx context.Context, actionCode string, update model.CPSAction) error {
	r.logger.Infof("Updating CPSAction with ActionCode: %s, Update: %+v", actionCode, update)
	filterMap := BuildCPSActionFilter(update)
	updateMap := BuildCPSActionUpdateMap(update)

	_, err := r.dal.UpdateOne(ctx, filterMap, updateMap)
	if err != nil {

		r.logger.Errorf("Error updating CPSAction: %v", err)
		code, _ := local_utils.HandleMongoError(err)
		return errors.New(code)
	}
	r.logger.Infof("Successfully updated CPSAction with ActionCode: %s", actionCode)
	return nil
}

func (r *CPSActionStorage) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting CPSAction with ID: %s", id)
	idObj, ok := local_utils.StringToObjectID(id)
	if !ok {
		r.logger.Errorf("Invalid ObjectID for deletion: %s", id)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	filterMap := BuildCPSActionFilter(model.CPSAction{ID: idObj})

	err := r.dal.DeleteOne(ctx, filterMap)
	if err != nil {
		r.logger.Errorf("Error deleting CPSAction: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Message)
	}
	r.logger.Infof("Successfully deleted CPSAction with ID: %s", id)
	return nil
}
