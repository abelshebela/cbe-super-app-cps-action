package ad

type AdvertService interface {
	FindAll() ([]*AdvertResponse, error)
	FindOne(id string) (*AdvertResponse, error)
	InsertOne(req CreateAdvertRequest) (*Advert, error)
	UpdateOne(id string, update UpdateAdvertRequest) (*Advert, error)
	DeleteOne(id string) error
}

type Service struct {
	repository Repository
}

func NewSerice(repository Repository) (AdvertService, error) {
	return &Service{
		repository: repository,
	}, nil
}

func (s *Service) FindAll() ([]*AdvertResponse, error) {
	adverts, err := s.repository.GetAllAdvert()
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (s *Service) FindOne(id string) (*AdvertResponse, error) {
	advert, err := s.repository.GetOneAdvert(id)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (s *Service) InsertOne(req CreateAdvertRequest) (*Advert, error) {
	advert, err := s.repository.CreateOneAdvert(req)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (s *Service) UpdateOne(id string, req UpdateAdvertRequest) (*Advert, error) {
	advert, err := s.repository.UpdateOneAdvert(id, req)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (s *Service) DeleteOne(id string) error {
	if err := s.repository.DeleteOneAdvert(id); err != nil {
		return err
	}

	return nil
}
