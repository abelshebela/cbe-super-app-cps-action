package cpsaction

import (
	"context"
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	infra_mongo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/mongo"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	repo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/repository"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cpsActionStore struct {
	MongoCPSAction *infra_mongo.MongoDal[model.CPSAction, model.CPSAction]
	logger         utils.Logger
}

func NewOutBoundStore(client *mongo.Client, dbName string, collectionName string, logger utils.Logger) repo.CPSActionRepository {
	MongoCPSAction := infra_mongo.NewMongoDal[model.CPSAction, model.CPSAction](client, dbName, collectionName)
	return &cpsActionStore{
		MongoCPSAction: MongoCPSAction,
		logger:         logger,
	}
}

func (o *cpsActionStore) CreateCPSAction(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	return nil, nil
}
func (o *cpsActionStore) UpdateCPSAction(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	return nil, nil
}

func (o *cpsActionStore) CPSActionExists(ctx context.Context, uniqueID string) (bool, error) {
	if uniqueID == "" {
		o.logger.Warnf("CPSActionExists: uniqueID is empty")
		return false, fmt.Errorf(common_util.InvalidID)
	}

	filter := bson.M{"unique_id": uniqueID}
	projection := bson.M{"_id": 1}

	_, err := o.MongoCPSAction.FindOne(ctx, filter, projection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		o.logger.Errorf("CPSActionExists query failed for unique_id=%s: %v", uniqueID, err)
		return false, fmt.Errorf("FAILED_TO_QUERY_CPS_ACTION: %w", err)
	}

	return true, nil
}

func (o *cpsActionStore) ApproveCPSAction(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	return nil, nil
}

func (o *cpsActionStore) RejectCPSAction(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	return nil, nil
}

func (o *cpsActionStore) GetCPSActionsByDepartment(ctx context.Context, department string) (*common_util.PaginatedResponse[[]*entity.CPSAction], error) {
	return nil, nil
}

func (o *cpsActionStore) GetCPSActionByID(ctx context.Context, id string) (*entity.CPSAction, error) {
	return nil, nil
}

func (o *cpsActionStore) GetCPSActionByUniqueID(ctx context.Context, uniqueID string) (*entity.CPSAction, error) {
	return nil, nil
}
