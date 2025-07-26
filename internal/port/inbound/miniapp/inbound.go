package miniapp

import "net/http"

type MiniAppInbound interface {
	MakerCreateMiniApp(w http.ResponseWriter, r *http.Request)

	MakerUpdateMiniApp(w http.ResponseWriter, r *http.Request)
	MakerDeleteMiniApp(w http.ResponseWriter, r *http.Request)
	ListMiniApp(w http.ResponseWriter, r *http.Request)
	DetailMiniAppByID(w http.ResponseWriter, r *http.Request)
	EnableMiniAppByID(w http.ResponseWriter, r *http.Request)
	DisableMiniAppByID(w http.ResponseWriter, r *http.Request)

}
