package event_merchant_repository

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"

	// "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/event_merchant/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type EventMerchantRepository struct {
	dal    dal.MongoDal[model.EventMerchant, model.EventMerchant]
	logger utils.Logger
}

func (m *EventMerchantRepository) Create(ctx context.Context, merchant model.EventMerchant) error {
	if merchant.ID.IsZero() {
		merchant.ID = bson.ObjectID(primitive.NewObjectID())
	}
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Time{}
	merchant.DeletedAt = time.Time{}
	_, err := m.dal.InsertOne(ctx, merchant)
	if err != nil {
		m.logger.Errorf("Failed to create event merchant: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *EventMerchantRepository) Update(ctx context.Context, id string, merchant model.EventMerchant) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	updateData := core.EventMerchantMapper(merchant)

	_, err = m.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Event merchant not found for update, id: %s", id)
			return errors.New(localization.ErrorEventMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to update event merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EventMerchantRepository) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"is_deleted": true,
		"deleted_at": time.Now(),
		"updated_at": time.Now(),
	}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Event merchant not found for deletion, id: %s", id)
			return errors.New(localization.ErrorEventMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to delete event merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EventMerchantRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("Event merchant not found for enable/disable, id: %s", id)
			return errors.New(localization.ErrorEventMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to enable/disable event merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *EventMerchantRepository) FindByID(ctx context.Context, id string) (*model.EventMerchant, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	result, err := m.dal.FindOne(context.Background(), filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Event merchant not found, id: %s", id)
			return nil, errors.New(localization.ErrorEventMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to find event merchant, id: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (s *EventMerchantRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EventMerchant], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "merchant_type", "merchant_id", "merchant_name", "email", "phone_number", "enabled", "bank_account_number"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"merchant_id": searchRegex},
			{"bank_account_number": searchRegex},
			{"merchant_name": searchRegex},
			{"merchant_id": searchRegex},
			{"phone_number": searchRegex},
			{"email": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("Failed to fetch paginated event merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("Failed to count total event merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.EventMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *EventMerchantRepository) FindOne(ctx context.Context, filter bson.M) (*model.EventMerchant, error) {
	result, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Event merchant not found")
			return nil, errors.New(localization.ErrorEventMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to find event merchant, filter: %v, error: %v", filter, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func NewEventMerchantRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.EventMerchantRepository {
	return &EventMerchantRepository{
		dal:    dal.NewMongoDal[model.EventMerchant, model.EventMerchant](client, dbName, collection),
		logger: logger,
	}
}
