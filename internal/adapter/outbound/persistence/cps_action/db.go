package cpsaction

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/common/mappers"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/port/outbound"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/action"
	error_codes "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cpsActionStore struct {
	MongoDalCPSAction *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
	logger            utils.Logger
}

func NewOutBoundStore(client *mongo.Client, dbName string, collectionNames string, logger utils.Logger) outbound.CPSOutboundInfra {
	mongoDalCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionNames)
	return &cpsActionStore{
		MongoDalCPSAction: mongoDalCPSAction,
		logger:            logger,
	}
}

func (o *cpsActionStore) CreateCPSAction(ctx context.Context, action domain.CPSAction) (*domain.CPSAction, error) {
	action.ActionCode = utils.RandomGenerator(20)
	action.CreatedAt = time.Now()
	action.LastModifiedAt = time.Now()
	action.MakerActionTime = time.Now()
	modelAction := mappers.DomainToModelCPSAction(action)
	data, err := o.MongoDalCPSAction.InsertOne(ctx, *modelAction)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBQueryFailed)
	}
	return mappers.ModelToDomainCPSAction(data), nil
}
func (o *cpsActionStore) UpdateCPSAction(ctx context.Context, action domain.CPSAction) (*domain.CPSAction, error) {
	modelAction := mappers.DomainToModelCPSAction(action)
	filter := bson.M{
		"action_code": modelAction.ActionCode,
	}
	update := bson.M{
		"maker_id":             modelAction.MakerID,
		"maker_name":           modelAction.MakerName,
		"maker_phone_number":   modelAction.MakerPhoneNumber,
		"checker_id":           modelAction.CheckerID,
		"checker_name":         modelAction.CheckerName,
		"checker_phone_number": modelAction.CheckerPhoneNumber,
		"unique_id":            modelAction.UniqueId,
		"department":           modelAction.Department,
		"rejection_reason":     modelAction.RejectionReason,
		"previous_action":      modelAction.PreviosAction,
		"current_action":       modelAction.CurrentAction,
		"action_status":        modelAction.ActionStatus,
		"action_type":          modelAction.ActionType,
		"request_action":       modelAction.RequestAction,
		"last_modified_at":     modelAction.LastModifiedAt,
	}
	updatedAction, err := o.MongoDalCPSAction.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf(error_codes.GeneralDBUpdateFailed)
	}

	return mappers.ModelToDomainCPSAction(updatedAction), nil
}

func (o *cpsActionStore) CPSActionExists(ctx context.Context, uniqueID string) (bool, error) {
	if uniqueID == "" {
		o.logger.Warnf("CPSActionExists: uniqueID is empty")
		return false, fmt.Errorf(error_codes.InvalidID)
	}

	filter := bson.M{"unique_id": uniqueID}
	projection := bson.M{"_id": 1}

	_, err := o.MongoDalCPSAction.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		o.logger.Errorf("CPSActionExists query failed for unique_id=%s: %v", uniqueID, err)
		return false, fmt.Errorf("FAILED_TO_QUERY_CPS_ACTION: %w", err)
	}

	return true, nil
}
