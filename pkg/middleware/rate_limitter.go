package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	requestCounts           = make(map[string][]int64)
	blockedHeaders          = make(map[string]int64)
	blockedHeadersFile      = "blockedHeaders.json"
	storeBlockedHeadersFile = "storedBlockedHeaders.json"
	mu                      sync.Mutex
)

func init() {
	if _, err := os.Stat(blockedHeadersFile); err == nil {
		data, err := ioutil.ReadFile(blockedHeadersFile)
		if err == nil {
			json.Unmarshal(data, &blockedHeaders)
		}
	}

	// Clean expired headers every minute
	go func() {
		for range time.Tick(time.Minute) {
			clearExpiredBlocks()
		}
	}()
}

func HeaderKeyGenerator(r *http.Request) string {
	relevantHeaders := map[string]string{
		"host":                      r.Host,
		"connection":                r.Header.Get("Connection"),
		"sec-ch-ua":                 r.Header.Get("Sec-Ch-Ua"),
		"sec-ch-ua-mobile":          r.Header.Get("Sec-Ch-Ua-Mobile"),
		"sec-ch-ua-platform":        r.Header.Get("Sec-Ch-Ua-Platform"),
		"upgrade-insecure-requests": r.Header.Get("Upgrade-Insecure-Requests"),
		"user-agent":                r.Header.Get("User-Agent"),
		"accept":                    r.Header.Get("Accept"),
		"sec-fetch-site":            r.Header.Get("Sec-Fetch-Site"),
		"sec-fetch-mode":            r.Header.Get("Sec-Fetch-Mode"),
		"sec-fetch-user":            r.Header.Get("Sec-Fetch-User"),
		"sec-fetch-dest":            r.Header.Get("Sec-Fetch-Dest"),
		"accept-encoding":           r.Header.Get("Accept-Encoding"),
		"accept-language":           r.Header.Get("Accept-Language"),
	}
	key, _ := json.Marshal(relevantHeaders)
	return string(key)
}

func saveBlockedHeaders() {
	data, err := json.Marshal(blockedHeaders)
	if err == nil {
		_ = ioutil.WriteFile(blockedHeadersFile, data, 0644)
		_ = ioutil.WriteFile(storeBlockedHeadersFile, data, 0644)
	}
}

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		headersKey := HeaderKeyGenerator(r)
		now := time.Now().UnixMilli()

		if expiry, exists := blockedHeaders[headersKey]; exists && now < expiry {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Too many requests, please try again later",
			})
			return
		}

		timestamps := requestCounts[headersKey]
		timestamps = append(timestamps, now)

		tenSecondsAgo := now - 10_000
		var recent []int64
		for _, ts := range timestamps {
			if ts >= tenSecondsAgo {
				recent = append(recent, ts)
			}
		}
		requestCounts[headersKey] = recent

		if len(recent) > 10 {
			blockedHeaders[headersKey] = now + 60_000
			delete(requestCounts, headersKey)
			saveBlockedHeaders()

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Too many requests, please try again later",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clearExpiredBlocks() {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now().UnixMilli()
	for key, expiry := range blockedHeaders {
		if now >= expiry {
			delete(blockedHeaders, key)
		}
	}
	saveBlockedHeaders()
}
