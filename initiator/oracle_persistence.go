package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	"cbe-super-app-cps-action/internal/storage/persistance/bank/gen/sqlc"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	services_repo "cbe-super-app-cps-action/internal/storage/persistance/services"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	BankOracle          storage.BankOracleRepository
	BankVault           storage.BankVaultRepository
	vaultCategory       storage.VaultCategoryRepository
	Sitota              storage.SitotaRepository
	Transaction         storage.TransactionRepository
	ServicesPersistence storage.ServicesRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, redisRepository storage.RedisRepository, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		BankOracle:          sqlc.NewBankRepository(db, log),
		BankVault:           bankvault.NewBankVaultRepository(db, log),
		vaultCategory:       vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		Sitota:              sitota.NewSitotaRepository(db, log),
		Transaction:         transaction_repo.NewTransactionRepository(db, log),
		ServicesPersistence: services_repo.NewServicesRepository(db, cfg, clientOrchestrationProducer, redisRepository, log),
	}
}
