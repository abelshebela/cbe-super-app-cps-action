package miniapp_oracle

import (
	constants "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"encoding/hex"
	"strings"
	"time"

	local_model "cbe-super-app-cps-action/internal/constants/model"
	querypkg "cbe-super-app-cps-action/internal/storage/persistance/mini_app/oracle/query"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"database/sql"
)



type miniAppOraclePersistence struct {
	logger shared_utils.Logger
	db     *sql.DB
}

func NewMiniAppOracleRepository(logger shared_utils.Logger, db *sql.DB) storage.MiniAppRepository {
	return &miniAppOraclePersistence{logger: logger, db: db}
}


func (m *miniAppOraclePersistence) Create(ctx context.Context, miniApp *local_model.MiniApp) error {
	
	var err error
	if miniApp.AppType == local_model.NonFinancial{
	q := querypkg.MiniAppsInsert

	_, err = m.db.ExecContext(
		ctx, q,
		miniApp.CategoryID,
		nil,
		miniApp.AppName,
		miniApp.AppIcon,
		string(miniApp.AppType),
		miniApp.BannerImage,
		miniApp.URL,
		string(miniApp.AppViewType),
		boolToInt64(miniApp.IsFeatured),
		boolToInt64(miniApp.Enabled),
		boolToInt64(miniApp.IsDeleted),
		miniApp.CreatedAt,
		miniApp.UpdatedAt,
		miniApp.Credential.MerchantAppID,
		miniApp.Credential.FabricAppID,
		miniApp.Credential.ShortCode,
		miniApp.Credential.AppSecret,
		miniApp.Credential.PrivateKey,
		miniApp.Credential.PublicKey,
		miniApp.ServiceID, // SERVICE_ID
	)}else{
		q := querypkg.MiniAppsInsert

		_, err = m.db.ExecContext(
			ctx, q,
			miniApp.CategoryID,
			miniApp.MerchantID,
			miniApp.AppName,
			miniApp.AppIcon,
			string(miniApp.AppType),
			miniApp.BannerImage,
			miniApp.URL,
			string(miniApp.AppViewType),
			boolToInt64(miniApp.IsFeatured),
			boolToInt64(miniApp.Enabled),
			boolToInt64(miniApp.IsDeleted),
			miniApp.CreatedAt,
			miniApp.UpdatedAt,
			miniApp.Credential.MerchantAppID,
			miniApp.Credential.FabricAppID,
			miniApp.Credential.ShortCode,
			miniApp.Credential.AppSecret,
			miniApp.Credential.PrivateKey,
			miniApp.Credential.PublicKey,
			miniApp.ServiceID, // SERVICE_ID
		)
	}

	if err != nil {
		m.logger.Errorf("[create mini app cps]failed to create mini app db error for app type: %v", "err", miniApp.AppType, err)

		return constants.ErrDatabaseError
	}
	return nil
}

func (m *miniAppOraclePersistence) Update(ctx context.Context, id string, miniApp *local_model.MiniApp) error {
	if !IsValidRaw16ID(id) {
		return constants.ErrInvalidID
	}
	var res sql.Result
	var err error
	if local_model.AppType((strings.ToUpper(string(miniApp.AppType)))) == local_model.NonFinancial {

	q := querypkg.MiniAppsUpdate
	res, err = m.db.ExecContext(
		ctx, q,
		miniApp.CategoryID[:], nil, miniApp.AppName, miniApp.AppIcon, string(miniApp.AppType), miniApp.BannerImage, miniApp.URL, string(miniApp.AppViewType),
		boolToInt64(miniApp.IsFeatured), boolToInt64(miniApp.Enabled), miniApp.UpdatedAt,
		miniApp.Credential.MerchantAppID, miniApp.Credential.FabricAppID, miniApp.Credential.ShortCode, miniApp.Credential.AppSecret, miniApp.Credential.PrivateKey, miniApp.Credential.PublicKey,
		miniApp.ServiceID, id,
	)} else{
			q := querypkg.MiniAppsUpdate
	res, err = m.db.ExecContext(
		ctx, q,
		miniApp.CategoryID[:], miniApp.MerchantID, miniApp.AppName, miniApp.AppIcon, string(miniApp.AppType), miniApp.BannerImage, miniApp.URL, string(miniApp.AppViewType),
		boolToInt64(miniApp.IsFeatured), boolToInt64(miniApp.Enabled), miniApp.UpdatedAt,
		miniApp.Credential.MerchantAppID, miniApp.Credential.FabricAppID, miniApp.Credential.ShortCode, miniApp.Credential.AppSecret, miniApp.Credential.PrivateKey, miniApp.Credential.PublicKey,
		miniApp.ServiceID, id,
	)
	}
	if err != nil {
		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrMiniAppNotFound
	}
	return nil
}

func (m *miniAppOraclePersistence) Delete(ctx context.Context, id string) error {
	if !IsValidRaw16ID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppsSoftDelete
	res, err := m.db.ExecContext(ctx, q, time.Now(), time.Now(), id)
	if err != nil {
		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrMiniAppNotFound
	}
	return nil
}

func (m *miniAppOraclePersistence) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	if !IsValidRaw16ID(id) {
		return constants.ErrInvalidID
	}
	q := querypkg.MiniAppsToggleEnable
	res, err := m.db.ExecContext(ctx, q, boolToInt64(enable), time.Now(), id)
	if err != nil {
		return constants.ErrDatabaseError
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return constants.ErrMiniAppNotFound
	}
	return nil
}

func (m *miniAppOraclePersistence) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return constants.ErrDatabaseError
	}
	txCtx := context.WithValue(ctx, "oracle_tx", tx)
	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return constants.ErrDatabaseError
	}
	return nil
}


func IsValidRaw16ID(id string) bool {
	if len(id) != 32 {
		return false
	}

	_, err := hex.DecodeString(id)
	return err == nil
}