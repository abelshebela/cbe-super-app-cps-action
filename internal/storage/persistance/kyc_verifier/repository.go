package kyc_verifier

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type KYCVerifierStorage struct {
	dal           dal.MongoDal[model.CustomerKYC, model.CustomerKYC]
	client        *mongo.Client
	logger        utils.Logger
	coll          *mongo.Collection
	kafkaProducer kafka.ClientOrchestrationProducer
}

// FindByIDPopulated returns a populated response with user details
func (s *KYCVerifierStorage) FindByIDPopulated(ctx context.Context, id string) (*dto.KYCVerifierResponse, error) {
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"_id": idObj, "is_deleted": false}}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"user_id_obj": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$eq": []any{bson.M{"$type": "$user_id"}, "string"}},
					"then": bson.M{"$toObjectId": "$user_id"},
					"else": "$user_id",
				},
			},
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "members",
			"localField":   "user_id_obj",
			"foreignField": "_id",
			"as":           "user_doc",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$project", Value: bson.M{
					"full_name":       1,
					"phone_number":    1,
					"customer_number": 1,
					"user_code":       1,
				}}},
			},
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$user_doc", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$project", Value: bson.M{
			"id":                       "$_id",
			"user_full_name":           "$user_doc.full_name",
			"phone_number":             "$user_doc.phone_number",
			"customer_number":          "$user_doc.customer_number",
			"user_code":                "$user_doc.user_code",
			"kyc_data":                 "$kyc_data",
			"kyc_status":               "$kyc_status",
			"kyc_reject_reason":        "$kyc_reject_reason",
			"kyc_reject_reason_failed": "$kyc_reject_reason_failed",
			"kyc_approved":             "$kyc_approved",
			"kyc_activity_by":          "$kyc_activity_by",
			"kyc_level":                "$kyc_level",
			"created_at":               "$created_at",
		}}},
	}
	cur, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("aggregate kyc by id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cur.Close(ctx)
	if !cur.Next(ctx) {
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}
	var resp dto.KYCVerifierResponse
	if err := cur.Decode(&resp); err != nil {
		s.logger.Errorf("decode kyc populated by id: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	return &resp, nil
}

// FindAllWithPaginationPopulated returns populated list with pagination
func (s *KYCVerifierStorage) FindAllWithPaginationPopulated(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*dto.KYCVerifierResponse], error) {
	searchKeys := bson.M{}
	allowedKeys := []string{"kyc_status", "kyc_level", "kyc_approved", "enabled"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{{"kyc_status": searchRegex}}
	}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	filter["is_deleted"] = false

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: filter}},
		bson.D{{Key: "$addFields", Value: bson.M{
			"user_id_obj": bson.M{
				"$cond": bson.M{
					"if":   bson.M{"$eq": []any{bson.M{"$type": "$user_id"}, "string"}},
					"then": bson.M{"$toObjectId": "$user_id"},
					"else": "$user_id",
				},
			},
		}}},
		bson.D{{Key: "$lookup", Value: bson.M{
			"from":         "members",
			"localField":   "user_id_obj",
			"foreignField": "_id",
			"as":           "user_doc",
			"pipeline": mongo.Pipeline{
				bson.D{{Key: "$project", Value: bson.M{
					"full_name":       1,
					"phone_number":    1,
					"customer_number": 1,
					"user_code":       1,
				}}},
			},
		}}},
		bson.D{{Key: "$unwind", Value: bson.M{"path": "$user_doc", "preserveNullAndEmptyArrays": true}}},
		bson.D{{Key: "$project", Value: bson.M{
			"id":                       "$_id",
			"user_full_name":           "$user_doc.full_name",
			"phone_number":             "$user_doc.phone_number",
			"customer_number":          "$user_doc.customer_number",
			"user_code":                "$user_doc.user_code",
			"kyc_data":                 "$kyc_data",
			"kyc_status":               "$kyc_status",
			"kyc_reject_reason":        "$kyc_reject_reason",
			"kyc_reject_reason_failed": "$kyc_reject_reason_failed",
			"kyc_approved":             "$kyc_approved",
			"kyc_activity_by":          "$kyc_activity_by",
			"kyc_level":                "$kyc_level",
			"created_at":               "$created_at",
		}}},
		bson.D{{Key: "$facet", Value: bson.M{
			"data":  []bson.D{{{Key: "$skip", Value: skip}}, {{Key: "$limit", Value: limit}}},
			"total": []bson.D{{{Key: "$count", Value: "count"}}},
		}}},
	}

	cur, err := s.coll.Aggregate(ctx, pipeline)
	if err != nil {
		s.logger.Errorf("aggregate kyc list: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	defer cur.Close(ctx)

	var result []struct {
		Data  []*dto.KYCVerifierResponse `bson:"data"`
		Total []struct {
			Count int64 `bson:"count"`
		} `bson:"total"`
	}
	if err := cur.All(ctx, &result); err != nil {
		s.logger.Errorf("decode kyc list: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	var data []*dto.KYCVerifierResponse
	var total int64
	if len(result) > 0 {
		data = result[0].Data
		if len(result[0].Total) > 0 {
			total = result[0].Total[0].Count
		}
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*dto.KYCVerifierResponse]{
		Data: data,
		Meta: meta,
	}, nil
}

func NewKYCVerifierRepository(client *mongo.Client, dbName string, collection string, clientOrchestrationProducer kafka.ClientOrchestrationProducer, logger utils.Logger) storage.KYCVerifierRepository {
	return &KYCVerifierStorage{
		dal:           dal.NewMongoDal[model.CustomerKYC, model.CustomerKYC](client, dbName, collection),
		client:        client,
		logger:        logger,
		coll:          client.Database(dbName).Collection(collection),
		kafkaProducer: clientOrchestrationProducer,
	}
}

func (s *KYCVerifierStorage) Update(ctx context.Context, id string, kyc *model.CustomerKYC) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false, "kyc_status": "PENDING"}
	data, err := local_util.JsonUnmarshal[bson.M](kyc)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	update := data
	updatedKycVerifier, err := s.dal.UpdateOne(ctx, filter, *update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(ctx, updatedKycVerifier, string(constants.ClientOrchestrationKycTopic), string(constants.ClientOrchestrationKycTopic), "update kyc-verifier")

	return nil
}

func (s *KYCVerifierStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New(localization.ErrorUnexpectedError.Code)
	}
	filter := bson.M{"_id": objID, "is_deleted": false, "kyc_status": "PENDING"}
	update := bson.M{"$set": bson.M{"enabled": enable}}
	updatedKycVerifier, err := s.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New(localization.ErrorFileNotFound.Code)
		}
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	s.kafkaProducer.PublishMessage(ctx, updatedKycVerifier, string(constants.ClientOrchestrationKycTopic), string(constants.ClientOrchestrationKycTopic), "enable/disable kyc-verifier")

	return nil
}

func (s *KYCVerifierStorage) FindByID(ctx context.Context, id string) (*model.CustomerKYC, error) {
	idObj, ok := local_util.StringToObjectID(id)
	if !ok {
		s.logger.Errorf("Invalid ObjectID for fetch by id: %s", id)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}
	filter := bson.M{"_id": idObj, "is_deleted": false, "kyc_status": "PENDING"}
	result, err := s.dal.FindOne(ctx, filter, nil)
	if err != nil {
		s.logger.Errorf("Error finding kyc verifier: %v", err)
		code, _ := local_util.HandleMongoError(err)
		return nil, errors.New(code)
	}
	return result, nil
}

func (s *KYCVerifierStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CustomerKYC], error) {
	filter := bson.M{"is_deleted": false, "kyc_status": "PENDING"}
	searchKeys := bson.M{}
	allowedKeys := []string{"kyc_status", "kyc_level", "kyc_approved", "enabled"}
	if filterParam.Search != "" {
		searchRegex := bson.M{"$regex": filterParam.Search, "$options": "i"}
		searchKeys["$or"] = []bson.M{{"kyc_status": searchRegex}}
	}
	projection := bson.M{}
	filter, skip, limit := lib.FilterBuilder(filterParam, searchKeys, allowedKeys)
	data, err := s.dal.FindAllWithPagination(ctx, filter, projection, skip, limit)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	total, err := s.dal.TotalCount(ctx, filter)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Message)
	}
	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)
	return &types.PaginatedResponse[[]*model.CustomerKYC]{
		Data: data,
		Meta: meta,
	}, nil
}
