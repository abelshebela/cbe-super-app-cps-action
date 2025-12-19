package external_call

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ExternalCallServices struct {
	SMS *SMSPersistence
}

func InitExternalCallServices(smsBaseURL string, logger utils.Logger) *ExternalCallServices {
	return &ExternalCallServices{
		SMS: NewSMSPersistence(smsBaseURL, logger),
	}
}
