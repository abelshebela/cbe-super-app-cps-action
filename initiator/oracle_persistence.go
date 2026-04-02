package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	access_list_segmentation_oracle "cbe-super-app-cps-action/internal/storage/persistance/access_list_segmentation/oracle"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	AccessListSegmentaion storage.AccessListSegmentationRepositoryOracle
	BankOracle            storage.BankOracleRepository
	BankVault             storage.BankVaultRepository
	vaultCategory         storage.VaultCategoryRepository
	Sitota                storage.SitotaRepository
	Transaction           storage.TransactionRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, accessListSegmentationProducer *kafka.AccessListSegmentationProducer, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		AccessListSegmentaion: access_list_segmentation_oracle.NewAccessListSegmentationOracle(db, *cfg, accessListSegmentationProducer, log),
		BankOracle:            sqlc.NewBankRepository(db, log),
		BankVault:             bankvault.NewBankVaultRepository(db, log),
		vaultCategory:         vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		Sitota:                sitota.NewSitotaRepository(db, log),
		Transaction:           transaction_repo.NewTransactionRepository(db, log),
	}
}
