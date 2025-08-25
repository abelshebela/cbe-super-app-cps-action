package validation_rule

// import (
// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	"cbe-super-app-cps-action/internal/constants/model"
// 	"cbe-super-app-cps-action/internal/constants/types"
// 	"cbe-super-app-cps-action/internal/storage"
// 	"context"
// 	"errors"

// 	local_util "cbe-super-app-cps-action/pkgs/utils"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
// 	"go.mongodb.org/mongo-driver/v2/bson"
// 	"go.mongodb.org/mongo-driver/v2/mongo"
// )

// type ValidationRuleStorage struct {
// 	dal    dal.MongoDal[model.ValidationRule, model.ValidationRule]
// 	client *mongo.Client
// 	logger utils.Logger
// }

// func NewValidationRuleRepository(client *mongo.Client, dbName string, collection string, logger utils.Logger) storage.ValidationRuleRepository {
// 	return &ValidationRuleStorage{
// 		dal:    dal.NewMongoDal[model.ValidationRule, model.ValidationRule](client, dbName, collection),
// 		client: client,
// 		logger: logger,
// 	}
// }

// func (v *ValidationRuleStorage) Create(ctx context.Context, rule *model.ValidationRule) error {
// 	_, err := v.dal.InsertOne(ctx, *rule)
// 	if err != nil {
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	return nil
// }

// func (v *ValidationRuleStorage) Update(ctx context.Context, id string, rule *model.ValidationRule) error {
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	filter := bson.M{"_id": objID, "is_deleted": false}
// 	updateData := ValidationRuleMapper(*rule)

// 	_, err = v.dal.UpdateOne(ctx, filter, updateData)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			return errors.New(localization.ErrorFileNotFound.Code)
// 		}
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	return nil
// }

// func (v *ValidationRuleStorage) Delete(ctx context.Context, id string) error {
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	filter := bson.M{"_id": objID, "is_deleted": false}
// 	return v.dal.DeleteOne(ctx, filter)
// }

// func (v *ValidationRuleStorage) FindByID(ctx context.Context, id string) (*model.ValidationRule, error) {
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, errors.New(localization.ErrorUnexpectedError.Code)
// 	}
// 	filter := bson.M{"_id": objID, "is_deleted": false}

// 	result, err := v.dal.FindOne(ctx, filter, nil)

// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }

// func (v *ValidationRuleStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error) {
// 	filter := bson.M{
// 		"is_deleted": false,
// 	}

// 	if filterParam.Search != "" {
// 		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
// 		filter["$or"] = []bson.M{
// 			{"entity_type": searchRegex},
// 			{"validation_for": searchRegex},
// 		}
// 	}

// 	skip := int64((filterParam.Page - 1) * filterParam.PerPage)
// 	limit := int64(filterParam.PerPage)

// 	data, err := v.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
// 	if err != nil {
// 		return nil, err
// 	}

// 	total, err := v.dal.TotalCount(ctx, filter)
// 	if err != nil {
// 		return nil, err
// 	}

// 	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

// 	return &types.PaginatedResponse[[]*model.ValidationRule]{
// 		Data: data,
// 		Meta: meta,
// 	}, nil
// }
