package utils

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse structure for standard error responses
type ErrorRes struct {
	Data    []string `json:"data"`
	Message string   `json:"message"`
	Status  int      `json:"status"`
}

func CheckRequiredHeaders(w http.ResponseWriter, r *http.Request, requiredHeaders []string) bool {
	for _, header := range requiredHeaders {
		if r.Header.Get(header) == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorRes{
				Data:    []string{},
				Message: "Missing required header: " + header,
				Status:  400,
			})
			return false
		}
	}
	return true
}

func CheckRequiredQueriesNotEmpty(w http.ResponseWriter, r *http.Request, requiredQueries []string) bool {
	q := r.URL.Query()
	for _, param := range requiredQueries {
		value := q.Get(param)
		if value == "" || len(value) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorRes{
				Data:    []string{},
				Message: "Missing or empty required query parameter: " + param,
				Status:  400,
			})
			return false
		}
	}
	return true
}

func CheckRequiredQueries(w http.ResponseWriter, r *http.Request, requiredQueries []string) bool {
	q := r.URL.Query()
	for _, param := range requiredQueries {
		if q.Get(param) == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorRes{
				Data:    []string{},
				Message: "Missing required query parameter: " + param,
				Status:  400,
			})
			return false
		}
	}
	return true
}
