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

	"github.com/go-chi/cors"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/golang-jwt/jwt/v5"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants"
	customErr "github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/errors"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/response"
)

func CORS() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept", "Authorization", "Content-Type", "X-CSRF-Token",
			"access-control-allow-origin", "x-api-applicationid",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

type JSONResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := JSONResponse{
		Status:  status,
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

type UserPayload struct {
	PhoneNumber string   `json:"phone_number,omitempty"`
	UserRole    string   `json:"user_role,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	BranchCode  []string `json:"branch_code,omitempty"`
	UserCode    string   `json:"user_code,omitempty"`
	FullName    string   `json:"full_name,omitempty"`
	Department  string   `json:"department,omitempty"`
	NextStep    string   `json:"next_step,omitempty"`
}

type authMiddleware struct {
	logger       utils.Logger
	JWTSecretKey string
	Key          string
	IV           string
}

type AuthMiddleware interface {
	AuthenticateToken(next http.Handler) http.Handler
	AuthenticateTempToken(next http.Handler) http.Handler
}

func InitAuthMiddleware(secretKey, key, iv string, logger utils.Logger) AuthMiddleware {
	return &authMiddleware{
		JWTSecretKey: secretKey,
		Key:          key,
		IV:           iv,
		logger:       logger,
	}
}

func (a *authMiddleware) AuthenticateTempToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		bearer := "Bearer "
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, bearer) {
			a.logger.Warnf("bearer token is not provided")
			response.SendErrorResponse(w, customErr.ErrUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, bearer)
		if tokenString == "" {
			a.logger.Warnf("empty token string provided")
			response.SendErrorResponse(w, customErr.ErrUnauthorized)
			return
		}

		data, err := a.validateToken(r.Context(), tokenString)
		if err != nil {
			response.SendErrorResponse(w, err)
			return
		}

		userPayload, err := a.extractUserPayload(r.Context(), data)
		if err != nil {
			response.SendErrorResponse(w, customErr.ErrUnauthorized)
			return
		}

		ctx := a.setUserPayload(r.Context(), userPayload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *authMiddleware) AuthenticateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		bearer := "Bearer "

		if !strings.HasPrefix(authHeader, bearer) {
			a.logger.Warnf("bearer token is not present")
			response.SendErrorResponse(w, customErr.ErrUnauthorized)
			return
		}

		tokenString := authHeader[len(bearer):]

		if tokenString == "" {
			a.logger.Warnf("empty token string provided")
			response.SendErrorResponse(w, customErr.ErrUnauthorized)
			return
		}

		data, err := a.validateToken(r.Context(), tokenString)
		if err != nil {
			response.SendErrorResponse(w, err)
			return
		}

		userPayload, err := a.extractUserPayload(r.Context(), data)
		if err != nil {
			response.SendErrorResponse(w, customErr.ErrUnauthorized)
			return
		}

		ctx := a.setUserPayload(r.Context(), userPayload)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (a *authMiddleware) validateToken(ctx context.Context, tokenString string) (string, error) {
	jwtSecret := []byte(a.JWTSecretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			a.logger.Errorf("unexpected signing method used")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		a.logger.Errorf("invalid or expired token: %v", err)
		return "", customErr.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		a.logger.Errorf("failed to cast to map claims")
		return "", customErr.ErrUnauthorized
	}

	data, ok := claims["data"].(string)
	if !ok || data == "" {
		a.logger.Errorf("invalid token or data not present")
		return "", customErr.ErrUnauthorized
	}

	return data, nil
}

func (a *authMiddleware) extractUserPayload(ctx context.Context, data string) (UserPayload, error) {
	decryptedUser, err := a.decryptUserData(ctx, data)
	if err != nil {
		return UserPayload{}, customErr.ErrUnauthorized
	}

	var userPayload UserPayload
	err = json.Unmarshal([]byte(decryptedUser), &userPayload)
	if err != nil {
		a.logger.Errorf("failed to unmarshal user payload: %v", err)
		return UserPayload{}, customErr.ErrUnauthorized
	}

	return userPayload, nil
}

func (a *authMiddleware) setUserPayload(ctx context.Context, userPayload UserPayload) context.Context {
	ctx = context.WithValue(ctx, constants.ContextKey("branch_code"), userPayload.BranchCode)
	ctx = context.WithValue(ctx, constants.ContextKey("user_role"), userPayload.UserRole)
	ctx = context.WithValue(ctx, constants.ContextKey("user_id"), userPayload.UserID)
	ctx = context.WithValue(ctx, constants.ContextKey("phone_number"), userPayload.PhoneNumber)
	ctx = context.WithValue(ctx, constants.ContextKey("user_code"), userPayload.UserCode)
	ctx = context.WithValue(ctx, constants.ContextKey("full_name"), userPayload.FullName)
	ctx = context.WithValue(ctx, constants.ContextKey("department"), userPayload.Department)
	ctx = context.WithValue(ctx, constants.ContextKey("next_step"), userPayload.NextStep)
	return ctx
}

func (a *authMiddleware) decryptUserData(ctx context.Context, data string) (string, error) {
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
		a.logger.Errorf("failed to decode hex: %v", err)
		return "", fmt.Errorf("hex decode failed: %w", err)
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		a.logger.Errorf("new cipher failed: %v", err)
		return "", fmt.Errorf("NewCipher failed: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, ivByte)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// Remove PKCS7 padding
	decrypted, err = a.pkcs7Unpad(ctx, decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad failed: %w", err)
	}

	result := string(decrypted)
	return result, nil
}

func (a *authMiddleware) pkcs7Unpad(ctx context.Context, data []byte, blockSize int) ([]byte, error) {
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
