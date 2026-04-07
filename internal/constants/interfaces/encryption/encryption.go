package encryption

import "net/http"

type EncryptionAdapter interface {
	Encrypt(w http.ResponseWriter, r *http.Request)
}
