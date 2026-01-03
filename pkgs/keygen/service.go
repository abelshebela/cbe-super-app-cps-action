package keygen

import (
	"cbe-super-app-cps-action/internal/constants/localization"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type Ed25519KeyGen struct {
	logger shared_utils.Logger
	cfg    *config.VaultConfig
}

type KeyPair struct {
	PublicKey  string
	PrivateKey string
}

type KeyGeneratorService interface {
	GenerateKeyPair() (KeyPair, error)
	Sign(data []byte, privateKey string) ([]byte, error)
	Verify(data, signature []byte, publicKey string) (bool, error)
}

var secretKey = []byte("01234567890123456789012345678901")

func NewKeyGenerator(
	logger shared_utils.Logger,
	cfg *config.VaultConfig,
) KeyGeneratorService {
	logger.Infof("Initializing Ed25519KeyGen")
	return &Ed25519KeyGen{
		logger: logger,
		cfg:    cfg,
	}
}

func (g *Ed25519KeyGen) GenerateKeyPair() (KeyPair, error) {
	g.logger.Infof("Generating Ed25519 key pair")
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		g.logger.Errorf("Failed to generate Ed25519 key pair", "error", err)
		return KeyPair{}, errors.New(localization.ErrorUnhandledServer.Code)
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
		return nil, errors.New(localization.ErrorUnhandledServer.Code)
	}
	if len(privateKeyBytes) != ed25519.PrivateKeySize {
		g.logger.Errorf("Invalid private key size", "expected", ed25519.PrivateKeySize, "got", len(privateKeyBytes))
		return nil, errors.New(localization.ErrorUnhandledServer.Code)
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
		return false, errors.New(localization.ErrorUnhandledServer.Code)
	}
	if len(publicKeyBytes) != ed25519.PublicKeySize {
		g.logger.Errorf("Invalid public key size", "expected", ed25519.PublicKeySize, "got", len(publicKeyBytes))
		return false, errors.New(localization.ErrorUnhandledServer.Code)
	}
	isValid := ed25519.Verify(publicKeyBytes, data, signature)
	if isValid {
		g.logger.Infof("Signature verification successful")
	} else {
		g.logger.Infof("Signature verification failed")
	}
	return isValid, nil
}
