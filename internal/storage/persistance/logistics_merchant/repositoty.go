package logistics_merchant_repository

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/logistics_merchant/core"
	"context"
	"errors"

	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LogisticsMerchantRepository struct {
	dal    dal.MongoDal[local_model.LogisticsMerchant, local_model.LogisticsMerchant]
	logger utils.Logger
}

func (m *LogisticsMerchantRepository) Create(ctx context.Context, merchant local_model.LogisticsMerchant) error {
	if merchant.ID.IsZero() {
		merchant.ID = bson.ObjectID(primitive.NewObjectID())
	}
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Time{}
	merchant.DeletedAt = time.Time{}
	_, err := m.dal.InsertOne(ctx, merchant)
	if err != nil {
		m.logger.Errorf("Failed to create logistics merchant: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	return nil
}

func (m *LogisticsMerchantRepository) Update(ctx context.Context, id string, merchant local_model.LogisticsMerchant) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	updateData := core.LogisticsMerchantMapper(merchant)

	_, err = m.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		return local_util.HandleDBError(err)
	}

	return nil
}

func (m *LogisticsMerchantRepository) Delete(ctx context.Context, id string) error {
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
		return local_util.HandleDBError(err)
	}

	return nil
}

func (m *LogisticsMerchantRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "updated_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (m *LogisticsMerchantRepository) FindByID(ctx context.Context, id string) (*local_model.LogisticsMerchant, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		m.logger.Errorf("Invalid ID format: %s, error: %v", id, err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	result, err := m.dal.FindOne(context.Background(), filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func (s *LogisticsMerchantRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.LogisticsMerchant], error) {
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

	data, err := s.dal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		s.logger.Errorf("Failed to fetch paginated logistics merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		s.logger.Errorf("Failed to count total logistics merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]local_model.LogisticsMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *LogisticsMerchantRepository) FindOne(ctx context.Context, filter bson.M) (*local_model.LogisticsMerchant, error) {
	result, err := m.dal.FindOne(ctx, filter, nil)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	return result, nil
}

func NewLogisticsMerchantRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.LogisticsMerchantRepository {
	return &LogisticsMerchantRepository{
		dal:    dal.NewMongoDal[local_model.LogisticsMerchant, local_model.LogisticsMerchant](client, cfg, dbName, collection),
		logger: logger,
	}
}
