package department

import "net/http"

type DepartmentPortHandler interface {
	CreateDepartment(w http.ResponseWriter, r *http.Request)
	ApproveDepartmentRequest(w http.ResponseWriter, r *http.Request)
	RejectDepartmentRequest(w http.ResponseWriter, r *http.Request)
	UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request)
}
