package fayda

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FaydaStorage struct {
	dal    dal.MongoDal[model.User, model.User]
	client *mongo.Client
	logger utils.Logger
}

func InitFaydaAccountPersistence(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.FaydaRepository {
	return &FaydaStorage{
		dal:    dal.NewMongoDal[model.User, model.User](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (f *FaydaStorage) FindByUserCode(ctx context.Context, user_code string) (*model.User, error) {
	filter := bson.M{"user_code": user_code}
	user, err := f.dal.FindOne(ctx, filter, bson.M{})

	if err != nil {

	}

	return user, nil
}

func (f *FaydaStorage) AuthorizeEnableFaydaUser(ctx context.Context, user *model.User) error {

}

func (f *FaydaStorage) AuthorizeDisableFaydaUser(ctx context.Context, user *model.User) error {

}
