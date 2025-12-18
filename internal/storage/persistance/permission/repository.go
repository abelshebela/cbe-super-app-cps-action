package permission

import (
	cps_user_dto "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PermissionPersistence struct {
	permissionGroupsDal   dal.MongoDal[model.PermissionGroup, model.PermissionGroup]
	permissionCategoryDal dal.MongoDal[model.PermissionCategory, model.PermissionCategory]
	permissionDal         dal.MongoDal[model.Permission, model.Permission]
	collections           []mongo.Collection
	cpsdal                dal.MongoDal[model.CPSAction, model.CPSAction]
	timeout               time.Duration
	logger                utils.Logger
}

var _ storage.PermissionRepository = (*PermissionPersistence)(nil)

func InitPermission(
	client *mongo.Client,
	dbName string,
	collectionNames []string,
	timeout time.Duration,
	logger utils.Logger,
) *PermissionPersistence {

	permissionGroupsDal := dal.NewMongoDal[model.PermissionGroup, model.PermissionGroup](client, dbName, collectionNames[0])
	permissionCategoryDal := dal.NewMongoDal[model.PermissionCategory, model.PermissionCategory](client, dbName, collectionNames[1])
	permissionDal := dal.NewMongoDal[model.Permission, model.Permission](client, dbName, collectionNames[2])
	cpsdal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames[3])

	// build []mongo.Collection
	var cols []mongo.Collection
	for _, name := range collectionNames {

		cols = append(cols, *client.Database(dbName).Collection(name))
	}

	return &PermissionPersistence{
		permissionGroupsDal:   permissionGroupsDal,
		permissionCategoryDal: permissionCategoryDal,
		permissionDal:         permissionDal,
		collections:           cols,
		cpsdal:                cpsdal,
		timeout:               timeout,
		logger:                logger,
	}
}

// Basic CRUD operations
func (r *PermissionPersistence) Create(ctx context.Context, permissionGroup *model.PermissionGroup) error {
	_, err := r.collections[0].InsertOne(ctx, *permissionGroup)
	return err
}

func (r *PermissionPersistence) Update(ctx context.Context, id string, permissionGroup *model.PermissionGroup) error {
	r.logger.Infof("[Update] updating permission group for id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[Update] invalid object id: %v", err)
		return err
	}

	filter := bson.M{"_id": objectID}
	update := PermissionGroupUpdateMapper(permissionGroup)
	update = bson.M{"$set": update}

	_, err = r.collections[0].UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[Update] failed to update permission group: %v", err)
		return err
	}
	r.logger.Infof("[Update] permission group updated successfully")
	return nil
}

func (r *PermissionPersistence) Delete(ctx context.Context, id string) error {
	r.logger.Infof("[Delete] deleting permission group for id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[Delete] invalid object id: %v", err)
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"is_deleted": true}}

	_, err = r.collections[0].UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("[Delete] failed to delete permission group: %v", err)
		return err
	}
	r.logger.Infof("[Delete] permission group deleted successfully")
	return nil
}

func (r *PermissionPersistence) FindByID(ctx context.Context, id string) (*model.PermissionGroup, error) {
	r.logger.Infof("[FindByID] fetching permission group by id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objectID, "is_deleted": bson.M{"$ne": true}}

	var result model.PermissionGroup
	err = r.collections[0].FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			r.logger.Errorf("[FindByID] permission group not found")
			return nil, errors.New(localization.ErrorPermissionGroupNotFound.Code)
		}
		r.logger.Errorf("[FindByID] failed to find permission group: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	r.logger.Infof("[FindByID] permission group retrieved successfully")
	return &result, nil
}
func (r *PermissionPersistence) FindByIDPopulated(ctx context.Context, id string) (cps_user_dto.PermissionGroupResponse, error) {
	r.logger.Infof("[FindByIDPopulated] fetching populated permission group by id: %s", id)
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.logger.Errorf("[FindByIDPopulated] invalid object id: %v", err)
		return cps_user_dto.PermissionGroupResponse{}, errors.New(localization.ErrorInvalidID.Code)
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"_id":        objectID,
			"is_deleted": bson.M{"$ne": true},
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": "permission_category",
			"let":  bson.M{"categoryIds": "$permission_category"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
					"$in": []interface{}{"$_id", bson.M{
						"$map": bson.M{
							"input": "$$categoryIds",
							"as":    "catId",
							"in": bson.M{
								"$cond": bson.M{
									"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$$catId"}, "string"}},
									"then": bson.M{"$toObjectId": "$$catId"},
									"else": "$$catId",
								},
							},
						},
					}},
				}}}},
				bson.D{{Key: "$match", Value: bson.M{
					"is_deleted": bson.M{"$ne": true},
				}}},
				bson.D{{Key: "$project", Value: bson.M{
					"category_name": 1,
					"access":        1,
					"permissions":   1,
				}}},
			},
			"as": "permission_category_docs",
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"id":                  "$_id",
			"group_name":          1,
			"permission_category": "$permission_category_docs",
			"created_at":          1,
			"updated_at":          1,
		}}},
	}

	cursor, err := r.collections[0].Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("[FindByIDPopulated] failed to aggregate permission group: %v", err)
		return cps_user_dto.PermissionGroupResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var groups []cps_user_dto.PermissionGroupResponse
	if err := cursor.All(ctx, &groups); err != nil {
		r.logger.Errorf("[FindByIDPopulated] failed to decode aggregation results: %v", err)
		return cps_user_dto.PermissionGroupResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(groups) == 0 {
		r.logger.Errorf("[FindByIDPopulated] permission group not found")
		return cps_user_dto.PermissionGroupResponse{}, errors.New(localization.ErrorPermissionGroupNotFound.Code)
	}
	r.logger.Infof("[FindByIDPopulated] populated permission group retrieved successfully")
	return groups[0], nil
}

func (s *PermissionPersistence) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": bson.M{"$ne": true}}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"enabled", "is_deleted", "department_id", "role", "realm", "group_name"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["group_name"] = searchRegex
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	pipeline := PermissionGroupsPipeline(filter, skip, limit)

	s.logger.Infof("[FindAllWithPagination] fetching permission groups with pagination")
	cursor, err := s.collections[0].Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to aggregate permission groups: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	var data []*model.PermissionGroup
	if err := cursor.All(ctx, &data); err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to decode aggregation results: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.collections[0].CountDocuments(ctx, filter)
	if err != nil {
		s.logger.Errorf("[FindAllWithPagination] failed to count permission groups: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	s.logger.Infof("[FindAllWithPagination] retrieved %d permission groups", len(data))

	// 8. Return standard paginated response
	return &types.PaginatedResponse[[]*model.PermissionGroup]{
		Data: data,
		Meta: meta,
	}, nil
}

func (s *PermissionPersistence) FindAllGroupsWithPagination(ctx context.Context, departmentId string, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	filter := bson.M{
		"department_id": departmentId,
		"is_deleted":    bson.M{"$ne": true},
	}

	page := filterParam.Page
	perPage := filterParam.PerPage

	skip := int64((page - 1) * perPage)
	limit := int64(perPage)

	pipeline := PermissionGroupsPipeline(filter, skip, limit)

	cursor, err := s.collections[0].Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	var data []*model.PermissionGroup
	if err := cursor.All(ctx, &data); err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	total, err := s.collections[0].CountDocuments(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	meta := local_util.BuildPaginationMeta(total, page, int(limit))

	return &types.PaginatedResponse[[]*model.PermissionGroup]{
		Data: data,
		Meta: meta,
	}, nil
}

// Permission category operations
func (r *PermissionPersistence) ValidatePermissionCategories(ctx context.Context, categoryIDs []string) ([]string, error) {
	if len(categoryIDs) == 0 {
		return []string{}, nil
	}

	objectIDs := make([]bson.ObjectID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"is_deleted": bson.M{"$ne": true},
	}

	var categories []model.PermissionCategory
	cursor, err := r.collections[1].Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	validIDs := make([]string, 0, len(categories))
	for _, category := range categories {
		validIDs = append(validIDs, category.ID.Hex())
	}

	return validIDs, nil
}

func (r *PermissionPersistence) GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error) {
	filter := bson.M{"is_deleted": bson.M{"$ne": true}}

	var categories []model.PermissionCategory
	cursor, err := r.collections[1].Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	result := make([]*model.PermissionCategory, len(categories))
	for i, cat := range categories {
		result[i] = &cat
	}
	return result, nil
}

// Permission group operations
func (r *PermissionPersistence) ValidatePermissionGroups(ctx context.Context, groupIDs []string) ([]string, error) {
	if len(groupIDs) == 0 {
		return []string{}, nil
	}

	objectIDs := make([]bson.ObjectID, 0, len(groupIDs))
	for _, id := range groupIDs {
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"is_deleted": bson.M{"$ne": true},
	}

	var groups []model.PermissionGroup
	cursor, err := r.collections[0].Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &groups); err != nil {
		return nil, err
	}

	validIDs := make([]string, 0, len(groups))
	for _, group := range groups {
		validIDs = append(validIDs, group.ID.Hex())
	}

	return validIDs, nil
}

func (r *PermissionPersistence) CheckPermissionGroupExists(groupName string) bool {
	groupName = strings.ToUpper(groupName)
	ctx := context.Background()

	filter := bson.M{"group_name": groupName, "is_deleted": bson.M{"$ne": true}}

	var result model.PermissionGroup
	err := r.collections[0].FindOne(ctx, filter).Decode(&result)
	if err != nil {
		r.logger.Errorf("CheckPermissionGroupExists failed: %v", err)
		return false
	}

	return true
}

func (r *PermissionPersistence) GetPermissionGroup(groupName string) (*model.PermissionGroup, error) {
	groupName = strings.ToUpper(groupName)
	ctx := context.Background()

	filter := bson.M{"group_name": groupName, "is_deleted": bson.M{"$ne": true}}

	var group model.PermissionGroup
	err := r.collections[0].FindOne(ctx, filter).Decode(&group)
	if err != nil {
		return nil, err
	}
	if group.ID.IsZero() {
		return nil, nil
	}

	// Populate full permission categories if PermissionCategory holds IDs
	var categoryIDs []bson.ObjectID
	switch v := group.PermissionCategory.(type) {
	case []bson.ObjectID:
		categoryIDs = v
	case []string:
		for _, s := range v {
			if oid, err := bson.ObjectIDFromHex(s); err == nil {
				categoryIDs = append(categoryIDs, oid)
			}
		}
	case bson.A:
		for _, raw := range v {
			switch val := raw.(type) {
			case bson.ObjectID:
				categoryIDs = append(categoryIDs, val)
			case string:
				if oid, err := bson.ObjectIDFromHex(val); err == nil {
					categoryIDs = append(categoryIDs, oid)
				}
			}
		}
	}

	if len(categoryIDs) > 0 {
		catFilter := bson.M{"_id": bson.M{"$in": categoryIDs}}
		var cats []model.PermissionCategory
		catCursor, err := r.collections[1].Find(ctx, catFilter)
		if err == nil {
			defer catCursor.Close(ctx)
			if catErr := catCursor.All(ctx, &cats); catErr == nil && len(cats) > 0 {
				// Replace the interface field with the full documents
				group.PermissionCategory = cats
			}
		}
	}

	return &group, nil
}

func (r *PermissionPersistence) GetPermissionGroupById(ctx context.Context, id string) (*model.PermissionGroup, error) {
	group, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, nil
	}

	var categoryIDs []bson.ObjectID
	switch v := group.PermissionCategory.(type) {
	case []bson.ObjectID:
		categoryIDs = v
	case []string:
		for _, s := range v {
			if oid, err := bson.ObjectIDFromHex(s); err == nil {
				categoryIDs = append(categoryIDs, oid)
			}
		}
	case bson.A:
		for _, raw := range v {
			switch val := raw.(type) {
			case bson.ObjectID:
				categoryIDs = append(categoryIDs, val)
			case string:
				if oid, err := bson.ObjectIDFromHex(val); err == nil {
					categoryIDs = append(categoryIDs, oid)
				}
			}
		}
	}

	if len(categoryIDs) > 0 {
		catFilter := bson.M{"_id": bson.M{"$in": categoryIDs}}
		var cats []model.PermissionCategory
		catCursor, err := r.collections[1].Find(ctx, catFilter)
		if err == nil {
			defer catCursor.Close(ctx)
			if catErr := catCursor.All(ctx, &cats); catErr == nil && len(cats) > 0 {
				group.PermissionCategory = cats
			}
		}
	}

	return group, nil
}
func (p *PermissionPersistence) ValidatePermissionGroupByID(ctx context.Context, ids []string) (bool, error) {
	if len(ids) == 0 {
		return false, fmt.Errorf("PERMISSION_GROUP_ARRAY_EMPTY")
	}

	var cleaned []string
	for _, id := range ids {
		if trimmed := strings.TrimSpace(id); trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}

	if len(cleaned) == 0 {
		return false, fmt.Errorf("NO_VALID_PERMISSION_GROUP_ID")
	}
	// Query DB for all given ids
	objectIDs := make([]bson.ObjectID, 0, len(cleaned))
	for _, idStr := range cleaned {
		oid, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			return false, fmt.Errorf("invalid id: %s: %w", idStr, err)
		}
		objectIDs = append(objectIDs, oid)
	}
	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	count, err := p.collections[0].CountDocuments(ctx, filter)
	if err != nil {
		return false, fmt.Errorf("DB_ERROR: %w", err)
	}

	if count != int64(len(cleaned)) {
		return false, fmt.Errorf("PERMISSION_GROUP_NOT_FOUND")
	}

	return true, nil
}

func (p *PermissionPersistence) GetAllPermissionCategories(
	ctx context.Context,
	card string,
) ([]*model.PermissionCategory, error) {
	mongoFilter := bson.M{
		"is_deleted":  bson.M{"$ne": true},
		"portal_card": card,
	}
	var categories []model.PermissionCategory
	cursor, err := p.collections[1].Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	result := make([]*model.PermissionCategory, len(categories))
	for i, cat := range categories {
		result[i] = &cat
	}
	return result, nil
}

func (p *PermissionPersistence) GetPopulatedPermissionCategories(ctx context.Context, categoryIDs []string) ([]cps_user_dto.PermissionCategoryResponse, error) {
	if len(categoryIDs) == 0 {
		return []cps_user_dto.PermissionCategoryResponse{}, nil
	}
	objectIDs := make([]bson.ObjectID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	filter := bson.M{
		"_id":        bson.M{"$in": objectIDs},
		"is_deleted": bson.M{"$ne": true},
	}

	var categories []model.PermissionCategory
	cursor, err := p.collections[1].Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	result := make([]cps_user_dto.PermissionCategoryResponse, 0, len(categories))
	for _, cat := range categories {
		var permissions []cps_user_dto.PermissionResponse
		if cat.Permissions != nil {
			switch perms := cat.Permissions.(type) {
			case []interface{}:
				for _, p := range perms {
					if permMap, ok := p.(map[string]interface{}); ok {
						perm := cps_user_dto.PermissionResponse{}
						if id, ok := permMap["_id"].(bson.ObjectID); ok {
							perm.ID = id
						}
						if name, ok := permMap["permission_name"].(string); ok {
							perm.Name = name
						}
						permissions = append(permissions, perm)
					}
				}
			}
		}

		result = append(result, cps_user_dto.PermissionCategoryResponse{
			ID:           cat.ID,
			Access:       cat.Access,
			CategoryName: cat.CategoryName,
			Permissions:  permissions,
		})
	}

	return result, nil
}

func (p *PermissionPersistence) GetPopulatedPermissionGroups(ctx context.Context, groupIDs []string) ([]cps_user_dto.PermissionGroupResponse, error) {
	if len(groupIDs) == 0 {
		return []cps_user_dto.PermissionGroupResponse{}, nil
	}

	objectIDs := make([]bson.ObjectID, 0, len(groupIDs))
	for _, id := range groupIDs {
		objectID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"_id":        bson.M{"$in": objectIDs},
			"is_deleted": bson.M{"$ne": true},
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from": p.collections[1].Name(),
			"let":  bson.M{"categoryIds": "$permission_category"},
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.M{"$expr": bson.M{
					"$in": []interface{}{
						"$_id",
						bson.M{
							"$map": bson.M{
								"input": "$$categoryIds",
								"as":    "catId",
								"in": bson.M{
									"$cond": bson.M{
										"if":   bson.M{"$eq": []interface{}{bson.M{"$type": "$$catId"}, "string"}},
										"then": bson.M{"$toObjectId": "$$catId"},
										"else": "$$catId",
									},
								},
							},
						},
					},
				}}}},
				bson.D{{Key: "$match", Value: bson.M{"is_deleted": bson.M{"$ne": true}}}},
				bson.D{{Key: "$project", Value: bson.M{
					"category_name": 1,
					"access":        1,
				}}},
			},
			"as": "permission_category",
		}}},
		bson.D{{Key: "$project", Value: bson.M{
			"id":                  "$_id",
			"group_name":          1,
			"permission_category": 1,
			"department_id":       1,
			"role":                1,
			"realm":               1,
			"enabled":             1,
			"created_at":          1,
			"updated_at":          1,
		}}},
	}

	cursor, err := p.collections[0].Aggregate(ctx, pipeline)
	if err != nil {
		p.logger.Errorf("Failed to aggregate Permission Groups, groupIDs: %v, error: %v", groupIDs, err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var groups []cps_user_dto.PermissionGroupResponse
	if err := cursor.All(ctx, &groups); err != nil {
		p.logger.Errorf("Failed to decode Permission Groups aggregation, groupIDs: %v, error: %v", groupIDs, err)
		return nil, err
	}

	return groups, nil
}
