package merchant_lookup

import "cbe-super-app-cps-action/internal/constants/types"

type MerchantLookUpResponse struct {
	Status   string              `json:"status"`
	Company  types.Company       `json:"company"`
	Branches []types.Branch      `json:"branches"`
	Users    []types.UserAccount `json:"users"`
}
