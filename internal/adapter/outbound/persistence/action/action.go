package action

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	action_entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ActionRepository interface {
	CreateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error)
	UpdateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error)
	FetchCpsActionById(ctx context.Context, Action_Id string) (*action.CPSAction, error)
	FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (*action.CPSAction, error)
}

type ActionRepo struct {
	client *mongo.Client
	logger utils.Logger

	cpsDal dal.MongoDal[action_entity.CPSAction, action_entity.CPSAction]
}

func NewActionRepo(client *mongo.Client, dbName string, collections []string, logger utils.Logger) action.IActionRepository {
	dalCPS := dal.NewMongoDal[action_entity.CPSAction, action_entity.CPSAction](client, dbName, "cps_actions")
	return &ActionRepo{
		cpsDal: dalCPS,
		logger: logger,
		client: client,
	}
}

func (r *ActionRepo) CreateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error) {
	return r.cpsDal.InsertOne(ctx, Action)
}

func (r *ActionRepo) UpdateCpsAction(ctx context.Context, Action action.CPSAction) (action.CPSAction, error) {
	update := bson.M{
		"action_code":      Action.ActionCode,
		"maker":            Action.Maker,
		"checker":          Action.Checker,
		"department":       Action.Department,
		"action_type":      Action.ActionType,
		"request_action":   Action.RequestAction,
		"action_status":    Action.ActionStatus,
		"current_action":   Action.CurrentAction,
		"previos_action":   Action.PreviosAction,
		"rejection_reason": Action.RejectionReason,
		"last_modified_at": Action.LastModifiedAt,
	}
	return r.cpsDal.UpdateOne(ctx, bson.M{"_id": Action.ID}, update)
}

func (r *ActionRepo) FetchCpsActionById(ctx context.Context, actionId string) (*action.CPSAction, error) {
	oid, err := bson.ObjectIDFromHex(actionId)
	if err != nil {
		r.logger.Errorf("failed to convert action ID to ObjectID: %v", err)
		return nil, err
	}

	r.logger.Infof("action ID: %v", oid)
	return r.cpsDal.FindOne(ctx, bson.M{"_id": oid}, nil)
}

func (r *ActionRepo) FetchLastCpsActionByMakerID(ctx context.Context, makerId string) (*action.CPSAction, error) {
	_, err := bson.ObjectIDFromHex(makerId)
	if err != nil {
		r.logger.Errorf("failed to convert maker ID to ObjectID: %v", err)
		return nil, err
	}
	return nil, nil
}
