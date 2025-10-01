package initiator

import (
	"cbe-super-app-cps-action/internal/storage"
	"database/sql"

	"cbe-super-app-cps-action/internal/storage/persistance/bankvault"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type OraclePersistence struct {
	BankVault storage.BankVaultRepository
}

func InitOraclePersistence(db *sql.DB, log utils.Logger) OraclePersistence {
	return OraclePersistence{
		BankVault: bankvault.NewBankVaultRepository(db, log),
	}
}
