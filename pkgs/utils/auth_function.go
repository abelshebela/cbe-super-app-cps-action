package utils

import (
	"bytes"
	"cbe-super-app-member-users/pkgs/entities"
	"fmt"
	"strconv"

	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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

func TempTokenMaker(user *entities.User, permissions []string, otpFor string, additional map[string]interface{}, action string, env *config.VaultConfig) (string, error) {
	timeSession := 0
	sessionExpiry := time.Now().Unix() + int64(timeSession)*60
	attributes := map[string]interface{}{
		"user_id":                 user.ID,
		"user_code":               user.UserCode,
		"full_name":               user.FullName,
		"phone_number":            user.PhoneNumber,
		"user_email":              user.Email,
		"user_realm":              user.Realm,
		"ifb_member":              user.MemberType == "ifb",
		"device_uuid":             user.Device.DeviceUUID,
		"user_device_linked_date": "",
		"permissions":             permissions,
		"session_expiry":          sessionExpiry,
	}

	if additional != nil {
		for k, v := range additional {
			attributes[k] = v
		}
	}

	attrBytes, err := json.Marshal(attributes)
	if err != nil {
		return "", err
	}
	encrypted, _, _ := LocalEncryptPassword(string(attrBytes), "token", "token", "token", env)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": encrypted,
	})

	return token.SignedString([]byte(env.JwtSecretKey))
}

func TokenMaker(user entities.User, permissions []string, otpFor string) (string, error) {
	env, _ := config.Load()
	timeSession, _ := strconv.Atoi(env.TempSessionTimeout)
	sessionExpiry := time.Now().Unix() + int64(timeSession)*60

	attributes := map[string]interface{}{
		"user_id":                 user.ID,
		"user_code":               user.UserCode,
		"full_name":               user.FullName,
		"phone_number":            user.PhoneNumber,
		"user_email":              user.Email,
		"user_realm":              user.Realm,
		"ifb_member":              user.MemberType == "ifb",
		"device_uuid":             user.Device.DeviceUUID,
		"user_device_linked_date": "",
		"permissions":             permissions,
		"session_expiry":          sessionExpiry,
	}

	attrBytes, err := json.Marshal(attributes)
	if err != nil {
		return "", err
	}

	encrypted, _, _ := LocalEncryptPassword(string(attrBytes), "token", "token", "token", nil)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": encrypted,
	})
	fmt.Println("Encrypted Data:", env.JwtSecretKey)
	return token.SignedString([]byte(env.JwtSecretKey))
}

func mustJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// Helper function to check if permission exists in permissions slice
func containsPermission(permissions []string, permission string) bool {
	for _, p := range permissions {
		if p == permission {
			return true
		}
	}
	return false
}
