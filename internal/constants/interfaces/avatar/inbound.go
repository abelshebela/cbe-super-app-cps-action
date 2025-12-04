package avatar

import "net/http"

type AvatarInbound interface {
	CreateAvatar(w http.ResponseWriter, r *http.Request)
	DeleteAvatar(w http.ResponseWriter, r *http.Request)
	Enable(w http.ResponseWriter, r *http.Request)
	Disable(w http.ResponseWriter, r *http.Request)
	FetchAvatar(w http.ResponseWriter, r *http.Request)
	FetchAvatars(w http.ResponseWriter, r *http.Request)
	UpdateAvatar(w http.ResponseWriter, r *http.Request)
}
