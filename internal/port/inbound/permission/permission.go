package permission

import "net/http"

type PermissionPortHandler interface {
	CreatePermissionGroup(w http.ResponseWriter, r *http.Request)
	ApprovePermissionGroup(w http.ResponseWriter, r *http.Request)
}
