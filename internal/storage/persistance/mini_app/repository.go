package mini_app

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MiniAppStorage struct {
	dal        dal.MongoDal[model.MiniApp, model.MiniApp]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewMiniAppRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.MiniAppRepository {
	return &MiniAppStorage{
		dal:        dal.NewMongoDal[model.MiniApp, model.MiniApp](client, dbName, collection),
		client:     client,
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

func (m *MiniAppStorage) Create(ctx context.Context, miniApp *model.MiniApp) error {

	miniAppDoc := MiniAppDocumentMapper(*miniApp)

	_, err := m.dal.InsertOne(ctx, *miniAppDoc)
	if err != nil {
		m.logger.Errorf("Create MiniApp failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppStorage) Update(ctx context.Context, id string, miniApp *model.MiniApp) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}

	update := MiniAppDocumentToBsonM(*miniApp)
	if len(update) == 1 {
		m.logger.Warnf("No data provided for MiniApp update, ID: %s", id)
		return errors.New(localization.ErrorUpdateMiniAppEmptyPayload.Code)
	}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("MiniApp not found for ID: %s", id)
			return errors.New(localization.ErrorMiniAppNotFound.Code)
		}
		m.logger.Errorf("Update MiniApp failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateFields := bson.M{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}

	_, err = m.dal.UpdateOne(ctx, filter, updateFields)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("MiniApp not found for ID: %s", id)
			return errors.New(localization.ErrorMiniAppNotFound.Code)
		}
		m.logger.Errorf("Delete MiniApp failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{
		"enabled":          enable,
		"last_modified_at": time.Now(),
	}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("MiniApp not found for ID: %s", id)
			return errors.New(localization.ErrorMiniAppNotFound.Code)
		}
		m.logger.Errorf("EnableOrDisable MiniApp failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppStorage) DisableManyByMerchantIDs(ctx context.Context, merchantID string) error {
	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		m.logger.Errorf("Invalid Merchant ID format: %s, error: %v", merchantID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{
		"merchant_id": objID,
		"is_deleted":  false,
	}
	update := bson.M{
		"enabled":          false,
		"last_modified_at": time.Now(),
	}
	_, err = m.collection.UpdateMany(ctx, filter, bson.M{"$set": update})
	if err != nil {
		m.logger.Errorf("DisableManyByMerchantIDs failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	m.logger.Infof("Successfully disabled mini apps for merchant ID: %s", merchantID)
	return nil

}

func (m *MiniAppStorage) DeleteManyByMerchantIDs(ctx context.Context, merchantID string) error {
	objID, err := bson.ObjectIDFromHex(merchantID)
	if err != nil {
		m.logger.Errorf("Invalid Merchant ID format: %s, error: %v", merchantID, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{
		"merchant_id": objID,
		"is_deleted":  false,
	}
	update := bson.M{
		"is_deleted":       true,
		"deleted_at":       time.Now(),
		"last_modified_at": time.Now(),
	}
	_, err = m.collection.UpdateMany(ctx, filter, bson.M{"$set": update})
	if err != nil {
		m.logger.Errorf("DeleteManyByMerchantIDs failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}
