package unlink

import "net/http"

type UnlinkPortHandler interface {
	UnlinkDevice(w http.ResponseWriter, r *http.Request)
}

type UnlinkHandler interface {
	GetArchivedUser(w http.ResponseWriter, r *http.Request)
	GetUserByAccount(w http.ResponseWriter, r *http.Request)
	UnlinkUserCif(w http.ResponseWriter, r *http.Request)
}
