package utils

import (
	//"cbe-super-app-member-auth/pkg/config"
	//"cbe-super-app-member-auth/pkg/config"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

var (
	encryptionSecretPath = os.Getenv("VAULT_ENCRYPTION_PATH")
	vaultAddr            = os.Getenv("VAULT_ADDR")
	vaultToken           = os.Getenv("VAULT_TOKEN")
	vaultClient          *vault.Client
)

func init() {
	config := vault.DefaultConfig()
	config.Address = vaultAddr
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("failed to create vault client: %v", err)
	}
	client.SetToken(vaultToken)
	vaultClient = client
}

func writeToVault(path string, data map[string]interface{}) error {
	_, err := vaultClient.KVv2("secret").Put(context.Background(), path, data)
	return err
}

func readFromVault(path string) (map[string]interface{}, error) {
	secret, err := vaultClient.KVv2("secret").Get(context.Background(), path)
	if err != nil {
		return nil, err
	}
	return secret.Data, nil
}

func GenerateAndSaveKeyPair() error {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	privBytes := priv.D.Bytes()
	pubBytes := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	var data map[string]interface{}
	data = map[string]interface{}{
		"SERVER_PRIVATE_KEY": hex.EncodeToString(privBytes),
		"SERVER_PUBLIC_KEY":  hex.EncodeToString(pubBytes),
	}
	if err := writeToVault(encryptionSecretPath, data); err != nil {
		return err
	}

	return nil
}

func GetPublicKey() (string, error) {
	// Read public key from a local file instead of Vault
	_, err := config.Load()
	publicKey := "" //env.ServerPublicKey

	if err != nil {
		return "", err
	}
	return publicKey, nil
}

func GetPrivateKey() (string, error) {
	keys, err := readFromVault(encryptionSecretPath)
	if err != nil {
		return "", err
	}
	if v, ok := keys["SERVER_PRIVATE_KEY"].(string); ok {
		return v, nil
	}
	return "", nil
}
