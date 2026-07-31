package accountlookup

import "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"

type AccountResponse struct {
	Data types.AccountInfo `json:"data"`
}

type CreateAccountResponse struct {
	Data types.Account `json:"data"`
}
