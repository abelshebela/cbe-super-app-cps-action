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
	constant "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/common"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type UserPayload struct {
	PhoneNumber string   `json:"phone_number,omitempty"`
	UserRole    string   `json:"user_role,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	BranchCode  []string `json:"branch_code,omitempty"`
	UserCode    string   `json:"user_code,omitempty"`
	FullName    string   `json:"full_name,omitempty"`
	Department  string   `json:"department,omitempty"`
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
		authHeader := r.Header.Get("Authorization")
		bearer := "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			a.logger.Warnf("bearer token is not present")
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

		tokenString := authHeader[len(bearer):]

		jwtSecret := []byte(a.JWTSecretKey)

		if tokenString == "" {
			a.logger.Warnf("empty token string provided")
			res := common.Response[constant.ErrorDefinition]{
				ResponseWriter: w,
				Status:         http.StatusUnauthorized,
				Data: constant.ErrorDefinition{
					Code:    http.StatusUnauthorized,
					Message: "access token required",
				},
			}
			res.SendJSON()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				a.logger.Errorf("unexpected signing method used")
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			a.logger.Errorf("invalid or expired token", err)
			res := common.Response[constant.ErrorDefinition]{
				ResponseWriter: w,
				Status:         http.StatusUnauthorized,
				Data: constant.ErrorDefinition{
					Code:    http.StatusUnauthorized,
					Message: "invalid or expired token",
				},
			}
			res.SendJSON()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			a.logger.Errorf("failed to cast to map claims")
			res := common.Response[constant.ErrorDefinition]{
				ResponseWriter: w,
				Status:         http.StatusUnauthorized,
				Data: constant.ErrorDefinition{
					Code:    http.StatusUnauthorized,
					Message: "invalid token",
				},
			}
			res.SendJSON()
			return
		}

		data, ok := claims["data"].(string)
		if !ok || data == "" {
			a.logger.Errorf("invalid token or data not present")
			res := common.Response[constant.ErrorDefinition]{
				ResponseWriter: w,
				Status:         http.StatusUnauthorized,
				Data: constant.ErrorDefinition{
					Code:    http.StatusUnauthorized,
					Message: "invalid token",
				},
			}
			res.SendJSON()
			return
		}

		decryptedUser, err := a.decryptUserData(data)
		if err != nil {
			res := common.Response[constant.ErrorDefinition]{
				ResponseWriter: w,
				Status:         http.StatusUnauthorized,
				Data: constant.ErrorDefinition{
					Code:    http.StatusUnauthorized,
					Message: "invalid token",
				},
			}
			res.SendJSON()
			return
		}

		var userPayload UserPayload
		err = json.Unmarshal([]byte(decryptedUser), &userPayload)
		if err != nil {
			a.logger.Errorf("failed to unmarshal user payload", err)
			res := common.Response[constant.ErrorDefinition]{
				ResponseWriter: w,
				Status:         http.StatusUnauthorized,
				Data: constant.ErrorDefinition{
					Code:    http.StatusUnauthorized,
					Message: "invalid token payload",
				},
			}
			res.SendJSON()
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, constant.ContextKey("branch_code"), userPayload.BranchCode)
		ctx = context.WithValue(ctx, constant.ContextKey("user_role"), userPayload.UserRole)
		ctx = context.WithValue(ctx, constant.ContextKey("user_id"), userPayload.UserID)
		ctx = context.WithValue(ctx, constant.ContextKey("phone_number"), userPayload.PhoneNumber)
		ctx = context.WithValue(ctx, constant.ContextKey("user_code"), userPayload.UserCode)
		ctx = context.WithValue(ctx, constant.ContextKey("full_name"), userPayload.FullName)
		ctx = context.WithValue(ctx, constant.ContextKey("department"), userPayload.Department)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
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
