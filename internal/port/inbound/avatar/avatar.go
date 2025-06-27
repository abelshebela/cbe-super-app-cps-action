package avatar

import "net/http"

type AvatarInbound interface {
	CreateAvatar(w http.ResponseWriter, r *http.Request)
}
