package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
)

// Required when clients send X-Api-Key:
//   - MERCHANT_INTEGRATION_API_KEY
//   - MERCHANT_INTEGRATION_ROLE_CODE
//   - MERCHANT_INTEGRATION_USER_ID
//   - MERCHANT_INTEGRATION_FULL_NAME
//   - MERCHANT_INTEGRATION_PHONE
//   - MERCHANT_INTEGRATION_DEPARTMENT
//
// Optional
//   - MERCHANT_INTEGRATION_USER_CODE
//   - MERCHANT_INTEGRATION_USERNAME
//   - MERCHANT_INTEGRATION_USER_ROLE
const (
	envIntegrationAPIKey   = "MERCHANT_INTEGRATION_API_KEY"
	envIntegrationRoleCode = "MERCHANT_INTEGRATION_ROLE_CODE"
	envIntegrationUserID   = "MERCHANT_INTEGRATION_USER_ID"
	envIntegrationUserCode = "MERCHANT_INTEGRATION_USER_CODE"
	envIntegrationUserName = "MERCHANT_INTEGRATION_USERNAME"
	envIntegrationFullName = "MERCHANT_INTEGRATION_FULL_NAME"
	envIntegrationPhone    = "MERCHANT_INTEGRATION_PHONE"
	envIntegrationDept     = "MERCHANT_INTEGRATION_DEPARTMENT"
	envIntegrationUserRole = "MERCHANT_INTEGRATION_USER_ROLE"
)

func (a *authMiddleware) AuthenticateTokenOrMerchantIntegrationAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimSpace(r.Header.Get("X-Api-Key"))
		if apiKey == "" {
			a.AuthenticateToken(next).ServeHTTP(w, r)
			return
		}

		configured := strings.TrimSpace(os.Getenv(envIntegrationAPIKey))
		if configured == "" {
			if a.logger != nil {
				a.logger.Warnf("[AuthMW][EcomIntegration] X-Api-Key present but %s is not set", envIntegrationAPIKey)
			}
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(configured)) != 1 {
			if a.logger != nil {
				a.logger.Warnf("[AuthMW][EcomIntegration] invalid X-Api-Key")
			}
			localization.SendUnauthorizedResponse(w, localization.ErrorUserUnauthorized.Message)
			return
		}

		roleCode := strings.TrimSpace(os.Getenv(envIntegrationRoleCode))
		userID := strings.TrimSpace(os.Getenv(envIntegrationUserID))
		fullName := strings.TrimSpace(os.Getenv(envIntegrationFullName))
		phone := strings.TrimSpace(os.Getenv(envIntegrationPhone))
		dept := strings.TrimSpace(os.Getenv(envIntegrationDept))

		if roleCode == "" || userID == "" || fullName == "" || phone == "" || dept == "" {
			if a.logger != nil {
				a.logger.Errorf("[AuthMW][EcomIntegration] missing required integration identity env (role, user id, name, phone, or department)")
			}
			localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
			return
		}

		userCode := strings.TrimSpace(os.Getenv(envIntegrationUserCode))
		if userCode == "" {
			userCode = userID
		}
		userName := strings.TrimSpace(os.Getenv(envIntegrationUserName))
		if userName == "" {
			userName = "ecommerce-integration"
		}
		userRole := strings.TrimSpace(os.Getenv(envIntegrationUserRole))
		if userRole == "" {
			userRole = "SERVICE"
		}

		payload := UserPayload{
			PhoneNumber: phone,
			UserName:    userName,
			UserRole:    userRole,
			RoleId:      roleCode,
			UserID:      userID,
			UserCode:    userCode,
			FullName:    fullName,
			Department:  dept,
			Environment: a.cfg.GoEnv,
			Permission:  nil,
		}

		ctx := a.setUserPayload(r.Context(), payload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
