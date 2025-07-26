package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Ed25519KeyGen struct{}
type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

type KeyGeneratorService interface {
	GenerateKeyPair() (KeyPair, error)
	Sign(data []byte, privateKey string) ([]byte, error)
	Verify(data, signature []byte, publicKey string) (bool, error)
	GenerateFabricID() string
	GenerateNumericCode(length int) (string, error)
	GenerateAppSecret(byteLength int) (string, error)
	HashAppSecret(secret string) (string, error)
	VerifyAppSecret(hashedSecret, inputSecret string) error
}

func NewKeyGenerator() KeyGeneratorService {
	return &Ed25519KeyGen{}
}

func (g *Ed25519KeyGen) GenerateKeyPair() (KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, fmt.Errorf(utils.UnhandledServerError)
	}
	return KeyPair{
		PublicKey:  base64.StdEncoding.EncodeToString(publicKey),
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
	}, nil
}

func (g *Ed25519KeyGen) Sign(data []byte, privateKey string) ([]byte, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf(utils.UnhandledServerError)
	}
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf(utils.UnhandledServerError)
	}
	return ed25519.Sign(privateKeyBytes, data), nil
}

func (g *Ed25519KeyGen) Verify(data, signature []byte, publicKey string) (bool, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, fmt.Errorf(utils.UnhandledServerError)
	}
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return false, fmt.Errorf(utils.UnhandledServerError)
	}
	return ed25519.Verify(publicKeyBytes, data, signature), nil
}

func (g *Ed25519KeyGen) GenerateFabricID() string {
	return uuid.New().String()
}

func (g *Ed25519KeyGen) GenerateNumericCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf(utils.UnhandledServerError)
	}

	digits := "0123456789"
	max := big.NewInt(int64(len(digits)))
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf(utils.UnhandledServerError)
		}
		code[i] = digits[n.Int64()]
	}

	return string(code), nil
}

func (g *Ed25519KeyGen) GenerateAppSecret(byteLength int) (string, error) {
	if byteLength < 32 {
		return "", fmt.Errorf(utils.UnhandledServerError)
	}

	secret := make([]byte, byteLength)
	_, err := rand.Read(secret)
	if err != nil {
		return "", fmt.Errorf(utils.UnhandledServerError)
	}

	return base64.StdEncoding.EncodeToString(secret), nil
}

// HashAppSecret hashes the secret before storing it
func (g *Ed25519KeyGen) HashAppSecret(secret string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	return string(hashed), err
}

// VerifyAppSecret compares the hashed and raw secret
func (g *Ed25519KeyGen) VerifyAppSecret(hashedSecret, inputSecret string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(inputSecret))
}
