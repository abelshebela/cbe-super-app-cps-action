package hq

import (
	"context"
	"time"

	models "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HQPersistence struct {
	hqDal   dal.MongoDal[models.HQ, models.HQ]
	timeout time.Duration
	logger  utils.Logger
}

func NewHQPersistence(client *mongo.Client, dbName string, timeout time.Duration, logger utils.Logger) *HQPersistence {
	hqDal := dal.NewMongoDal[models.HQ, models.HQ](client, dbName, "hq")
	return &HQPersistence{
		hqDal:   hqDal,
		timeout: timeout,
		logger:  logger,
	}
}

func (p *HQPersistence) GetHQByID(ctx context.Context, id string) (hq.HQ, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if id == "" {
		p.logger.Errorf("invalid HQ ID: empty")
		return hq.HQ{}, mongo.ErrNoDocuments
	}

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		p.logger.Errorf("invalid HQ ObjectID: %v", err)
		return hq.HQ{}, err
	}

	filter := bson.M{"_id": objID}

	result, err := p.hqDal.FindOne(ctx, filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			p.logger.Errorf("HQ not found: id=%s", id)
			return hq.HQ{}, err
		}
		p.logger.Errorf("failed to fetch HQ: %v", err)
		return hq.HQ{}, err
	}
	if result == nil {
		p.logger.Errorf("HQ not found: id=%s", id)
		return hq.HQ{}, mongo.ErrNoDocuments
	}
	return modelToDomainHQ(*result), nil
}

func modelToDomainHQ(m models.HQ) hq.HQ {
	return hq.HQ{
		ID:           m.ID.Hex(),
		Name:         m.Name,
		BlockTime:    m.BlockTime,
		ArchiveTime:  m.ArchiveTime,
		CreatedAt:    m.CreatedAt,
		LastModified: m.LastModifiedAt,
	}
}

func (p *HQPersistence) UpdateHQ(ctx context.Context, id string, update hq.HQ) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	updateDoc := bson.M{
		"block_time":       update.BlockTime,
		"archive_time":     update.ArchiveTime,
		"last_modified_at": time.Now(),
	}
	_, err := p.hqDal.UpdateOne(ctx, bson.M{"_id": id}, updateDoc)
	if err != nil {
		p.logger.Errorf("failed to update HQ: %v", err)
		return err
	}
	return nil
}
