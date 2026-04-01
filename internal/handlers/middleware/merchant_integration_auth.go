package middleware

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"crypto/subtle"
	"net/http"
	"os"
	"strings"
)

const (
	envIntegrationAPIKey = "MERCHANT_INTEGRATION_API_KEY"
	// envIntegrationRoleCode = "MERCHANT_INTEGRATION_ROLE_CODE"
	// envIntegrationUserID   = "MERCHANT_INTEGRATION_USER_ID"
	// envIntegrationUserCode = "MERCHANT_INTEGRATION_USER_CODE"
	// envIntegrationUserName = "MERCHANT_INTEGRATION_USERNAME"
	// envIntegrationFullName = "MERCHANT_INTEGRATION_FULL_NAME"
	// envIntegrationPhone    = "MERCHANT_INTEGRATION_PHONE"
	// envIntegrationDept     = "MERCHANT_INTEGRATION_DEPARTMENT"
	// envIntegrationUserRole = "MERCHANT_INTEGRATION_USER_ROLE"
)

func (a *authMiddleware) AuthenticateTokenOrMerchantIntegrationAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimSpace(r.Header.Get("X-Api-Key"))
		if apiKey == "" {
			a.AuthenticateToken(next).ServeHTTP(w, r)
			return
		}

		configured := strings.TrimSpace(os.Getenv(envIntegrationAPIKey))
		// configured := a.cfg.MerchantIntegrationAPIKey
		if configured == "" {
			if a.logger != nil {
				a.logger.Warnf("[AuthMW][EcomIntegration] X-Api-Key present but %s is not set", envIntegrationAPIKey)
			}
			next.ServeHTTP(w, r)
			return
		}

		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(configured)) != 1 {
			if a.logger != nil {
				a.logger.Warnf("[AuthMW][EcomIntegration] invalid X-Api-Key")
			}
			next.ServeHTTP(w, r)
			return
		}

		username := strings.TrimSpace(r.Header.Get("username"))
		if username == "" {
			a.AuthenticateTempToken(next).ServeHTTP(w, r)
			return
		}

		fullname := strings.TrimSpace(r.Header.Get("fullname"))
		if fullname == "" {
			a.AuthenticateTempToken(next).ServeHTTP(w, r)
			return
		}

		phone_number := strings.TrimSpace(r.Header.Get("phone_number"))
		if phone_number == "" {
			a.AuthenticateTempToken(next).ServeHTTP(w, r)
			return
		}

		// roleCode := strings.TrimSpace(os.Getenv(envIntegrationRoleCode))
		roleCode := "ERP"
		// userID := strings.TrimSpace(os.Getenv(envIntegrationUserID))
		// fullName := strings.TrimSpace(os.Getenv(envIntegrationFullName))
		// phone := strings.TrimSpace(os.Getenv(envIntegrationPhone))
		// dept := strings.TrimSpace(os.Getenv(envIntegrationDept))

		if username == "" || fullname == "" || phone_number == "" {
			if a.logger != nil {
				a.logger.Errorf("[AuthMW][EcomIntegration] missing required headers (username, fullname, or phone_number)")
			}
			localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
			return
		}

		// userCode := strings.TrimSpace(os.Getenv(envIntegrationUserCode))
		// if userCode == "" {
		// 	userCode = userID
		// }
		// userName := strings.TrimSpace(os.Getenv(envIntegrationUserName))
		// if userName == "" {
		// 	userName = "ecommerce-integration"
		// }
		// userRole := strings.TrimSpace(os.Getenv(envIntegrationUserRole))
		// if userRole == "" {
		// 	userRole = "SERVICE"
		// }

		payload := UserPayload{
			IsERP:       true,
			PhoneNumber: phone_number,
			UserName:    username,
			// UserRole:    userRole,
			RoleId: roleCode,
			// UserID:      userID,
			// UserCode:    userCode,
			FullName: fullname,
			// Department:  dept,
			Environment: a.cfg.GoEnv,
			Permission:  nil,
		}

		ctx := a.setUserPayload(r.Context(), payload)
		localization.UpdateWriterContext(w, ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
