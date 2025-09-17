package utils

/*
import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"strconv"
	"time"
	//"cbe-super-app-member-auth/pkg/config"

	"gitlab.com/bersufekadgetachew/cbe-super-app-member-auth/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-member-auth/internal/domain/auth/entities"
	//"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	//"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/budget/entities"
)

type Permission struct {
	ID             string `bson:"_id" json:"_id"`
	PermissionName string `bson:"permission_name" json:"permission_name"`
}

type PermissionGroup struct {
	ID          string       `bson:"_id" json:"_id"`
	GroupName   string       `bson:"group_name" json:"group_name"`
	Permissions []Permission `bson:"permissions" json:"permissions"`
}

var valueMapping = map[string]string{
	"sms":   "phone_number",
	"email": "email",
	"both":  "email_and_phone",
}

func VerifyPassword(password, inputPassword string, key, iv []byte) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in verifyPassword:", r)
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("AES cipher error:", err)
		return false
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	ciphertext, err := hex.DecodeString(password)
	if err != nil {
		log.Println("Hex decode error:", err)
		return false
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		log.Println("Ciphertext is not a multiple of the block size")
		return false
	}
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)
	decrypted = bytes.Trim(decrypted, "\x00")
	plain := string(decrypted)
	return plain == inputPassword
}

func CompressJSON(data interface{}) (string, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err = w.Write(b)
	if err != nil {
		return "", err
	}
	w.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

/*
func TempTokenMaker(user entities.User, permissions []string, otpFor string, additional map[string]interface{}) (string, error) {
	env, _ := config.Load()
	timeSession, _ := strconv.Atoi(env.TempSessionTimeout)
	sessionExpiry := time.Now().Unix() + int64(timeSession)*60

	attributes := map[string]interface{}{
		"userId":        user.ID,
		"userRealm":     user.Realm,
		"fullname":      user.FullName,
		"username":      user.Username,
		"branchCode":    user.BranchCode,
		"userCode":      user.UserCode,
		"phoneNumber":   user.PhoneNumber,
		"sourceApp":     "memberapp",
		"sessionexpiry": sessionExpiry,
	}

	for k, v := range additional {
		attributes[k] = v
	}

	attrBytes, err := json.Marshal(attributes)
	if err != nil {
		return "", err
	}
	encrypted, _, _ := LocalEncryptPassword(string(attrBytes), "token", "token", "token")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": encrypted,
		// "exp":  jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	})

	return token.SignedString([]byte(env.JWTSecretKey))
}
// type User struct {
// 	ID           string `json:"_id"`
// 	UserCode     string `json:"userCode"`
// 	FullName     string `json:"fullName"`
// 	PhoneNumber  string `json:"phoneNumber"`
// 	Realm        string `json:"realm"`
// 	DeviceUUID   string `json:"deviceUUID"`
// }

// // Replace this with your config structure
// var Config = struct {
// 	JWTSecret          string
// 	JWTExpiry          time.Duration // Example: time.Minute * 30
// 	TempSessionTimeout time.Duration // Example: time.Minute * 10
// }{
// 	JWTSecret:          "your-secret",
// 	JWTExpiry:          time.Minute * 30,
// 	TempSessionTimeout: time.Minute * 10,
// }

// // Encrypt function placeholder
// func LocalEncryptPassword(input string) string {
// 	// Add your encryption logic here
// 	return input
// }

// // Handle session logic
// func HandleSession(user User) {
// 	// Implement your session handler here
// }

func TempTokenMaker(user *entities.User, permissions []string, publicKey string) (string, error) {
	env, _ := config.Load()
	fromLinkAccount := len(user.ID) == 0
	realm := "member"
	if !fromLinkAccount {
		realm = string(user.Realm)
	}
	tempSession, _ := strconv.Atoi(env.TempSessionTimeout)
	attributes := map[string]interface{}{
		"userid":          user.ID,
		"usercode":        user.UserCode,
		"fullname":        user.FullName,
		"phonenumber":     user.PhoneNumber,
		"userpermissions": permissions,
		"userrealm":       realm,
		"fromlinkaccount": fromLinkAccount,
		"deviceuuid":      user.Device.DeviceUUID,

		"sessionexpiry": time.Now().Add(time.Duration(tempSession)).Format(time.RFC3339),
		"publicKey":     publicKey,
	}

	attrBytes, err := json.Marshal(attributes)
	if err != nil {
		return "", err
	}

	encrypted, _, _ := LocalEncryptPassword(string(attrBytes), "token", "", "")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": encrypted,
	})

	signedToken, err := token.SignedString([]byte(env.JWTSecretKey))
	if err != nil {
		return "", err
	}

	// HandleSession(user)
	return signedToken, nil
}

// type SessionQuery struct {
// 	UserID string `json:"userID"`
// }

// type SessionUpdate struct {
// 	LastActivity  time.Time `json:"lastActivity"`
// 	SessionExpiry time.Time `json:"sessionExpiry"`
// }

// type SessionData struct {
// 	UserID        string    `json:"userID"`
// 	LastActivity  time.Time `json:"lastActivity"`
// 	SessionExpiry time.Time `json:"sessionExpiry"`
// }

// type SessionResponse struct {
// 	StatusCode int
// 	Body       struct {
// 		Data  interface{}
// 		Error string
// 	}
// }

// func SessionDal(method string, query SessionQuery, update *SessionUpdate, create *SessionData) (SessionResponse, error) {
// 	// Placeholder logic
// 	return SessionResponse{StatusCode: 404}, nil
// }

func TokenMaker(user *entities.User, tokentype, sourceapp, publicKey, deviceuuid string, setpin bool, platform string, fcmtoken string) (string, map[string]interface{}, error) {
	env, _ := config.Load()
	secret := []byte(env.JWTSecretKey)
	permissions := []string{}
	primaryauth := ""
	if user.Realm == "member" {
		for k, v := range valueMapping {
			if v == string(user.PrimaryAuthentication) {
				primaryauth = k
				break
			}
		}
	}
	sessionExp, _ := strconv.Atoi(env.AppSessionExpiry)
	claims := map[string]interface{}{
		"user_id":                 user.ID,
		"user_code":               user.UserCode,
		"full_name":               user.FullName,
		"organization_id":         user.OrganizationID,
		"phone_number":            user.PhoneNumber,
		"user_email":              user.Email,
		"user_realm":              user.Realm,
		"ifb_member":              user.MemberType == "ifb",
		"device_uuid":             user.Device.DeviceUUID,
		"user_device_linked_date": "",
		// "user_device_linked_date": user.DeviceLinkedDate,
		"permissions":    permissions,
		"primary_auth":   primaryauth,
		"session_expiry": time.Now().Add(time.Duration(sessionExp) * time.Minute).Format(time.RFC3339),
		"public_key":     publicKey,
	}
	jsonBytes, _ := json.Marshal(claims)
	compressed, err := CompressAndEncrypt(jsonBytes)
	if err != nil {
		return "", map[string]interface{}{}, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": compressed,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // adjust as needed
	})
	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", map[string]interface{}{}, err
	}
	updateUser := updateUser(user, signedToken, sourceapp, deviceuuid, setpin, platform, fcmtoken)
	return signedToken, updateUser, nil
}
func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func CompressAndEncrypt(data []byte) (string, error) {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return "", err
	}
	w.Close()
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

func updateUser(user *entities.User, token, sourceapp, deviceuuid string, setpin bool, platform string, fcmtoken string) map[string]interface{} {
	updatedata := map[string]interface{}{
		"session_expires_on":  user.SessionExpirsOn,
		"login_attempt_count": 0,
		"last_login":          time.Now().UTC().Format(time.RFC3339),
		"last_modified":       time.Now().UTC().Format(time.RFC3339),
	}

	if user.Realm == "merchant" && sourceapp == "agentapp" {
		updatedata["device_uuid"] = deviceuuid
	}
	if user.Realm == "member" {
		if setpin {
			updatedata["device_platform"] = stringToUpper(platform)
		}
		if fcmtoken != "" {
			updatedata["push_token"] = fcmtoken
		}
	}
	return updatedata
}

func stringToUpper(s string) string {
	return string(bytes.ToUpper([]byte(s)))
}
*/
