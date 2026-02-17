package erp_merchant_update_dto

type ERPBranch struct {
	Merchant         string `json:"merchant"`
	CPSAccountNumber string `json:"cps_account_number"`
	CpsEnabled       *bool  `json:"cps_enabled"`
}
type ERPUpdateRequest struct {
	MainAccountNumber string      `json:"cps_account_number"`
	CpsEnabled        *bool       `json:"cps_enabled"`
	Branches          []ERPBranch `json:"branches"`
}
