package common

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	ResponseWriter http.ResponseWriter
	Status         int
	Data           T
}

func (r *Response[T]) SendJSON() {
	r.ResponseWriter.Header().Set("Content-Type", "application/json")
	r.ResponseWriter.WriteHeader(r.Status)
	if err := json.NewEncoder(r.ResponseWriter).Encode(r.Data); err != nil {
		http.Error(
			r.ResponseWriter,
			`{"error": "failed to encode response"}`,
			http.StatusInternalServerError,
		)
	}
}
