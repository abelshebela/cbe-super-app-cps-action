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
	"github.com/spf13/viper"
)

type UserPayload struct {
	PhoneNumber string   `json:"phoneNumber,omitempty"`
	UserRole    string   `json:"userrole,omitempty"`
	UserID      string   `json:"userId,omitempty"`
	BranchCode  []string `json:"branchCode,omitempty"`
	FullName    string   `json:"fullname,omitempty"`
	HomeBranch  string   `json:"homeBranch,omitempty"`
	Username    string   `json:"username,omitempty"`
	SourceApp   string   `json:"sourceApp,omitempty"`
}

type ContextKey string

func AccessControl(allowedRoles []string) func(http.Handler) http.Handler {
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := r.Context().Value(ContextKey("user_role")).(string)
			if !ok || role == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if _, allowed := roleSet[strings.ToUpper(role)]; !allowed {
				http.Error(w, "Action not allowed", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthenticateToken(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")

        tokenString := ""
        if strings.HasPrefix(authHeader, "Bearer ") {
            tokenString = strings.TrimPrefix(authHeader, "Bearer ")
        }

        jwtSecret := []byte(viper.GetString("JwtSecretKey"))

        if tokenString == "" {
            fmt.Println("Access token required: token string is empty")
            http.Error(w, "Access token required", http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })
        if err != nil {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }
        if !token.Valid {
            http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        data, ok := claims["data"].(string)
        if !ok || data == "" {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        decryptedUser, err := decryptUserData(data)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        var userPayload UserPayload
        err = json.Unmarshal([]byte(decryptedUser), &userPayload)
        if err != nil {
            http.Error(w, "Invalid token payload", http.StatusUnauthorized)
            return
        }


        ctx := context.WithValue(r.Context(), ContextKey("user_payload"), userPayload)
        ctx = context.WithValue(ctx, ContextKey("user_role"), userPayload.UserRole)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}
func decryptUserData(data string) (string, error) {
	keyByte := []byte(viper.GetString("Key"))
	ivByte := []byte(viper.GetString("IV"))

	if len(keyByte) != 32 {
		return "", errors.New("key must be 32 bytes for AES-256")
	}
	if len(ivByte) != aes.BlockSize {
		return "", errors.New("IV must be 16 bytes for AES-256-CBC")
	}

	ciphertext, err := hex.DecodeString(data)
	if err != nil {
		return "", fmt.Errorf("hex decode failed: %w", err)
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", fmt.Errorf("NewCipher failed: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, ivByte)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// Remove PKCS7 padding
	decrypted, err = pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpad failed: %w", err)
	}

	result := string(decrypted)
	return result, nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}
	if len(data)%blockSize != 0 {
		return nil, errors.New("input length is not a multiple of block size")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("invalid padding bytes")
		}
	}
	return data[:len(data)-padding], nil
}
