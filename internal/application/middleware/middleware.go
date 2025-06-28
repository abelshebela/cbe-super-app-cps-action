package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	constant "cbe-super-app-member-users/pkgs/utils"
)

type UserPayload struct {
	UserID               string      `json:"user_id,omitempty"`
	UserCode             string      `json:"user_code,omitempty"`
	FullName             string      `json:"full_name,omitempty"`
	PhoneNumber          string      `json:"phone_number,omitempty"`
	Email                string      `json:"user_email,omitempty"`
	Realm                string      `json:"user_realm,omitempty"`
	MemberType           bool        `json:"ifb_member,omitempty"`
	DeviceUUID           string      `json:"device_uuid,omitempty"`
	UserDeviceLinkedDate string      `json:"user_device_linked_date,omitempty"`
	Permissions          interface{} `json:"permissions,omitempty"`
	SessionExpiry        interface{} `json:"session_expiry,omitempty"`
}

type authMiddleware struct {
	logger       utils.Logger
	JWTSecretKey string
	Key          string
	IV           string
}

type AuthMiddleware interface {
	AccessControl(allowedRoles []string) func(http.Handler) http.Handler
	AuthenticateToken(next http.Handler) http.Handler
}

func InitAuthMiddleware(secretKey, key, iv string, logger utils.Logger) AuthMiddleware {
	return &authMiddleware{
		JWTSecretKey: secretKey,
		Key:          key,
		IV:           iv,
		logger:       logger,
	}
}

func (a *authMiddleware) AccessControl(allowedRoles []string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[strings.ToUpper(r)] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(constant.ContextKey("user_role")).(string)
			if !ok || role == "" {
				res := common.Response[constant.ErrorDefinition]{
					ResponseWriter: w,
					Status:         http.StatusUnauthorized,
					Data: constant.ErrorDefinition{
						Code:    http.StatusUnauthorized,
						Message: "unauthorized",
					},
				}
				res.SendJSON()
				return
			}

			if _, allowed := roleSet[strings.ToUpper(role)]; !allowed {
				res := common.Response[constant.ErrorDefinition]{
					ResponseWriter: w,
					Status:         http.StatusForbidden,
					Data: constant.ErrorDefinition{
						Code:    http.StatusForbidden,
						Message: "action not allowed",
					},
				}
				res.SendJSON()
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (a *authMiddleware) AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const bearer = "Bearer "
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, bearer) {
			a.sendUnauthorizedResponse(w, "bearer token is not present", "unauthorized")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)
		if tokenString == "" {
			a.sendUnauthorizedResponse(w, "empty token string provided", "access token required")
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				a.logger.Errorf("unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(a.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			a.sendUnauthorizedResponse(w, "invalid or expired token", "invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			a.sendUnauthorizedResponse(w, "failed to cast to map claims", "invalid token")
			return
		}

		data, ok := claims["data"].(string)
		if !ok || data == "" {
			a.sendUnauthorizedResponse(w, "invalid token or data not present", "invalid token")
			return
		}

		decryptedUser, err := a.decryptUserData(data)
		if err != nil {
			a.sendUnauthorizedResponse(w, "failed to decrypt user data", "invalid token")
			return
		}

		var userPayload UserPayload
		if err := json.Unmarshal([]byte(decryptedUser), &userPayload); err != nil {
			a.sendUnauthorizedResponse(w, "failed to unmarshal user payload", "invalid token payload")
			return
		}

		a.logger.Warnf("user payload: %+v", userPayload)
		ctx := a.populateContextWithUserPayload(r.Context(), userPayload)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *authMiddleware) sendUnauthorizedResponse(w http.ResponseWriter, logMessage, responseMessage string) {
	a.logger.Warnf(logMessage)
	res := common.Response[constant.ErrorDefinition]{
		ResponseWriter: w,
		Status:         http.StatusUnauthorized,
		Data: constant.ErrorDefinition{
			Code:    http.StatusUnauthorized,
			Message: responseMessage,
		},
	}
	res.SendJSON()
}

func (a *authMiddleware) populateContextWithUserPayload(ctx context.Context, userPayload UserPayload) context.Context {
	// ctx = context.WithValue(ctx, constant.ContextKey("user"), userPayload)
	ctx = context.WithValue(ctx, constant.ContextKey("user_id"), userPayload.UserID)
	ctx = context.WithValue(ctx, constant.ContextKey("user_code"), userPayload.UserCode)
	ctx = context.WithValue(ctx, constant.ContextKey("full_name"), userPayload.FullName)
	ctx = context.WithValue(ctx, constant.ContextKey("phone_number"), userPayload.PhoneNumber)
	ctx = context.WithValue(ctx, constant.ContextKey("user_email"), userPayload.Email)
	ctx = context.WithValue(ctx, constant.ContextKey("user_realm"), userPayload.Realm)
	ctx = context.WithValue(ctx, constant.ContextKey("ifb_member"), userPayload.MemberType)
	ctx = context.WithValue(ctx, constant.ContextKey("device_uuid"), userPayload.DeviceUUID)
	ctx = context.WithValue(ctx, constant.ContextKey("user_device_linked_date"), userPayload.UserDeviceLinkedDate)
	ctx = context.WithValue(ctx, constant.ContextKey("permissions"), userPayload.Permissions)
	ctx = context.WithValue(ctx, constant.ContextKey("session_expiry"), userPayload.SessionExpiry)
	return ctx
}

func (a *authMiddleware) decryptUserData(data string) (string, error) {
	keyByte := []byte(a.Key)
	ivByte := []byte(a.IV)

	if len(keyByte) != 32 {
		a.logger.Warnf("invalid key byte provided")
		return "", errors.New("key must be 32 bytes for AES-256")
	}
	if len(ivByte) != aes.BlockSize {
		a.logger.Warnf("invalid iv bytes provided must be 16 bytes")
		return "", errors.New("IV must be 16 bytes for AES-256-CBC")
	}

	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		a.logger.Errorf("failed to decode hex", err)
		return "", fmt.Errorf("hex decode failed: %w", err)
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		a.logger.Errorf("new cipher failed", err)
		return "", fmt.Errorf("NewCipher failed: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, ivByte)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// Remove PKCS7 padding
	decrypted, err = a.pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad failed: %w", err)
	}

	result := string(decrypted)
	return result, nil
}

func (a *authMiddleware) pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		a.logger.Warnf("input data is empty")
		return nil, errors.New("input data is empty")
	}
	if len(data)%blockSize != 0 {
		a.logger.Warnf("input length is not multiple of block size")
		return nil, errors.New("input length is not a multiple of block size")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		a.logger.Warnf("invalid padding")
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			a.logger.Warnf("invalid padding bytes")
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:len(data)-padding], nil
}
