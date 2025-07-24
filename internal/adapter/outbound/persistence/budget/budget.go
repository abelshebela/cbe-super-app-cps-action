package budget

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	contexts "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"encoding/json"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type BudgetPersistence struct {
	iconDal  dal.MongoDal[entities.Icon, entities.Icon]
	colorDal dal.MongoDal[entities.Color, entities.Color]
	cpsDal   dal.MongoDal[entities.CPSAction, entities.CPSAction]
	logger   utils.Logger
}

var (
	_ budget.Repository = (*BudgetPersistence)(nil)
)

func InitBudget(client *mongo.Client, dbName string, collections []string, logger utils.Logger) *BudgetPersistence {
	return &BudgetPersistence{
		iconDal:  dal.NewMongoDal[entities.Icon, entities.Icon](client, dbName, collections[0]),
		colorDal: dal.NewMongoDal[entities.Color, entities.Color](client, dbName, collections[1]),
		cpsDal:   dal.NewMongoDal[entities.CPSAction, entities.CPSAction](client, dbName, collections[2]),
		logger:   logger,
	}
}

func (b *BudgetPersistence) CreateIconAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {

	makerData := contexts.ExtractContext(ctx)
	filter := bson.M{
		"maker_id":      makerData.UserCode,
		"department":    makerData.Department,
		"action_status": "PENDING",
	}
	_, err := b.cpsDal.FindOne(ctx, filter, bson.M{})
	if err == nil {
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	// if existingAction != nil {
	// 	return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	// }

	cpsAction.MakerActionTime = time.Now()
	cpsAction.CreatedAt = time.Now()
	objID := bson.NewObjectID()
	cpsAction.ID = objID
	cps, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		return nil, fmt.Errorf(error_codes.FailedToCreateAction)
	}

	return &cps, nil
}

func (b *BudgetPersistence) FetchIcons(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Icon], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParams.Search != "" {
		filter["icon"] = bson.M{
			"$regex":   filterParams.Search,
			"$options": "i",
		}
	}

	if filterParams.Filters != "" {
		switch filterParams.Filters {
		case "enabled":
			filter["enabled"] = true
		case "disabled":
			filter["enabled"] = false
		}
	}

	skip := (filterParams.Page - 1) * filterParams.PerPage
	limit := filterParams.PerPage

	icons, err := b.iconDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	total, err := b.iconDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := common_util.BuildPaginationMeta(total, limit, filterParams.Page)

	return &common_util.PaginatedResponse[[]*entities.Icon]{
		Data: icons,
		Meta: meta,
	}, nil
}

func (b *BudgetPersistence) UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	makerData := contexts.ExtractContext(ctx)
	filter := bson.M{
		"maker_id":       makerData.UserCode,
		"department":     makerData.Department,
		"action_status":  "PENDING",
		"request_action": cpsAction.RequestAction,
	}
	_, err := b.cpsDal.FindOne(ctx, filter, nil)

	if err == nil {
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	cpsAction.MakerActionTime = time.Now()
	cpsAction.ID = bson.NewObjectID()
	cpsAction.LastModifiedAt = time.Now()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID")
	}

	filter = bson.M{
		"_id":        objectID,
		"is_deleted": false,
	}
	existingIcon, err := b.iconDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("icon not found: %s", id)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		b.logger.Errorf("icon not found or db error: %v", err)
		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}

	cpsAction.PreviousAction = map[string]interface{}{
		"icon_id":  existingIcon.ID,
		"icon_url": existingIcon.Icon,
	}
	// cpsAction.CurrentAction = map[string]interface{}{
	// 	"icon_id": existingIcon.ID,
	// }
	currentMap, ok := cpsAction.CurrentAction.(map[string]interface{})
	if !ok {
		// If it's not a map yet, initialize it
		currentMap = make(map[string]interface{})
	}
	currentMap["icon_id"] = existingIcon.ID
	cpsAction.CurrentAction = currentMap

	createdAction, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create CPSAction update request: %v", err)
		err = fmt.Errorf("UNHANDLED_SERVER_ERROR")
		return nil, err
	}

	return &createdAction, nil
}

func (b *BudgetPersistence) CreateColor(ctx context.Context, color string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {

	makerData := contexts.ExtractContext(ctx)
	filter := bson.M{
		"maker_id":      makerData.UserCode,
		"department":    makerData.Department,
		"action_status": "PENDING",
	}
	_, err := b.cpsDal.FindOne(ctx, filter, nil)

	if err == nil {
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	filter = bson.M{
		"color": color,
	}
	existing, err := b.colorDal.FindOne(ctx, filter, nil)
	if err != nil {
		if err.Error() != "mongo: no documents in result" {
			return nil, err
		}
	}
	if existing != nil {
		return nil, fmt.Errorf("COLOR_ALREADY_EXISTED")
	}
	cpsAction.CreatedAt = time.Now()
	cpsAction.MakerActionTime = time.Now()
	cpsAction.ID = bson.NewObjectID()
	cpsAction.LastModifiedAt = time.Now()
	createdAction, err := b.cpsDal.InsertOne(ctx, cpsAction)

	if err != nil {
		b.logger.Errorf("failed to create CPSAction update request: %v", err)
		return nil, fmt.Errorf("GENERAL_DB_INSERT_FAILED")
	}

	return &createdAction, nil
}

func (b *BudgetPersistence) ListAllColor(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Color], error) {
	filter := bson.M{
		"is_deleted": false,
	}

	if filterParams.Search != "" {
		filter["color"] = bson.M{
			"$regex":   filterParams.Search,
			"$options": "i",
		}
	}

	if filterParams.Filters != "" {
		switch filterParams.Filters {
		case "enabled":
			filter["enabled"] = true
		case "disabled":
			filter["enabled"] = false
		}
	}

	page := filterParams.Page
	limit := filterParams.PerPage
	skip := (page - 1) * limit

	colors, err := b.colorDal.FindAllWithPagination(ctx, filter, bson.M{}, int64(skip), int64(limit))
	if err != nil {
		return nil, err
	}

	totalDocs, err := b.colorDal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	meta := common_util.BuildPaginationMeta(totalDocs, skip, limit)

	return &common_util.PaginatedResponse[[]*entities.Color]{
		Data: colors,
		Meta: meta,
	}, nil
}

func (b *BudgetPersistence) GetByIDColor(ctx context.Context, id string) (*entities.Color, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid object ID")
	}

	colors, err := b.colorDal.FindOne(ctx, bson.M{"_id": objectId}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("color not found: %s", id)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		b.logger.Errorf("color not found or db error: %v", err)
		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}

	return colors, nil
}

func (b *BudgetPersistence) CheckColorExist(ctx context.Context, color string) (bool, error) {
	colors, err := b.colorDal.FindOne(ctx, bson.M{"color": color}, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		b.logger.Errorf("failed to check color existence: %v", err)
		return false, fmt.Errorf("failed to check color existence: %w", err)
	}

	return colors != nil, nil
}

func (b *BudgetPersistence) UpdateColor(ctx context.Context, color entities.CPSAction) (*entities.CPSAction, error) {

	makerData := contexts.ExtractContext(ctx)
	filter := bson.M{
		"maker_id":      makerData.UserCode,
		"department":    makerData.Department,
		"action_status": "PENDING",
	}
	_, err := b.cpsDal.FindOne(ctx, filter, nil)

	if err == nil {
		return nil, fmt.Errorf("PENDING_ACTION_EXISTS")
	}

	update, err := bson.Marshal(color)
	if err != nil {
		return nil, err
	}
	var updateMap bson.M
	if err := bson.Unmarshal(update, &updateMap); err != nil {
		return nil, err
	}

	action, err := b.cpsDal.UpdateOne(ctx, bson.M{"_id": color.ID}, bson.M{"$set": updateMap})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("action not found: %s", color.ID.Hex())
			return nil, fmt.Errorf("NOT_FOUND")
		}
		b.logger.Errorf("failed to update CPSAction: %v", err)
		return nil, fmt.Errorf("GENERAL_DB_UPDATE_FAILED")
	}
	return &action, nil
}

func (b *BudgetPersistence) CreateAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction.ID = bson.NewObjectID()
	createdAction, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create CPSAction update request: %v", err)
		return nil, err
	}

	return &createdAction, nil
}

func (b *BudgetPersistence) ApproveAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	action, err := b.cpsDal.FindOne(ctx, bson.M{"action_code": cpsAction.ActionCode}, bson.M{})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			b.logger.Errorf("action not found: %s", cpsAction.ActionCode)
			return nil, fmt.Errorf("NOT_FOUND")
		}
		b.logger.Errorf("action not found: %s, error: %v", cpsAction.ActionCode, err)
		return nil, fmt.Errorf("GENERAL_DB_QUERY_FAILED")
	}

	if action.ActionStatus != "PENDING" {
		b.logger.Errorf("action already processed: %s", action.ActionCode)
		return nil, fmt.Errorf("action already processed")
	}

	action.CheckerActionTime = time.Now()
	action.ActionStatus = entities.ActionApproved

	castToBsonM := func(input interface{}) (bson.M, error) {
		raw, err := bson.Marshal(input)
		if err != nil {
			return nil, err
		}
		var out bson.M
		err = bson.Unmarshal(raw, &out)
		return out, err
	}

	switch action.RequestAction {
	case "BUDGET_ICON":
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			b.logger.Errorf("invalid currentAction format for icon approval: %v", err)
			return nil, fmt.Errorf("invalid currentAction format for icon approval")
		}
		iconURL, ok := current["icon_url"].(string)
		if !ok || iconURL == "" {
			b.logger.Errorf("missing icon_url in currentAction")
			return nil, fmt.Errorf("missing icon_url in currentAction")
		}

		switch action.ActionType {
		case entities.ActionCreate:
			icon := entities.Icon{
				Icon:         iconURL,
				Enabled:      true,
				IsDeleted:    false,
				CreatedAt:    time.Now(),
				LastModified: time.Now(),
			}
			_, err := b.iconDal.InsertOne(ctx, icon)
			if err != nil {
				b.logger.Errorf("failed to approve icon creation: %v", err)
				return nil, fmt.Errorf("icon creation failed: %w", err)
			}

		case entities.ActionUpdate:
			fmt.Println("----------KKKKKKKKKKKK-----------")

			prev, err := castToBsonM(action.PreviousAction)
			if err != nil {
				b.logger.Errorf("invalid previousAction format for icon update: %v", err)
				return nil, fmt.Errorf("invalid previousAction format for icon update: %w", err)
			}
			iconID, ok := prev["icon_id"].(string)
			fmt.Println("---------------------")
			fmt.Println(iconID)
			if !ok {
				b.logger.Errorf("missing icon_id in previousAction")
				return nil, fmt.Errorf("missing icon_id in previousAction")
			}

			filter := bson.M{"_id": iconID, "is_deleted": false}
			update := bson.M{
				"icon":          iconURL,
				"last_modified": time.Now(),
			}
			_, err = b.iconDal.UpdateOne(ctx, filter, update)
			if err != nil {
				b.logger.Errorf("failed to approve icon update: %v", err)
				return nil, fmt.Errorf("icon update failed: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown action type for icon")
		}

	case "BUDGET_COLOR":
		current, err := castToBsonM(action.CurrentAction)
		if err != nil {
			return nil, fmt.Errorf("invalid currentAction format for color: %w", err)
		}
		colorName, ok := current["color"].(string)
		if !ok || colorName == "" {
			return nil, fmt.Errorf("missing color name in currentAction")
		}

		switch action.ActionType {
		case entities.ActionCreate:
			color := entities.Color{
				Color:     colorName,
				Enabled:   true,
				IsDeleted: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			_, err := b.colorDal.InsertOne(ctx, color)
			if err != nil {
				b.logger.Errorf("failed to approve color creation: %v", err)
				return nil, fmt.Errorf("color creation failed: %w", err)
			}

		case entities.ActionUpdate:
			prev, err := castToBsonM(action.PreviousAction)
			if err != nil {
				return nil, fmt.Errorf("invalid previousAction format for color update: %w", err)
			}
			colorID, ok := prev["color_id"].(string)
			if !ok || colorID == "" {
				return nil, fmt.Errorf("missing color_id in previousAction")
			}

			filter := bson.M{"_id": colorID, "is_deleted": false}
			update := bson.M{
				"color":      colorName,
				"updated_at": time.Now(),
			}
			_, err = b.colorDal.UpdateOne(ctx, filter, update)
			if err != nil {
				b.logger.Errorf("failed to approve color update: %v", err)
				return nil, fmt.Errorf("color update failed: %w", err)
			}
		default:
			return nil, fmt.Errorf("unknown action type for color")
		}

	default:
		return nil, fmt.Errorf("unknown request action type")
	}

	filter := bson.M{"_id": action.ID}
	update := bson.M{
		"action_status":        entities.ActionApproved,
		"checker_action_time":  action.CheckerActionTime,
		"checker_id":           cpsAction.CheckerID,
		"checker_name":         cpsAction.CheckerName,
		"checker_phone_number": cpsAction.CheckerPhoneNumber,
		"last_modified_at":     time.Now(),
	}
	res, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to update CPSAction status: %v", err)
		return nil, fmt.Errorf("CPSAction status update failed: %w", err)
	}

	return &res, nil
}

func (b *BudgetPersistence) Authorize(ctx context.Context, cpsAction *cps_entities.CPSAction) (*cps_entities.CPSAction, error) {
	if cpsAction == nil {
		return nil, errors.New("cpsAction is required")
	}

	castToBsonM := func(input interface{}) (bson.M, error) {
		raw, err := bson.Marshal(input)
		if err != nil {
			return nil, err
		}
		var out bson.M
		err = bson.Unmarshal(raw, &out)
		return out, err
	}
	switch cpsAction.RequestAction {
	case "BUDGET_CREATE_ICON":

		var actionData map[string]interface{}
		if cpsAction.CurrentAction == nil {
			return nil, fmt.Errorf("missing action data")
		}
		if err := UnmarshalMap(cpsAction.CurrentAction, &actionData); err != nil {
			return nil, fmt.Errorf("INVALID_ACTION_FORMAT")
		}
		iconURL, ok := actionData["icon_url"].(string)
		if !ok || iconURL == "" {
			return nil, fmt.Errorf("MISSING_ICON")
		}
		icon := entities.Icon{
			Icon:         iconURL,
			Enabled:      true,
			IsDeleted:    false,
			CreatedAt:    time.Now(),
			LastModified: time.Now(),
		}
		_, err := b.iconDal.InsertOne(ctx, icon)
		if err != nil {
			b.logger.Errorf("failed to approve icon creation: %v", err)
			return nil, fmt.Errorf("icon creation failed: %w", err)
		}

	case "BUDGET_UPDATE_ICON":

		current, err := castToBsonM(cpsAction.CurrentAction)
		if err != nil {
			b.logger.Errorf("invalid currentAction format for icon approval: %v", err)
			return nil, fmt.Errorf("invalid currentAction format for icon approval")

		}
		iconURL, ok := current["icon_url"].(string)
		if !ok || iconURL == "" {
			b.logger.Errorf("missing icon_url in currentAction")
			return nil, fmt.Errorf("missing icon_url in currentAction")
		}
		prev, err := castToBsonM(cpsAction.PreviousAction)
		if err != nil {
			b.logger.Errorf("invalid previousAction format for icon update: %v", err)
			return nil, fmt.Errorf("invalid previousAction format for icon update: %w", err)
		}
		iconID := prev["icon_id"]

		filter := bson.M{"_id": iconID, "is_deleted": false}
		update := bson.M{
			"icon":          iconURL,
			"last_modified": time.Now(),
		}
		_, err = b.iconDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to approve icon update: %v", err)
			return nil, fmt.Errorf("icon update failed: %w", err)
		}

	case "BUDGET_DELETE_ICON":
		// For icon deletion, mark the icon as deleted
		var actionData map[string]interface{}
		if cpsAction.CurrentAction == nil {
			return nil, fmt.Errorf("missing action data")
		}
		if err := UnmarshalMap(cpsAction.CurrentAction, &actionData); err != nil {
			return nil, fmt.Errorf("invalid action data format: %w", err)
		}
		iconID, ok := actionData["icon_id"].(string)
		if !ok || iconID == "" {
			return nil, fmt.Errorf("missing icon_id in action data")
		}
		filter := bson.M{"_id": iconID, "is_deleted": false}
		update := bson.M{
			"is_deleted":    true,
			"last_modified": time.Now(),
		}
		_, err := b.iconDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to authorize icon deletion: %v", err)
			return nil, fmt.Errorf("failed to authorize icon deletion: %w", err)
		}

	case "BUDGET_CREATE_COLOR":
		// For color creation, mark the color as enabled (or perform any other necessary business logic)
		var actionData map[string]interface{}
		if cpsAction.CurrentAction == nil {
			return nil, fmt.Errorf("missing action data")
		}
		if err := UnmarshalMap(cpsAction.CurrentAction, &actionData); err != nil {
			return nil, fmt.Errorf("invalid action data format: %w", err)
		}

		colorName, ok := actionData["color"].(string)
		if !ok || colorName == "" {
			return nil, fmt.Errorf("missing color in action data")
		}

		color := entities.Color{
			Color:     colorName,
			Enabled:   true,
			IsDeleted: false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := b.colorDal.InsertOne(ctx, color)
		if err != nil {
			b.logger.Errorf("failed to approve color creation: %v", err)
			return nil, fmt.Errorf("color creation failed: %w", err)
		}

	case "BUDGET_UPDATE_COLOR":
		// For color update, update the color fields as needed
		current, err := castToBsonM(cpsAction.CurrentAction)
		if err != nil {
			return nil, fmt.Errorf("invalid currentAction format for color: %w", err)
		}
		colorName, ok := current["color"].(string)
		if !ok || colorName == "" {
			return nil, fmt.Errorf("missing color name in currentAction")
		}

		prev, err := castToBsonM(cpsAction.PreviousAction)
		if err != nil {
			return nil, fmt.Errorf("invalid previousAction format for color update: %w", err)
		}
		colorID, ok := prev["color_id"]

		filter := bson.M{"_id": colorID, "is_deleted": false}
		update := bson.M{
			"color":      colorName,
			"updated_at": time.Now(),
		}
		_, err = b.colorDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to approve color update: %v", err)
			return nil, fmt.Errorf("FAILED_COLOR_UPDATE")
		}

	case "BUDGET_DELETE_COLOR":
		// For color deletion, mark the color as deleted
		var actionData map[string]interface{}
		if cpsAction.CurrentAction == nil {
			return nil, fmt.Errorf("MISSING_ACTION_DATA")
		}
		if err := UnmarshalMap(cpsAction.CurrentAction, &actionData); err != nil {
			return nil, fmt.Errorf("INVALID_ACTION_FORMAT")
		}
		colorID, ok := actionData["color_id"].(string)
		if !ok || colorID == "" {
			return nil, fmt.Errorf("MISSING_COLOR_ID")
		}
		filter := bson.M{"_id": colorID, "is_deleted": false}
		update := bson.M{
			"is_deleted":    true,
			"last_modified": time.Now(),
		}
		_, err := b.colorDal.UpdateOne(ctx, filter, update)
		if err != nil {
			b.logger.Errorf("failed to authorize color deletion: %v", err)
			return nil, fmt.Errorf("FAILED_TO_AUTHORIZE")
		}

	default:
		b.logger.Errorf("failed to authorize action")

		return nil, fmt.Errorf("FAILED_TO_AUTHORIZE")

	}

	return cpsAction, nil
}

// UnmarshalMap converts a map[string]interface{} to a struct.
func UnmarshalMap(input interface{}, output interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, output)
}
