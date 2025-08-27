package accountvalidation

import (
	"net/http"
)

type AccountValidation interface {
	FindById(w http.ResponseWriter, r *http.Request)
	FindAllWithPagination(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}
