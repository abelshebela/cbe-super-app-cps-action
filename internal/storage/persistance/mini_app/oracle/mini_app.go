package miniapp_oracle

import (
	constants "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/storage"
	"context"
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
	// var merchantID any = nil
	// if !miniApp.MerchantID.IsZero() {
	// 	merchantID = miniApp.MerchantID[:]
	// }

	q := querypkg.MiniAppsInsert

_, err := m.db.ExecContext(
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

	// fmt.Printf("Err :%v", err)
	if err != nil {
		return constants.ErrDatabaseError
	}
	return nil
}

func (m *miniAppOraclePersistence) Update(ctx context.Context, id string, miniApp *local_model.MiniApp) error {
	// if !isHexID(id) {
	// 	return constants.ErrInvalidID
	// }
	// var merchantID any = nil
	// if !miniApp.MerchantID.IsZero() {
	// 	merchantID = miniApp.MerchantID[:]
	// }
	q := querypkg.MiniAppsUpdate
	res, err := m.db.ExecContext(
		ctx, q,
		miniApp.CategoryID[:], miniApp.MerchantID, miniApp.AppName, miniApp.AppIcon, string(miniApp.AppType), miniApp.BannerImage, miniApp.URL, string(miniApp.AppViewType),
		boolToInt64(miniApp.IsFeatured), boolToInt64(miniApp.Enabled), miniApp.UpdatedAt,
		miniApp.Credential.MerchantAppID, miniApp.Credential.FabricAppID, miniApp.Credential.ShortCode, miniApp.Credential.AppSecret, miniApp.Credential.PrivateKey, miniApp.Credential.PublicKey,
		miniApp.ServiceID, id,
	)
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
	if !isHexID(id) {
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
	if !isHexID(id) {
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
