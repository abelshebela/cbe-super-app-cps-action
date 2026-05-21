package superapprole

import "net/http"

type SuperAppRole interface {
	GetAllSuperAppRoles(w http.ResponseWriter, r *http.Request)
	GetTransferLimitByRole(w http.ResponseWriter, r *http.Request)
	EnableByRole(w http.ResponseWriter, r *http.Request)
	DisableByRole(w http.ResponseWriter, r *http.Request)
	DeleteByRole(w http.ResponseWriter, r *http.Request)
	GetAccessListsByRole(w http.ResponseWriter, r *http.Request)
	BulkDisableAccessLists(w http.ResponseWriter, r *http.Request)
	BulkEnableAccessLists(w http.ResponseWriter, r *http.Request)
}
