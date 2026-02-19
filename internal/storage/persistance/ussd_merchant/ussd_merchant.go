package ussd_merchant

import (
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"regexp"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UssdMerchantRepository struct {
	dal             dal.MongoDal[imodel.UssdMerchant, imodel.UssdMerchant]
	dbName          string
	cfg             config.VaultConfig
	collection      string
	mongoCollection *mongo.Collection
	logger          utils.Logger
}

func NewUssdMerchant(client *mongo.Client, dbName string, collection string, cfg *config.VaultConfig, logger utils.Logger) storage.UssdMerchantRepository {
	return &UssdMerchantRepository{
		dal:             dal.NewMongoDal[imodel.UssdMerchant, imodel.UssdMerchant](client, cfg, dbName, collection),
		mongoCollection: client.Database(dbName).Collection(collection),
		cfg:             *cfg,
		logger:          logger,
	}
}

func (u *UssdMerchantRepository) Create(ctx context.Context, data imodel.UssdMerchant) error {

	data.ID = bson.NewObjectID()
	data.CreatedAt = time.Now()
	_, err := u.dal.InsertOne(ctx, data)
	if err != nil {
		u.logger.Errorf("[UssdDalCreateMerchant] Error creating ussd_merchant: %v", err)
		if err == mongo.ErrNoDocuments {
			return localization.ErrorResourceNotFound
		}
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	return nil
}
func (u *UssdMerchantRepository) Update(ctx context.Context, id string, update bson.M) error {

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		u.logger.Errorf("[UssdDalUpdateMerchant] Error parsing id: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	update["updated_at"] = time.Now()

	_, err = u.dal.UpdateOne(ctx, bson.M{"_id": objID}, update)
	if err != nil {
		u.logger.Errorf("[UssdDalUpdateMerchant] Error creating ussd_merchant: %v", err)
		if err == mongo.ErrNoDocuments {
			return localization.ErrorResourceNotFound
		}
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	return nil
}
func (u *UssdMerchantRepository) FindById(ctx context.Context, id string) (ussd_merchant_dto.UssdMerchantResponse, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		u.logger.Errorf("[UssdDalUpdateMerchant] Error parsing id: %v", err)
		return ussd_merchant_dto.UssdMerchantResponse{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	data, err := u.dal.FindOne(ctx, bson.M{"_id": objID}, nil)
	if err != nil {
		u.logger.Errorf("[UssdDalFindMerchant] Error fetching ussd_merchant: %v", err)
		if err == mongo.ErrNoDocuments {
			return ussd_merchant_dto.UssdMerchantResponse{}, localization.ErrorResourceNotFound
		}
		return ussd_merchant_dto.UssdMerchantResponse{}, errors.New(localization.ErrorUnhandledServer.Code)
	}

	res := ResponseMapper(*data)

	return res, nil
}

func (u *UssdMerchantRepository) Find(ctx context.Context, filter bson.M) (ussd_merchant_dto.UssdMerchantResponse, error) {

	data, err := u.dal.FindOne(ctx, filter, nil)
	if err != nil {
		u.logger.Errorf("[UssdDalFindMerchant] Error fetching ussd_merchant: %v", err)
		if err == mongo.ErrNoDocuments {
			return ussd_merchant_dto.UssdMerchantResponse{}, localization.ErrorResourceNotFound
		}
		return ussd_merchant_dto.UssdMerchantResponse{}, errors.New(localization.ErrorUnhandledServer.Code)
	}

	res := ResponseMapper(*data)

	return res, nil
}

func (u *UssdMerchantRepository) FindByOr(ctx context.Context, phone, email, account_number string) (imodel.UssdMerchant, error) {

	// Build conditions dynamically, only for non-empty parameters
	conditions := []bson.M{}

	if phone != "" {
		conditions = append(conditions, bson.M{
			"phone_number": bson.M{"$regex": "^" + regexp.QuoteMeta(phone) + "$", "$options": "i"},
		})
	}

	if account_number != "" {
		conditions = append(conditions, bson.M{
			"account_number": account_number,
		})
	}

	if email != "" {
		conditions = append(conditions, bson.M{
			"email": bson.M{"$regex": "^" + regexp.QuoteMeta(email) + "$", "$options": "i"},
		})
	}

	filter := bson.M{"$or": conditions}
	data, err := u.dal.FindOne(ctx, filter, nil)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Infof("[FindByOr] no merchant found matching the criteria")
			return imodel.UssdMerchant{}, nil
		}
		u.logger.Errorf("[FindByOr] failed to find merchant: %v", err)
		return imodel.UssdMerchant{}, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return *data, nil
}

func (u *UssdMerchantRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse], error) {
	filter := bson.M{}
	searchKeys := bson.M{}

	allowedKeys := []string{"enabled", "merchant_code", "name", "settlement_method", "phone_number", "service", "account_number"}

	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["enabled"] = searchRegex
		searchKeys["merchant_code"] = searchRegex
		searchKeys["name"] = searchRegex
		searchKeys["settlement_method"] = searchRegex
		searchKeys["phone_number"] = searchRegex
		searchKeys["service"] = searchRegex
		searchKeys["account_number"] = searchRegex
	}

	if filterParam.Page == 0 {
		filterParam.Page = 1
	}
	if filterParam.PerPage == 0 {
		filterParam.PerPage = 10
	}

	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}},
		bson.D{{Key: "$project", Value: bson.M{
			"_id":               1,
			"merchant_code":     1,
			"name":              1,
			"settlement_method": 1,
			"phone_number":      1,
			"service":           1,
			"email":             1,
			"enabled":           1,
			"credential":        1,
			"account_number":    1,
			"logo":              1,
			"updated_at":        1,
			"created_at":        1,
		}}},

		// Safe string → ObjectId
		bson.D{{Key: "$addFields", Value: bson.M{
			"service_obj_id": bson.M{
				"$cond": bson.M{
					"if": bson.M{
						"$regexMatch": bson.M{
							"input": "$service",
							"regex": "^[a-fA-F0-9]{24}$",
						},
					},
					"then": bson.M{"$toObjectId": "$service"},
					"else": nil,
				},
			},
		}}},

		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "services",
			"localField":   "service_obj_id",
			"foreignField": "_id",
			"as":           "service_info",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$project", Value: bson.M{
					"_id":          1,
					"service_name": 1,
				}}},
			},
		}}},

		bson.D{{Key: "$unwind", Value: bson.M{
			"path":                       "$service_info",
			"preserveNullAndEmptyArrays": true,
		}}},

		bson.D{{Key: "$addFields", Value: bson.M{
			"service_name": "$service_info.service_name",
		}}},

		bson.D{{Key: "$project", Value: bson.M{
			"service_info":   0,
			"service_obj_id": 0,
		}}},

		bson.D{{Key: "$facet", Value: bson.M{
			"data": []bson.D{
				{{Key: "$skip", Value: skip}},
				{{Key: "$limit", Value: limit}},
			},
			"total": []bson.D{
				{{Key: "$count", Value: "count"}},
			},
		}}},
	}

	u.logger.Infof("[FindAllWithPagination] fetching Ussd Merchant with pagination")
	cursor, err := u.mongoCollection.Aggregate(ctx, pipeline)
	if err != nil {
		u.logger.Errorf("[FindAllWithPagination] failed to execute aggregation pipeline: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}
	defer cursor.Close(ctx)

	var results []struct {
		Data  []ussd_merchant_dto.UssdMerchantResponse `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}

	if err := cursor.All(ctx, &results); err != nil {
		u.logger.Errorf("[FindAllWithPagination] failed to decode aggregation results: %v", err)
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{}, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if len(results) == 0 {
		return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{
			Data: []ussd_merchant_dto.UssdMerchantResponse{},
			Meta: local_util.BuildPaginationMeta(0, filterParam.Page, filterParam.PerPage),
		}, nil
	}

	var total int64
	if len(results[0].Total) > 0 {
		total = results[0].Total[0].Count
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	u.logger.Infof("[FindAllWithPagination] retrieved %d Ussd Merchant", len(results[0].Data))
	return types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse]{
		Data: results[0].Data,
		Meta: meta,
	}, nil

}
