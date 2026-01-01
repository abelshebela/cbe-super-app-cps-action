package ecommercemerchant

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"

	// "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

type EcommerceMerchantStorage struct {
	dal        dal.MongoDal[model.EcommerceMerchant, model.EcommerceMerchant]
	client     *mongo.Client
	logger     utils.Logger
	dbName     string
	collection string
}

func NewEcommerceMerchantRepository(client *mongo.Client,cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.EcommerceMerchantRepository {
	return &EcommerceMerchantStorage{
		dal:        dal.NewMongoDal[model.EcommerceMerchant, model.EcommerceMerchant](client,cfg, dbName, collection),
		client:     client,
		logger:     logger,
		dbName:     dbName,
		collection: collection,
	}
}

func (m *EcommerceMerchantStorage) Create(ctx context.Context, merchant *model.EcommerceMerchant) (*model.EcommerceMerchant, error) {
	createdMerchant, err := m.dal.InsertOne(ctx, *merchant)
	if err != nil {
		m.logger.Errorf("Failed to create mini app merchant: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &createdMerchant, nil
}

func (m *EcommerceMerchantStorage) Update(ctx context.Context, id string, merchant *model.EcommerceMerchant) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	updateData := MiniAppMerchantMapper(*merchant)

	_, err = m.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found for update, id: %s", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to update mini app merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) Delete(ctx context.Context, id string) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	update := bson.M{
		"is_deleted":       true,
		"deleted_at":       time.Now(),
		"last_modified_at": time.Now(),
	}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found for deletion, id: %s", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to delete mini app merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (m *EcommerceMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			m.logger.Warnf("Mini app merchant not found for enable/disable, id: %s", id)
			return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to enable/disable mini app merchant, id: %s, error: %v", id, err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *EcommerceMerchantStorage) FindByID(ctx context.Context, id string) (*model.EcommerceMerchant, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	result, err := m.dal.FindOne(context.Background(), filter, bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found, id: %s", id)
			return nil, errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to find mini app merchant, id: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}

func (s *EcommerceMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EcommerceMerchant], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"search", "merchant_code", "merchant_name", "email", "phone_number", "enabled", "bank_account_number"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"merchant_id": searchRegex},
			{"bank_account_number": searchRegex},
			{"merchant_name": searchRegex},
			{"merchant_code": searchRegex},
			{"phone_number": searchRegex},
			{"email": searchRegex},
		}
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	data, err := s.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("Failed to fetch paginated mini app merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("Failed to count total mini app merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]model.EcommerceMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *EcommerceMerchantStorage) FindOne(ctx context.Context, filter bson.M) (*model.EcommerceMerchant, error) {
	result, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			m.logger.Warnf("Mini app merchant not found")
			return nil, errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
		}
		m.logger.Errorf("Failed to find mini app merchant, id: %s, error: %v", filter, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return result, nil
}
