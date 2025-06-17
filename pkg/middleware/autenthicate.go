package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cbe-super-app-member-auth/pkg/common"
	"cbe-super-app-member-auth/pkg/config"

	// "your_project/dal/sessiondal"
	// "your_project/dal/userdal"
	"cbe-super-app-member-auth/pkg/utils"

	"context"

	"github.com/golang-jwt/jwt/v5"
)

type DecodedUser struct {
	FullName        string   `json:"fullname"`
	PhoneNumber     string   `json:"phoneNumber"`
	SessionExpiry   int64    `json:"sessionexpiry"`
	SourceApp       string   `json:"source_app"`
	UserID          string   `json:"userId"`
	UserName        string   `json:"username"`
	UserRole        string   `json:"role"`
	UserPermissions []string `json:"userpermissions"`
	UserRealm       string   `json:"userrealm"`
	BranchCode      []string `json:"branch_code"`
	HomeBranch      string   `json:"home_branch"`
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg, _ := config.Load()
		authorization := r.Header.Get("Authorization")
		returndata := make(map[string]interface{})
		data := make(map[string]interface{})
		body, err := io.ReadAll(r.Body)

		if err != nil {
			returndata["message"] = common.DefineError.General["INVALID_REQ"].Message
			returndata["code"] = common.DefineError.General["INVALID_REQ"].Code
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseMaker(returndata, w)
			return
		}

		if err := json.Unmarshal([]byte(body), &data); err != nil {
			returndata["message"] = common.DefineError.General["INVALID_JSON_PAYLOAD"].Message
			returndata["code"] = common.DefineError.General["INVALID_JSON_PAYLOAD"].Code
			w.WriteHeader(http.StatusBadRequest)
			utils.ResponseMaker(returndata, w)
			return
		}

		if authorization == "" {
			returndata["message"] = common.DefineError.General["UNAUTHORIZED"].Message
			returndata["code"] = common.DefineError.General["UNAUTHORIZED"].Code
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseMaker(returndata, w)
			return
		}

		splitToken := strings.Split(authorization, " ")
		if len(splitToken) != 2 || strings.ToLower(splitToken[0]) != "bearer" {
			returndata["message"] = common.DefineError.Auth["INVALID_BEARER"].Message
			returndata["code"] = common.DefineError.Auth["INVALID_BEARER"].Code
			w.WriteHeader(http.StatusUnauthorized)
			utils.ResponseMaker(returndata, w)
			return
		}

		tokenStr := splitToken[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			returndata["message"] = common.DefineError.Auth["USE_RIGHT_AUTH"].Message
			returndata["code"] = common.DefineError.Auth["USE_RIGHT_AUTH"].Code
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseMaker(returndata, w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			returndata["message"] = common.DefineError.Auth["INVALID_CLAIM"].Message
			returndata["code"] = common.DefineError.Auth["INVALID_CLAIM"].Code
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseMaker(returndata, w)
			return
		}

		dataStr, ok := claims["data"].(string)
		if !ok {
			returndata["message"] = common.DefineError.Auth["INVALID_TOKEN_DATA"].Message
			returndata["code"] = common.DefineError.Auth["INVALID_TOKEN_DATA"].Code
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseMaker(returndata, w)
			return
		}

		decrypted, err := utils.LocalDecryptPassword(dataStr)
		if err != nil {
			returndata["message"] = common.DefineError.Auth["UNABLE_TO_DYCRYPT_TOKEN"].Message
			returndata["code"] = common.DefineError.Auth["UNABLE_TO_DYCRYPT_TOKEN"].Code
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseMaker(returndata, w)
			return

		}

		var decoded DecodedUser
		err = json.Unmarshal([]byte(decrypted), &decoded)

		if err != nil {
			returndata["message"] = common.DefineError.General["INVALID_TOKEN_FORMAT"].Message
			returndata["code"] = common.DefineError.General["INVALID_TOKEN_FORMAT"].Code
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseMaker(returndata, w)
			return
		}

		if decoded.UserName != data["username"] {
			returndata["message"] = common.DefineError.Auth["INVALID_BEARER"].Message
			returndata["code"] = common.DefineError.Auth["INVALID_BEARER"].Code
			w.WriteHeader(http.StatusForbidden)
			utils.ResponseMaker(returndata, w)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(body))
		r = r.WithContext(context.WithValue(r.Context(), "user", &decoded))
		next.ServeHTTP(w, r)
	})
}
