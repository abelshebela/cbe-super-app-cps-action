package account

import "net/http"

type AccountPortHandler interface {
    CreateAccount(w http.ResponseWriter, r *http.Request)
}