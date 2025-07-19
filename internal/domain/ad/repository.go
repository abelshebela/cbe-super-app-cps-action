package ad
import (
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/ad/entity"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
)

type Repository interface {
	GetAllAdvert() (*common_util.PaginatedResponse[[]*entity.Advert], error)
	GetOneAdvert(id string) (*AdvertResponse, error)
	CreateOneAdvert(req CreateAdvertRequest) (*Advert, error)
	UpdateOneAdvert(id string, update UpdateAdvertRequest) (*Advert, error)
	DeleteOneAdvert(id string) error
}
