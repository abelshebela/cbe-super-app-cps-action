package customer_oracle

import (
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_constants "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/constants"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type customerOracleRepository struct {
	db            DBTX
	cfg           config.VaultConfig
	kafkaProducer *kafka.AccessListSegmentationProducer
	accBlock      storage.AccountBlockRepository
	logger        utils.Logger
}

// EnableOrDisable implements [storage.CustomerRepository].
func (c *customerOracleRepository) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	const userQ = `
	UPDATE USERS
	SET
		IS_SUPERAPP_ACTIVE = :1,
		LAST_MODIFIED_AT = SYSTIMESTAMP
	WHERE USER_CODE = :2 AND IS_DELETED = 0`

	res, err := c.db.ExecContext(ctx, userQ, local_util.BoolToOracleNumber(enable), id)
	if err != nil {
		c.logger.Errorf("[CustomerRepository][EnableOrDisable] user update failed: %v", err)
		return local_util.HandleDBError(err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return localization.ErrorResourceNotFound
	}

	return nil
}

// FetchLinkedAccount implements [storage.CustomerRepository].
func (c *customerOracleRepository) FetchLinkedAccount(ctx context.Context, id string) ([]model.LinkedAccount, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*customer.CustomerListResponse], error) {
	panic("unimplemented")
}

func (c *customerOracleRepository) FindByID(ctx context.Context, id string) (*member.User, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[customerOracleRepository][FindByID] fetching customer by id: %s", id)

	const query = `
		SELECT
			RAWTOHEX(ID),
			USER_CODE
		FROM USERS
		WHERE ID = HEXTORAW(:1)
	`

	user := &member.User{}

	err := c.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.UserCode,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[customerOracleRepository][FindByID] customer not found for id: %s", id)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}

		log.Errorf("[customerOracleRepository][FindByID] failed to fetch customer: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	log.Infof("[customerOracleRepository][FindByID] customer retrieved successfully")
	return user, nil
}

func (c *customerOracleRepository) FindUserByUserCode(ctx context.Context, userCode string) (*member.User, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[customerOracleRepository][FindUserByUserCode] fetching customer by userCode: %s", userCode)

	const query = `
		SELECT
			USER_CODE,
			FULL_NAME,
			USERNAME,
			CONTACT_PHONE,
			CONTACT_EMAIL,
			CUSTOMER_NUMBER,
			IS_SUPERAPP_ENABLED,
			IS_USSD_ENABLED,
			IS_BLOCKED
		FROM USERS
		WHERE USER_CODE = :1
	`

	user := &member.User{}

	err := c.db.QueryRowContext(ctx, query, userCode).Scan(
		&user.UserCode,
		&user.FullName,
		&user.Username,
		&user.PhoneNumber,
		&user.Email,
		&user.CustomerNumber,
		&user.ISuperappEnabled,
		&user.IsUSSDEnabled,
		&user.IsBlocked,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[customerOracleRepository][FindUserByUserCode] customer not found for userCode: %s", userCode)
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}

		log.Errorf("[customerOracleRepository][FindUserByUserCode] failed to fetch customer: %v", err)
		return nil, local_util.HandleDBError(err)
	}

	log.Infof("[customerOracleRepository][FindUserByUserCode] customer retrieved successfully")
	return user, nil
}

// FindCustomerByID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerByID(ctx context.Context, id string) (*member.User, error) {
	panic("unimplemented")
}

// FindCustomerByIDs implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerByIDs(ctx context.Context, ids []string) ([]member.User, error) {
	panic("unimplemented")
}

// FindCustomerByUserCode implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerByUserCode(ctx context.Context, usercode string) (*member.User, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[CustomerRepository][FindCustomerByUserCode][oracle] fetching feedback user fields by user_code: %s", usercode)

	// Only fetch fields needed for feedback
	userQuery := `
	SELECT
	  u.full_name,
	  u.contact_phone,
	  u.contact_email,
	  ld.platform
	FROM users u
	join linked_devices ld on ld.user_code = u.user_code
	WHERE u.user_code = :1`

	var (
		platform, fullName, phone, email string
	)
	err := c.db.QueryRowContext(ctx, userQuery, usercode).Scan(&fullName, &phone, &email, &platform)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, localization.ErrorResourceNotFound
		}
		log.Errorf("[CustomerRepository][FindCustomerByUserCode] user query failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}
	response := &member.User{
		UserCode:    usercode,
		FullName:    fullName,
		PhoneNumber: phone,
		Email:       email,
		Platform:    shared_constants.Platform(platform),
	}
	return response, nil
}

// // FindCustomerDetailByID implements [storage.CustomerRepository].
// func (c *customerOracleRepository) FindCustomerDetailByID(ctx context.Context, id string) (*customer.CustomerDetailResponse, error) {

// 	log.Infof("[CustomerRepository][FindCustomerDetailByID] fetching customer detail by user_code: %s", id)

// 	// 1. Fetch user info by user_code
// 	userQuery := `
// 	SELECT
// 	  RAWTOHEX(u.id),
// 	  u.full_name,
// 	  u.gender,
// 	  u.contact_phone,
// 	  u.contact_email,
// 	  u.customer_number,
// 	  u.is_superapp_active,
// 	  u.birth_of_date
// 	FROM users u
// 	WHERE u.user_code = :1`

// 	var (
// 		userID, fullName, gender, phone, email, customerNumber string
// 		isActive                                               int
// 		birthOfDate                                            sql.NullTime
// 	)
// 	err := c.db.QueryRowContext(ctx, userQuery, id).Scan(&userID, &fullName, &gender, &phone, &email, &customerNumber, &isActive, &birthOfDate)
// 	if err != nil {
// 		if err.Error() == "sql: no rows in result set" {
// 			return nil, localization.ErrorResourceNotFound
// 		}
// 		log.Errorf("[CustomerRepository][FindCustomerDetailByID] user query failed: %v", err)
// 		return nil, localization.ErrorUnexpectedError
// 	}

// 	// 2. Fetch linked accounts with account and branch info
// 	linkedQuery := `
// 	SELECT
// 	  a.account_number,
// 	  a.account_holder_name,
// 	  a.account_type,
// 	  ab.code,
// 	  la.is_active,
// 	  ab.name
// 	FROM linked_accounts la
// 	JOIN accounts a ON a.id = la.account_id
// 	JOIN account_blocks ab ON ab.id = a.bank_id
// 	WHERE la.user_code = :1`

// 	rows, err := c.db.QueryContext(ctx, linkedQuery, id)
// 	if err != nil {
// 		log.Errorf("[CustomerRepository][FindCustomerDetailByID] linked accounts query failed: %v", err)
// 		return nil, localization.ErrorUnexpectedError
// 	}
// 	defer rows.Close()

// 	var linkedAccounts []customer.LinkedAccount
// 	for rows.Next() {
// 		var accNum, accHolder, accType, branchCode, branchName sql.NullString
// 		var isActiveAcc int
// 		if err := rows.Scan(&accNum, &accHolder, &accType, &branchCode, &isActiveAcc, &branchName); err != nil {
// 			log.Errorf("[CustomerRepository][FindCustomerDetailByID] scan failed: %v", err)
// 			return nil, localization.ErrorUnexpectedError
// 		}
// 		linkedAccounts = append(linkedAccounts, customer.LinkedAccount{
// 			AccountNumber:     accNum.String,
// 			AccountHolderName: accHolder.String,
// 			AccountType:       accType.String,
// 			AccountBranchCode: branchCode.String,
// 			IsActive:          isActiveAcc == 1,
// 			AccountBranchName: branchName.String,
// 		})
// 	}

// 	response := &customer.CustomerDetailResponse{
// 		ID:            userID,
// 		LinkedAccount: linkedAccounts,
// 		PersonalInfo: customer.PersonalInfo{
// 			FullName:       fullName,
// 			Gender:         gender,
// 			PhoneNumber:    phone,
// 			Email:          email,
// 			CustomerNumber: customerNumber,
// 			IsActivated:    isActive == 1,
// 			DateOfBirth:    "",
// 		},
// 	}
// 	if birthOfDate.Valid {
// 		response.PersonalInfo.DateOfBirth = birthOfDate.Time.Format("2006-01-02")
// 	}

// 	return response, nil
// }

// FindCustomerDetailByID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerDetailByID(ctx context.Context, id string) (*customer.CustomerDetailResponse, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[CustomerRepository][FindCustomerDetailByID] fetching customer detail by user_code: %s", id)

	// New query: join users, account_blocks, linked_accounts, accounts, fetch all required fields
	query := `
		       SELECT
			       RAWTOHEX(u.id),
			       u.full_name,
			       u.gender,
			       u.contact_phone,
			       u.contact_email,
			       u.customer_number,
			       u.is_superapp_active,
				   u.is_ussd_active,
				   u.is_superapp_enabled,
				   u.is_ussd_enabled,
				   u.is_blocked,
			       u.birth_of_date,
			       a.account_holder_name,
			       a.account_type,
			       a.account_number,
			       ab.BRANCH_NAME,
			       ab.BRANCH_CODE,
			       la.is_active
		       FROM linked_accounts la
		       LEFT JOIN accounts a ON a.id = la.account_id
		       LEFT JOIN users u ON la.user_code = u.user_code
		       LEFT JOIN COMPANY ab ON ab.DAO_CODE = u.branch_code
		       WHERE la.user_code = :1`

	rows, err := c.db.QueryContext(ctx, query, id)
	if err != nil {
		log.Errorf("[CustomerRepository][FindCustomerDetailByID] query failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}
	defer rows.Close()

	var (
		userID, fullName, gender, phone, email, customerNumber                      string
		isSuperAppActive, isUSSDActive, isBlocked, isSuperAppEnabled, isUSSDEnalbed int
		birthOfDate                                                                 sql.NullTime
		linkedAccounts                                                              []customer.LinkedAccount
		fetchedFirstRow                                                             bool
	)

	for rows.Next() {
		var (
			accHolder, accType, accNum, branchName, branchCode                                          sql.NullString
			isActiveAcc                                                                                 int
			rowUserID, rowFullName, rowGender, rowPhone, rowEmail, rowCustomerNumber                    string
			rowIsSuperAppActive, rowIsUSSDActive, rowIsUSSDEnabled, rowIsSupperAppEnabled, rowIsBlocked int
			rowBirthOfDate                                                                              sql.NullTime
		)
		if err := rows.Scan(
			&rowUserID, &rowFullName, &rowGender, &rowPhone, &rowEmail, &rowCustomerNumber,
			&rowIsSuperAppActive, &rowIsUSSDActive, &rowIsSupperAppEnabled, &rowIsUSSDEnabled, &rowIsBlocked, &rowBirthOfDate,
			&accHolder, &accType, &accNum, &branchName, &branchCode, &isActiveAcc,
		); err != nil {
			log.Errorf("[CustomerRepository][FindCustomerDetailByID] scan failed: %v", err)
			return nil, localization.ErrorUnexpectedError
		}
		// Set user info from first row
		if !fetchedFirstRow {
			userID = rowUserID
			fullName = rowFullName
			gender = rowGender
			phone = rowPhone
			email = rowEmail
			customerNumber = rowCustomerNumber
			isBlocked = rowIsBlocked
			isSuperAppActive = rowIsSuperAppActive
			isSuperAppEnabled = rowIsSupperAppEnabled
			isUSSDActive = rowIsUSSDActive
			isUSSDEnalbed = rowIsUSSDEnabled
			isUSSDActive = rowIsUSSDActive
			birthOfDate = rowBirthOfDate
			fetchedFirstRow = true
		}
		linkedAccounts = append(linkedAccounts, customer.LinkedAccount{
			AccountNumber:     accNum.String,
			AccountHolderName: accHolder.String,
			AccountType:       accType.String,
			AccountBranchCode: branchCode.String,
			IsActive:          isActiveAcc == 1 || isUSSDActive == 1,
			AccountBranchName: branchName.String,
		})
	}

	if !fetchedFirstRow {
		// No rows found
		return nil, localization.ErrorResourceNotFound
	}

	response := &customer.CustomerDetailResponse{
		ID:            userID,
		LinkedAccount: linkedAccounts,
		PersonalInfo: customer.PersonalInfo{
			FullName:            fullName,
			Gender:              gender,
			PhoneNumber:         phone,
			Email:               email,
			CustomerNumber:      customerNumber,
			IsActivated:         isSuperAppActive == 1,
			IsSupperAppActivate: isSuperAppActive == 1,
			IsUSSDActivate:      isUSSDActive == 1,
			IsSupperAppEnabled:  isSuperAppEnabled == 1,
			IsUSSDEnabled:       isUSSDEnalbed == 1,
			IsBlocked:           isBlocked == 1,
			DateOfBirth:         birthOfDate.Time.Format("2006-01-02"),
		},
	}
	if birthOfDate.Valid {
		response.PersonalInfo.DateOfBirth = birthOfDate.Time.Format("2006-01-02")
	}

	return response, nil
}

// FindCustomerLinkedAccountByUserID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerLinkedAccountByUserID(ctx context.Context, userID string) (*model.LinkedAccount, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[CustomerRepository][FindCustomerLinkedAccountByUserID] fetching linked account for user ID: %s", userID)

	// 1. Find account_id from linked_accounts where user_id = :1 and is_main = 1
	var accountID string
	queryLinked := `SELECT account_id FROM linked_accounts WHERE user_code = :1 AND is_main_account = 1 AND is_deleted = 0`
	err := c.db.QueryRowContext(ctx, queryLinked, userID).Scan(&accountID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Errorf("[CustomerRepository][FindCustomerLinkedAccountByUserID] linked account not found for user ID: %s", userID)
			return nil, localization.ErrorResourceNotFound
		}
		log.Errorf("[CustomerRepository][FindCustomerLinkedAccountByUserID] failed to find linked account: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	// 2. Find account_number from accounts where id = account_id
	var accountNumber, branchCode, branchName string
	queryAccount := `SELECT account_number, branch_code, branch_name FROM accounts WHERE id = :1 AND is_deleted = 0`
	err = c.db.QueryRowContext(ctx, queryAccount, accountID).Scan(&accountNumber, &branchCode, &branchName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			log.Errorf("[CustomerRepository][FindCustomerLinkedAccountByUserID] account not found for account_id: %s", accountID)
			return nil, localization.ErrorResourceNotFound
		}
		log.Errorf("[CustomerRepository][FindCustomerLinkedAccountByUserID] failed to find account: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	linkedAccount := &model.LinkedAccount{
		AccountNumber: accountNumber,
		BranchCode:    branchCode,
		BranchName:    branchName,
	}
	return linkedAccount, nil
}

// SearchCustomerByCIForAccountNumber implements [storage.CustomerRepository].
func (c *customerOracleRepository) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer.CustomerListResponse, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[CustomerRepository][SearchCustomerByCIForAccountNumber] searching members by value: %s", number)

	// Oracle SQL: join USERS and LINKED_ACCOUNTS, search by phone, customer number, user code, or account id (as hex string)
	query := `
			 SELECT
				 RAWTOHEX(u.id) AS id,
				 u.user_code,
				 u.contact_email,
				 u.customer_number,
				 u.full_name,
				 u.contact_phone,
				 NVL(co.BRANCH_CODE, '') AS branch_code,
				 u.gender,
				 u.is_blocked,
				 NVL(ac.account_number, '') AS account_number,
				 u.IS_SUPERAPP_ENABLED,
				 u.IS_USSD_ENABLED,
				 NVL(co.BRANCH_NAME, '') AS branch_name,
				 NVL(co.ACCOUNT_TYPE, '') AS account_type,
				 NVL(co.DISTRICT_NAME, '') AS district_name,
				 NVL(co.REGION_NAME, '') AS region_name,
				 NVL(co.FEDERAL_REGION_NAME, '') AS federal_region_name,
				 NVL(TO_CHAR(co.DAO_CODE), '') AS dao_code
			 FROM users u
			 LEFT JOIN linked_accounts la ON la.user_code = u.user_code
			 LEFT JOIN accounts ac ON ac.id = la.account_id
			 LEFT JOIN COMPANY co ON TO_CHAR(co.DAO_CODE) = u.branch_code
			 WHERE (
				 u.contact_phone = :1
				 OR u.customer_number = :1
				 OR u.user_code = :1
				 OR ac.account_number = :1
			 )
			 AND u.is_superapp_active = 1
			 AND (la.is_active = 1 OR la.is_active IS NULL)
			 FETCH FIRST 1 ROWS ONLY`

	row := c.db.QueryRowContext(ctx, query, number, number, number, number)
	var (
		id, userCode, email, customerNumber, fullName, phoneNumber, branchCode, gender, accountNumber string
		branchName, accountType, districtName, regionName, federalRegionName, daoCode                 string
		createdAt                                                                                     time.Time
		isBlocked, isSupperAppEnabled, isUssdEnabled                                                  int
	)
	err := row.Scan(
		&id, &userCode, &email, &customerNumber, &fullName, &phoneNumber, &branchCode,
		&gender, &isBlocked, &accountNumber, &isSupperAppEnabled, &isUssdEnabled,
		&branchName, &accountType, &districtName, &regionName, &federalRegionName, &daoCode,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, localization.ErrorCustomerNotFound
		}
		log.Errorf("[CustomerRepository][SearchCustomerByCIForAccountNumber] query failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	return &customer.CustomerListResponse{
		ID:                 id,
		UserID:             id,
		UserCode:           userCode,
		Email:              email,
		CustomerNumber:     customerNumber,
		FullName:           fullName,
		PhoneNumber:        phoneNumber,
		BranchCode:         branchCode,
		BranchName:         branchName,
		AccountType:        accountType,
		DistrictName:       districtName,
		RegionName:         regionName,
		FederalRegionName:  federalRegionName,
		DaoCode:            daoCode,
		Gender:             gender,
		CreatedAt:          createdAt.Format(time.RFC3339),
		IsBlocked:          isBlocked == 1,
		AccountNumber:      accountNumber,
		IsSupperAppEnabled: isSupperAppEnabled == 1,
		IsUssdEnabled:      isUssdEnabled == 1,
	}, nil
}

// BlockCustomerByUserCode implements [storage.CustomerRepository].
func (c *customerOracleRepository) BlockCustomerByUserCode(ctx context.Context, userCode string) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	query := `UPDATE users SET is_blocked = 1 WHERE user_code = :1`
	result, err := c.db.ExecContext(ctx, query, userCode)
	if err != nil {
		log.Errorf("[CustomerRepository][BlockCustomerByUserCode] failed to block customer with user_code %s: %v", userCode, err)
		return localization.ErrorUnexpectedError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("[CustomerRepository][BlockCustomerByUserCode] failed to get rows affected: %v", err)
		return localization.ErrorUnexpectedError
	}
	if rowsAffected == 0 {
		log.Errorf("[CustomerRepository][BlockCustomerByUserCode] no user found with user_code: %s", userCode)
		return localization.ErrorCustomerNotFound
	}
	log.Infof("[CustomerRepository][BlockCustomerByUserCode] successfully blocked customer with user_code: %s", userCode)
	return nil
}

func (c *customerOracleRepository) UNBlockCustomerByUserCode(ctx context.Context, userCode string) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	query := `UPDATE users SET is_blocked = 0 WHERE user_code = :1`
	result, err := c.db.ExecContext(ctx, query, userCode)
	if err != nil {
		log.Errorf("[CustomerRepository][BlockCustomerByUserCode] failed to unblock customer with user_code %s: %v", userCode, err)
		return localization.ErrorUnexpectedError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("[CustomerRepository][BlockCustomerByUserCode] failed to get rows affected: %v", err)
		return localization.ErrorUnexpectedError
	}
	if rowsAffected == 0 {
		log.Errorf("[CustomerRepository][BlockCustomerByUserCode] no user found with user_code: %s", userCode)
		return localization.ErrorCustomerNotFound
	}
	log.Infof("[CustomerRepository][BlockCustomerByUserCode] successfully unblocked customer with user_code: %s", userCode)
	return nil
}

// Update implements [storage.CustomerRepository].
func (c *customerOracleRepository) Update(ctx context.Context, id string, data member.User) error {
	panic("unimplemented")
}

func (c *customerOracleRepository) DisableCustomerByChannel(ctx context.Context, userCode, channel string) error {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	var query string
	switch channel {
	case "BOTH":
		query = `UPDATE users SET IS_SUPERAPP_ENABLED = 0, IS_USSD_ENABLED = 0 WHERE user_code = :1`
	case "SUPERAPP":
		query = `UPDATE users SET IS_SUPERAPP_ENABLED = 0 WHERE user_code = :1`
	case "USSD":
		query = `UPDATE users SET IS_USSD_ENABLED = 0 WHERE user_code = :1`
	default:
		log.Errorf("[CustomerRepository][DisableCustomerByChannel] unsupported channel: %s", channel)
		return fmt.Errorf("%s", localization.ErrorInvalidInputParameter.Code)
	}

	result, err := c.db.ExecContext(ctx, query, userCode)
	if err != nil {
		log.Errorf("[CustomerRepository][DisableCustomerByChannel] failed for user_code %s channel %s: %v", userCode, channel, err)
		return localization.ErrorUnexpectedError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("[CustomerRepository][DisableCustomerByChannel] failed to get rows affected: %v", err)
		return localization.ErrorUnexpectedError
	}
	if rowsAffected == 0 {
		log.Errorf("[CustomerRepository][DisableCustomerByChannel] no user found with user_code: %s", userCode)
		return localization.ErrorCustomerNotFound
	}
	log.Infof("[CustomerRepository][DisableCustomerByChannel] disabled channel %s for user_code: %s", channel, userCode)
	return nil
}

func (c *customerOracleRepository) SaveBarUnBarReason(ctx context.Context, entry *imodel.CustomerBarUnBarReason) error {
	panic("unimplemented")
}

func (c *customerOracleRepository) GetBarUnBarReasons(ctx context.Context, userID string) ([]*imodel.CustomerBarUnBarReason, error) {
	panic("unimplemented")
}

func NewCustomerOracleRepository(db DBTX, cfg config.VaultConfig, logger utils.Logger) storage.CustomerRepository {
	return &customerOracleRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}
