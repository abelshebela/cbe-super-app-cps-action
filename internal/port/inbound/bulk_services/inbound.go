package inbound

import "net/http"

type Inbound interface {
	FetchServices(w http.ResponseWriter, r *http.Request)
	EnableDisableServicesMaker(w http.ResponseWriter, r *http.Request)
	EnableDisableServicesChecker(w http.ResponseWriter, r *http.Request)
}
