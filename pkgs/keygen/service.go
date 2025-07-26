package keygen

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Ed25519KeyGen struct {
	logger shared_utils.Logger
}

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

func NewKeyGenerator(
	logger shared_utils.Logger,
) KeyGeneratorService {
	logger.Infof("Initializing Ed25519KeyGen")
	return &Ed25519KeyGen{
		logger: logger,
	}
}

func (g *Ed25519KeyGen) GenerateKeyPair() (KeyPair, error) {
	g.logger.Infof("Generating Ed25519 key pair")
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		g.logger.Errorf("Failed to generate Ed25519 key pair", "error", err)
		return KeyPair{}, fmt.Errorf(utils.UnhandledServerError)
	}
	keyPair := KeyPair{
		PublicKey:  base64.StdEncoding.EncodeToString(publicKey),
		PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
	}
	g.logger.Infof("Successfully generated Ed25519 key pair", "publicKeyLength", len(publicKey), "privateKeyLength", len(privateKey))
	return keyPair, nil
}

func (g *Ed25519KeyGen) Sign(data []byte, privateKey string) ([]byte, error) {
	g.logger.Infof("Signing data", "dataLength", len(data))
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		g.logger.Errorf("Failed to decode private key", "error", err)
		return nil, fmt.Errorf(utils.UnhandledServerError)
	}
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		g.logger.Errorf("Invalid private key size", "expected", ed25519.PrivateKeySize, "got", len(privateKeyBytes))
		return nil, fmt.Errorf(utils.UnhandledServerError)
	}
	signature := ed25519.Sign(privateKeyBytes, data)
	g.logger.Infof("Successfully signed data", "signatureLength", len(signature))
	return signature, nil
}

func (g *Ed25519KeyGen) Verify(data, signature []byte, publicKey string) (bool, error) {
	g.logger.Infof("Verifying signature", "dataLength", len(data), "signatureLength", len(signature))
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		g.logger.Errorf("Failed to decode public key", "error", err)
		return false, fmt.Errorf(utils.UnhandledServerError)
	}
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		g.logger.Errorf("Invalid public key size", "expected", ed25519.PublicKeySize, "got", len(publicKeyBytes))
		return false, fmt.Errorf(utils.UnhandledServerError)
	}
	isValid := ed25519.Verify(publicKeyBytes, data, signature)
	if isValid {
		g.logger.Infof("Signature verification successful")
	} else {
		g.logger.Infof("Signature verification failed")
	}
	return isValid, nil
}

func (g *Ed25519KeyGen) GenerateFabricID() string {
	g.logger.Infof("Generating Fabric ID")
	fabricID := uuid.New().String()
	g.logger.Infof("Successfully generated Fabric ID", "fabricID", fabricID)
	return fabricID
}

func (g *Ed25519KeyGen) GenerateNumericCode(length int) (string, error) {
	g.logger.Infof("Generating numeric code", "length", length)
	if length <= 0 {
		g.logger.Errorf("Invalid code length", "length", length)
		return "", fmt.Errorf(utils.UnhandledServerError)
	}

	digits := "0123456789"
	max := big.NewInt(int64(len(digits)))
	code := make([]byte, length)

	for i := range code {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			g.logger.Errorf("Failed to generate random digit", "error", err)
			return "", fmt.Errorf(utils.UnhandledServerError)
		}
		code[i] = digits[n.Int64()]
	}

	g.logger.Infof("Successfully generated numeric code", "code", string(code))
	return string(code), nil
}

func (g *Ed25519KeyGen) GenerateAppSecret(byteLength int) (string, error) {
	g.logger.Infof("Generating app secret", "byteLength", byteLength)
	if byteLength < 32 {
		g.logger.Errorf("Invalid byte length for app secret", "byteLength", byteLength)
		return "", fmt.Errorf(utils.UnhandledServerError)
	}

	secret := make([]byte, byteLength)
	_, err := rand.Read(secret)
	if err != nil {
		g.logger.Errorf("Failed to generate app secret", "error", err)
		return "", fmt.Errorf(utils.UnhandledServerError)
	}

	encodedSecret := base64.StdEncoding.EncodeToString(secret)
	g.logger.Infof("Successfully generated app secret", "secretLength", len(encodedSecret))
	return encodedSecret, nil
}

func (g *Ed25519KeyGen) HashAppSecret(secret string) (string, error) {
	g.logger.Infof("Hashing app secret")
	hashed, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		g.logger.Errorf("Failed to hash app secret", "error", err)
		return "", err
	}
	g.logger.Infof("Successfully hashed app secret", "hashedLength", len(hashed))
	return string(hashed), nil
}

func (g *Ed25519KeyGen) VerifyAppSecret(hashedSecret, inputSecret string) error {
	g.logger.Infof("Verifying app secret")
	err := bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(inputSecret))
	if err != nil {
		g.logger.Errorf("Failed to verify app secret", "error", err)
		return err
	}
	g.logger.Infof("App secret verification successful")
	return nil
}