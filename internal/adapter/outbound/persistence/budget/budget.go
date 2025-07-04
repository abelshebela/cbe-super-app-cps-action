package budget

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
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
	cpsAction.MakerActionTime = time.Now()
	cps, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create cps action", err)
		err = fmt.Errorf("failed to create cps action %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &cps, nil
}

func (b *BudgetPersistence) FetchIcons(ctx context.Context) ([]*entities.Icon, error) {
	icons, err := b.iconDal.FindAll(ctx, bson.M{}, bson.M{})
	if err != nil {
		return nil, err
	}

	return icons, nil
}

func (b *BudgetPersistence) UpdateIcon(ctx context.Context, id string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	cpsAction.MakerActionTime = time.Now()

	filter := bson.M{
		"_id":        id,
		"is_deleted": false,
	}
	existingIcon, err := b.iconDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		b.logger.Errorf("icon not found or db error: %v", err)
		err = fmt.Errorf("icon not found: %w", constant.ErrorDefinition{
			Code:    http.StatusNotFound,
			Message: "icon not found",
		})
		return nil, err
	}

	cpsAction.PreviousAction = map[string]interface{}{
		"icon_id":  existingIcon.ID,
		"icon_url": existingIcon.Icon,
	}

	createdAction, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create CPSAction update request: %v", err)
		err = fmt.Errorf("failed to queue icon update: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &createdAction, nil
}

func (b *BudgetPersistence) CreateColor(ctx context.Context, color string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	createdAction, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create CPSAction update request: %v", err)
		err = fmt.Errorf("failed to queue icon update: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &createdAction, nil
}

func (b *BudgetPersistence) ListAllColor(ctx context.Context) ([]*entities.Color, error) {
	colors, err := b.colorDal.FindAll(ctx, bson.M{"isDeleted": false}, bson.M{})
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (b *BudgetPersistence) GetByIDColor(ctx context.Context, id string) (*entities.Color, error) {
	colors, err := b.colorDal.FindOne(ctx, bson.M{"_id": id}, bson.M{})
	if err != nil {
		return nil, err
	}

	return colors, nil
}

func (b *BudgetPersistence) UpdateColor(ctx context.Context, color entities.CPSAction) (*entities.CPSAction, error) {
	update, err := bson.Marshal(color)
	if err != nil {
		return nil, err
	}
	var updateMap bson.M
	if err := bson.Unmarshal(update, &updateMap); err != nil {
		return nil, err
	}
	action, err := b.cpsDal.UpdateOne(ctx, bson.M{"_id": color.ID}, bson.M{"$set": updateMap})
	return &action, err
}

func (b *BudgetPersistence) CreateAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	createdAction, err := b.cpsDal.InsertOne(ctx, cpsAction)
	if err != nil {
		b.logger.Errorf("failed to create CPSAction update request: %v", err)
		err = fmt.Errorf("failed to queue icon update: %w", constant.ErrorDefinition{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		})
		return nil, err
	}

	return &createdAction, nil
}

func (b *BudgetPersistence) ApproveAction(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	action, err := b.cpsDal.FindOne(ctx, bson.M{"action_code": cpsAction.ActionCode}, bson.M{})
	if err != nil || action == nil {
		return nil, errors.New("action not found")
	}
	if action.ActionStatus != "PENDING" {
		return nil, errors.New("action already processed")
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
			return nil, fmt.Errorf("invalid currentAction format for icon: %w", err)
		}
		iconURL, ok := current["icon_url"].(string)
		if !ok || iconURL == "" {
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
			prev, err := castToBsonM(action.PreviousAction)
			if err != nil {
				return nil, fmt.Errorf("invalid previousAction format for icon update: %w", err)
			}
			iconID, ok := prev["icon_id"].(string)
			if !ok || iconID == "" {
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
		"action_status":       entities.ActionApproved,
		"checker_action_time": action.CheckerActionTime,
		"last_modified_at":    time.Now(),
	}
	res, err := b.cpsDal.UpdateOne(ctx, filter, update)
	if err != nil {
		b.logger.Errorf("failed to update CPSAction status: %v", err)
		return nil, fmt.Errorf("CPSAction status update failed: %w", err)
	}

	return &res, nil
}
