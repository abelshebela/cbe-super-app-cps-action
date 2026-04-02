package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	access_list_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_oracle"
	access_list_segmentation_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation/oracle"
	amount_based_auth_oracle "cbe-super-app-cps-action/internal/storage/persistance/amount_based_auth_oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	budget_category_oracle "cbe-super-app-cps-action/internal/storage/persistance/budget_category_oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	// DB is the shared Oracle *sql.DB used by bulk-service relation proxy and other Oracle repos.
	DB *sql.DB

	BankOracle            storage.BankOracleRepository
	BankVault             storage.BankVaultRepository
	vaultCategory         storage.VaultCategoryRepository
	BudgetCategoryOracle  storage.BudgetCategoryOracleRepository
	AmountBasedAuthOracle storage.AmountBasedAuthOracleRepository
	AccessListOracle      storage.BulkServiceRepository
	Sitota                storage.SitotaRepository
	Transaction           storage.TransactionRepository
	AccessListSegmentaion storage.AccessListSegmentationRepositoryOracle
	BankOracle            storage.BankOracleRepository
	BankVault             storage.BankVaultRepository
	vaultCategory         storage.VaultCategoryRepository
	Sitota                storage.SitotaRepository
	Transaction           storage.TransactionRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, accessListSegmentationProducer *kafka.AccessListSegmentationProducer, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		DB: db,

		BankOracle:            sqlc.NewBankRepository(db, log),
		BankVault:             bankvault.NewBankVaultRepository(db, log),
		vaultCategory:         vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		BudgetCategoryOracle:  budget_category_oracle.NewBudgetCategoryOracleRepository(db, clientOrchestrationProducer, log),
		AmountBasedAuthOracle: amount_based_auth_oracle.NewAmountBasedAuthOracleRepository(db, log),
		AccessListOracle:      access_list_oracle.NewAccessListOracleRepository(db, log),
		Sitota:                sitota.NewSitotaRepository(db, log),
		Transaction:           transaction_repo.NewTransactionRepository(db, log),
		AccessListSegmentaion: access_list_segmentation_oracle.NewAccessListSegmentationOracle(db, *cfg, accessListSegmentationProducer, log),
		BankOracle:            sqlc.NewBankRepository(db, log),
		BankVault:             bankvault.NewBankVaultRepository(db, log),
		vaultCategory:         vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		Sitota:                sitota.NewSitotaRepository(db, log),
		Transaction:           transaction_repo.NewTransactionRepository(db, log),
	}
}
