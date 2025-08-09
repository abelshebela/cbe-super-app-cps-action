package encryption

import "net/http"

type Encryption interface {
	Encrypt(w http.ResponseWriter, r *http.Request)
}
