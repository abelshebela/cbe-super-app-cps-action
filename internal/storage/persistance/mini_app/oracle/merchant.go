package miniapp_oracle

import (
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	querypkg "cbe-super-app-cps-action/internal/storage/persistance/mini_app/oracle/query"

	// sharedModel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	local_model "cbe-super-app-cps-action/internal/constants/model"

	constants "cbe-super-app-cps-action/internal/constants/localization"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	// mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"
)
type miniAppMerchantOraclePersistence struct {
	logger shared_utils.Logger
	db     *sql.DB
}

func NewMiniAppMerchantOracleRepository(logger shared_utils.Logger, db *sql.DB) storage.MiniAppMerchant {
	return &miniAppMerchantOraclePersistence{logger: logger, db: db}
}


func (m *miniAppMerchantOraclePersistence) Create(ctx context.Context, merchant *local_model.MiniAppMerchant) (*local_model.MiniAppMerchant, error) {
	q := querypkg.MiniAppMerchantInsert
	var err error
	if merchant.SettlementMethod == local_model.SettlementMethodDirect || merchant.SettlementMethod == local_model.SettlementMethodGL{

	_, err = m.db.ExecContext(ctx, q, merchant.MerchantName, merchant.MerchantCode, fmt.Sprintf("%s", merchant.SettlementMethod), merchant.BankAccountNumber, boolToInt64(merchant.Enabled), boolToInt64(merchant.IsDeleted), merchant.CreatedAt, merchant.UpdatedAt, merchant.PhoneNumber, merchant.Email)
	}else{
		_, err = m.db.ExecContext(ctx, q, merchant.MerchantName, merchant.MerchantCode, string(merchant.SettlementMethod), nil, boolToInt64(merchant.Enabled), boolToInt64(merchant.IsDeleted), merchant.CreatedAt, merchant.UpdatedAt, merchant.PhoneNumber, merchant.Email)
	}
	if err != nil {
		m.logger.Errorf("[miniapp_merchant CREATE] got error while creating merchant %w", err)
		return nil, constants.ErrDatabaseError
	}


	return &local_model.MiniAppMerchant{
		ID:                strings.ToUpper(merchant.ID),
		MerchantName:      merchant.MerchantName,
		MerchantCode:      merchant.MerchantCode,
		SettlementMethod:  local_model.SettlementMethod(merchant.SettlementMethod),
		BankAccountNumber: merchant.BankAccountNumber,
		Enabled:           merchant.Enabled,
		IsDeleted:         merchant.IsDeleted,
		CreatedAt:         merchant.CreatedAt,
		UpdatedAt:         merchant.UpdatedAt,
	}, nil
}

func (m *miniAppMerchantOraclePersistence) Update(ctx context.Context, id string, merchant *local_model.MiniAppMerchant) error {
	if !IsValidRaw16ID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppMerchantUpdate
	var res sql.Result
	var  err error
	if merchant.SettlementMethod == local_model.SettlementMethodDirect || merchant.SettlementMethod == local_model.SettlementMethodGL{

	res, err = m.db.ExecContext(ctx, q, merchant.MerchantName, merchant.MerchantCode, string(merchant.SettlementMethod), merchant.BankAccountNumber, merchant.UpdatedAt,merchant.PhoneNumber,merchant.Email, id)
	} else {
		res, err = m.db.ExecContext(ctx, q, merchant.MerchantName, merchant.MerchantCode, string(merchant.SettlementMethod), nil, merchant.UpdatedAt,merchant.PhoneNumber,merchant.Email, id)
	}
	if err != nil {
		m.logger.Errorf("[miniapp_merchant UPDATE] got error while update merchant %w", err)

		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrMiniAppMerchantNotFound
	}

	return nil
}

func (m *miniAppMerchantOraclePersistence) Delete(ctx context.Context, id string) error {
	if !isHexID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppMerchantSoftDelete
	res, err := m.db.ExecContext(ctx, q, time.Now(), time.Now(), id)
	if err != nil {
		m.logger.Errorf("[miniapp_merchant Delete] got error while Delete merchant %w", err)

		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrMiniAppMerchantNotFound
	}
	return nil
}

func (m *miniAppMerchantOraclePersistence) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	if !isHexID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppMerchantToggleEnabled
	res, err := m.db.ExecContext(ctx, q, boolToInt64(enable), time.Now(), id)
	if err != nil {
		m.logger.Errorf("[miniapp_merchant status update ] got error while  merchant %w", err)

		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrMiniAppMerchantNotFound
	}
	return nil
}

//for test 

func (m *miniAppMerchantOraclePersistence) FindByID(ctx context.Context, id string) (*local_model.MiniAppMerchant, error) {
	if !isHexID(id) {
		return nil, constants.ErrInvalidID
	}
	q := querypkg.MiniAppMerchantByID
	var (
		ID                                         string
		name, code, settlement, bank, phone, email sql.NullString
		enabled, deleted                           sql.NullInt64
		createdAt, updatedAt, deletedAt            sql.NullTime
	)
	err := m.db.QueryRowContext(ctx, q, id).Scan(&ID, &name, &code, &settlement, &bank, &enabled, &deleted, &createdAt, &updatedAt, &deletedAt, &phone, &email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrMiniAppMerchantNotFound
		}
		return nil, constants.ErrDatabaseError
	}
	var deletedPtr *time.Time
	if deletedAt.Valid {
		t := deletedAt.Time
		deletedPtr = &t
	}
	return &local_model.MiniAppMerchant{
		ID:                ID,
		MerchantName:      nullString(name),
		MerchantCode:      nullString(code),
		PhoneNumber:       nullString(phone),
		Email:             nullString(email),
		SettlementMethod:  local_model.SettlementMethod(nullString(settlement)),
		BankAccountNumber: nullString(bank),
		Enabled:           enabled.Valid && enabled.Int64 == 1,
		IsDeleted:         deleted.Valid && deleted.Int64 == 1,
		CreatedAt:         nullTime(createdAt),
		UpdatedAt:         nullTime(updatedAt),
		DeletedAt:         deletedPtr,
	}, nil
}
func nullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}


func nullTime(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}
