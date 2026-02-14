package cpsroles

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cpsRoleStorage struct {
	dal        dal.MongoDal[model.CPSRoles, model.CPSRoles]
	client     *mongo.Client
	dbName     string
	collection string
	col        *mongo.Collection
	logger     utils.Logger
}

func NewCPSRolesStorage(client *mongo.Client, cfg *config.VaultConfig, dbName, collection string, logger utils.Logger) storage.CPSRolesRepository {
	return &cpsRoleStorage{
		dal:        dal.NewMongoDal[model.CPSRoles, model.CPSRoles](client, cfg, dbName, collection),
		client:     client,
		dbName:     dbName,
		collection: collection,
		logger:     logger,
	}
}

func (m *cpsRoleStorage) Create(ctx context.Context, req model.CPSRoles) error {
	_, err := m.dal.InsertOne(ctx, req)
	if err != nil {
		m.logger.Errorf("[Create] failed to create cps role: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *cpsRoleStorage) Update(ctx context.Context, id string, req model.CPSRoles) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[Update] invalid id format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"name":       req.Name,
		"enabled":    req.Enabled,
		"updated_at": time.Now(),
	}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[Update] cps role not found, id: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[Update] failed to update cps role, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *cpsRoleStorage) FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.CPSRoles], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "enabled"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"name": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(*filterParam, searchKeys, allowedKeys)

	data, err := m.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		m.logger.Errorf("[FindAllWithPagination] failed to fetch paginated cps roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := m.dal.TotalCount(ctx, filter)
	if err != nil {
		m.logger.Errorf("[FindAllWithPagination] failed to count total cps roles: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.CPSRoles]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *cpsRoleStorage) FindById(ctx context.Context, id string) (*model.CPSRoles, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[FindById] invalid id format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	result, err := m.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[FindById] cps role not found, id: %s", id)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[FindById] failed to find cps role, id: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// Use MongoDB pipeline to group action_names by maker/checker/auditor index
	approverCol := m.client.Database(m.dbName).Collection("cps_action_approver_index")
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{{"role_id", result.RoleCode}}}},
		bson.D{{Key: "$facet", Value: bson.M{
			"maker": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{{"maker_index", bson.D{{"$ne", nil}}}}}},
				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
			},
			"checker": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{{"checker_index", bson.D{{"$ne", nil}}}}}},
				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
			},
			"auditor": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{{"auditor_index", bson.D{{"$ne", nil}}}}}},
				bson.D{{Key: "$group", Value: bson.D{{"_id", nil}, {"actions", bson.D{{"$addToSet", "$action_name"}}}}}},
			},
		}}},
	}
	cursor, err := approverCol.Aggregate(ctx, pipeline)
	if err != nil {
		m.logger.Errorf("Error aggregating cps_action_approver_index: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var facetResult []bson.M
	if err := cursor.All(ctx, &facetResult); err != nil {
		m.logger.Errorf("Error decoding facet result: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	var maker, checker, auditor []string
	if len(facetResult) > 0 {
		if arr, ok := facetResult[0]["maker"].([]any); ok && len(arr) > 0 {
			if doc, ok := arr[0].(bson.M); ok {
				if actions, ok := doc["actions"].([]string); ok {
					maker = actions
				} else if actions, ok := doc["actions"].([]any); ok {
					for _, a := range actions {
						if s, ok := a.(string); ok {
							maker = append(maker, s)
						}
					}
				}
			}
		}
		if arr, ok := facetResult[0]["checker"].([]any); ok && len(arr) > 0 {
			if doc, ok := arr[0].(bson.M); ok {
				if actions, ok := doc["actions"].([]string); ok {
					checker = actions
				} else if actions, ok := doc["actions"].([]any); ok {
					for _, a := range actions {
						if s, ok := a.(string); ok {
							checker = append(checker, s)
						}
					}
				}
			}
		}
		if arr, ok := facetResult[0]["auditor"].([]any); ok && len(arr) > 0 {
			if doc, ok := arr[0].(bson.M); ok {
				if actions, ok := doc["actions"].([]string); ok {
					auditor = actions
				} else if actions, ok := doc["actions"].([]any); ok {
					for _, a := range actions {
						if s, ok := a.(string); ok {
							auditor = append(auditor, s)
						}
					}
				}
			}
		}
	}

	result.MakerActions = maker
	result.CheckerActions = checker
	result.AuditorActions = auditor

	return result, nil
}

func (m *cpsRoleStorage) FindByNameOrRoleCode(ctx context.Context, name, roleCode string) (*model.CPSRoles, error) {
	filter := []bson.M{}
	if roleCode == "" {
		filter = append(filter, bson.M{"name": bson.M{"$regex": "^" + name, "$options": "i"}})
	}
	if name == "" {
		filter = append(filter, bson.M{"role_code": bson.M{"$regex": "^" + roleCode, "$options": "i"}})
	}
	orFilter := bson.M{"$or": filter}
	result, err := m.dal.FindOne(ctx, orFilter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[FindByNameOrRoleCode] cps role not found, name: %s, roleCode: %s", name, roleCode)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[FindByNameOrRoleCode] failed to find cps role, name: %s, roleCode: %s, error: %v", name, roleCode, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (m *cpsRoleStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("[EnableOrDisable] invalid id format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[EnableOrDisable] cps role not found, id: %s", id)
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[EnableOrDisable] failed to enable/disable cps role, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *cpsRoleStorage) FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*model.CPSRoles, error) {
	filter := bson.M{"name": customerSegment, "enabled": true}
	result, err := m.dal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("[FindByCustomerSegmentation] cps role not found for customer segment: %s", customerSegment)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		m.logger.Errorf("[FindByCustomerSegmentation] failed to find cps role for customer segment: %s, error: %v", customerSegment, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}
