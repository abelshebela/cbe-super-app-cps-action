package mini_app

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"
	// "fmt"

	mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"

	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	local_model "cbe-super-app-cps-action/internal/constants/model"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MiniAppMerchantStorage struct {
	dal        dal.MongoDal[mini_model.MiniAppMerchant, mini_model.MiniAppMerchant]
	client     *mongo.Client
	logger     utils.Logger
	dbName     string
	collection string
}

func NewMiniAppMerchantRepository(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, logger utils.Logger) storage.MiniAppMerchant {
	return &MiniAppMerchantStorage{
		dal:        dal.NewMongoDal[mini_model.MiniAppMerchant, mini_model.MiniAppMerchant](client, cfg, dbName, collection),
		client:     client,
		logger:     logger,
		dbName:     dbName,
		collection: collection,
	}
}

func (m *MiniAppMerchantStorage) Create(ctx context.Context, merchant *local_model.MiniAppMerchant) (*local_model.MiniAppMerchant, error) {
	// log := local_util.LoggerFromCtx(ctx, m.logger)

	// if merchant.ID.IsZero() {
	// 	merchant.ID = bson.ObjectID(primitive.NewObjectID())
	// }
	// createdMerchant, err := m.dal.InsertOne(ctx, *merchant)
	// if err != nil {
	// 	log.Errorf("Failed to create mini app merchant: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	return &local_model.MiniAppMerchant{}, nil
}

func (m *MiniAppMerchantStorage) Update(ctx context.Context, id string, merchant *local_model.MiniAppMerchant) error {
	// log := local_util.LoggerFromCtx(ctx, m.logger)

	// objID, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	log.Errorf("Invalid ID format: %s, error: %v", id, err)
	// 	return errors.New(localization.ErrorInvalidID.Code)
	// }

	// filter := bson.M{"_id": objID}
	// updateData := MiniAppMerchantMapper(*merchant)

	// _, err = m.dal.UpdateOne(ctx, filter, updateData)
	// if err != nil {
	// 	return local_util.HandleDBError(err)
	// }

	return nil
}

func (m *MiniAppMerchantStorage) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Invalid ID format: %s, error: %v", id, err)
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
		return local_util.HandleDBError(err)
	}

	return nil
}

func (m *MiniAppMerchantStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, m.logger)

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("Invalid ID format: %s, error: %v", id, err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID, "is_deleted": false}
	update := bson.M{"enabled": enable, "last_modified_at": time.Now()}

	_, err = m.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		return local_util.HandleDBError(err)
	}
	return nil
}

func (m *MiniAppMerchantStorage) FindByID(ctx context.Context, id string) (*local_model.MiniAppMerchant, error) {
	// log := local_util.LoggerFromCtx(ctx, m.logger)

	// fmt.Println("=============I ahvwbfkjsBDFKBFRKDBSSFSHF", id)
	// objID, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	log.Errorf("Invalid ID format: %s, error: %v", id, err)
	// 	return nil, errors.New(localization.ErrorInvalidID.Code)
	// }

	// filter := bson.M{"_id": objID, "is_deleted": false}
	// result, err := m.dal.FindOne(ctx, filter, bson.M{})
	// if err != nil {
	// 	return nil, local_util.HandleDBError(err)
	// }
	// return result, nil

	return &local_model.MiniAppMerchant{},nil
}

func (s *MiniAppMerchantStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]mini_model.MiniAppMerchant], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

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
		log.Errorf("Failed to fetch paginated mini app merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("Failed to count total mini app merchants: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]mini_model.MiniAppMerchant]{
		Data: data,
		Meta: meta,
	}, nil
}

func (m *MiniAppMerchantStorage) FindOne(ctx context.Context, filter bson.M) (*local_model.MiniAppMerchant, error) {
	// result, err := m.dal.FindOne(ctx, filter, nil)
	// if err != nil {
	// 	return nil, local_util.HandleDBError(err)
	// }
	// return result, nil
	return &local_model.MiniAppMerchant{},nil
}
