package logistics_merchant_repository

import (
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	local_model "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage/persistance/logistics_merchant/core"
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
	client     *mongo.Client
	dbName     string
	collection string
	dal        dal.MongoDal[local_model.LogisticsMerchant, local_model.LogisticsMerchant]
	logger     utils.Logger
}

func NewLogisticsMerchantRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.LogisticsMerchantRepository {
	return &LogisticsMerchantRepository{
		client:     client,
		dbName:     dbName,
		collection: collection,
		dal:        dal.NewMongoDal[local_model.LogisticsMerchant, local_model.LogisticsMerchant](client, cfg, dbName, collection),
		logger:     logger,
	}
}

func (m *LogisticsMerchantRepository) Create(ctx context.Context, merchant local_model.LogisticsMerchant) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	if merchant.ID == "" {
		merchant.ID = bson.ObjectID(primitive.NewObjectID()).Hex()
	}
	merchant.CreatedAt = time.Now()
	merchant.UpdatedAt = time.Time{}
	merchant.DeletedAt = time.Time{}

	coll := m.client.Database(m.dbName).Collection(m.collection)
	res, err := coll.InsertOne(ctx, merchant)
	if err != nil {
		log.Errorf("Failed to create logistics merchant: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	merchant.ID = res.InsertedID.(string)
	// types.SetId(ctx, merchant.ID.Hex())

	return nil
}

func (m *LogisticsMerchantRepository) Update(ctx context.Context, id string, merchant local_model.LogisticsMerchant) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Invalid ID format: %s, error: %v", id, err)
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
	log := local_util.LoggerFromCtx(ctx, m.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Invalid ID format: %s, error: %v", id, err)
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

func (m *LogisticsMerchantRepository) EnableOrDisable(ctx context.Context, ids []string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	for _, id := range ids {
		objID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			log.Errorf("Invalid ID format: %s, error: %v", id, err)
			return errors.New(localization.ErrorInvalidID.Code)
		}

		filter := bson.M{"_id": objID, "is_deleted": false}
		update := bson.M{"enabled": enable, "updated_at": time.Now()}

		_, err = m.dal.UpdateOne(ctx, filter, update)
		if err != nil {
			return local_util.HandleDBError(err)
		}
	}
	return nil
}

func (m *LogisticsMerchantRepository) FindByID(ctx context.Context, id string) (*local_model.LogisticsMerchant, error) {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Invalid ID format: %s, error: %v", id, err)
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
	log := local_util.LoggerFromCtx(ctx, s.logger)

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
		log.Errorf("Failed to fetch paginated logistics merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("Failed to count total logistics merchants: %v", err)
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
