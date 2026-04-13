package customer_oracle

import (
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"database/sql"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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
	panic("unimplemented")
}

// FetchLinkedAccount implements [storage.CustomerRepository].
func (c *customerOracleRepository) FetchLinkedAccount(ctx context.Context, id string) ([]model.LinkedAccount, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*customer.CustomerListResponse], error) {
	panic("unimplemented")
}

// FindByID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindByID(ctx context.Context, id string) (*member.User, error) {
	panic("unimplemented")
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
	panic("unimplemented")
}

// FindCustomerDetailByID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerDetailByID(ctx context.Context, id string) (*customer.CustomerDetailResponse, error) {
	c.logger.Infof("[CustomerRepository][FindCustomerDetailByID] fetching customer detail by user_code: %s", id)

	// 1. Fetch user info by user_code
	userQuery := `
	SELECT
	  RAWTOHEX(u.id),
	  u.full_name,
	  u.gender,
	  u.contact_phone,
	  u.contact_email,
	  u.customer_number,
	  u.is_active,
	  u.birth_of_date
	FROM users u
	WHERE u.user_code = :1`

	var (
		userID, fullName, gender, phone, email, customerNumber string
		isActive                                               int
		birthOfDate                                            sql.NullTime
	)
	err := c.db.QueryRowContext(ctx, userQuery, id).Scan(&userID, &fullName, &gender, &phone, &email, &customerNumber, &isActive, &birthOfDate)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, localization.ErrorResourceNotFound
		}
		c.logger.Errorf("[CustomerRepository][FindCustomerDetailByID] user query failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	// 2. Fetch linked accounts with account and branch info
	linkedQuery := `
	SELECT
	  a.account_number,
	  a.account_holder_name,
	  a.account_type,
	  ab.code AS account_branch_code,
	  la.is_active,
	  ab.name AS account_branch_name
	FROM linked_accounts la
	JOIN accounts a ON a.id = la.account_id
	LEFT JOIN account_blocks ab ON ab.id = a.bank_id
	WHERE la.user_code = :1`

	rows, err := c.db.QueryContext(ctx, linkedQuery, id)
	if err != nil {
		c.logger.Errorf("[CustomerRepository][FindCustomerDetailByID] linked accounts query failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}
	defer rows.Close()

	var linkedAccounts []customer.LinkedAccount
	for rows.Next() {
		var accNum, accHolder, accType, branchCode, branchName sql.NullString
		var isActiveAcc int
		if err := rows.Scan(&accNum, &accHolder, &accType, &branchCode, &isActiveAcc, &branchName); err != nil {
			c.logger.Errorf("[CustomerRepository][FindCustomerDetailByID] scan failed: %v", err)
			return nil, localization.ErrorUnexpectedError
		}
		linkedAccounts = append(linkedAccounts, customer.LinkedAccount{
			AccountNumber:     accNum.String,
			AccountHolderName: accHolder.String,
			AccountType:       accType.String,
			AccountBranchCode: branchCode.String,
			IsActive:          isActiveAcc == 1,
			AccountBranchName: branchName.String,
		})
	}

	response := &customer.CustomerDetailResponse{
		ID:            userID,
		LinkedAccount: linkedAccounts,
		PersonalInfo: customer.PersonalInfo{
			FullName:       fullName,
			Gender:         gender,
			PhoneNumber:    phone,
			Email:          email,
			CustomerNumber: customerNumber,
			IsActivated:    isActive == 1,
			DateOfBirth:    "",
		},
	}
	if birthOfDate.Valid {
		response.PersonalInfo.DateOfBirth = birthOfDate.Time.Format("2006-01-02")
	}

	return response, nil
}

// FindCustomerLinkedAccountByUserID implements [storage.CustomerRepository].
func (c *customerOracleRepository) FindCustomerLinkedAccountByUserID(ctx context.Context, userID string) (*model.LinkedAccount, error) {
	panic("unimplemented")
}

// SearchCustomerByCIForAccountNumber implements [storage.CustomerRepository].
func (c *customerOracleRepository) SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer.CustomerListResponse, error) {
	c.logger.Infof("[CustomerRepository][SearchCustomerByCIForAccountNumber] searching members by value: %s", number)

	// Oracle SQL: join USERS and LINKED_ACCOUNTS, search by phone, customer number, user code, or account id (as hex string)
	query := `
			 SELECT
				 RAWTOHEX(u.id) AS id,
				 u.user_code,
				 u.contact_email,
				 u.customer_number,
				 u.full_name,
				 u.contact_phone,
				 '' AS branch_code,
				 u.gender,
				 u.created_at,
				 u.is_blocked,
				 ac.account_number
			 FROM users u
			 LEFT JOIN linked_accounts la ON la.user_code = u.user_code
			 Left JOIN accounts ac ON ac.id = la.account_id 
			 WHERE (
				 u.contact_phone = ':1'
				 OR u.customer_number = ':1'
				 OR u.user_code = ':1'
				 OR ac.account_number = ':1'
			 )
			 AND u.is_active = 1
			 AND (la.is_active = 1 OR la.is_active IS NULL)
			 FETCH FIRST 1 ROWS ONLY`

	row := c.db.QueryRowContext(ctx, query, number, number, number, number)
	var (
		id, userCode, email, customerNumber, fullName, phoneNumber, branchCode, gender, accountNumber string
		createdAt                                                                                     time.Time
		isBlocked                                                                                     int
	)
	err := row.Scan(&id, &userCode, &email, &customerNumber, &fullName, &phoneNumber, &branchCode, &gender, &createdAt, &isBlocked, &accountNumber)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, localization.ErrorCustomerNotFound
		}
		c.logger.Errorf("[CustomerRepository][SearchCustomerByCIForAccountNumber] query failed: %v", err)
		return nil, localization.ErrorUnexpectedError
	}

	return &customer.CustomerListResponse{
		ID:             id,
		UserID:         id,
		UserCode:       userCode,
		Email:          email,
		CustomerNumber: customerNumber,
		FullName:       fullName,
		PhoneNumber:    phoneNumber,
		BranchCode:     branchCode,
		Gender:         gender,
		CreatedAt:      createdAt.Format(time.RFC3339),
		IsBlocked:      isBlocked == 1,
		AccountNumber:  accountNumber,
	}, nil
}

// Update implements [storage.CustomerRepository].
func (c *customerOracleRepository) Update(ctx context.Context, id string, data member.User) error {
	panic("unimplemented")
}

func NewCustomerOracleRepository(db DBTX, cfg config.VaultConfig, logger utils.Logger) storage.CustomerRepository {
	return &customerOracleRepository{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}
