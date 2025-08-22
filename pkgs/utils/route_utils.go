package utils

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetParam(r *http.Request, key string) (string, bool) {
	value := chi.URLParam(r, key)
	if value == "" {
		return "", false
	}
	return value, true
}
