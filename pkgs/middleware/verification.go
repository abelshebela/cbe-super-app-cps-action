package middleware

// import (
// 	"context"
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/elliptic"
// 	"crypto/hmac"
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"net/http"
// 	"time"

// )

// type IncomingPayload struct {
// 	EncryptedPayload string `json:"encryptedPayload"`
// 	IV               string `json:"iv"`
// 	AuthTag          string `json:"authTag"`
// 	Timestamp        int64  `json:"timestamp"`
// 	ClientPublicKey  string `json:"clientPublicKey"`
// }

// type DecryptedPayload struct {
// 	Payload  string `json:"payload"`
// 	Checksum string `json:"checksum"`
// }

// func getPrivateKeyHex() (string, error) {
// 	return "YOUR_PRIVATE_KEY_HEX", nil // Replace with your actual key retrieval
// }

// func DecryptAndVerifyMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		var body IncomingPayload
// 		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
// 			http.Error(w, "Invalid JSON", http.StatusBadRequest)
// 			return
// 		}

// 		if body.EncryptedPayload == "" || body.IV == "" || body.AuthTag == "" ||
// 			body.Timestamp == 0 || body.ClientPublicKey == "" {
// 			http.Error(w, "Missing required fields", http.StatusBadRequest)
// 			return
// 		}

// 		loc, _ := time.LoadLocation("Africa/Nairobi")
// 		now := time.Now().In(loc)
// 		clientTime := time.UnixMilli(body.Timestamp).In(loc)

// 		if now.Sub(clientTime) > 4*time.Minute {
// 			http.Error(w, "Expired request", http.StatusBadRequest)
// 			return
// 		}

// 		privHex, _ := getPrivateKeyHex()
// 		privBytes, _ := hex.DecodeString(privHex)
// 		clientPubBytes, _ := hex.DecodeString(body.ClientPublicKey)

// 		curve := elliptic.P256()
// 		x, y := elliptic.Unmarshal(curve, clientPubBytes)
// 		if x == nil {
// 			http.Error(w, "Invalid client public key", http.StatusBadRequest)
// 			return
// 		}

// 		d := new(big.Int).SetBytes(privBytes)
// 		sharedX, _ := curve.ScalarMult(x, y, d.Bytes())
// 		sharedKey := sha256.Sum256(sharedX.Bytes())

// 		ivBytes, _ := hex.DecodeString(body.IV)

// 		mac := hmac.New(sha256.New, sharedKey[:])
// 		mac.Write(append(ivBytes, []byte(body.EncryptedPayload)...))

// 		if hex.EncodeToString(mac.Sum(nil)) != body.AuthTag {
// 			http.Error(w, "Auth tag mismatch", http.StatusBadRequest)
// 			return
// 		}

// 		cipherText, _ := hex.DecodeString(body.EncryptedPayload)
// 		block, _ := aes.NewCipher(sharedKey[:])
// 		mode := cipher.NewCBCDecrypter(block, ivBytes)
// 		decrypted := make([]byte, len(cipherText))
// 		mode.CryptBlocks(decrypted, cipherText)

// 		padLen := int(decrypted[len(decrypted)-1])
// 		decrypted = decrypted[:len(decrypted)-padLen]

// 		var payloadParsed DecryptedPayload
// 		if err := json.Unmarshal(decrypted, &payloadParsed); err != nil {
// 			http.Error(w, "Invalid decrypted JSON", http.StatusBadRequest)
// 			return
// 		}

// 		checksum := sha256.Sum256([]byte(payloadParsed.Payload))
// 		if hex.EncodeToString(checksum[:]) != payloadParsed.Checksum {
// 			http.Error(w, "Checksum mismatch", http.StatusBadRequest)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "decrypted", payloadParsed.Payload)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

// // func secureHandler(w http.ResponseWriter, r *http.Request) {
// // 	decrypted := r.Context().Value("decrypted")
// // 	fmt.Fprintf(w, "Decrypted payload: %s", decrypted)
// // }

// // func main() {
// // 	r := chi.NewRouter()
// // 	r.With(DecryptAndVerifyMiddleware).Post("/secure", secureHandler)

// // 	fmt.Println("Server listening on http://localhost:8080")
// // 	log.Fatal(http.ListenAndServe(":8080", r))
// // }
