package customersegmentaion

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerStorage struct {
	dal        dal.MongoDal[imodel.CustomerSegmentation, imodel.CustomerSegmentation]
	client     *mongo.Client
	dbName     string
	collection string
	logger     utils.Logger
}

func NewCustomerSegmentationRepository(client *mongo.Client, dbName, collection string, logger utils.Logger) storage.CustomerSegmentationRepository {
	return &customerStorage{
		dal:        dal.NewMongoDal[imodel.CustomerSegmentation, imodel.CustomerSegmentation](client, dbName, collection),
		client:     client,
		dbName:     dbName,
		collection: collection,
		logger:     logger,
	}
}

func (r *customerStorage) Create(ctx context.Context, seg *imodel.CustomerSegmentation) error {
	_, err := r.dal.InsertOne(ctx, *seg)
	if err != nil {
		r.logger.Errorf("Unable to create customer segmentation with error: %s", err)
		return err
	}
	return nil
}
