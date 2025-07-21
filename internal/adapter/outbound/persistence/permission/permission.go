package permission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repository "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"time"
)

type PermissionPersistence struct {
	permissionGroupsDal   dal.MongoDal[entities.PermissionGroup, entities.PermissionGroup]
	permissionCategoryDal dal.MongoDal[entities.PermissionCategory, entities.PermissionCategory]
	cpsdal                dal.MongoDal[model.CPSAction, model.CPSAction]
	timeout               time.Duration
	logger                utils.Logger
}

var _ repository.PermissionGroupRepository = (*PermissionPersistence)(nil)
var _ repository.CPSActionRepository = (*PermissionPersistence)(nil)
var _ repository.PermissionCategoryRepository = (*PermissionPersistence)(nil)

func InitPermission(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *PermissionPersistence {
	permissionGroupsDal := dal.NewMongoDal[entities.PermissionGroup, entities.PermissionGroup](client, dbName, "permission_groups")

	permissionCategoryDal := dal.NewMongoDal[entities.PermissionCategory, entities.PermissionCategory](client, dbName, "permission_category")

	cpsdal := dal.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, "cps_actions")
	return &PermissionPersistence{
		permissionGroupsDal:   permissionGroupsDal,
		permissionCategoryDal: permissionCategoryDal,
		cpsdal:                cpsdal,
		timeout:               timeout,
		logger:                logger,
	}
}

func (r *PermissionPersistence) CheckPendingRequest(userCode string, status model.ActionStatus, action model.RequestAction) error {
	ctx := context.Background()

	filter := bson.M{
		"maker_id":       userCode,
		"action_status":  status,
		"request_action": action,
	}

	result, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		r.logger.Errorf("Check Pending Request failed:", err)
		return err
	}

	if result != nil {
		return fmt.Errorf("user already has a pending request for this action")
	}
	return nil
}

func (r *PermissionPersistence) CheckPermissionGroupExists(groupName string) bool {
	ctx := context.Background()

	filter := bson.M{"group_name": groupName}
	r.logger.Infof("here create PermissionPersistence", groupName)
	result, err := r.permissionGroupsDal.FindOne(ctx, filter, bson.M{})

	fmt.Println("check permission", result)
	if err != nil {
		r.logger.Errorf("CheckPermissionGroupExists failed:", err)
		return false
	}

	if result == nil {
		r.logger.Infof("no permission Group found ")
		return false
	}
	return true
}

// InvalidIDsError is a structured error for reporting missing IDs in validation
// Field is the type of entity (e.g., "permission categories"), IDs is the list of missing IDs
type InvalidIDsError struct {
	Field string
	IDs   []string
}

func (e InvalidIDsError) Error() string {
	return fmt.Sprintf("%s not found: %v", e.Field, e.IDs)
}

func (r *PermissionPersistence) ValidatePermissionCategories(ids []string) ([]string, error) {
	ctx := context.Background()

	var validObjectIDs []bson.ObjectID
	var invalidIDs []string

	for _, id := range ids {
		if id == "" {
			invalidIDs = append(invalidIDs, "<empty>")
			continue
		}

		if len(id) != 24 {
			invalidIDs = append(invalidIDs, fmt.Sprintf("%s (length %d)", id, len(id)))
			continue
		}

		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			invalidIDs = append(invalidIDs, fmt.Sprintf("%s (invalid format)", id))
			continue
		}
		validObjectIDs = append(validObjectIDs, objID)
	}

	if len(invalidIDs) > 0 {
		return nil, fmt.Errorf("invalid permission category IDs: %v", invalidIDs)
	}

	categories, err := r.permissionCategoryDal.FindAll(ctx,
		bson.M{"_id": bson.M{"$in": validObjectIDs}},
		bson.M{"_id": 1},
	)
	if err != nil {
		r.logger.Errorf("failed to fetch permission categories: %v", err)
		return nil, fmt.Errorf("database error while validating permissions")
	}

	foundIDs := make(map[bson.ObjectID]bool)
	for _, cat := range categories {
		if !cat.ID.IsZero() {
			foundIDs[cat.ID] = true
		}
	}

	var missingIDs []string
	var validIDs []string
	for _, objID := range validObjectIDs {
		if foundIDs[objID] {
			validIDs = append(validIDs, objID.Hex())
		} else {
			missingIDs = append(missingIDs, objID.Hex())
		}
	}

	if len(missingIDs) > 0 {
		return nil, InvalidIDsError{Field: "permission categories", IDs: missingIDs}
	}

	return validIDs, nil
}

func (r *PermissionPersistence) ValidatePermissionGroups(ids []string) ([]string, error) {
	ctx := context.Background()

	var validObjectIDs []bson.ObjectID
	var invalidIDs []string

	for _, id := range ids {
		if id == "" {
			invalidIDs = append(invalidIDs, "<empty>")
			continue
		}

		if len(id) != 24 {
			invalidIDs = append(invalidIDs, fmt.Sprintf("%s (length %d)", id, len(id)))
			continue
		}

		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			invalidIDs = append(invalidIDs, fmt.Sprintf("%s (invalid format)", id))
			continue
		}
		validObjectIDs = append(validObjectIDs, objID)
	}

	if len(invalidIDs) > 0 {
		return nil, fmt.Errorf("invalid permission group IDs: %v", invalidIDs)
	}

	groups, err := r.permissionGroupsDal.FindAll(ctx,
		bson.M{"_id": bson.M{"$in": validObjectIDs}},
		bson.M{"_id": 1},
	)
	if err != nil {
		r.logger.Errorf("failed to fetch permission groups: %v", err)
		return nil, fmt.Errorf("database error while validating permission groups")
	}

	foundIDs := make(map[bson.ObjectID]bool)
	for _, group := range groups {
		if !group.ID.IsZero() {
			foundIDs[group.ID] = true
		}
	}

	var missingIDs []string
	var validIDs []string
	for _, objID := range validObjectIDs {
		if foundIDs[objID] {
			validIDs = append(validIDs, objID.Hex())
		} else {
			missingIDs = append(missingIDs, objID.Hex())
		}
	}

	if len(missingIDs) > 0 {
		return nil, InvalidIDsError{Field: "permission groups", IDs: missingIDs}
	}

	return validIDs, nil
}

func (r *PermissionPersistence) CreatePermissionGroup(group model.CPSAction) (model.CPSAction, error) {
	ctx := context.Background()
	group.ID = bson.NewObjectID()

	cps_action, err := r.cpsdal.InsertOne(ctx, group)
	return cps_action, err
}

func (r *PermissionPersistence) ValidateActionRequest(actionCode, department string) (model.CPSAction, error) {
	ctx := context.Background()

	filter := bson.M{
		"action_code":   actionCode,
		"action_status": "PENDING",
	}

	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.CPSAction{}, errors.New("no pending request found for this action")
		}
		return model.CPSAction{}, err
	}

	if action == nil {
		return model.CPSAction{}, errors.New("no pending request found for this action")
	}

	if action.Department != department {
		return model.CPSAction{}, errors.New("you are not allowed to approve this request")
	}

	return *action, nil
}

func (r *PermissionPersistence) CreatePermissionGroupFromAction(ctx context.Context, action *cps_entities.CPSAction) (*cps_entities.CPSAction, error) {

	var actionData map[string]interface{}
	bytes, err := json.Marshal(action.CurrentAction)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}
	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal CurrentAction: %v", err)
	}

	groupName, ok := actionData["group_name"].(string)
	if !ok {
		return nil, errors.New("invalid group_name")
	}

	permissionCategoriesIface, ok := actionData["permission_categories"].([]interface{})
	if !ok {
		return nil, errors.New("invalid permission_categories")
	}

	permissionCategories := make([]entities.PermissionCategory, len(permissionCategoriesIface))
	for i, v := range permissionCategoriesIface {
		idStr, ok := v.(string)
		if !ok {
			return nil, errors.New("permission_categories contains non-string value")
		}
		objID, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, fmt.Errorf("invalid permission category id: %s", idStr)
		}
		permissionCategories[i] = entities.PermissionCategory{
			ID: objID,
		}
	}

	role, ok := actionData["role"].(string)
	if !ok {
		return nil, errors.New("invalid role")
	}

	realm, ok := actionData["realm"].(string)
	if !ok {
		return nil, errors.New("invalid realm")
	}

	newPermissionGroup := entities.PermissionGroup{
		GroupName:          groupName,
		PermissionCategory: permissionCategories,
		Role:               role,
		Realm:              realm,
		CreatedAt:          time.Now(),
		LastModified:       time.Now(),
	}

	_, err = r.permissionGroupsDal.InsertOne(ctx, newPermissionGroup)
	if err != nil {
		r.logger.Errorf("Error creating permission group from action: %v", err)
		return nil, err
	}

	return action, nil
}

func (r *PermissionPersistence) RejectActionRequest(actionCode string, action model.CPSAction, rejectedReason string) (model.CPSAction, error) {
	ctx := context.Background()
	filter := bson.M{"action_code": actionCode}
	update := bson.M{
		"checker_name":         action.CheckerName,
		"checker_id":           action.CheckerID,
		"rejection_reason":     rejectedReason,
		"checker_phone_number": action.CheckerPhoneNumber,
		"action_status":        model.ActionRejected,
		"checker_action_time":  time.Now(),
		"last_modified_at":     time.Now(),
	}
	rejectedAction, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return model.CPSAction{}, err
	}
	return rejectedAction, nil
}

func (r *PermissionPersistence) UpdatePermissionGroupFromAction(ctx context.Context, action *cps_entities.CPSAction) (*cps_entities.CPSAction, error) {

	// 1. Handle different possible input types
	var actionData map[string]interface{}

	switch v := action.CurrentAction.(type) {
	case map[string]interface{}:
		actionData = v
	case bson.D:
		actionData = make(map[string]interface{})
		// Convert bson.D to map
		for _, elem := range v {
			actionData[elem.Key] = elem.Value
		}
	case bson.M:
		actionData = v
	case string:
		// Try to unmarshal JSON string
		if err := json.Unmarshal([]byte(v), &actionData); err != nil {
			r.logger.Errorf("failed to unmarshal action data: %v", err)
			return nil, errors.New("invalid action data format")
		}
	default:
		r.logger.Errorf("unexpected action data type: %T", action.CurrentAction)
		return nil, errors.New("invalid action data format")
	}

	groupName, ok := actionData["group_name"].(string)
	if !ok || groupName == "" {
		r.logger.Errorf("missing or invalid group_name")
		return nil, errors.New("group_name is required and must be a string")
	}

	// 3. Handle permission categories more robustly
	var permissionCategories []string
	if pc, ok := actionData["permission_categories"]; ok && pc != nil {
		switch v := pc.(type) {
		case []interface{}:
			permissionCategories = make([]string, 0, len(v))
			for i, item := range v {
				if s, ok := item.(string); ok {
					permissionCategories = append(permissionCategories, s)
				} else {
					r.logger.Errorf("permission_categories[%d] is not a string", i)
					return nil, fmt.Errorf("permission_categories must contain only strings")
				}
			}
		case []string:
			permissionCategories = v
		case string:
			permissionCategories = strings.Split(v, ",")
			for i := range permissionCategories {
				permissionCategories[i] = strings.TrimSpace(permissionCategories[i])
			}
		case bson.A:
			permissionCategories = make([]string, 0, len(v))
			for i, val := range v {
				switch strVal := val.(type) {
				case string:
					permissionCategories = append(permissionCategories, strVal)
				default:
					r.logger.Errorf("permission_categories[%d] is not a string, got: %T", i, val)
					return nil, fmt.Errorf("permission_categories must contain only strings")
				}
			}
		default:
			r.logger.Errorf("unexpected permission_categories type: %T", pc)
			return nil, errors.New("invalid permission_categories format")
		}

	}

	oldGroup := actionData["old_group"].(string)
	// 4. Prepare update document
	update := bson.M{
		"group_name":          groupName,
		"permission_category": permissionCategories,
		"updated_at":          time.Now(),
		"last_modified":       time.Now(),
	}

	// Include role if provided
	if role, ok := actionData["role"].(string); ok && role != "" {
		update["role"] = role
	}

	// 5. Execute update
	filter := bson.M{"group_name": oldGroup}
	_, err := r.permissionGroupsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("error updating permission group: %v", err)
		return nil, fmt.Errorf("failed to update permission group")
	}

	return action, nil
}

func (r *PermissionPersistence) UpdatePermissionGroup(groupName string, permissionCategoryLists []string) (entities.PermissionGroup, error) {
	ctx := context.Background()

	filter := bson.M{"group_name": groupName}
	update := bson.M{
		"permission_category": permissionCategoryLists,
		"updated_at":          time.Now(),
	}
	_, err := r.permissionGroupsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("Error updating permission group: %v", err)
		return entities.PermissionGroup{}, err
	}
	return entities.PermissionGroup{}, nil
}

func (r *PermissionPersistence) GetPermissionGroup(groupName string) (entities.PermissionGroup, error) {
	filter := bson.M{"group_name": groupName}
	result, err := r.permissionGroupsDal.FindOne(context.Background(), filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.PermissionGroup{}, fmt.Errorf("NOT_FOUND")
		}
		return entities.PermissionGroup{}, err
	}
	r.logger.Infof("Permission group found: %v", result)
	return *result, nil
}

func (r *PermissionPersistence) GetPermissionGroups(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.PermissionGroup], error) {
	filter := bson.M{
		// "is_deleted": false,
	}
	projection := bson.M{}

	if filterParams.Search != "" {
		filter["icon"] = bson.M{
			"$regex":   filterParams.Search,
			"$options": "i",
		}
	}

	if filterParams.Filters != "" {
		filter["group_name"] = filterParams.Filters
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	permissionGroups, err := r.permissionGroupsDal.FindAllWithPagination(context.Background(), filter, projection, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	total, err := r.permissionGroupsDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := common_util.BuildPaginationMeta(total, limit, filterParams.Page)

	return &common_util.PaginatedResponse[[]*entities.PermissionGroup]{
		Data: permissionGroups,
		Meta: meta,
	}, nil
}
