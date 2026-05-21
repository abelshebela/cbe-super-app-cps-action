package superapprole

import "net/http"

type SuperAppRole interface {
	GetAllSuperAppRoles(w http.ResponseWriter, r *http.Request)
	EnableByRole(w http.ResponseWriter, r *http.Request)
	DisableByRole(w http.ResponseWriter, r *http.Request)
	DeleteByRole(w http.ResponseWriter, r *http.Request)
}
