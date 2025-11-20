package permission

import "net/http"

type PermissionHandler interface {
	CreatePermissionGroup(w http.ResponseWriter, r *http.Request)
	GetPermissionGroups(w http.ResponseWriter, r *http.Request)
	GetPermissionGroup(w http.ResponseWriter, r *http.Request)
	GetPermissionGroupById(w http.ResponseWriter, r *http.Request)
	UpdatePermissionGroup(w http.ResponseWriter, r *http.Request)
	GetAllPermissionCategoriesWithPermissions(w http.ResponseWriter, r *http.Request)
	GetPermissionGroupsByDepartment(w http.ResponseWriter, r *http.Request)
}
