package middleware

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"context"
// )

// // User represents the authenticated user structure.
// // type User struct {
// // 	ID             string   `json:"_id"`
// // 	FullName       string   `json:"fullName"`
// // 	Role           string   `json:"role"`
// // 	UserCode       string   `json:"userCode"`
// // 	OrganizationID string   `json:"organizationID"`
// // 	PhoneNumber    string   `json:"phoneNumber"`
// // 	Email          string   `json:"email"`
// // 	Permissions    []string `json:"permissions"`
// // 	Realm          string   `json:"realm"`
// // 	IsMaker        bool     `json:"isMaker"`
// // 	IsChecker      bool     `json:"isChecker"`
// // 	SessionExpires string   `json:"sessionExpiresOn"`
// // 	DeviceUUID     string   `json:"deviceUUID"`
// // }

// type contextKey string

// func Authorization(realms []string, permissions []string, userRoles ...[]string) func(http.Handler) http.Handler {

// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			userVal := r.Context().Value("_user")

// 			if userVal == nil {
// 				writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
// 					"status":  401,
// 					"type":    "AUTHORIZATION_ERROR",
// 					"message": "Missing Authenticated User",
// 				})
// 				return
// 			}
// 			fmt.Println(userVal)
// 			// user, ok := userVal.(User)
// 			fmt.Println("=============user===========")
// 			// fmt.Println(user)
// 			var user User
// 			var ok bool
// 			if m, isMap := userVal.(map[string]interface{}); isMap {
// 				// Marshal then unmarshal to User struct
// 				b, err := json.Marshal(m)
// 				if err == nil {
// 					err = json.Unmarshal(b, &user)
// 					ok = err == nil
// 				}
// 			} else {
// 				user, ok = userVal.(User)
// 			}

// 			if !ok {
// 				writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
// 					"status":  401,
// 					"type":    "AUTHORIZATION_ERROR",
// 					"message": "Invalid User Type",
// 				})
// 				return
// 			}
// 			// ctx := context.WithValue(r.Context(), "_user", user)
// 			// r = r.WithContext(ctx)
// 			// fmt.Println(user.Realm)
// 			// Check realm
// 			realmFound := false
// 			for _, realm := range realms {
// 				if realm == "*" || user.Realm == realm {
// 					realmFound = true
// 					break
// 				}
// 			}
// 			fmt.Println(user.Permissions)

// 			// Check permissions
// 			permissionFound := false
// 			if contains(permissions, "*") {
// 				permissionFound = true
// 			} else {
// 				for _, p := range permissions {
// 					if contains(user.Permissions, p) {
// 						permissionFound = true
// 						break
// 					}
// 				}
// 			}

// 			// Check user role if realm is bank
// 			bankFound := true
// 			if user.Realm == "bank" {
// 				bankFound = false
// 				if len(userRoles) > 0 {
// 					for _, r := range userRoles[0] {
// 						if user.Role == r {
// 							bankFound = true
// 							break
// 						}
// 					}
// 				}
// 			}
// 			fmt.Println(realmFound)
// 			fmt.Println(permissionFound)
// 			fmt.Println(bankFound)
// 			isAuthorized := realmFound && bankFound && permissionFound

// 			if !isAuthorized {
// 				writeJSON(w, http.StatusBadRequest, map[string]interface{}{
// 					"status":  400,
// 					"type":    "AUTHORIZATION_ERROR",
// 					"message": "Action Not Allowed",
// 				})
// 				return
// 			}

// 			// Session expiry check (uncomment if needed)
// 			// if user.SessionExpires == "" || isSessionExpired(user.SessionExpires) {
// 			// 	writeJSON(w, http.StatusForbidden, map[string]interface{}{
// 			// 		"message": "session expired",
// 			// 	})
// 			// 	return
// 			// }

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // contains checks if a string slice contains a value.
// func contains(slice []string, val string) bool {
// 	for _, s := range slice {
// 		if s == val {
// 			return true
// 		}
// 	}
// 	return false
// }

// // isSessionExpired checks if the session is expired.
// func isSessionExpired(sessionExpires string) bool {
// 	t, err := time.Parse(time.RFC3339, sessionExpires)
// 	if err != nil {
// 		return true
// 	}
// 	return time.Now().After(t)
// }

// // Helper to write JSON responses
// func writeJSON(w http.ResponseWriter, status int, data interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	_ = json.NewEncoder(w).Encode(data)
// }

// // Helper to set user in context (to be used in your authentication middleware)
// // func WithUser(next http.Handler) http.Handler {
// // 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// // 		// Example: get user from request (replace with your logic)
// // 		var user User
// // 		// ... populate user ...
// // 		ctx := context.WithValue(r.Context(), userContextKey, user)
// // 		next.ServeHTTP(w, r.WithContext(ctx))
// // 	})
// // }
