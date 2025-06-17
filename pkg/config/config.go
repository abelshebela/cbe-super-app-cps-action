package config

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"

	// "github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type VaultConfig struct {
	ServerPort         string `mapstructure:"SERVER_PORT"`
	ServerTimeout      string `mapstructure:"SERVER_TIMEOUT"`
	MongoDBURI         string `mapstructure:"MONGODB_URI"`
	MongoDBDatabase    string `mapstructure:"MONGODB_DATABASE"`
	MinioAccessKey     string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioEndPoint      string `mapstructure:"MINIO_ENDPOINT"`
	MinioSecretKey     string `mapstructure:"MINIO_SECRET_KEY"`
	JWTSecretKey       string `mapstructure:"JWT_SECRET_KEY"`
	JwtTTL             string `mapstructure:"JWT_TTL"`
	TempSessionTimeout string `mapstructure:"TEMP_SESSION_TIMEOUT`
	DashSessionExpiry  string `mapstructure:"DASH_SESSION_TIMEOUT"`
	LDAPJWTExpire      string `mapstructure:"LDAP_JWT_EXPIRE"`
	Key                string `mapstructure:"KEY"`
	Iv                 string `mapstructure:"IV"`
	ApproovSecret      string `mapstructure:"APPROVE_SECRET"`
	ServerPublicKey    string `mapstructure:"SERVER_PUBLIC_KEY"`
	AppSessionExpiry   string `mapstructure:"APP_SESSION_EXPIRY"`
}

func Load() (*VaultConfig, error) {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	})

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Warn().Err(err).Msg("Error reading config file")
			return nil, err
		}
		log.Info().Msg("Config file not found; using environment variables")
	}

	VaultAddress := viper.GetString("VAULT_ADDR")
	VaultToken := viper.GetString("VAULT_TOKEN")
	VaultPath := viper.GetString("VAULT_PATH")

	client, err := api.NewClient(&api.Config{Address: VaultAddress})
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create Vault client")
		return nil, err
	}

	client.SetToken(VaultToken)

	secret, err := client.Logical().Read(VaultPath)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to read secret from Vault")
		return nil, err
	}

	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no secret found at path: %s", VaultPath)
	}

	vaultData, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Warn().Msg("Unexpected Vault secret format — expected KV v2 at `data.data`")
		return nil, fmt.Errorf("unexpected secret format at path: %s", VaultPath)
	}

	log.Info().Msg("Vault Secrets Successfully Retrieved!")

	var config VaultConfig
	if err := mapstructure.Decode(vaultData, &config); err != nil {
		log.Warn().
			Err(err).
			Interface("rawVaultData", vaultData).
			Msg("Failed to decode Vault secret into config struct")
		return nil, err
	}

	log.Info().Msg("Vault Secrets Decode Successfully!")
	return &config, nil
}
