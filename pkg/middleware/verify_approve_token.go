package middleware

import (
	"cbe-super-app-member-auth/pkg/config"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyApproovToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("========== Verifying Approov Token ==========")

		cfg, err := config.Load()
		if err != nil {
			log.Println("Failed to initialize config:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		approovToken := r.Header.Get("approov-token")
		if approovToken == "" {
			log.Println("Approov Token is missing in request headers")
			http.Error(w, "Missing Approov Token", http.StatusBadRequest)
			return
		}

		approovSecret := cfg.ApproovSecret // Adjust field access if needed
		if approovSecret == "" {
			log.Println("Approov Secret is missing in the configuration")
			http.Error(w, "Missing Approov Secret", http.StatusBadRequest)
			return
		}

		// Decode the base64 secret
		decodedSecret, err := base64.StdEncoding.DecodeString(approovSecret)
		if err != nil {
			log.Println("Failed to decode Approov secret:", err)
			http.Error(w, "Invalid Approov Secret", http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(approovToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != "HS256" {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return decodedSecret, nil
		})

		if err != nil || !token.Valid {
			log.Println("Approov Token verification failed:", err)
			http.Error(w, "Invalid Approov Token", http.StatusBadRequest)
			return
		}

		log.Println("========== Approov Token Verified ==========")
		next.ServeHTTP(w, r)
	})
}
