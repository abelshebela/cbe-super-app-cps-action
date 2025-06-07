package event

type Repository interface {
	GetAllEvent() ([]*EventResponse, error)
	GetOneEvent(id string) (*EventResponse, error)
	CreateOneEvent(req CreateEventRequest) (*Event, error)
	UpdateOneEvent(id string, update UpdateEventRequest) (*Event, error)
	DeleteOneEvent(id string) error
}
