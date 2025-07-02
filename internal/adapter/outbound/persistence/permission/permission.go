package permission

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	repository "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"time"
)

type PermissionPersistence struct {
	permissionGroupsDal   dal.MongoDal[entities.PermissionGroup, entities.PermissionGroup]
	permissionCategoryDal dal.MongoDal[entities.PermissionCategory, entities.PermissionCategory]
	cpsdal                dal.MongoDal[entities.CPSAction, entities.CPSAction]
	timeout               time.Duration
	logger                utils.Logger
}

var _ repository.PermissionGroupRepository = (*PermissionPersistence)(nil)
var _ repository.CPSActionRepository = (*PermissionPersistence)(nil)
var _ repository.PermissionCategoryRepository = (*PermissionPersistence)(nil)

func InitPermission(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *PermissionPersistence {
	permissionGroupsDal := dal.NewMongoDal[entities.PermissionGroup, entities.PermissionGroup](client, dbName, "permission_groups")
	permissionCategoryDal := dal.NewMongoDal[entities.PermissionCategory, entities.PermissionCategory](client, dbName, "permission_categories")
	cpsdal := dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, "cps_actions")
	return &PermissionPersistence{
		permissionGroupsDal:   permissionGroupsDal,
		permissionCategoryDal: permissionCategoryDal,
		cpsdal:                cpsdal,
		timeout:               timeout,
		logger:                logger,
	}
}

func (r *PermissionPersistence) CheckPendingRequest(userCode string, status entities.ActionStatus, action entities.RequestAction) error {
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

	filter := bson.M{"group_name": groupName, "is_deleted": false}
	result, err := r.permissionGroupsDal.FindOne(ctx, filter, bson.M{})
	if err != nil || result == nil {
		r.logger.Errorf("CheckPermissionGroupExists failed:", err)
		return false
	}
	return true
}

func (r *PermissionPersistence) ValidatePermissionCategories(ids []string) ([]string, error) {
	ctx := context.Background()

	filter := bson.M{}
	categories, err := r.permissionCategoryDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		r.logger.Errorf("failed to fetch permission categories: ", err)
		return nil, fmt.Errorf("failed to fetch permission categories")
	}

	var allowedPermissionCategories []bson.ObjectID

	for _, permissionID := range ids {
		objID, err := bson.ObjectIDFromHex(permissionID)
		if err != nil {
			r.logger.Errorf("invalid permission category ID format: ", err)
			continue
		}

		for _, cat := range categories {
			if cat.ID == objID {
				allowedPermissionCategories = append(allowedPermissionCategories, objID)
				break
			}
		}
	}

	// if len(allowedPermissionCategories) == 0 {
	//     return nil, fmt.Errorf("no valid permission categories found")
	// }

	var allowedStrings []string
	for _, obj := range allowedPermissionCategories {
		allowedStrings = append(allowedStrings, obj.Hex())
	}
	return allowedStrings, nil

}

func (r *PermissionPersistence) CreatePermissionGroup(group entities.CPSAction) error {
	ctx := context.Background()
	_, err := r.cpsdal.InsertOne(ctx, group)
	return err
}

func (r *PermissionPersistence) ValidateActionRequest(actionCode, department string) (entities.CPSAction, error) {
	ctx := context.Background()

	filter := bson.M{
		"action_code":   actionCode,
		"action_status": "PENDING",
	}

	action, err := r.cpsdal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entities.CPSAction{}, errors.New("no pending request found for this action")
		}
		return entities.CPSAction{}, err
	}

	if action == nil {
		return entities.CPSAction{}, errors.New("no pending request found for this action")
	}

	if action.Department != department {
		return entities.CPSAction{}, errors.New("you are not allowed to approve this request")
	}

	return *action, nil
}

func (r *PermissionPersistence) CreatePermissionGroupFromAction(action entities.CPSAction) error {
	ctx := context.Background()

	// Step 1: Normalize action.CurrentAction into map[string]interface{}
	var actionData map[string]interface{}
	bytes, err := json.Marshal(action.CurrentAction)
	if err != nil {
		return fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}
	if err := json.Unmarshal(bytes, &actionData); err != nil {
		return fmt.Errorf("failed to unmarshal CurrentAction: %v", err)
	}

	// Step 2: Extract fields
	groupName, ok := actionData["group_name"].(string)
	if !ok {
		return errors.New("invalid group_name")
	}

	permissionCategoriesIface, ok := actionData["permission_categories"].([]interface{})
	if !ok {
		return errors.New("invalid permission_categories")
	}

	permissionCategories := make([]entities.PermissionCategory, len(permissionCategoriesIface))
	for i, v := range permissionCategoriesIface {
		idStr, ok := v.(string)
		if !ok {
			return errors.New("permission_categories contains non-string value")
		}
		objID, err := bson.ObjectIDFromHex(idStr)
		if err != nil {
			return fmt.Errorf("invalid permission category id: %s", idStr)
		}
		permissionCategories[i] = entities.PermissionCategory{
			ID: objID,
		}
	}

	role, ok := actionData["role"].(string)
	if !ok {
		return errors.New("invalid role")
	}

	realm, ok := actionData["realm"].(string)
	if !ok {
		return errors.New("invalid realm")
	}

	// Step 3: Save to DB
	newPermissionGroup := entities.PermissionGroup{
		GroupName:          groupName,
		PermissionCategory: permissionCategories,
		Role:               role,
		Realm:              realm,
		CreatedAt:          time.Now(),
	}

	_, err = r.permissionGroupsDal.InsertOne(ctx, newPermissionGroup)
	if err != nil {
		r.logger.Errorf("Error creating permission group from action: %v", err)
		return err
	}

	return nil
}

func (r *PermissionPersistence) ApproveActionRequest(actionCode string, action entities.CPSAction) error {
	ctx := context.Background()
	filter := bson.M{"action_code": actionCode}
	update := bson.M{
		"checker_name":         action.CheckerName,
		"checker_id":           action.CheckerID,
		"checker_phone_number": action.CheckerPhoneNumber,
		"action_status":        entities.ActionApproved,
		"checker_action_time":  time.Now(),
	}
	_, err := r.cpsdal.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *PermissionPersistence) UpdatePermissionGroupFromAction(action entities.CPSAction) error {
	ctx := context.Background()

	actionData, ok := action.CurrentAction.(map[string]interface{})
	if !ok {
		return errors.New("invalid action data format")
	}

	groupID, ok := actionData["permission_group_id"].(string)
	if !ok {
		return errors.New("invalid permission_group_id")
	}

	objectID, err := bson.ObjectIDFromHex(groupID)
	if err != nil {
		return errors.New("invalid objectID format")
	}

	groupName, ok := actionData["group_name"].(string)
	if !ok {
		return errors.New("invalid group_name")
	}

	permissionCategoriesIface, ok := actionData["permission_categories"].([]interface{})
	if !ok {
		return errors.New("invalid permission_categories")
	}
	permissionCategories := make([]string, len(permissionCategoriesIface))
	for i, v := range permissionCategoriesIface {
		permissionCategories[i], ok = v.(string)
		if !ok {
			return errors.New("permission_categories contains non-string value")
		}
	}

	role, ok := actionData["role"].(string)
	if !ok {
		return errors.New("invalid role")
	}

	update := bson.M{
		"$set": bson.M{
			"groupName":          groupName,
			"permissionCategory": permissionCategories,
			"role":               role,
			"updatedAt":          time.Now(),
		},
	}

	filter := bson.M{"_id": objectID}
	_, err = r.permissionGroupsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Errorf("Error updating permission group from action: %v", err)
		return err
	}

	return nil
}
