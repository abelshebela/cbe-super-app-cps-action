package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cbe-super-app-member-users/pkgs/common"

	"cbe-super-app-member-users/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"context"

	"github.com/golang-jwt/jwt/v5"
)

// type DecodedUser struct {
// 	FullName        string   `json:"full_name"`
// 	PhoneNumber     string   `json:"phone_number"`
// 	SessionExpiry   int64    `json:"session_expiry"`
// 	SourceApp       string   `json:"source_app"`
// 	UserID          string   `json:"userId"`
// 	UserName        string   `json:"username"`
// 	UserRole        string   `json:"role"`
// 	UserPermissions []string `json:"user_permissions"`
// 	UserRealm       string   `json:"user_realm"`
// 	BranchCode      []string `json:"branch_code"`
// 	HomeBranch      string   `json:"home_branch"`
// }

type DecodedUser struct {
	UserID               string   `json:"user_id"`
	UserCode             string   `json:"user_code"`
	FullName             string   `json:"full_name"`
	OrganizationID       string   `json:"organization_id"`
	PhoneNumber          string   `json:"phone_number"`
	UserEmail            string   `json:"user_email"`
	UserRealm            string   `json:"user_realm"`
	IFBMember            bool     `json:"ifb_member"`
	DeviceUUID           string   `json:"device_uuid"`
	UserDeviceLinkedDate string   `json:"user_device_linked_date"` // Consider using time.Time if not always empty
	Permissions          []string `json:"permissions"`
	PrimaryAuth          bool     `json:"primary_auth"`
	SessionExpiry        string   `json:"session_expiry"` // Or time.Time if parsed
	UserName             string   `json:"username"`
	PublicKey            string   `json:"public_key"`
}

func Authenticate(next http.Handler, cfg *config.VaultConfig) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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
			return []byte(cfg.JwtSecretKey), nil
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

		decrypted, err := utils.LocalDecryptPassword(dataStr, cfg)
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
