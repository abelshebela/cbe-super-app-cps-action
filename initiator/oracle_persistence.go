package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/kafka"
	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault"
	vaultamounttiers "cbe-super-app-cps-action/internal/storage/persistance/vault_amount_tiers"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	BankVault       storage.BankVaultRepository
	vaultCategory   storage.VaultCategoryRepository
	Sitota          storage.SitotaRepository
	Transaction     storage.TransactionRepository
	VaultAmountTier storage.VaultAmountTierRepository
}

func InitOraclePersistence(db *sql.DB, cfg *config.VaultConfig, clientOrchestrationProducer kafka.ClientOrchestrationProducer, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		BankVault:       bankvault.NewBankVaultRepository(db, log),
		vaultCategory:   vaultCategory.NewVaultCategoryRepository(db, cfg, clientOrchestrationProducer, log),
		Sitota:          sitota.NewSitotaRepository(db, log),
		Transaction:     transaction_repo.NewTransactionRepository(db, log),
		VaultAmountTier: vaultamounttiers.NewVaultAmountTierRepository(db, log),
	}
}
