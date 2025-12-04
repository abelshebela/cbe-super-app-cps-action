package unlink

import "net/http"

type UnlinkAdapter interface {
	GetArchivedUser(w http.ResponseWriter, r *http.Request)
	GetUserByAccount(w http.ResponseWriter, r *http.Request)
	UnlinkUserCif(w http.ResponseWriter, r *http.Request)
}
