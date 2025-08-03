package hq

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/model"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type hqRepository struct {
	hqDal  dal.MongoDal[model.HQ, model.HQ]
	logger utils.Logger
}

func NewHQRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.HQRepository {
	return &hqRepository{
		hqDal:  dal.NewMongoDal[model.HQ, model.HQ](client, dbName, collection),
		logger: logger,
	}
}

func (h *hqRepository) FindOne(ctx context.Context, filter bson.M) (*model.HQ, error) {
	if filter == nil {
		filter = bson.M{}
	}
	projection := HQProjection()

	h.logger.Infof("Finding one HQ with filter: %v and projection: %v", filter, projection)

	hq, err := h.hqDal.FindAll(ctx, filter, projection)
	if err != nil {
		h.logger.Errorf("Failed to find HQ: %v", err)
		return nil, err
	}
	if len(hq) == 0 {
		h.logger.Warnf("No HQ found for filter: %v", filter)
		return nil, mongo.ErrNoDocuments
	}
	h.logger.Infof("Successfully found HQ: %+v", hq[0])
	return hq[0], nil
}
