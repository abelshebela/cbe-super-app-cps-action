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
	"regexp"
	"strings"
	"time"

	coreio "github.com/hugokessem/coreio/core"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/dal"
	sharedconst "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerKYCRepository struct {
	oracleDB         *sql.DB
	dal              dal.MongoDal[imodel.CustomerKYC, imodel.CustomerKYC]
	logger           utils.Logger
	coll             *mongo.Collection
	membersCol       *mongo.Collection
	linkedAccountCol *mongo.Collection
}

func NewCustomerKYCRepository(client *mongo.Client, oracleDB *sql.DB, cfg *config.VaultConfig, dbName string, custKycCollection string, logger utils.Logger) storage.CustomerKYCRepository {
	return &customerKYCRepository{
		oracleDB:         oracleDB,
		dal:              dal.NewMongoDal[imodel.CustomerKYC, imodel.CustomerKYC](client, cfg, dbName, custKycCollection),
		logger:           logger,
		coll:             client.Database(dbName).Collection(custKycCollection),
		membersCol:       client.Database(dbName).Collection("members"),
		linkedAccountCol: client.Database(dbName).Collection("linked_account"),
	}
}

func (r *customerKYCRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[CustomerKYC][FindAllWithPagination] fetching kyc requests")

	allowed := []string{"review.status", "enabled", "created_at"}
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
		"kyc_status": "review.status",
	}

	for param, mongoField := range nestedFieldMap {
		if v, ok := filterParam.Filters[param]; ok && v != nil && v != "" {
			filter[mongoField] = v
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

func (r *customerKYCRepository) CreateUser(ctx context.Context, userAccount *coreio.CusteomerAccountCreationResponse, userData imodel.CustomerKYC) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	detail := userAccount.CustomerCreationDetail.Detail
	accountDetail := userAccount.AccountCreationDetail.Detail
	customerNumber := detail.CustomerNumber
	userCode := "SA" + customerNumber

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
	username := core.GenerateUsername(firstName, lastName, middleName)

	tx, err := r.oracleDB.BeginTx(ctx, nil)
	if err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] begin tx: %v", err)
		return local_util.HandleDBError(err)
	}
	defer func() { _ = tx.Rollback() }()

	// ── 1. Resolve CBE bank ID ────────────────────────────────────────────
	const bankQ = `SELECT RAWTOHEX(id) FROM banks WHERE is_cbe = 1 AND is_enabled = 1 AND is_deleted = 0 AND ROWNUM = 1`
	var bankID string
	if err := tx.QueryRowContext(ctx, bankQ).Scan(&bankID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("cbe bank not found")
		}
		log.Errorf("[customerKYCRepository][CreateUser] fetch bank id: %v", err)
		return local_util.HandleDBError(err)
	}

	// ── 2. INSERT into USERS ──────────────────────────────────────────────
	const insertUserQ = `
INSERT INTO USERS (
	USER_CODE, USERNAME, CONTACT_EMAIL, CONTACT_PHONE, CUSTOMER_NUMBER,
	SECTOR, OWNERSHIP, INDUSTRY,
	FIRST_NAME, LAST_NAME, MIDDLE_NAME, FULL_NAME,
	GENDER, BIRTH_OF_DATE,
	PIN, PIN_HISTORY, FAILED_LOGIN_ATTEMPT, IS_LOCKED,
	IS_SUPERAPP_ENABLED, IS_USSD_ENABLED, IS_USSD_ACTIVE, IS_BLOCKED, IS_SUPERAPP_ACTIVE,
	LANGUAGE, PUSH_TOKEN, BRANCH_CODE, IS_BUDGET_ENABLED, EXPIRY_AT,
	PIN_CREATED_AT, CREATED_AT, LAST_MODIFIED_AT
) VALUES (
	:1, :2, :3, :4, :5, :6, :7, :8, :9,
	:10, :11, :12, :13, :14, :15, :16, :17, :18, :19,
	:20, :21, :22, :23, :24, :25, :26, :27, :28,
	SYSTIMESTAMP, SYSTIMESTAMP, SYSTIMESTAMP
)`

	currency := detail.Currency
	if strings.TrimSpace(currency) == "" {
		currency = "ETB"
	}

	// Parse T24 YYYYMMDD date string to time.Time for Oracle DATE column
	var dob *time.Time
	if t, parseErr := time.Parse("20060102", detail.DateOfBirth); parseErr == nil {
		dob = &t
	} else {
		log.Warnf("[customerKYCRepository][CreateUser] unparseable DateOfBirth %q: %v", detail.DateOfBirth, parseErr)
	}

	_, err = tx.ExecContext(ctx, insertUserQ,
		userCode,           // :1  USER_CODE
		username,           // :2  USERNAME
		detail.Email,       // :3  CONTACT_EMAIL
		detail.PhoneNumber, // :4  CONTACT_PHONE
		customerNumber,     // :5  CUSTOMER_NUMBER

		detail.Industry,  // :6  SECTOR
		detail.Ownership, // :7  OWNERSHIP
		detail.Industry,  // :8  INDUSTRY

		firstName,       // :9  FIRST_NAME
		lastName,        // :10 LAST_NAME
		middleName,      // :11 MIDDLE_NAME
		detail.FullName, // :12 FULL_NAME

		detail.Gender, // :13 GENDER
		dob,           // :14 BIRTH_OF_DATE — time.Time for Oracle DATE (was string "YYYYMMDD" → ORA-01861)

		"UNSET", // :15 PIN — NOT NULL; user sets PIN later via app
		"",      // :16 PIN_HISTORY
		0,       // :17 FAILED_LOGIN_ATTEMPT
		0,       // :18 IS_LOCKED

		1, // :19 IS_SUPERAPP_ENABLED
		0, // :20 IS_USSD_ENABLED
		0, // :21 IS_USSD_ACTIVE
		0, // :22 IS_BLOCKED
		1, // :23 IS_SUPERAPP_ACTIVE

		"EN", // :24 LANGUAGE
		"",   // :25 PUSH_TOKEN

		detail.AccountOfficer, // :26 BRANCH_CODE
		0,                     // :27 IS_BUDGET_ENABLED
		(*time.Time)(nil),     // :28 EXPIRY_AT — NULL (was 0, caused ORA-00932)
	)
	if err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] user insert: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[customerKYCRepository][CreateUser] user created: %s", userCode)

	// ── 3. INSERT into ACCOUNTS ───────────────────────────────────────────
	const insertAccountQ = `
INSERT INTO ACCOUNTS (
	BANK_ID, ACCOUNT_HOLDER_NAME, ACCOUNT_NUMBER,
	ACCOUNT_CURRENCY, ACCOUNT_TYPE, ACCOUNT_BRANCH, CUSTOMER_NUMBER
) VALUES (
	HEXTORAW(:bank_id), :holder_name, :account_number,
	:currency, :account_type, :branch, :customer_number
) RETURNING RAWTOHEX(ID) INTO :id`

	var accountID string
	_, err = tx.ExecContext(ctx, insertAccountQ,
		sql.Named("bank_id", bankID),
		sql.Named("holder_name", detail.FullName),
		sql.Named("account_number", accountDetail.AccountNumber),
		sql.Named("currency", currency),
		sql.Named("account_type", userData.KYCData.AccountType),
		sql.Named("branch", detail.AccountOfficer),
		sql.Named("customer_number", customerNumber),
		sql.Named("id", sql.Out{Dest: &accountID}),
	)
	if err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] account insert: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[customerKYCRepository][CreateUser] account created: id=%s number=%s", accountID, accountDetail.AccountNumber)

	// ── 4. INSERT into LINKED_ACCOUNTS ────────────────────────────────────
	const insertLinkedQ = `
INSERT INTO LINKED_ACCOUNTS (
	USER_CODE, ACCOUNT_ID, IS_MAIN_ACCOUNT, IS_SUPERAPP_ENABLED, IS_USSD_ENABLED, IS_ACTIVE
) VALUES (
	:user_code, HEXTORAW(:account_id), 1, 1, 0, 1
)`

	_, err = tx.ExecContext(ctx, insertLinkedQ,
		sql.Named("user_code", userCode),
		sql.Named("account_id", accountID),
	)
	if err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] linked account insert: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[customerKYCRepository][CreateUser] linked account created: user=%s account=%s", userCode, accountID)

	if err := tx.Commit(); err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] commit: %v", err)
		return local_util.HandleDBError(err)
	}

	// ── 5. INSERT member into MongoDB (members collection) ────────────────
	// Required so SearchCustomerByCIForAccountNumber can find this customer.
	memberID := bson.NewObjectID()
	now := time.Now()

	newMember := member.User{
		ID:               memberID,
		UserCode:         userCode,
		FullName:         detail.FullName,
		PhoneNumber:      detail.PhoneNumber,
		Email:            detail.Email,
		CustomerNumber:   customerNumber,
		BranchCode:       detail.AccountOfficer,
		Gender:           t24Gender(detail.Gender),
		BirthOfDate:      t24Date(detail.DateOfBirth),
		Industry:         detail.Industry,
		Sector:           detail.Industry,
		Ownership:        detail.Ownership,
		Language:         "EN",
		ISuperappEnabled: true,
		IsBlocked:        false,
		IsLocked:         false,
		Enabled:          true,
		IsActivated:      true,
		CreatedAt:        now,
		LastModifiedAt:   now,
	}

	if _, err := r.membersCol.InsertOne(ctx, newMember); err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] members insert: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[customerKYCRepository][CreateUser] member created in MongoDB: userCode=%s customerNumber=%s", userCode, customerNumber)

	// ── 6. INSERT linked_account into MongoDB ─────────────────────────────
	linkedAccount := member.LinkedAccount{
		ID:                bson.NewObjectID(),
		UserID:            memberID,
		AccountNumber:     accountDetail.AccountNumber,
		AccountHolderName: detail.FullName,
		AccountBranchCode: detail.AccountOfficer,
		AccountType:       userData.KYCData.SubAccountType,
		CustomerNumber:    customerNumber,
		Currency:          currency,
		ISuperappEnabled:  true,
		IsUSSDEnabled:     false,
		LinkedStatus:      true,
		IsMain:            true,
		CreatedAt:         now,
		LastModifiedAt:    now,
	}

	if _, err := r.linkedAccountCol.InsertOne(ctx, linkedAccount); err != nil {
		log.Errorf("[customerKYCRepository][CreateUser] linked_account insert: %v", err)
		return local_util.HandleDBError(err)
	}
	log.Infof("[customerKYCRepository][CreateUser] linked_account created: accountNumber=%s", accountDetail.AccountNumber)

	return nil
}

func (r *customerKYCRepository) ApproveOrReject(ctx context.Context, id, status, rejectionReason string, approved bool) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[customerKYCRepository][ApproveOrReject] updating review for id: %s to %s", id, status)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[customerKYCRepository][ApproveOrReject] invalid object id: %v", err)
		return errors.New(localization.ErrorInvalidID.Code)
	}

	now := time.Now()
	update := bson.M{
		"review.status": status,
		"enabled":       approved,
		"updated_at":    now,
	}

	if status != string(imodel.KYCStatusCancelled) {
		update["review.decision.reviewed_at"] = now
	}

	if rejectionReason != "" {
		update["review.decision.rejection_reason"] = rejectionReason
	}

	filter := bson.M{"_id": objID}

	if _, err := r.dal.UpdateOne(ctx, filter, update); err != nil {
		log.Errorf("[customerKYCRepository][ApproveOrReject] failed to update review: %v", err)
		return errors.New(localization.ErrorUnexpectedError.Code)
	}

	return nil
}

func (r *customerKYCRepository) UpdateKYC(ctx context.Context, id string, data *imodel.CustomerKYC) (*imodel.CustomerKYC, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[customerKYCRepository][UpdateKYC] updating kyc for id: %s", id)
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		log.Errorf("[customerKYCRepository][UpdateKYC] invalid object id: %v", err)
		return nil, errors.New(localization.ErrorInvalidID.Code)
	}

	update := bson.M{}

	if data.Review != nil {
		if data.Review.Status != "" {
			update["review.status"] = string(data.Review.Status)
		}
		if data.Review.ReviewStartedAt != nil {
			update["review.review_started_at"] = data.Review.ReviewStartedAt
		}
		if data.Review.ReviewExpiresAt != nil {
			update["review.review_expires_at"] = data.Review.ReviewExpiresAt
		}
		if data.Review.PickCount != 0 {
			update["review.pick_count"] = data.Review.PickCount
		}
		if data.Review.Assignments != nil {
			update["review.assignments"] = data.Review.Assignments
		}
		if data.Review.Decision != nil {
			if data.Review.Decision.RejectionReason != "" {
				update["review.decision.rejection_reason"] = data.Review.Decision.RejectionReason
			}
			if data.Review.Decision.ReviewedAt != nil {
				update["review.decision.reviewed_at"] = data.Review.Decision.ReviewedAt
			}
			if data.Review.Decision.Reviewer != nil {
				update["review.decision.reviewer"] = data.Review.Decision.Reviewer
			}
		}
	}

	filter := bson.M{"_id": objID}

	result, err := r.dal.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Errorf("[customerKYCRepository][UpdateKYC] failed to update kyc: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return &result, err
}

func t24Gender(g string) sharedconst.Gender {
	if strings.ToUpper(strings.TrimSpace(g)) == "FEMALE" {
		return sharedconst.Female
	}
	return sharedconst.Male
}

func t24Date(yyyymmdd string) time.Time {
	t, _ := time.Parse("20060102", yyyymmdd)
	return t
}

func (r *customerKYCRepository) FindKYCOnboardingForExport(ctx context.Context, from, to time.Time, customerName string) ([]imodel.ExportSelfActivationRequest, error) {
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
			Key: "$project",
			Value: bson.M{
				"_id": 0,

				"customer_name": "$kyc_data.full_name",
				"phone_number":  "$kyc_data.phone_number",
				"gender":        "$kyc_data.gender",
				"date_of_birth": bson.M{
					"$dateToString": bson.M{
						"format": "%Y-%m-%d",
						"date":   "$kyc_data.birth_date",
					},
				},
				"region":            "$kyc_data.address.region",
				"registration_date": "$last_modified_at",
				"rejection_reason":  "$kyc_reject_reason_failed",
				"customer_status":   "NEW",
				"kyc_status":        "$review.status",
			},
		}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		r.logger.Errorf("failed to get kyc onboarding data: %v", err)
		return nil, local_util.HandleDBError(err)
	}
	defer cursor.Close(ctx)

	var result []imodel.ExportSelfActivationRequest
	if err := cursor.All(ctx, &result); err != nil {
		return nil, local_util.HandleDBError(err)
	}

	return result, nil
}
