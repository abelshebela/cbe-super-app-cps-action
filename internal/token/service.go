package token

// import (
// 	"encoding/json"
// 	"errors"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"cbe-super-app-cps-action/internal/constants"
// 	"cbe-super-app-cps-action/internal/constants/localization"
// 	"cbe-super-app-cps-action/internal/constants/model"
// 	"cbe-super-app-cps-action/pkgs/utils"

// 	"github.com/golang-jwt/jwt/v5"

// 	"crypto/aes"
// 	"crypto/cipher"
// 	"encoding/hex"

// 	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
// )

// type TokenService struct {
// 	cfg *config.VaultConfig
// }

// func NewTokenService(config *config.VaultConfig) *TokenService {
// 	return &TokenService{cfg: config}
// }

// func (t *TokenService) TempTokenMaker(user *model.User, additional map[string]interface{}, action string) (string, error) {

// 	timeSession := 0
// 	sessionExpiry := time.Now().Unix() + int64(timeSession)*60

// 	attributes := ClaimBuilder(user, sessionExpiry, action, additional)

// 	attrBytes, err := json.Marshal(attributes)
// 	if err != nil {
// 		return "", err
// 	}

// 	encrypted, _, _ := t.LocalEncryptPassword(string(attrBytes), "token", "token", "token")

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"data": encrypted,
// 	})

// 	return token.SignedString([]byte(t.cfg.JwtSecretKey))
// }

// func (t *TokenService) TokenMaker(user *model.User, env *config.VaultConfig, tokenType ...string) (string, error) {

// 	timeSession, _ := strconv.Atoi(env.TempSessionTimeout)
// 	sessionExpiry := time.Now().Unix() + int64(timeSession)*60

// 	attributes := PermanentClaimBuilder(user, sessionExpiry, nil)

// 	typeVal := constants.Permanent
// 	if len(tokenType) > 0 && tokenType[0] != "" {
// 		typeVal = tokenType[0]
// 	}
// 	attributes[constants.TokenType] = typeVal

// 	attrBytes, err := json.Marshal(attributes)
// 	if err != nil {
// 		return "", err
// 	}

// 	encrypted, _, _ := t.LocalEncryptPassword(string(attrBytes), constants.Token, constants.Token, constants.Token)

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"data": encrypted,
// 	})
// 	return token.SignedString([]byte(env.JwtSecretKey))
// }

// func (t *TokenService) LocalEncryptPassword(password string, dataType string, userSalt string, action string) (string, string, error) {

// 	var signedPass, salt string
// 	if dataType == constants.Password {
// 		salt, _ = utils.GenerateSalt(20)
// 		signedPass, _ = utils.SignWithHS256(password, salt)
// 	} else {
// 		signedPass = password
// 	}

// 	if action == constants.Login || action == constants.Change {
// 		salt = userSalt
// 		signedPass, _ = utils.SignWithHS256(password, userSalt)
// 	}

// 	key := []byte(t.cfg.Key)
// 	iv := []byte(t.cfg.IV)
// 	if len(key) != 32 || len(iv) != aes.BlockSize {
// 		return "", salt, errors.New(localization.ErrorInvalidKey.Code)
// 	}

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", salt, err
// 	}

// 	mode := cipher.NewCBCEncrypter(block, iv)
// 	padLen := aes.BlockSize - len(signedPass)%aes.BlockSize
// 	padding := strings.Repeat(string(byte(padLen)), padLen)
// 	padded := []byte(signedPass + padding)
// 	encrypted := make([]byte, len(padded))
// 	mode.CryptBlocks(encrypted, padded)

// 	return hex.EncodeToString(encrypted), salt, nil
// }

// func (t *TokenService) LocalDecryptPassword(encryptedHex string, env *config.VaultConfig) (string, error) {

// 	key := []byte(t.cfg.Key)
// 	iv := []byte(t.cfg.IV)

// 	if len(key) != 32 || len(iv) != aes.BlockSize {
// 		return "", errors.New(localization.ErrorInvalidKey.Code)
// 	}

// 	encrypted, err := hex.DecodeString(encryptedHex)
// 	if err != nil {
// 		return "", err
// 	}
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	if len(encrypted)%aes.BlockSize != 0 {
// 		return "", errors.New(localization.ErrorInvalidEncData.Code)
// 	}

// 	mode := cipher.NewCBCDecrypter(block, iv)
// 	decrypted := make([]byte, len(encrypted))
// 	mode.CryptBlocks(decrypted, encrypted)
// 	// Remove PKCS#7 padding
// 	padLen := int(decrypted[len(decrypted)-1])
// 	if padLen > aes.BlockSize || padLen == 0 {
// 		return "", errors.New(localization.ErrorInvalidPadding.Code)
// 	}
// 	return string(decrypted[:len(decrypted)-padLen]), nil
// }
