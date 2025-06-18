package utils

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenVals struct {
	SessionExpiry string                 `json:"session_expiry"`
	OtherFields   map[string]interface{} `json:"-"`
}

func (t *TokenVals) UnmarshalJSON(data []byte) error {
	type Alias TokenVals
	aux := &struct {
		*Alias
	}{Alias: (*Alias)(t)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	delete(m, "sessionexpiry")
	t.OtherFields = m
	return nil
}

func debuggingLoggerV2(args ...interface{}) {
	// Implement your debug logger here
}

func localEncryptPassword(s string) string {
	// Implement your encryption logic here
	return s
}

func CheckSession(tokenVals *TokenVals, tokenExpiryTime int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			debuggingLoggerV2("----- * * * * session check * * * * ----")
			debuggingLoggerV2("tokenvals: ", tokenVals)

			currentTime := time.Now().UnixMilli()
			sessionExpiry := tokenVals.SessionExpiry
			sessionExpiryTime, err := time.Parse(time.RFC3339, sessionExpiry)
			if err != nil {
				http.Error(w, "invalid sessionexpiry", http.StatusUnauthorized)
				return
			}
			sessionExpiryMillis := sessionExpiryTime.UnixMilli()

			secretKey := os.Getenv("JWTSECRET")
			realm := r.Context().Value("realm")
			var expiry int64
			if realm == "member" {
				expiry, _ = strconv.ParseInt(os.Getenv("APP_SESSIONEXPIREY"), 10, 64)
			} else {
				expiry, _ = strconv.ParseInt(os.Getenv("DASH_SESSIONEXPIREY"), 10, 64)
			}
			sessionThreshold, _ := strconv.ParseInt(os.Getenv("SESSION_THRESHOLD"), 10, 64)
			sessionThreshold = sessionThreshold * 60 * 1000
			tempSessionTimeout, _ := strconv.ParseInt(os.Getenv("_TEMPSESSIONTIMEOUT"), 10, 64)
			tempSessionTimeout = tempSessionTimeout * 60 * 1000
			remainingSessionTime := sessionExpiryMillis - currentTime

			debuggingLoggerV2("token_expiry_time: ", tokenExpiryTime)
			debuggingLoggerV2("_expiry: ", expiry)
			debuggingLoggerV2("currentTime: ", currentTime)
			debuggingLoggerV2("_sessionexpiry: ", sessionExpiry)
			debuggingLoggerV2("_sessionexpiry_time: ", sessionExpiryMillis)
			debuggingLoggerV2("_sessionthreshold: ", sessionThreshold)
			debuggingLoggerV2("_remaining_session_time: ", remainingSessionTime)

			if remainingSessionTime <= 0 {
				http.Error(w, "expired session", http.StatusForbidden)
				return
			}

			if remainingSessionTime <= sessionThreshold {
				newSessionExpiry := time.UnixMilli(currentTime + tempSessionTimeout).UTC()
				newTokenData := make(map[string]interface{})
				for k, v := range tokenVals.OtherFields {
					newTokenData[k] = v
				}
				newTokenData["sessionexpiry"] = newSessionExpiry.Format(time.RFC3339)
				buf, _ := json.Marshal(newTokenData)

				var b bytes.Buffer
				zw := zlib.NewWriter(&b)
				_, _ = zw.Write(buf)
				zw.Close()
				compressed := base64.StdEncoding.EncodeToString(b.Bytes())

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"data": localEncryptPassword(compressed),
					"exp":  time.Now().Unix() + int64Abs(tokenExpiryTime-time.Now().Unix()),
				})
				newToken, _ := token.SignedString([]byte(secretKey))
				w.Header().Set("x_new_token", newToken)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func int64Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
