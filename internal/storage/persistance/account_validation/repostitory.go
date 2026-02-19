package accountvalidation

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AccountValidationStore struct {
	dal           dal.MongoDal[model.ValidationRule, model.ValidationRule]
	client        *mongo.Client
	collection    *mongo.Collection
	kafkaProducer kafka.ClientOrchestrationProducer
	logger        utils.Logger
}

// NewAccountValidationStore returns a ValidationRuleRepository
func NewAccountValidationStore(client *mongo.Client, cfg *config.VaultConfig, dbName string, collection string, kafkaProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.ValidationRuleRepository {
	return &AccountValidationStore{
		dal:           dal.NewMongoDal[model.ValidationRule, model.ValidationRule](client, cfg, dbName, collection),
		client:        client,
		logger:        logger,
		kafkaProducer: kafkaProducer,
		collection:    client.Database(dbName).Collection(collection),
	}
}

// GetAccountValidationByID implements ValidationRuleRepository
func (l *AccountValidationStore) FindByID(ctx context.Context, id string) (*model.ValidationRule, error) {
	l.logger.Infof("[FindByID] fetching account validation rule by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		l.logger.Errorf("[FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": objID}

	result, err := l.dal.FindOne(ctx, filter, nil)
	if err != nil {
		l.logger.Errorf("[FindByID] failed to find account validation rule: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	l.logger.Infof("[FindByID] account validation rule retrieved successfully")
	return result, nil
}

// UpdateAccountValidation implements ValidationRuleRepository
func (a *AccountValidationStore) Update(ctx context.Context, id string, rule *model.ValidationRule) error {
	a.logger.Infof("[Update] updating account validation rule for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		a.logger.Errorf("[Update] invalid object id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false}
	updateData := AccountValidationMapper(*rule)

	updateAccountValidation, err := a.dal.UpdateOne(ctx, filter, updateData)
	if err != nil {
		a.logger.Errorf("[Update] failed to update account validation rule: %v", err)
		return local_util.HandleDBError(err)
	}

	a.kafkaProducer.PublishMessage(ctx, updateAccountValidation, string(constants.ClientOrchestrationAccountValidationTopic), string(constants.ClientOrchestrationAccountValidationTopic), "update account validation rule")

	a.logger.Infof("[Update] account validation rule updated successfully")
	return nil
}

// GetAllAccountValidation implements ValidationRuleRepository
func (l *AccountValidationStore) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ValidationRule], error) {
	filter := bson.M{"is_deleted": false}

	searchKeys := bson.M{}

	// 2. Allowed filterable/searchable fields
	allowedKeys := []string{"enabled", "entity_type", "validation_for"}
	// 3. Add search (if provided)
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{
			{"service_id": searchRegex},
			{"enabled": searchRegex},
			{"validation_for": searchRegex},
			{"entity_type": searchRegex},
		} // choose your searchable field(s)
	}

	// 4. Build filter, skip, limit
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	// 5. Fetch data
	data, err := l.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		l.logger.Errorf("[FindAllWithPagination] failed to fetch account validation rules: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 6. Count total
	total, err := l.dal.TotalCount(ctx, filter)
	if err != nil {
		l.logger.Errorf("[FindAllWithPagination] failed to count account validation rules: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// 7. Build pagination metadata
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	l.logger.Infof("[FindAllWithPagination] retrieved %d account validation rules", len(data))

	return &types.PaginatedResponse[[]model.ValidationRule]{
		Data: data,
		Meta: meta,
	}, nil
}
