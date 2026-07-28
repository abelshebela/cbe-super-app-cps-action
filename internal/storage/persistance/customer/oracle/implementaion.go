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
	"strings"
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

//		return response, nil
//	}
//

// FindCustomerDetailByID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerDetailByID(ctx context.Context, id string) (*customer.CustomerDetailResponse, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[CustomerRepository][FindCustomerDetailByID] fetching customer detail by user_code: %s", id)

	query := `
		SELECT
			RAWTOHEX(u.id),
			u.full_name,
			u.gender,
			u.contact_phone,
			u.contact_email,
			u.customer_number,
			NVL(u.is_superapp_active, 0),
			NVL(u.is_ussd_active, 0),
			NVL(u.is_superapp_enabled, 0),
			NVL(u.is_ussd_enabled, 0),
			NVL(u.is_blocked, 0),
			u.birth_of_date,
			NVL(u.is_self_activated, 0),
			u.superapp_role,
			u.created_at,
			u.expiry_at,
			NVL(a.account_holder_name, '') AS account_holder_name,
			NVL(a.account_type, '') AS account_type,
			NVL(a.account_number, '') AS account_number,
			NVL(ab.BRANCH_NAME, '') AS branch_name,
			NVL(ab.BRANCH_CODE, '') AS branch_code,
			NVL(la.is_active, 0)
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
		userID, fullName, gender, phone, email, customerNumber, superAppRole                         string
		isSuperAppActive, isUSSDActive, isSuperAppEnabled, isUSSDEnabled, isBlocked, isSelfActivated int
		birthOfDate                                                                                  sql.NullTime
		createdAt, expiryAt                                                                          time.Time
		linkedAccounts                                                                               []customer.LinkedAccount
		fetchedFirstRow                                                                              bool
	)

	for rows.Next() {
		var (
			accHolder, accType, accNum, branchName, branchCode sql.NullString
			isActiveAcc                                        int

			rowUserID, rowFullName, rowGender, rowPhone, rowEmail, rowCustomerNumber, rowSuperAppRole                      string
			rowIsSuperAppActive, rowIsUSSDActive, rowIsSuperAppEnabled, rowIsUSSDEnabled, rowIsBlocked, rowIsSelfActivated int
			rowBirthOfDate                                                                                                 sql.NullTime
			rowCreatedAt, rowExpiryAt                                                                                      sql.NullTime
		)

		if err := rows.Scan(
			&rowUserID,
			&rowFullName,
			&rowGender,
			&rowPhone,
			&rowEmail,
			&rowCustomerNumber,
			&rowIsSuperAppActive,
			&rowIsUSSDActive,
			&rowIsSuperAppEnabled,
			&rowIsUSSDEnabled,
			&rowIsBlocked,
			&rowBirthOfDate,
			&rowIsSelfActivated,
			&rowSuperAppRole,
			&rowCreatedAt,
			&rowExpiryAt,
			&accHolder,
			&accType,
			&accNum,
			&branchName,
			&branchCode,
			&isActiveAcc,
		); err != nil {
			log.Errorf("[CustomerRepository][FindCustomerDetailByID] scan failed: %v", err)
			return nil, localization.ErrorUnexpectedError
		}

		if !fetchedFirstRow {
			userID = rowUserID
			fullName = rowFullName
			gender = rowGender
			phone = rowPhone
			email = rowEmail
			customerNumber = rowCustomerNumber
			isSuperAppActive = rowIsSuperAppActive
			isUSSDActive = rowIsUSSDActive
			isSuperAppEnabled = rowIsSuperAppEnabled
			isUSSDEnabled = rowIsUSSDEnabled
			isBlocked = rowIsBlocked
			birthOfDate = rowBirthOfDate
			isSelfActivated = rowIsSelfActivated
			superAppRole = rowSuperAppRole

			if rowCreatedAt.Valid {
				createdAt = rowCreatedAt.Time
			}

			if rowExpiryAt.Valid {
				expiryAt = rowExpiryAt.Time
			}

			fetchedFirstRow = true
		}

		linkedAccounts = append(linkedAccounts, customer.LinkedAccount{
			AccountNumber:     accNum.String,
			AccountHolderName: accHolder.String,
			AccountType:       accType.String,
			AccountBranchCode: branchCode.String,
			AccountBranchName: branchName.String,
			IsActive:          isActiveAcc == 1 || isUSSDActive == 1,
		})
	}

	if err := rows.Err(); err != nil {
		log.Errorf("[CustomerRepository][FindCustomerDetailByID] rows iteration failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	if !fetchedFirstRow {
		return nil, localization.ErrorResourceNotFound
	}

	dateOfBirth := ""
	if birthOfDate.Valid {
		dateOfBirth = birthOfDate.Time.Format("2006-01-02")
	}

	return &customer.CustomerDetailResponse{
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
			IsUSSDEnabled:       isUSSDEnabled == 1,
			IsBlocked:           isBlocked == 1,
			DateOfBirth:         dateOfBirth,
			SuperAppID:          "SA" + customerNumber,
			SupperAppRole:       superAppRole,
			IsSelfActivated:     isSelfActivated == 1,
			ActivationDate:      createdAt,
			ExpiryDate:          expiryAt,
		},
	}, nil
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

func (c *customerOracleRepository) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer.CustomerListResponse, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	log.Infof("[CustomerRepository][SearchCustomerByCIForAccountNumber] searching members by value: %s", number)

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
				 u.is_self_activated,
				 u.superapp_role,
				 u.is_blocked,
				 u.created_at,
				 u.expiry_at,
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
			 FETCH FIRST 1 ROWS ONLY`

	row := c.db.QueryRowContext(ctx, query, number, number, number, number)
	var (
		id, userCode, email, customerNumber, fullName, phoneNumber, branchCode, gender, accountNumber string
		branchName, accountType, districtName, regionName, federalRegionName, daoCode, superappRole   string
		createdAt, expiryAt                                                                           time.Time
		isBlocked, isSelfActivated, isSupperAppEnabled, isUssdEnabled                                 int
	)
	err := row.Scan(
		&id, &userCode, &email, &customerNumber, &fullName, &phoneNumber, &branchCode,
		&gender, &isSelfActivated, &superappRole, &isBlocked, &createdAt, &expiryAt, &accountNumber, &isSupperAppEnabled, &isUssdEnabled,
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
		SuperAppID:         "SA" + customerNumber,
		UserCode:           userCode,
		Email:              email,
		CustomerNumber:     customerNumber,
		FullName:           fullName,
		PhoneNumber:        phoneNumber,
		IsSelfActivated:    isSelfActivated == 1,
		SupperAppRole:      superappRole,
		ActivationDate:     createdAt,
		ExpiryDate:         expiryAt,
		BranchCode:         branchCode,
		BranchName:         branchName,
		AccountType:        accountType,
		DistrictName:       districtName,
		RegionName:         regionName,
		FederalRegionName:  federalRegionName,
		DaoCode:            daoCode,
		Gender:             gender,
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

func (c *customerOracleRepository) GetCustomerByCifNumber(ctx context.Context, cif string) (customer.CustomerByCIFResponse, error) {
	log := local_util.LoggerFromCtx(ctx, c.logger)

	cif = strings.TrimSpace(cif)
	if cif == "" {
		return customer.CustomerByCIFResponse{}, localization.ErrorNoDataProvided
	}

	query := `
		SELECT
			RAWTOHEX(ID),
			USER_CODE,
			USERNAME,
			BRANCH_CODE,
			CONTACT_EMAIL,
			CONTACT_PHONE,
			CUSTOMER_NUMBER,
			FIRST_NAME,
			LAST_NAME,
			MIDDLE_NAME,
			FULL_NAME,
			GENDER,
			BIRTH_OF_DATE,
			AVATAR,
			SECTOR,
			OWNERSHIP,
			PUSH_TOKEN,
			INDUSTRY,
			LANGUAGE,
			PIN,
			PIN_HISTORY,
			EXCLUDED_ACCESS_LIST,
			FAILED_LOGIN_ATTEMPT,
			IS_LOCKED,
			IS_BUDGET_ENABLED,
			IS_SUPERAPP_ACTIVE,
			IS_SUPERAPP_ENABLED,
			IS_USSD_ACTIVE,
			IS_USSD_ENABLED,
			IS_BLOCKED,
			LOCK_EXPIRES_AT,
			PIN_CREATED_AT,
			CREATED_AT,
			LAST_LOGIN_AT,
			LAST_MODIFIED_AT,
			EXPIRY_AT,
			CUSTOMER_SEGMENTATION,
			ACCOUNT_TYPE,
			CUSTOMER_GROUP,
			CUSTOMER_SUBSEGMENT,
			SUPERAPP_ROLE,
			IS_RESET_PIN,
			IS_TERMINATED,
			TERMINATION_COUNTER,
			HAS_CHANGE,
			IS_SELF_ACTIVATED,
			SELF_ACTIVATION_RULE_EXPIRATION,
			IS_SELF_ACTIVATION_RULE_OVERIDED
		FROM USERS
		WHERE CUSTOMER_NUMBER = :1`

	var (
		id, userCode, username, branchCode, contactEmail, contactPhone, customerNumber             sql.NullString
		firstName, lastName, middleName, fullName, gender, avatar, sector, ownership, pushToken    sql.NullString
		industry, language, pin, pinHistory, excludedAccessList, customerSegmentation, accountType sql.NullString
		customerGroup, customerSubsegment, superappRole, hasChange                                 sql.NullString
		failedLoginAttempt, isLocked, isBudgetEnabled, isSuperappActive, isSuperappEnabled         sql.NullInt64
		isUssdActive, isUssdEnabled, isBlocked, isResetPin, isTerminated, isSelfActivated          sql.NullInt64
		terminationCounter, isSelfActivationRuleOverided                                           sql.NullInt64
		birthOfDate, lockExpiresAt, pinCreatedAt, createdAt, lastLoginAt, lastModifiedAt, expiryAt sql.NullTime
		selfActivationRuleExpiration                                                               sql.NullTime
	)

	err := c.db.QueryRowContext(ctx, query, cif).Scan(
		&id, &userCode, &username, &branchCode, &contactEmail, &contactPhone, &customerNumber,
		&firstName, &lastName, &middleName, &fullName, &gender, &birthOfDate, &avatar, &sector,
		&ownership, &pushToken, &industry, &language, &pin, &pinHistory, &excludedAccessList,
		&failedLoginAttempt, &isLocked, &isBudgetEnabled, &isSuperappActive, &isSuperappEnabled,
		&isUssdActive, &isUssdEnabled, &isBlocked, &lockExpiresAt, &pinCreatedAt, &createdAt,
		&lastLoginAt, &lastModifiedAt, &expiryAt, &customerSegmentation, &accountType, &customerGroup,
		&customerSubsegment, &superappRole, &isResetPin, &isTerminated, &terminationCounter,
		&hasChange, &isSelfActivated, &selfActivationRuleExpiration, &isSelfActivationRuleOverided,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Errorf("[CustomerRepository][GetCustomerByCifNumber] customer not found for cif: %s", cif)
			return customer.CustomerByCIFResponse{}, localization.ErrorResourceNotFound
		}
		log.Errorf("[CustomerRepository][GetCustomerByCifNumber] query failed: %v", err)
		return customer.CustomerByCIFResponse{}, local_util.HandleDBError(err)
	}

	toString := func(v sql.NullString) string {
		if v.Valid {
			return v.String
		}
		return ""
	}
	toTimePtr := func(v sql.NullTime) *time.Time {
		if v.Valid {
			t := v.Time
			return &t
		}
		return nil
	}
	toBool := func(v sql.NullInt64) bool {
		return v.Valid && v.Int64 == 1
	}
	toInt := func(v sql.NullInt64) int {
		if v.Valid {
			return int(v.Int64)
		}
		return 0
	}

	response := customer.CustomerByCIFResponse{
		ID:                           toString(id),
		UserCode:                     toString(userCode),
		Username:                     toString(username),
		BranchCode:                   toString(branchCode),
		ContactEmail:                 toString(contactEmail),
		ContactPhone:                 toString(contactPhone),
		CustomerNumber:               toString(customerNumber),
		FirstName:                    toString(firstName),
		LastName:                     toString(lastName),
		MiddleName:                   toString(middleName),
		FullName:                     toString(fullName),
		Gender:                       toString(gender),
		BirthOfDate:                  toTimePtr(birthOfDate),
		Avatar:                       toString(avatar),
		Sector:                       toString(sector),
		Ownership:                    toString(ownership),
		PushToken:                    toString(pushToken),
		Industry:                     toString(industry),
		Language:                     toString(language),
		Pin:                          toString(pin),
		PinHistory:                   toString(pinHistory),
		ExcludedAccessList:           toString(excludedAccessList),
		FailedLoginAttempt:           toInt(failedLoginAttempt),
		IsLocked:                     toBool(isLocked),
		IsBudgetEnabled:              toBool(isBudgetEnabled),
		IsSuperappActive:             toBool(isSuperappActive),
		IsSuperappEnabled:            toBool(isSuperappEnabled),
		IsUssdActive:                 toBool(isUssdActive),
		IsUssdEnabled:                toBool(isUssdEnabled),
		IsBlocked:                    toBool(isBlocked),
		LockExpiresAt:                toTimePtr(lockExpiresAt),
		PinCreatedAt:                 toTimePtr(pinCreatedAt),
		CreatedAt:                    createdAt.Time,
		LastLoginAt:                  toTimePtr(lastLoginAt),
		LastModifiedAt:               lastModifiedAt.Time,
		ExpiryAt:                     toTimePtr(expiryAt),
		CustomerSegmentation:         toString(customerSegmentation),
		AccountType:                  toString(accountType),
		CustomerGroup:                toString(customerGroup),
		CustomerSubsegment:           toString(customerSubsegment),
		SuperappRole:                 toString(superappRole),
		IsResetPin:                   toBool(isResetPin),
		IsTerminated:                 toBool(isTerminated),
		TerminationCounter:           toInt(terminationCounter),
		HasChange:                    toString(hasChange),
		IsSelfActivated:              toBool(isSelfActivated),
		SelfActivationRuleExpiration: toTimePtr(selfActivationRuleExpiration),
		IsSelfActivationRuleOverided: toBool(isSelfActivationRuleOverided),
	}

	return response, nil
}

func NewCustomerOracleRepository(db DBTX, cfg config.VaultConfig, logger utils.Logger) storage.CustomerRepository {
	return &customerOracleRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}
