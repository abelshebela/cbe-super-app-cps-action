package initiator

import (
	"os"
	"time"

	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key string, logger utils.Logger) string {
	if value := os.Getenv(key); value != "" {

		return value
	}
	logger.Errorf("the %s not found in the vault", key)
	return ""
}

// InitAccountLookupService initializes the account lookup service using configuration
func InitAccountLookupService(cbebaseurl string, logger utils.Logger) account_lookup.Account {
	timeout := 30 * time.Second

	return account_lookup.InitAccountAPIClient(cbebaseurl, timeout, logger)
}
