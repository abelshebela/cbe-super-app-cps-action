package department

import "net/http"

type DepartmentHandler interface {
	GetAllDepartments(w http.ResponseWriter, r *http.Request)
	CreateDepartment(w http.ResponseWriter, r *http.Request)
	UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request)
	GetDepartmentByID(w http.ResponseWriter, r *http.Request)
	EnableDepartment(w http.ResponseWriter, r *http.Request)
	DisableDepartment(w http.ResponseWriter, r *http.Request)
}
