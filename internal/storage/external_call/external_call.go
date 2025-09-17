package external_call

import (
	"time"

	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// ExternalCallServices contains all external service clients
type ExternalCallServices struct {
	AccountLookup account_lookup.Account
	SMS           *SMSPersistence
}

// InitExternalCallServices initializes all external service clients
func InitExternalCallServices(cbeBaseURL, smsBaseURL string, logger utils.Logger) *ExternalCallServices {
	return &ExternalCallServices{
		AccountLookup: account_lookup.InitAccountAPIClient(cbeBaseURL, 30*time.Second, logger),
		SMS:           NewSMSPersistence(smsBaseURL, logger),
	}
}
