package bank

type BankService interface {
	FindAll() ([]*BankResponse, error)
	FindOne(id string) (*BankResponse, error)
	InsertOne(req CreateBankRequest) (*Bank, error)
	UpdateOne(id string, update UpdateBankRequest) (*Bank, error)
	DeleteOne(id string) error
}

type Service struct {
	repository Repository
}

func NewSerice(repository Repository) (BankService, error) {
	return &Service{
		repository: repository,
	}, nil
}

func (s *Service) FindAll() ([]*BankResponse, error) {
	banks, err := s.repository.GetAllBank()
	if err != nil {
		return nil, err
	}
	return banks, nil
}

func (s *Service) FindOne(id string) (*BankResponse, error) {
	bank, err := s.repository.GetOneBank(id)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (s *Service) InsertOne(req CreateBankRequest) (*Bank, error) {
	bank, err := s.repository.CreateOneBank(req)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (s *Service) UpdateOne(id string, req UpdateBankRequest) (*Bank, error) {
	bank, err := s.repository.UpdateOneBank(id, req)
	if err != nil {
		return nil, err
	}

	return bank, nil
}

func (s *Service) DeleteOne(id string) error {
	if err := s.repository.DeleteOneBank(id); err != nil {
		return err
	}

	return nil
}
