package customer

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type selfActivationRepository struct {
	oracleDB          *sql.DB
	selfActivationDal dal.MongoDal[imodel.SelfActivationUser, imodel.SelfActivationUser]
	inReviewDal       dal.MongoDal[imodel.StartedKycReview, imodel.StartedKycReview]
	logger            utils.Logger
	coll              *mongo.Collection
}

func NewSelfActivationRepository(client *mongo.Client, oracleDB *sql.DB, cfg *config.VaultConfig, dbName, collection string, logger utils.Logger) storage.SelfActivationKYCRepository {
	return &selfActivationRepository{
		oracleDB:          oracleDB,
		selfActivationDal: dal.NewMongoDal[imodel.SelfActivationUser, imodel.SelfActivationUser](client, cfg, dbName, collection),
		inReviewDal:       dal.NewMongoDal[imodel.StartedKycReview, imodel.StartedKycReview](client, cfg, dbName, "started_kyc_reviews"),
		logger:            logger,
		coll:              client.Database(dbName).Collection(collection),
	}
}

func (r *selfActivationRepository) CheckIfUserOrAccountExists(ctx context.Context, userData *imodel.SelfActivationUser) (bool, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	// phone := strings.TrimSpace(userData.PhoneNumber)
	customerID := strings.TrimSpace(userData.CustomerNumber)

	if customerID == "" {
		return false, nil
	}

	const userExistsQuery = `
		SELECT COUNT(1)
		FROM USERS u
		WHERE TRIM(u.CUSTOMER_NUMBER) = :1
	`

	var userCount int
	if err := r.oracleDB.QueryRowContext(ctx, userExistsQuery, customerID).Scan(&userCount); err != nil {
		log.Errorf("[SelfActivationKYC][CheckIfUserOrAccountExists] failed to query USERS: %v", err)
		return false, local_util.HandleDBError(err)
	}

	if userCount == 0 {
		return false, nil
	}

	return true, nil

	// const linkedAccountQuery = `
	// 	SELECT COUNT(1)
	// 	FROM LINKED_ACCOUNTS la
	// 	JOIN USERS u
	// 	    ON u.USER_CODE = la.USER_CODE
	// 	WHERE TRIM(u.CUSTOMER_NUMBER) = :1
	// 	AND la.IS_ACTIVE = 1
	// `

	// var linkedCount int
	// if err := r.oracleDB.QueryRowContext(ctx, linkedAccountQuery, customerID).Scan(&linkedCount); err != nil {
	// 	log.Errorf("[SelfActivationKYC][CheckIfUserOrAccountExists] failed to query LINKED_ACCOUNTS: %v", err)
	// 	return false, local_util.HandleDBError(err)
	// }

	// return linkedCount > 0, nil
}

func (r *selfActivationRepository) FindAllWithPaginationSA(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.SelfActivationUser], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindAllWithPagination] fetching kyc requests")

	allowed := []string{"kyc.kyc_status", "enabled", "created_at"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)

	filter["is_deleted"] = bson.M{"$ne": true}

	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"name": q},
			{"phone_number": q},
		}
	}

	nestedFieldMap := map[string]string{
		"kyc_status": "kyc.kyc_status",
	}
	for param, mongoField := range nestedFieldMap {
		if v, ok := filterParam.Filters[param]; ok && v != nil {
			statuses := local_util.StringSliceFromFilterValue(v)
			if len(statuses) > 0 {
				filter[mongoField] = bson.M{
					"$in": statuses,
				}
			}
		}
	}

	results, err := r.selfActivationDal.FindAllWithPaginationE(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[CustomerKYC][FindAllWithPagination] failed to fetch data: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.selfActivationDal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[CustomerKYC][FindAllWithPagination] failed to get total count: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]imodel.SelfActivationUser]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *selfActivationRepository) FindByIDSA(ctx context.Context, id string) (*imodel.SelfActivationUser, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindByID] fetching kyc by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[CustomerKYC][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	result, err := r.selfActivationDal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CustomerKYC][FindByID] failed to find kyc: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}

func (r *selfActivationRepository) UpdateKYCStatusSA(ctx context.Context, id, status, rejectionReason string, approved bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][UpdateKYCStatus] updating kyc status for id: %s to %s", id, status)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[CustomerKYC][UpdateKYCStatus] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{
		"kyc.kyc_status":   status,
		"enabled":          approved,
		"last_modified_at": time.Now(),
		"updated_at":       time.Now(),
	}

	if rejectionReason != "" {
		update["kyc_reject_reason"] = rejectionReason
	}

	filter := bson.M{"_id": objID}

	if _, err := r.selfActivationDal.UpdateOne(ctx, filter, update); err != nil {
		log.Errorf("[CustomerKYC][UpdateKYCStatus] failed to update kyc status: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *selfActivationRepository) FindKycInReviewSA(ctx context.Context, kycID string) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindKycInReview] fetching started kyc by id: %s", kycID)
	objID, err := bson.ObjectIDFromHex(kycID)
	if err != nil {
		log.Errorf("[CustomerKYC][FindKycInReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"kyc_id": objID}
	result, err := r.inReviewDal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CustomerKYC][FindKycInReview] failed to find kyc: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}

func (r *selfActivationRepository) StartKycReviewSA(ctx context.Context, reviewData *imodel.StartedKycReview) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Debugf("[CustomerKYC][StartKycReview] inserting started review: %+v", reviewData)

	result, err := r.inReviewDal.InsertOne(ctx, *reviewData)
	if err != nil {
		log.Errorf("[CustomerKYC][StartKycReview] failed to insert started kyc review: %v", err)
		log.Debugf("[CustomerKYC][StartKycReview] payload: %+v", reviewData)
		return nil, local_util.HandleDBError(err)
	}

	return &result, nil
}

func (r *selfActivationRepository) UpdateKycReviewSA(ctx context.Context, kycID string, reviewData *imodel.StartedKycReview) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[UpdateKycReview] updating kyc review for id: %s", kycID)
	objID, err := bson.ObjectIDFromHex(kycID)
	if err != nil {
		log.Errorf("[UpdateKycReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"kyc_id": objID}
	update := bson.M{
		"review_status": reviewData.ReviewStatus,
		"picked_at":     reviewData.PickedAt,
		"started_at":    reviewData.StartedAt,
		"expires_at":    reviewData.ExpiresAt,
		"reviewer":      reviewData.Reviewer,
		"picked_by":     reviewData.PickedBy,
		"pick_reason":   reviewData.PickReason,
		"pick_count":    reviewData.PickCount,
	}

	result, err := r.inReviewDal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[UpdateKycReview] failed to update kyc review: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, nil
}

func (r *selfActivationRepository) FindForExport(ctx context.Context, from, to time.Time, customerName string) ([]imodel.ExportSelfActivationRequest, error) {
	matchFilter := bson.M{
		"created_at": bson.M{
			"$gte": from,
			"$lte": to,
		},
		"is_deleted": bson.M{"$ne": true},
	}

	if customerName != "" {
		matchFilter["name"] = bson.M{
			"$regex":   "^" + regexp.QuoteMeta(customerName) + "$",
			"$options": "i",
		}
	}

	pipeline := mongo.Pipeline{
		bson.D{{
			Key:   "$match",
			Value: matchFilter,
		}},
		bson.D{{
			Key: "$project",
			Value: bson.M{
				"_id": 0,

				"customer_name":     "$name",
				"phone_number":      1,
				"gender":            1,
				"date_of_birth":     "$birth_date",
				"region":            "$address.region",
				"registration_date": "$updated_at",
				"rejection_reason":  "$kyc_reject_reason",
				"customer_status":   "NEW",
				"kyc_status":        "$kyc.kyc_status",
			},
		}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var result []imodel.ExportSelfActivationRequest
	if err := cursor.All(ctx, &result); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}
