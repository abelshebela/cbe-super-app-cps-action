package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"strings"

	"cbe-super-app-cps-action/internal/constants/localization"
)

// Environment variables for server-to-server access to ecommerce-merchant APIs
// (Option A: X-Api-Key OR Bearer JWT). Configure the key and a CPS identity that
// exists in CPS action approve index for ecommerce merchant request actions.
//
// Required when clients send X-Api-Key:
//   - ECOMMERCE_MERCHANT_INTEGRATION_API_KEY
//   - ECOMMERCE_MERCHANT_INTEGRATION_ROLE_CODE
//   - ECOMMERCE_MERCHANT_INTEGRATION_USER_ID
//   - ECOMMERCE_MERCHANT_INTEGRATION_FULL_NAME
//   - ECOMMERCE_MERCHANT_INTEGRATION_PHONE
//   - ECOMMERCE_MERCHANT_INTEGRATION_DEPARTMENT
//
// Optional (defaults shown in implementation):
//   - ECOMMERCE_MERCHANT_INTEGRATION_USER_CODE
//   - ECOMMERCE_MERCHANT_INTEGRATION_USERNAME
//   - ECOMMERCE_MERCHANT_INTEGRATION_USER_ROLE
const (
	envEcommerceIntegrationAPIKey   = "ECOMMERCE_MERCHANT_INTEGRATION_API_KEY"
	envEcommerceIntegrationRoleCode = "ECOMMERCE_MERCHANT_INTEGRATION_ROLE_CODE"
	envEcommerceIntegrationUserID   = "ECOMMERCE_MERCHANT_INTEGRATION_USER_ID"
	envEcommerceIntegrationUserCode = "ECOMMERCE_MERCHANT_INTEGRATION_USER_CODE"
	envEcommerceIntegrationUserName = "ECOMMERCE_MERCHANT_INTEGRATION_USERNAME"
	envEcommerceIntegrationFullName = "ECOMMERCE_MERCHANT_INTEGRATION_FULL_NAME"
	envEcommerceIntegrationPhone    = "ECOMMERCE_MERCHANT_INTEGRATION_PHONE"
	envEcommerceIntegrationDept     = "ECOMMERCE_MERCHANT_INTEGRATION_DEPARTMENT"
	envEcommerceIntegrationUserRole = "ECOMMERCE_MERCHANT_INTEGRATION_USER_ROLE"
)

// AuthenticateTokenOrMerchantIntegrationAPIKey accepts either:
//   - Authorization: Bearer <JWT> (same behavior as AuthenticateToken), or
//   - X-Api-Key (or X-API-Key): must match ECOMMERCE_MERCHANT_INTEGRATION_API_KEY
//     and required integration identity env vars; context is populated like a CPS user
//     so handlers and CPS action creation receive role_code and maker fields.
//
// If the client sends a non-empty X-Api-Key header, JWT is not attempted for that request.
func (a *authMiddleware) AuthenticateTokenOrMerchantIntegrationAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := strings.TrimSpace(r.Header.Get("X-Api-Key"))
		if apiKey == "" {
			a.AuthenticateToken(next).ServeHTTP(w, r)
			return
		}

		configured := strings.TrimSpace(os.Getenv(envEcommerceIntegrationAPIKey))
		if configured == "" {
			if a.logger != nil {
				a.logger.Warnf("[AuthMW][EcomIntegration] X-Api-Key present but %s is not set", envEcommerceIntegrationAPIKey)
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

		roleCode := strings.TrimSpace(os.Getenv(envEcommerceIntegrationRoleCode))
		userID := strings.TrimSpace(os.Getenv(envEcommerceIntegrationUserID))
		fullName := strings.TrimSpace(os.Getenv(envEcommerceIntegrationFullName))
		phone := strings.TrimSpace(os.Getenv(envEcommerceIntegrationPhone))
		dept := strings.TrimSpace(os.Getenv(envEcommerceIntegrationDept))

		if roleCode == "" || userID == "" || fullName == "" || phone == "" || dept == "" {
			if a.logger != nil {
				a.logger.Errorf("[AuthMW][EcomIntegration] missing required integration identity env (role, user id, name, phone, or department)")
			}
			localization.SendErrorResponse(w, localization.ErrorUnexpectedError, nil, nil)
			return
		}

		userCode := strings.TrimSpace(os.Getenv(envEcommerceIntegrationUserCode))
		if userCode == "" {
			userCode = userID
		}
		userName := strings.TrimSpace(os.Getenv(envEcommerceIntegrationUserName))
		if userName == "" {
			userName = "ecommerce-integration"
		}
		userRole := strings.TrimSpace(os.Getenv(envEcommerceIntegrationUserRole))
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
