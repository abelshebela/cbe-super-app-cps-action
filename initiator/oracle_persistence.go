package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"
	wallet_oracle "cbe-super-app-cps-action/internal/storage/persistance/wallet/oracle"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	BankOracle    storage.BankOracleRepository
	WalletOracle  storage.WalletOracleRepository
	BankVault     storage.BankVaultRepository
	vaultCategory storage.VaultCategoryRepository
	Sitota        storage.SitotaRepository
	Transaction   storage.TransactionRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		BankOracle:    sqlc.NewBankRepository(db, log),
		WalletOracle:  wallet_oracle.NewWalletOracleRepository(db, log),
		BankVault:     bankvault.NewBankVaultRepository(db, log),
		vaultCategory: vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		Sitota:        sitota.NewSitotaRepository(db, log),
		Transaction:   transaction_repo.NewTransactionRepository(db, log),
	}
}
