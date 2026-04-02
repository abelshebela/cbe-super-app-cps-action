package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	budget_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/budget_category_oracle"
	amount_based_auth_oracle "cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth_oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	BankOracle              storage.BankOracleRepository
	BankVault               storage.BankVaultRepository
	vaultCategory           storage.VaultCategoryRepository
	BudgetCategoryOracle    storage.BudgetCategoryOracleRepository
	AmountBasedAuthOracle  storage.AmountBasedAuthOracleRepository
	Sitota                  storage.SitotaRepository
	Transaction             storage.TransactionRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		BankOracle:           sqlc.NewBankRepository(db, log),
		BankVault:            bankvault.NewBankVaultRepository(db, log),
		vaultCategory:        vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		BudgetCategoryOracle: budget_category_oracle.NewBudgetCategoryOracleRepository(db, clientOrchestrationProducer, log),
		AmountBasedAuthOracle: amount_based_auth_oracle.NewAmountBasedAuthOracleRepository(db, log),
		Sitota:               sitota.NewSitotaRepository(db, log),
		Transaction:          transaction_repo.NewTransactionRepository(db, log),
	}
}
