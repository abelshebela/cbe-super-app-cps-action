package mini_app

import (
	"context"
	"errors"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MiniAppRepositoryImpl struct {
	dal    dal.MongoDal[model.MiniApp, model.MiniApp]
	client *mongo.Client
	logger utils.Logger
}

func NewMiniAppRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.MiniAppRepository {
	return &MiniAppRepositoryImpl{
		dal:    dal.NewMongoDal[model.MiniApp, model.MiniApp](client, dbName, collection),
		client: client,
		logger: logger,
	}
}

func (r *MiniAppRepositoryImpl) Create(ctx context.Context, miniApp *model.MiniApp) error {
	_, err := r.dal.InsertOne(ctx, *miniApp)
	return err
}

func (r *MiniAppRepositoryImpl) Update(ctx context.Context, miniApp *model.MiniApp) error {
	if miniApp.ID.IsZero() {
		return mongo.ErrNilDocument
	}
	filter := bson.M{"_id": miniApp.ID}
	update := bson.M{"$set": miniApp}
	_, err := r.dal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MiniAppRepositoryImpl) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	err = r.dal.DeleteOne(ctx, filter)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return err
}

func (r *MiniAppRepositoryImpl) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	_, err = r.dal.UpdateOne(ctx, filter, update)
	return err
}

func (r *MiniAppRepositoryImpl) FindByID(ctx context.Context, id string) (*model.MiniApp, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	miniApp, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return miniApp, nil
}

func (r *MiniAppRepositoryImpl) FindByCode(ctx context.Context, code string) (*model.MiniApp, error) {
	filter := bson.M{"code": code}
	miniApp, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, err
	}
	return miniApp, nil
}

func (r *MiniAppRepositoryImpl) FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniApp], error) {
	filter := bson.M{
		"is_deleted": false,
		"department": department,
	}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["name"] = searchRegex
	}

	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
	limit := int64(filterParam.PerPage)

	data, err := r.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		return nil, err
	}

	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, err
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]*model.MiniApp]{
		Data: data,
		Meta: meta,
	}, nil
}
