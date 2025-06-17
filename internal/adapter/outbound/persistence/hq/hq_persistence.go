package hq

import (
	"context"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/action"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/hq/models"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HQPersistence struct {
	hqDal        dal.MongoDal[models.HQ, models.HQ]
	cpsActionDal dal.MongoDal[action.CPSAction, action.CPSAction]
	timeout      time.Duration
	logger       utils.Logger
}

func NewHQPersistence(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *HQPersistence {
	hqDal := dal.NewMongoDal[models.HQ, models.HQ](client, dbName, "hq")
	cpsActionDal := dal.NewMongoDal[action.CPSAction, action.CPSAction](client, dbName, "cps_actions")
	return &HQPersistence{
		hqDal:        hqDal,
		cpsActionDal: cpsActionDal,
		timeout:      timeout,
		logger:       logger,
	}
}

func (p *HQPersistence) GetHQByID(ctx context.Context, id string) (models.HQ, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	result, err := p.hqDal.FindOne(ctx, bson.M{"_id": id}, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			p.logger.Errorf("HQ not found", "id", id)
			return models.HQ{}, err
		}
		p.logger.Errorf("failed to fetch HQ: %v", err)
		return models.HQ{}, err
	}
	return *result, nil
}

func (p *HQPersistence) UpdateHQ(ctx context.Context, id string, update models.HQ) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	updateDoc := bson.M{
		"block_time":          update.BlockTime,
		"block_time_status":   update.BlockTimeStatus,
		"archive_time":        update.ArchiveTime,
		"archive_time_status": update.ArchiveTimeStatus,
		"last_modified_at":    time.Now(),
	}

	_, err := p.hqDal.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	if err != nil {
		p.logger.Errorf("failed to update HQ: %v", err)
		return err
	}
	return nil
}

func (p *HQPersistence) FetchPendingActionsByUniqueID(ctx context.Context, uniqueID string) ([]action.CPSAction, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	filter := bson.M{
		"department":    uniqueID,
		"action_status": action.ActionPending,
	}

	results, err := p.cpsActionDal.FindAll(ctx, filter, bson.M{})
	if err != nil {
		p.logger.Errorf("failed to fetch pending actions: %v", err)
		return nil, err
	}

	actions := make([]action.CPSAction, 0, len(results))
	for _, result := range results {
		if result != nil {
			actions = append(actions, *result)
		}
	}
	return actions, nil
}
