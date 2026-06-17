package customer

import (
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance/customer_kyc/core"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	coreio "github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerKYCRepository struct {
	oracleDB *sql.DB
	dal      dal.MongoDal[imodel.CustomerKYC, imodel.CustomerKYC]
	kycDal   dal.MongoDal[imodel.StartedKycReview, imodel.StartedKycReview]
	logger   utils.Logger
	coll     *mongo.Collection
}

func NewCustomerKYCRepository(client *mongo.Client, oracleDB *sql.DB, cfg *config.VaultConfig, dbName string, custKycCollection string, logger utils.Logger) storage.CustomerKYCRepository {
	return &customerKYCRepository{
		oracleDB: oracleDB,
		dal:      dal.NewMongoDal[imodel.CustomerKYC, imodel.CustomerKYC](client, cfg, dbName, custKycCollection),
		kycDal:   dal.NewMongoDal[imodel.StartedKycReview, imodel.StartedKycReview](client, cfg, dbName, "started_kyc_reviews"),
		logger:   logger,
		coll:     client.Database(dbName).Collection(custKycCollection),
	}
}

func (r *customerKYCRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindAllWithPagination] fetching kyc requests")

	allowed := []string{"kyc_status", "enabled", "created_at"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"kyc_data.full_name": q},
			{"kyc_data.phone_number": q},
			{"kyc_status": q},
			{"user_id": q},
		}
	}

	results, err := r.dal.FindAllWithPagination(ctx, filter, bson.M{}, skip, limit)
	if err != nil {
		log.Errorf("[CustomerKYC][FindAllWithPagination] failed to fetch data: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	total, err := r.dal.TotalCount(ctx, filter)
	if err != nil {
		log.Errorf("[CustomerKYC][FindAllWithPagination] failed to get total count: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	meta := local_util.BuildPaginationMeta(total, filterParam.Page, filterParam.PerPage)

	return &types.PaginatedResponse[[]imodel.CustomerKYC]{
		Data: results,
		Meta: meta,
	}, nil
}

func (r *customerKYCRepository) FindByID(ctx context.Context, id string) (*imodel.CustomerKYC, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindByID] fetching kyc by id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[CustomerKYC][FindByID] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"_id": objID}
	result, err := r.dal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CustomerKYC][FindByID] failed to find kyc: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}
func (r *customerKYCRepository) CreateUser(ctx context.Context, userAccount *coreio.CreateCustomerResult, userData imodel.CustomerKYC) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	const q = `
INSERT INTO USERS (
	USER_CODE,
	USERNAME,
	CONTACT_EMAIL,
	CONTACT_PHONE,
	CUSTOMER_NUMBER,
	SECTOR,
	OWNERSHIP,
	INDUSTRY,
	FIRST_NAME,
	LAST_NAME,
	MIDDLE_NAME,
	FULL_NAME,
	GENDER,
	BIRTH_OF_DATE,
	PIN,
	PIN_HISTORY,
	FAILED_LOGIN_ATTEMPT,
	IS_LOCKED,
	IS_SUPERAPP_ENABLED,
	IS_USSD_ENABLED,
	IS_USSD_ACTIVE,
	IS_BLOCKED,
	IS_SUPERAPP_ACTIVE,
	LANGUAGE,
	PUSH_TOKEN,
	BRANCH_CODE,
	IS_BUDGET_ENABLED,
	EXPIRY_AT,
	PIN_CREATED_AT,
	CREATED_AT,
	LAST_MODIFIED_AT
) VALUES (
	:1, :2, :3, :4, :5, :6, :7, :8, :9,
	:10, :11, :12, :13, :14, :15, :16, :17, :18, :19,
	:20, :21, :22, :23, :24, :25, :26, :27, :28,
	SYSTIMESTAMP, SYSTIMESTAMP, SYSTIMESTAMP
)`

	detail := userAccount.Detail

	var firstName, middleName, lastName string

	nameParts := strings.Fields(detail.FullName)

	switch len(nameParts) {
	case 1:
		firstName = nameParts[0]

	case 2:
		firstName = nameParts[0]
		lastName = nameParts[1]

	default:
		firstName = nameParts[0]
		middleName = strings.Join(nameParts[1:len(nameParts)-1], " ")
		lastName = nameParts[len(nameParts)-1]
	}

	userCode := "SA" + detail.Customer

	username := core.GenerateUsername(
		firstName,
		lastName,
		middleName,
	)

	_, err := r.oracleDB.ExecContext(
		ctx,
		q,

		userCode,           // USER_CODE
		username,           // USERNAME
		detail.Email,       // CONTACT_EMAIL
		detail.PhoneNumber, // CONTACT_PHONE
		"",                 // CUSTOMER_NUMBER

		detail.Industry,  // SECTOR
		detail.Ownership, // OWNERSHIP
		detail.Industry,  // INDUSTRY

		firstName,       // FIRST_NAME
		lastName,        // LAST_NAME
		middleName,      // MIDDLE_NAME
		detail.FullName, // FULL_NAME

		detail.Gender,      // GENDER
		detail.DateOfBirth, // BIRTH_OF_DATE

		"", // PIN
		"", // PIN_HISTORY
		0,  // FAILED_LOGIN_ATTEMPT
		0,  // IS_LOCKED

		0, // IS_SUPERAPP_ENABLED
		0, // IS_USSD_ENABLED
		0, // IS_USSD_ACTIVE
		0, // IS_BLOCKED
		0, // IS_SUPERAPP_ACTIVE

		"EN", // LANGUAGE
		"",   // PUSH_TOKEN

		detail.AccountOfficer, // BRANCH_CODE

		0, // IS_BUDGET_ENABLED
		0, // EXPIRY_AT
	)
	if err != nil {
		log.Errorf(
			"[customerKYCRepository][CreateUser] insert failed: %v",
			err,
		)

		return local_util.HandleDBError(err)
	}

	return nil
}

func (r *customerKYCRepository) UpdateKYCStatus(ctx context.Context, id, status, rejectionReason string, approved bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][UpdateKYCStatus] updating kyc status for id: %s to %s", id, status)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[CustomerKYC][UpdateKYCStatus] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{
		"kyc_status":       status,
		"kyc_approved":     approved,
		"enabled":          approved,
		"last_modified_at": time.Now(),
	}
	if rejectionReason != "" {
		update["kyc_reject_reason"] = rejectionReason
	}

	filter := bson.M{"_id": objID}

	if _, err := r.dal.UpdateOne(ctx, filter, update); err != nil {
		log.Errorf("[CustomerKYC][UpdateKYCStatus] failed to update kyc status: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

// func (r *customerKYCRepository) Delete(ctx context.Context, id string) error {

// 	log.Infof("[CustomerKYC][Delete] hard deleting kyc for id: %s", id)
// 	objID, err := bson.ObjectIDFromHex(id)
// 	if err != nil {
// 		log.Errorf("[CustomerKYC][Delete] invalid object id: %v", err)
// 		return errors.New(localization.ErrorInvalidID.Code)
// 	}

// 	filter := bson.M{"_id": objID}
// 	if err := r.dal.DeleteOneH(ctx, filter); err != nil {
// 		log.Errorf("[CustomerKYC][Delete] failed to delete kyc: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	return nil
// }

func (r *customerKYCRepository) FindKycInReview(ctx context.Context, kycID string) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindKycInReview] fetching started kyc by id: %s", kycID)
	objID, err := bson.ObjectIDFromHex(kycID)
	if err != nil {
		log.Errorf("[CustomerKYC][FindKycInReview] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	filter := bson.M{"kyc_id": objID}
	result, err := r.kycDal.FindOne(ctx, filter, nil)
	if err != nil {
		log.Errorf("[CustomerKYC][FindKycInReview] failed to find kyc: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}

func (r *customerKYCRepository) StartKycReview(ctx context.Context, reviewData *imodel.StartedKycReview) (*imodel.StartedKycReview, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Debugf("[CustomerKYC][StartKycReview] inserting started review: %+v", reviewData)

	result, err := r.kycDal.InsertOne(ctx, *reviewData)
	if err != nil {
		log.Errorf("[CustomerKYC][StartKycReview] failed to insert started kyc review: %v", err)
		log.Debugf("[CustomerKYC][StartKycReview] payload: %+v", reviewData)
		return nil, local_util.HandleDBError(err)
	}

	return &result, nil
}

func (r *customerKYCRepository) UpdateKycReview(ctx context.Context, kycID string, reviewData *imodel.StartedKycReview) (*imodel.StartedKycReview, error) {
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

	result, err := r.kycDal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[UpdateKycReview] failed to update kyc review: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, nil
}
