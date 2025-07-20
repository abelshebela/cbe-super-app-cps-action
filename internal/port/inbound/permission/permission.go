package permission

import "net/http"

type PermissionPortHandler interface {
	CreatePermissionGroup(w http.ResponseWriter, r *http.Request)
	GetPermissionGroups(w http.ResponseWriter, r *http.Request)
	GetPermissionGroup(w http.ResponseWriter, r *http.Request)
	UpdatePermissionGroup(w http.ResponseWriter, r *http.Request)
	// ApprovePermissionGroup(w http.ResponseWriter, r *http.Request)
	// RejectPermissionGroup(w http.ResponseWriter, r *http.Request)
}
