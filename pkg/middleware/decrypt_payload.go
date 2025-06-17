package middleware

import (
	"cbe-super-app-member-auth/pkg/utils"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"

	"context"
)

type DecryptedPayload struct {
	Payload  string `json:"payload"`
	Checksum string `json:"checksum"`
}

// ContextKey is used for storing values in context
type ContextKey string

const DecryptedBodyKey ContextKey = "decryptedBody"

func DecryptAndVerifyPayload(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(" - - - -  decrypt and verify payload - - - - ")
		isNewAPK := r.Header.Get("x-csg-cvd") == "yes"

		var (
			encryptedPayload string
			iv               string
			authTag          string
			timestamp        string
			clientPublicKey  string
		)

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		if isNewAPK {
			var body struct {
				X1 string `json:"x1"`
				Y2 string `json:"y2"`
				Z3 string `json:"z3"`
				A4 string `json:"a4"`
				B5 string `json:"b5"`
			}
			if err := decoder.Decode(&body); err != nil {
				http.Error(w, `{"status":400,"message":"Invalid payload: `+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			encryptedPayload = body.X1
			iv = body.Y2
			authTag = body.Z3
			timestamp = body.A4
			clientPublicKey = body.B5
		} else {
			var body struct {
				EncryptedPayload string `json:"encryptedPayload"`
				Iv               string `json:"iv"`
				AuthTag          string `json:"authTag"`
				Timestamp        string `json:"timestamp"`
				ClientPublicKey  string `json:"clientPublicKey"`
			}
			if err := decoder.Decode(&body); err != nil {
				http.Error(w, `{"status":400,"message":"Invalid payload: `+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			encryptedPayload = body.EncryptedPayload
			iv = body.Iv
			authTag = body.AuthTag
			timestamp = body.Timestamp
			clientPublicKey = body.ClientPublicKey
		}

		vaultPublicKey, err := utils.GetPublicKey()
		if err != nil {
			http.Error(w, `{"status":500,"message":"Failed to get public key"}`, http.StatusInternalServerError)
			return
		}

		// Check for _user in context (if set by previous middleware)
		if userVal := r.Context().Value(ContextKey("_user")); userVal != nil {
			if u, ok := userVal.(map[string]interface{}); ok {
				if pk, ok := u["publicKey"].(string); ok && pk != vaultPublicKey {
					log.Println("Public Key is not the same")
					log.Println("Public Key from request: ", pk)
					log.Println("Public Key from vault: ", vaultPublicKey)
					http.Error(w, `{"status":403,"message":"Invalid Public Key"}`, http.StatusForbidden)
					return
				}
			}
		}
		log.Println("----------Public Key: Verified----------")

		if encryptedPayload == "" || iv == "" || authTag == "" || timestamp == "" || clientPublicKey == "" {
			http.Error(w, `{"status":400,"message":"Invalid payload: Missing required fields."}`, http.StatusBadRequest)
			return
		}

		serverPrivateKey, err := utils.GetPrivateKey()
		if err != nil {
			http.Error(w, `{"status":500,"message":"Failed to get private key"}`, http.StatusInternalServerError)
			return
		}

		sharedSecret, err := computeSharedSecret(serverPrivateKey, clientPublicKey)
		if err != nil {
			http.Error(w, `{"status":500,"message":"Failed to compute shared secret"}`, http.StatusInternalServerError)
			return
		}

		computedAuthTag := computeHMACSHA256(sharedSecret, iv+encryptedPayload)

		log.Println("Computed Auth Tag >>>>>>>>>>>>:", computedAuthTag)
		log.Println("Received Auth Tag >>>>>>>>>>>>>:", authTag)

		if computedAuthTag != authTag {
			http.Error(w, `{"status":400,"message":"Sorry, we could not verify the authenticity of the request."}`, http.StatusBadRequest)
			return
		}
		log.Println("----------Authentication Tag: Verified----------")

		decryptedPayload, err := decryptAES256CBC(sharedSecret, iv, encryptedPayload)
		if err != nil {
			log.Println("Decryption error:", err)
			http.Error(w, `{"status":500,"message":"Failed to decrypt payload."}`, http.StatusInternalServerError)
			return
		}

		var dp DecryptedPayload
		if err := json.Unmarshal([]byte(decryptedPayload), &dp); err != nil {
			http.Error(w, `{"status":400,"message":"Invalid decrypted payload"}`, http.StatusBadRequest)
			return
		}

		computedChecksum := sha256Hex(dp.Payload)
		log.Println("Computed Checksum >>>>>>>>>>>>:", computedChecksum)
		log.Println("Received Checksum >>>>>>>>>>>>>:", dp.Checksum)
		if computedChecksum != dp.Checksum {
			http.Error(w, `{"status":400,"message":"Sorry, we could not verify the integrity of the request."}`, http.StatusBadRequest)
			return
		}
		log.Println("----------Checksum: Verified----------")

		var newBody map[string]interface{}
		if err := json.Unmarshal([]byte(dp.Payload), &newBody); err != nil {
			http.Error(w, `{"status":400,"message":"Invalid payload body"}`, http.StatusBadRequest)
			return
		}

		// Store decrypted body in context for downstream handlers
		ctx := context.WithValue(r.Context(), DecryptedBodyKey, newBody)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func computeSharedSecret(serverPrivHex, clientPubHex string) ([]byte, error) {
	curve := elliptic.P256() // Use secp256k1 if available, else P256
	serverPriv, err := hex.DecodeString(serverPrivHex)
	if err != nil {
		return nil, err
	}
	clientPub, err := hex.DecodeString(clientPubHex)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(curve, clientPub)
	if x == nil || y == nil {
		return nil, errors.New("invalid client public key")
	}
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(serverPriv)
	sharedX, _ := curve.ScalarMult(x, y, serverPriv)
	shared := sharedX.Bytes()
	return padTo32Bytes(shared), nil
}

func padTo32Bytes(b []byte) []byte {
	if len(b) >= 32 {
		return b[:32]
	}
	padded := make([]byte, 32)
	copy(padded[32-len(b):], b)
	return padded
}

func computeHMACSHA256(key []byte, data string) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func decryptAES256CBC(key []byte, ivHex, cipherHex string) (string, error) {
	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return "", err
	}
	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)
	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}
	paddingLen := int(data[len(data)-1])
	if paddingLen == 0 || paddingLen > blockSize {
		return nil, errors.New("invalid padding")
	}
	for i := 0; i < paddingLen; i++ {
		if data[len(data)-1-i] != byte(paddingLen) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-paddingLen], nil
}

func sha256Hex(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}
