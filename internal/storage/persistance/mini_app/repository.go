package mini_app

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	// "cbe-super-app-cps-action/internal/constants/model"

	mini_app "cbe-super-app-cps-action/internal/constants/dto/mini_app"

	mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"

	"cbe-super-app-cps-action/internal/storage"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type MiniAppStorage struct {
	dal        dal.MongoDal[mini_model.MiniApp, mini_model.MiniApp]
	client     *mongo.Client
	collection *mongo.Collection
	logger     utils.Logger
}

func NewMiniAppRepository(client *mongo.Client,cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.MiniAppRepository {
	return &MiniAppStorage{
		dal:        dal.NewMongoDal[mini_model.MiniApp, mini_model.MiniApp](client,cfg, dbName, collection),
		client:     client,
		collection: client.Database(dbName).Collection(collection),
		logger:     logger,
	}
}

func (m *MiniAppStorage) Create(ctx context.Context, miniApp *mini_model.MiniApp) error {

	miniAppDoc := MiniAppDocumentMapper(*miniApp)

	_, err := m.dal.InsertOne(ctx, *miniAppDoc)
	if err != nil {
		m.logger.Errorf("Create MiniApp failed: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *MiniAppStorage) Update(ctx context.Context, id string, miniApp *mini_model.MiniApp) error {
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

func (m *MiniAppStorage) FindByIDWithMerchant(ctx context.Context, id string) (*mini_app.MiniAppResponse, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	pipeline := mongo.Pipeline{
		// Match mini app by ID and not deleted
		{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: objID},
			{Key: "is_deleted", Value: false},
		}}},

		// Convert merchant_id string to ObjectID
		{{Key: "$addFields", Value: bson.D{
			{Key: "merchant_id_obj", Value: bson.D{
				{Key: "$toObjectId", Value: "$merchant_id"},
			}},
		}}},

		// Lookup merchant
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "mini_app_merchant"},
			{Key: "localField", Value: "merchant_id_obj"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "merchant"},
		}}},

		// Unwind merchant array
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$merchant"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},

		// Project only required fields
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 1},
			{Key: "app_name", Value: 1},
			{Key: "app_icon", Value: 1},
			{Key: "banner_image", Value: 1},
			{Key: "commison_gl_account", Value: 1},
			{Key: "app_type", Value: 1},
			{Key: "app_view_type", Value: 1},
			{Key: "url", Value: 1},
			{Key: "stage", Value: 1},
			{Key: "product_code", Value: 1},
			{Key: "credential", Value: 1},
			{Key: "is_event_mini_app", Value: 1},
			{Key: "is_three_click", Value: 1},
			{Key: "enabled", Value: 1},
			{Key: "created_at", Value: 1},
			{Key: "last_modified_at", Value: 1},
			{Key: "merchant._id", Value: 1},
			{Key: "merchant.merchant_name", Value: 1},
		}}},
	}

	cursor, err := m.collection.Aggregate(ctx, pipeline)
	if err != nil {
		m.logger.Errorf("aggregate mini app by id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) {
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	var resp mini_app.MiniAppResponse
	if err := cursor.Decode(&resp); err != nil {
		m.logger.Errorf("decode mini app response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	return &resp, nil
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
