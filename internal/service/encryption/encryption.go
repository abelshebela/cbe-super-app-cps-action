package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	dtoEncryption "cbe-super-app-cps-action/internal/constants/dto/encryption"
	"cbe-super-app-cps-action/internal/constants/localization"

	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type EncryptionService struct {
	cfg    *config.VaultConfig
	logger utils.Logger
}

func NewEncryptionService(cfg *config.VaultConfig, logger utils.Logger) *EncryptionService {
	return &EncryptionService{
		cfg:    cfg,
		logger: logger,
	}
}

func GenerateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func SignWithHS256(data string, saltHex string) (string, error) {
	key, err := hex.DecodeString(saltHex)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (e *EncryptionService) LocalEncryptPassword(req dtoEncryption.EncryptionRequest, dataType, userSalt, action string) (dtoEncryption.EncryptionResponse, string, error) {
	ctx := context.Background()
	ctx, span := local_util.TraceLogger(ctx, "service", "LocalEncryptPassword", "Encryption", "LocalEncryptPassword")
	defer span.End()

	if e.cfg == nil {
		e.logger.Errorf("[LocalEncryptPassword] config is empty")
		span.AddEvent("Config is empty", trace.WithAttributes(
			attribute.String("error", localization.ErrConfigIsEmpty.Code),
		))
		return dtoEncryption.EncryptionResponse{}, "", errors.New(localization.ErrConfigIsEmpty.Code)
	}

	key := []byte(e.cfg.Key)
	iv := []byte(e.cfg.IV)

	data := map[string]string{
		"username": req.Username,
		"password": req.Password,
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		e.logger.Errorf("[LocalEncryptPassword] failed to marshal data: %v", err)
		span.AddEvent("Failed to marshal data", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return dtoEncryption.EncryptionResponse{}, "", errors.New(localization.ErrMarshalingData.Code)
	}
	dataStr := string(dataBytes)

	var signedPass, salt string
	if dataType == "password" {
		salt, _ = GenerateSalt(20)
		signedPass, _ = SignWithHS256(dataStr, salt)
	} else {
		signedPass = dataStr
	}

	if action == "login" || action == "change" {
		salt = userSalt
		signedPass, _ = SignWithHS256(dataStr, userSalt)
	}

	if len(key) != 32 || len(iv) != aes.BlockSize {
		e.logger.Errorf("[LocalEncryptPassword] invalid key or IV length")
		span.AddEvent("Invalid key or IV length", trace.WithAttributes(
			attribute.String("error", localization.ErrInvalidKeyOrIv.Code),
		))
		return dtoEncryption.EncryptionResponse{}, salt, errors.New(localization.ErrInvalidKeyOrIv.Code)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		e.logger.Errorf("[LocalEncryptPassword] failed to create cipher: %v", err)
		span.AddEvent("Failed to create cipher", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return dtoEncryption.EncryptionResponse{}, salt, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	padLen := aes.BlockSize - len(signedPass)%aes.BlockSize
	padding := strings.Repeat(string(byte(padLen)), padLen)
	padded := []byte(signedPass + padding)
	encrypted := make([]byte, len(padded))
	mode.CryptBlocks(encrypted, padded)

	e.logger.Infof("[LocalEncryptPassword] password encrypted successfully for username (hashed)")
	return dtoEncryption.EncryptionResponse{Encryption: hex.EncodeToString(encrypted)}, salt, nil
}
