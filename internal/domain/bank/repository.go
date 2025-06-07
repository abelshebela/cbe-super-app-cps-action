package bank

type Repository interface {
	GetAllBank() ([]*BankResponse, error)
	GetOneBank(id string) (*BankResponse, error)
	CreateOneBank(req CreateBankRequest) (*Bank, error)
	UpdateOneBank(id string, update UpdateBankRequest) (*Bank, error)
	DeleteOneBank(id string) error
}
