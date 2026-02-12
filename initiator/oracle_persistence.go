package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"
	"cbe-super-app-cps-action/internal/storage/persistance/sitota"
	transaction_repo "cbe-super-app-cps-action/internal/storage/persistance/transaction"
	vaultamounttiers "cbe-super-app-cps-action/internal/storage/persistance/vault_amount_tiers"
	vaultCategory "cbe-super-app-cps-action/internal/storage/persistance/vault_category"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	BankVault       storage.BankVaultRepository
	vaultCategory   storage.VaultCategoryRepository
	Sitota          storage.SitotaRepository
	Transaction     storage.TransactionRepository
	VaultAmountTier storage.VaultAmountTierRepository
}

func InitOraclePersistence(db *sql.DB, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		BankVault:       bankvault.NewBankVaultRepository(db, log),
		vaultCategory:   vaultCategory.NewVaultCategoryRepository(db, log),
		Sitota:          sitota.NewSitotaRepository(db, log),
		Transaction:     transaction_repo.NewTransactionRepository(db, log),
		VaultAmountTier: vaultamounttiers.NewVaultAmountTierRepository(db, log),
	}
}
