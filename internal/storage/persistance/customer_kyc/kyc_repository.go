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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerKYCRepository struct {
	oracleDB *sql.DB
	dal      dal.MongoDal[imodel.CustomerKYC, imodel.CustomerKYC]
	logger   utils.Logger
	coll     *mongo.Collection
}

func NewCustomerKYCRepository(client *mongo.Client, oracleDB *sql.DB, cfg *config.VaultConfig, dbName string, custKycCollection string, logger utils.Logger) storage.CustomerKYCRepository {
	return &customerKYCRepository{
		oracleDB: oracleDB,
		dal:      dal.NewMongoDal[imodel.CustomerKYC, imodel.CustomerKYC](client, cfg, dbName, custKycCollection),
		logger:   logger,
		coll:     client.Database(dbName).Collection(custKycCollection),
	}
}

// func (r *customerKYCRepository) Create(ctx context.Context, kyc *imodel.CustomerKYC) error {

// 	log.Infof("[CustomerKYC][Create] creating kyc for customer: %s", kyc.CustomerCode)
// 	kyc.CreatedAt = time.Now()
// 	kyc.UpdatedAt = time.Now()
// 	kyc.KYCStatus = "PENDING"
// 	kyc.CustomerStatus = constants.CustomerPending

// 	if _, err := r.dal.InsertOne(ctx, *kyc); err != nil {
// 		log.Errorf("[CustomerKYC][Create] failed to create kyc: %v", err)
// 		return errors.New(localization.ErrorUnexpectedError.Code)
// 	}

// 	return nil
// }

func (r *customerKYCRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindAllWithPagination] fetching kyc requests")

	allowed := []string{"kyc_status"}
	filter, skip, limit := lib.FilterBuilder(filterParam, bson.M{}, allowed)
	if filterParam.Search != "" {
		q := bson.M{"$regex": filterParam.Search, "$options": "i"}
		filter["$or"] = []bson.M{
			{"service_name": q},
			{"service_code": q},
			{"service_type": q},
			{"kyc_status": q},
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

	// resultData := make([]*imodel.CustomerKYC, len(results))
	// for i := range data {
	// 	resultData[i] = &data[i]
	// }

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

func (r *customerKYCRepository) CreateUser(ctx context.Context, userAccount types.Account, userData imodel.CustomerKYC) error {
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

	var firstName, middleName, lastName string

	nameParts := strings.Fields(userAccount.CustomerName)
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}
	if len(nameParts) > 2 {
		middleName = nameParts[2]
	}

	userCode := "SA" + userAccount.CustomerNumber
	username := core.GenerateUsername(firstName, lastName, middleName)

	_, err := r.oracleDB.ExecContext(
		ctx,
		q,
		userCode,                      // :1  USER_CODE
		username,                      // :2  USERNAME
		"",                            // :3  CONTACT_EMAIL
		userData.KYCData.PhoneNumber,  // :4  CONTACT_PHONE
		userAccount.CustomerNumber,    // :5  CUSTOMER_NUMBER
		"",                            // :6  SECTOR
		"",                            // :7  OWNERSHIP
		"",                            // :8  INDUSTRY
		firstName,                     // :9  FIRST_NAME
		lastName,                      // :10 LAST_NAME
		middleName,                    // :11 MIDDLE_NAME
		userAccount.CustomerName,      // :12 FULL_NAME
		userAccount.Gender,            // :13 GENDER
		userData.KYCData.BirthDate,    // :14 BIRTH_OF_DATE
		"",                            // :15 PIN
		"",                            // :16 PIN_HISTORY
		0,                             // :17 FAILED_LOGIN_ATTEMPT
		0,                             // :18 IS_LOCKED
		0,                             // :19 IS_SUPERAPP_ENABLED
		0,                             // :20 IS_USSD_ENABLED
		0,                             // :21 IS_USSD_ACTIVE
		0,                             // :22 IS_BLOCKED
		0,                             // :23 IS_SUPERAPP_ACTIVE
		"EN",                          // :24 LANGUAGE
		"",                            // :25 PUSH_TOKEN
		userAccount.AccountBranchCode, // :26 BRANCH_CODE
		0,                             // :27 IS_BUDGET_ENABLED
		0,                             // :28 EXPIRY_AT
	)
	if err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] insert failed: %v", err)
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

	var update bson.M
	if rejectionReason != "" {
		update = bson.M{"kyc_status": status, "kyc_approved": approved, "kyc_reject_reason": rejectionReason}
	} else {
		update = bson.M{"kyc_status": status, "kyc_approved": approved}
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
