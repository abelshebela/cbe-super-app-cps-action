package permission

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PermissionPersistence struct {
	permissionGroupsDal   dal.MongoDal[model.PermissionGroup, model.PermissionGroup]
	permissionCategoryDal dal.MongoDal[model.PermissionCategory, model.PermissionCategory]
	permissionDal         dal.MongoDal[model.Permission, model.Permission]
	cpsdal                dal.MongoDal[model.CPSAction, model.CPSAction]
	timeout               time.Duration
	logger                utils.Logger
}

var _ storage.PermissionRepository = (*PermissionPersistence)(nil)

func InitPermission(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *PermissionPersistence {
	permissionGroupsDal := dal.NewMongoDal[model.PermissionGroup, model.PermissionGroup](client, dbName, "permission_groups")
	permissionCategoryDal := dal.NewMongoDal[model.PermissionCategory, model.PermissionCategory](client, dbName, "permission_category")
	permissionDal := dal.NewMongoDal[model.Permission, model.Permission](client, dbName, "permission")
	cpsdal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "cps_actions")

	return &PermissionPersistence{
		permissionGroupsDal:   permissionGroupsDal,
		permissionCategoryDal: permissionCategoryDal,
		permissionDal:         permissionDal,
		cpsdal:                cpsdal,
		timeout:               timeout,
		logger:                logger,
	}
}

// Basic CRUD operations
func (r *PermissionPersistence) Create(ctx context.Context, permissionGroup *model.PermissionGroup) error {
	_, err := r.permissionGroupsDal.InsertOne(ctx, *permissionGroup)
	return err
}

func (r *PermissionPersistence) Update(ctx context.Context, id string, permissionGroup *model.PermissionGroup) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": permissionGroup}

	_, err = r.permissionGroupsDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *PermissionPersistence) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": bson.M{"is_deleted": true}}

	_, err = r.permissionGroupsDal.UpdateOne(ctx, filter, update)
	return err
}

func (r *PermissionPersistence) FindByID(ctx context.Context, id string) (*model.PermissionGroup, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID, "is_deleted": bson.M{"$ne": true}}

	return r.permissionGroupsDal.FindOne(ctx, filter, bson.M{})
}

func (s *PermissionPersistence) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error) {
	// 1. Base filter (only active records)
	filter := bson.M{"is_deleted": false}
	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"enabled", "is_deleted", "role", "realm", "group_name"}

	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["group_name"] = searchRegex
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := s.permissionGroupsDal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 6. Count total
	total, err := s.permissionGroupsDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	// 8. Return standard paginated response
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

	categories, err := r.permissionCategoryDal.FindAllWithPagination(ctx, filter, bson.M{}, 0, 0)
	if err != nil {
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

	return r.permissionCategoryDal.FindAllWithPagination(ctx, filter, bson.M{}, 0, 0)
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

	groups, err := r.permissionGroupsDal.FindAllWithPagination(ctx, filter, bson.M{}, 0, 0)
	if err != nil {
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

	result, err := r.permissionGroupsDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		r.logger.Errorf("CheckPermissionGroupExists failed: %v", err)
		return false
	}

	return result != nil
}

func (r *PermissionPersistence) GetPermissionGroup(groupName string) (*model.PermissionGroup, error) {
	groupName = strings.ToUpper(groupName)
	ctx := context.Background()

	filter := bson.M{"group_name": groupName, "is_deleted": bson.M{"$ne": true}}

	group, err := r.permissionGroupsDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		return nil, err
	}
	if group == nil {
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
		cats, err := r.permissionCategoryDal.FindAll(ctx, bson.M{"_id": bson.M{"$in": categoryIDs}}, bson.M{})
		if err == nil && len(cats) > 0 {
			// Replace the interface field with the full documents
			group.PermissionCategory = cats
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
	filter := bson.M{"_id": bson.M{"$in": cleaned}}
	projection := bson.M{"_id": 1}

	groups, err := p.permissionGroupsDal.FindAll(ctx, filter, projection)
	if err != nil {
		return false, fmt.Errorf("DB_ERROR: %w", err)
	}

	if len(cleaned) != len(groups) {
		return false, fmt.Errorf("PERMISSION_GROUP_NOT_FOUND")
	}

	return true, nil
}
