package job_role_interface

import "net/http"

type RolesInbound interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	GetAllWithPagination(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
}
