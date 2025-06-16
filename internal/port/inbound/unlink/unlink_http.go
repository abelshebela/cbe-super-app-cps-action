package unlink

import "net/http"

type UnlinkPortHandler interface {
	UnlinkDevice(w http.ResponseWriter, r *http.Request)
    ApproveUnlinkDevice(w http.ResponseWriter, r *http.Request) 
}
