package miniapp_oracle

import (
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"
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
	_, err := m.db.ExecContext(ctx, q, merchant.MerchantName, merchant.MerchantCode, fmt.Sprintf("%s", merchant.SettlementMethod), merchant.BankAccountNumber, boolToInt64(merchant.Enabled), boolToInt64(merchant.IsDeleted), merchant.CreatedAt, merchant.UpdatedAt, merchant.PhoneNumber, merchant.Email)
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
	if !isHexID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppMerchantUpdate
	res, err := m.db.ExecContext(ctx, q, merchant.MerchantName, merchant.MerchantCode, string(merchant.SettlementMethod), merchant.BankAccountNumber, merchant.UpdatedAt,merchant.PhoneNumber,merchant.Email, id)
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
func (m *miniAppMerchantOraclePersistence) FindByID(ctx context.Context, id string) (*local_model.MiniAppMerchant, error){
	return &local_model.MiniAppMerchant{},nil
}
