package initiator

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitConfig(logger utils.Logger) *config.VaultConfig {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("failed to load config %v", err)
	}

	return cfg
}
