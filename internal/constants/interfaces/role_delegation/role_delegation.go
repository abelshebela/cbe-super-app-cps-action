package role_delegation_outbound

import "net/http"

type RoleDelegation interface {
	CreateWithNewUser(w http.ResponseWriter, r *http.Request)
	CreateWithExistingUser(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Export(w http.ResponseWriter, r *http.Request)
	FindAll(w http.ResponseWriter, r *http.Request)
	FindAllWithPagination(w http.ResponseWriter, r *http.Request)
	FindById(w http.ResponseWriter, r *http.Request)
	FindByUsername(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
