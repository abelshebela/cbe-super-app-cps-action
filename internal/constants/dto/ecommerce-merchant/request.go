package ecommercemerchant

type BranchInformation struct {
	BranchCode          string `json:"branch_code" bson:"branch_code"`
	BranchName          string `json:"branch_name" bson:"branch_name"`
	BranchAddress       string `json:"branch_address" bson:"branch_address"`
	BranchOwner         string `json:"branch_owner" bson:"branch_owner"`
	BranchAccountNumber string `json:"branch_account_number" bson:"branch_account_number"`
}

type EcommerceMerchant struct {
	MerchantName     string              `json:"merchant_name"`
	MerchantCode     string              `json:"merchant_code"`
	AccountNumber    string              `json:"account_number"`
	SettlementMethod string              `json:"settlement_method" bson:"settlement_method"`
	Branches         []BranchInformation `json:"branches"`
	// IsEcommerceMerchant *bool                     `json:"is_ecommerce_merchant"`
}

type UpdateEcommerceMerchant struct {
	MerchantName     *string             `json:"merchant_name"`
	MerchantCode     *string             `json:"merchant_code"`
	AccountNumber    *string             `json:"account_number"`
	SettlementMethod *string             `json:"settlement_method" bson:"settlement_method"`
	Branches         []BranchInformation `json:"branches"`
}

type KYCDTO struct {
	Status         string            `json:"status"`
	Representative RepresentativeDTO `json:"representative"`
}

type RepresentativeDTO struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type ERPUpdateMerchantRequest struct {
	CPSAccountNumber string            `json:"cps_account_number"`
	Branches         []ERPUpdateBranch `json:"branches"`
}

type ERPUpdateBranch struct {
	Merchant         string `json:"merchant"`
	CPSAccountNumber string `json:"cps_account_number"`
}
