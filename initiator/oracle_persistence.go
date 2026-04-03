package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	access_list_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_oracle"
	account_block_repo "cbe-super-app-cps-action/internal/storage/persistance/account_block"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	amount_based_auth_oracle "cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth_oracle"
	budget_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/budget_category_oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"
	wallet_oracle "cbe-super-app-cps-action/internal/storage/persistance/wallet/oracle"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	// DB is the shared Oracle *sql.DB used by bulk-service relation proxy and other Oracle repos.
	DB *sql.DB

	BankOracle            storage.BankOracleRepository
	BankVault             storage.BankVaultRepository
  WalletOracle  storage.WalletOracleRepository
	vaultCategory         storage.VaultCategoryRepository
	BudgetCategoryOracle  storage.BudgetCategoryOracleRepository
	AmountBasedAuthOracle storage.AmountBasedAuthOracleRepository
	AccessListOracle      storage.BulkServiceRepository
	Sitota                storage.SitotaRepository
	Transaction           storage.TransactionRepository
  AccountBlock  storage.AccountBlockRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		DB: db,
		BankOracle:            sqlc.NewBankRepository(db, log),
		BankVault:             bankvault.NewBankVaultRepository(db, log),
    WalletOracle:  wallet_oracle.NewWalletOracleRepository(db, log),
		vaultCategory:         vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		BudgetCategoryOracle:  budget_category_oracle.NewBudgetCategoryOracleRepository(db, clientOrchestrationProducer, log),
		AmountBasedAuthOracle: amount_based_auth_oracle.NewAmountBasedAuthOracleRepository(db, log),
		AccessListOracle:      access_list_oracle.NewAccessListOracleRepository(db, log),
		Sitota:                sitota.NewSitotaRepository(db, log),
		Transaction:           transaction_repo.NewTransactionRepository(db, log),
    AccountBlock:  account_block_repo.NewAccountBlockRepository(db, log),
	}
}
