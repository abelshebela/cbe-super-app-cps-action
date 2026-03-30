package sqlc

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"database/sql"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type WalletStorage struct {
	db     DBTX
	logger utils.Logger
}

func NewWalletOracleRepository(db *sql.DB, log utils.Logger) storage.WalletOracleRepository {
	return &WalletStorage{
		db:     db,
		logger: log,
	}
}

// Create implements [storage.WalletOracleRepository].
func (q *WalletStorage) Create(ctx context.Context, wallet *model.WalletOracle) error {
	panic("unimplemented")
}

// Delete implements [storage.WalletOracleRepository].
func (q *WalletStorage) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// EnableOrDisable implements [storage.WalletOracleRepository].
func (q *WalletStorage) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	panic("unimplemented")
}

// Find implements [storage.WalletOracleRepository].
func (q *WalletStorage) Find(ctx context.Context, key string, value string) (*model.WalletOracle, error) {
	panic("unimplemented")
}

// FindAllWithPagination implements [storage.WalletOracleRepository].
func (q *WalletStorage) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.WalletOracle], error) {
	panic("unimplemented")
}

// FindAllWithPaginationForGRPC implements [storage.WalletOracleRepository].
func (q *WalletStorage) FindAllWithPaginationForGRPC(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.WalletOracle], error) {
	panic("unimplemented")
}

// FindByID implements [storage.WalletOracleRepository].
func (q *WalletStorage) FindByID(ctx context.Context, id string) (*model.WalletOracle, error) {
	panic("unimplemented")
}

// FindByIDForGRPC implements [storage.WalletOracleRepository].
func (q *WalletStorage) FindByIDForGRPC(ctx context.Context, id string) (*model.GRPCWallet, error) {
	panic("unimplemented")
}

// Update implements [storage.WalletOracleRepository].
func (q *WalletStorage) Update(ctx context.Context, id string, wallet *model.WalletOracle) error {
	panic("unimplemented")
}
