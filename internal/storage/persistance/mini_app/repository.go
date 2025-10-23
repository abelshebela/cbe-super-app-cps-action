package mini_app

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"cbe-super-app-cps-action/internal/constants/dto/mini_app"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MiniAppStorage struct {
	dal    dal.MongoDal[model.MiniApp, model.MiniApp]
	client *mongo.Client
	collection *mongo.Collection
	logger utils.Logger
}

func NewMiniAppRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.MiniAppRepository {
	return &MiniAppStorage{
		dal:    dal.NewMongoDal[model.MiniApp, model.MiniApp](client, dbName, collection),
		client: client,
		collection: client.Database(dbName).Collection(collection),
		logger: logger,
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

	_, err = m.dal.UpdateOne(ctx, filter,update)
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

// cascading operations are handled in service layer

func (m *MiniAppStorage) FindByID(ctx context.Context, id string) (*model.MiniApp, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}

	doc, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("MiniApp not found for ID: %s,error ", id, err)
			return nil, errors.New(localization.ErrorMiniAppNotFound.Code)
		}
		m.logger.Errorf("FindByID MiniApp failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return doc, nil
}

func (m *MiniAppStorage) Find(ctx context.Context, name string) (*model.MiniApp, error) {
	if name == "" {
		m.logger.Warnf("Find called with empty name")
		return nil, errors.New(localization.ErrorMiniAppNameRequired.Code)
	}

	filter := bson.M{
		"app_name":   name,
		"is_deleted": false,
	}

	doc, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("No MiniApp found with name: %s", name)
			return nil, nil
		}
		m.logger.Errorf("Find MiniApp failed for name %s: %v", name, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return doc, nil
}

func (m *MiniAppStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniApp], error) {
	allowedKeys := []string{"app_type", "enabled", "is_event_mini_app", "is_three_click", "app_name"}
	searchKeys := bson.M{}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"app_name": searchRegex},
			{"commison_gl_account": searchRegex},
			{"app_type.uat": searchRegex},
			{"app_type.production": searchRegex},
			{"app_type.test": searchRegex},
			{"app_type.dev": searchRegex},
			{"merchant_id": searchRegex},
			{"product_code.product_code": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	docs, err := m.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		m.logger.Errorf("FindAllWithPagination MiniApp failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := m.dal.TotalCount(ctx, filter)
	if err != nil {
		m.logger.Errorf("Count MiniApp failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	m.logger.Infof("FindAllWithPagination returning %d MiniApps, total: %d", len(docs), total)

	return &types.PaginatedResponse[[]*model.MiniApp]{
		Data: docs,
		Meta: meta,
	}, nil
}

func (p *MiniAppStorage) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	p.logger.Debugf("Starting MongoDB session for transaction")

	session, err := p.client.StartSession()
	if err != nil {
		p.logger.Errorf("failed to start MongoDB session: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(txCtx context.Context) error {
		p.logger.Debugf("Starting MongoDB transaction")

		if err := session.StartTransaction(); err != nil {
			p.logger.Errorf("failed to start transaction: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}

		err := fn(txCtx)
		if err != nil {
			p.logger.Errorf("transaction logic failed: %v", err)
			if abortErr := session.AbortTransaction(txCtx); abortErr != nil {
				p.logger.Errorf("failed to abort transaction: %v", abortErr)
			} else {
				p.logger.Debugf("Transaction aborted successfully")
			}
			return err
		}

		if err := session.CommitTransaction(txCtx); err != nil {
			p.logger.Errorf("failed to commit transaction: %v", err)
			return errors.New(localization.ErrorUnhandledServer.Code)
		}

		p.logger.Debugf("Transaction committed successfully")
		return nil
	})
}



func (m *MiniAppStorage) FindByIDWithMerchant(ctx context.Context, id string) (*miniappdto.MiniAppResponse, error) {
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

	var resp miniappdto.MiniAppResponse
	if err := cursor.Decode(&resp); err != nil {
		m.logger.Errorf("decode mini app response: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}

	return &resp, nil
}

