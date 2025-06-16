package bank

type BankResponse struct {
	ID   string
	Name string
	Logo string
	Code string
	BIC  string
}

type CreateBankRequest struct {
	Name string
	Logo string
	Code string
	BIC  string
}

type UpdateBankRequest struct {
	Name *string
	Logo *string
	Code *string
	BIC  *string
}
