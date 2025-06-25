package account

import (
	"net/http"
)

type InBound interface {
	CreateAccount(w http.ResponseWriter, r *http.Request)
}

type AccountCreationResult struct {
	CustomerNumber string
	AccountNumber  string
}

