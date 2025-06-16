package service

type Repository interface {
	GetAllService() ([]*Service, error)
	GetOneService(id string) (Service, error)
	UpdateOneService(id string, update any) error
}
