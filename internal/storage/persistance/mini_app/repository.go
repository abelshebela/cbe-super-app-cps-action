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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MiniAppStorage struct {
	dal    dal.MongoDal[model.MiniApp, model.MiniApp]
	client *mongo.Client

	logger utils.Logger
}

func NewMiniAppRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.MiniAppRepository {
	return &MiniAppStorage{
		dal:    dal.NewMongoDal[model.MiniApp, model.MiniApp](client, dbName, collection),
		client: client,
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

	_, err = m.dal.UpdateOne(ctx, filter, bson.M{"$set": update})
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

func (m *MiniAppStorage) FindByID(ctx context.Context, id string) (*model.MiniApp, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
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
	allowedKeys := []string{"app_type", "enabled", "is_event_mini_app", "is_three_click"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{"is_deleted": false}, allowedKeys)

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
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

	meta := local_util.BuildPaginationMeta(total, filterParam.PerPage, filterParam.Page)

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
