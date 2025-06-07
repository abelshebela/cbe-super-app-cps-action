package event

type EventService interface {
	FindAll() ([]*EventResponse, error)
	FindOne(id string) (*EventResponse, error)
	InsertOne(req CreateEventRequest) (*Event, error)
	UpdateOne(id string, update UpdateEventRequest) (*Event, error)
	DeleteOne(id string) error
}

type Service struct {
	repository Repository
}

func NewSerice(repository Repository) (EventService, error) {
	return &Service{
		repository: repository,
	}, nil
}

func (s *Service) FindAll() ([]*EventResponse, error) {
	adverts, err := s.repository.GetAllEvent()
	if err != nil {
		return nil, err
	}
	return adverts, nil
}

func (s *Service) FindOne(id string) (*EventResponse, error) {
	advert, err := s.repository.GetOneEvent(id)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (s *Service) InsertOne(req CreateEventRequest) (*Event, error) {
	advert, err := s.repository.CreateOneEvent(req)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (s *Service) UpdateOne(id string, req UpdateEventRequest) (*Event, error) {
	advert, err := s.repository.UpdateOneEvent(id, req)
	if err != nil {
		return nil, err
	}

	return advert, nil
}

func (s *Service) DeleteOne(id string) error {
	if err := s.repository.DeleteOneEvent(id); err != nil {
		return err
	}

	return nil
}
