package ad

type Repository interface {
	GetAllAdvert() ([]*AdvertResponse, error)
	GetOneAdvert(id string) (*AdvertResponse, error)
	CreateOneAdvert(req CreateAdvertRequest) (*Advert, error)
	UpdateOneAdvert(id string, update UpdateAdvertRequest) (*Advert, error)
	DeleteOneAdvert(id string) error
}
