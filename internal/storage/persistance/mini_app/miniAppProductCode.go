package mini_app

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	// "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/handlers/middleware"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type miniAppProductCode struct {
	logger                shared_utils.Logger
	miniAppProductCodeDal dal.MongoDal[model.MiniAppProductCode, model.MiniAppProductCode]
	client                *mongo.Client
}

func NewMiniAppProdutCodeRepository(logger shared_utils.Logger, client *mongo.Client, dbName, collectionName string) storage.MiniAppProductCodeRepository {
	return &miniAppProductCode{
		logger:                logger,
		miniAppProductCodeDal: dal.NewMongoDal[model.MiniAppProductCode, model.MiniAppProductCode](client, dbName, collectionName),
		client:                client,
	}
}

func (a *miniAppProductCode) Create(ctx context.Context, productCode *model.MiniAppProductCode) error {
	_, err := a.miniAppProductCodeDal.InsertOne(ctx, *productCode)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *miniAppProductCode) Update(ctx context.Context, productCode *model.MiniAppProductCode, id string) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf(invalidCategoryID, err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	filter := bson.M{"_id": objId}
	update := buildProductCodeUpdate(*productCode)
	_, err = a.miniAppProductCodeDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error updating mini app product code in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (a *miniAppProductCode) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		a.logger.Errorf("Invalid mini app product code ID:", err)
		return middleware.NewBadRequestError(invalidCategoryID, err)
	}
	filter := bson.M{"_id": objId}
	update := bson.M{"is_enabled": enable, "updated_at": time.Now()}
	_, err = a.miniAppProductCodeDal.UpdateOne(ctx, filter, update)
	if err != nil {
		a.logger.Errorf("Error enabling or disabling mini app product code in database:", err)
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func buildProductCodeUpdate(updateFields model.MiniAppProductCode) bson.M {
	update := bson.M{}
	if updateFields.Name != "" {
		update["name"] = updateFields.Name
	}
	if updateFields.IsEnabled {
		update["is_enabled"] = updateFields.IsEnabled
	}

	if updateFields.ChargeCode != "" {
		update["charge_code"] = updateFields.ChargeCode
	}
	if updateFields.CommisionCode != "" {
		update["commission_code"] = updateFields.CommisionCode
	}
	if updateFields.ServiceCode != "" {
		update["service_code"] = updateFields.ServiceCode
	}

	update["updated_at"] = time.Now()

	return update
}
