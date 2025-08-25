package accountvalidation

import (
	"net/http"
)

type AccountValidation interface {
	FetchAccountValidation(w http.ResponseWriter, r *http.Request)
	FetchAllAccountValidation(w http.ResponseWriter, r *http.Request)
	UpdateAccountValidation(w http.ResponseWriter, r *http.Request)
}
