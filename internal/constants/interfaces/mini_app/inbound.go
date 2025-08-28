package miniapp

import "net/http"

type MiniAppInbound interface {
	CreateMiniApp(w http.ResponseWriter, r *http.Request)

	UpdateMiniApp(w http.ResponseWriter, r *http.Request)
	DeleteMiniApp(w http.ResponseWriter, r *http.Request)
	ListMiniApp(w http.ResponseWriter, r *http.Request)
	DetailMiniAppByID(w http.ResponseWriter, r *http.Request)
	EnableMiniAppByID(w http.ResponseWriter, r *http.Request)
	DisableMiniAppByID(w http.ResponseWriter, r *http.Request)
}
