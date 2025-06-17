package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

var specialCharRegex = regexp.MustCompile(`[<>/*@.]+`)

func SantizeString(s string) string {

	return specialCharRegex.ReplaceAllString(s, "")
}

func sanitizeMap(m map[string]interface{}) {
	for key, val := range m {
		// fmt.Println(val)
		switch v := val.(type) {
		case string:
			m[key] = SantizeString(v)
		case map[string]interface{}:
			sanitizeMap(v)
		case []interface{}:
			for i, item := range v {
				switch itemTyped := item.(type) {
				case string:
					v[i] = SantizeString(itemTyped)
				case map[string]interface{}:
					sanitizeMap(itemTyped)
				}
			}
		}
	}
}

func SanitizeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		for key, values := range query {
			for i, val := range values {
				values[i] = SantizeString(val)
			}
			query[key] = values
		}

		r.URL.RawQuery = query.Encode()

		if r.Header.Get("Content-Type") == "application/json" {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err == nil && len(bodyBytes) > 0 {
				var jsonData map[string]interface{}
				if json.Unmarshal(bodyBytes, &jsonData) == nil {
					sanitizeMap(jsonData)
					newBody, _ := json.Marshal(jsonData)
					r.Body = io.NopCloser(bytes.NewBuffer(newBody))
					r.ContentLength = int64(len(newBody))
				} else {
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
