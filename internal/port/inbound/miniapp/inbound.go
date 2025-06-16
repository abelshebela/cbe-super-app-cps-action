package miniapp

import "net/http"

type Inbound interface {
	MakerCreateMiniApp(w http.ResponseWriter, r *http.Request)
	CheckerMiniApp(w http.ResponseWriter, r *http.Request)
}
