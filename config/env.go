package config

import (
	"fmt"

	"cbe-super-app-member-auth/platform/logger"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/vault/api"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type VaultConfig struct {
	MinioAccessKey        string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioEndPoint         string `mapstructure:"MINIO_ENDPOINT"`
	MinioSecretKey        string `mapstructure:"MINIO_SECRET_KEY"`
	MongoDBConnectTimeout string `mapstructure:"MONGODB_CONNECT_TIMEOUT"`
	MongoDBDatabase       string `mapstructure:"MONGODB_DATABASE"`
	MongoDBQueryTimeout   string `mapstructure:"MONGODB_QUERY_TIMEOUT"`
	MongoDBURI            string `mapstructure:"MONGODB_URI"`
	ReadTimeout           uint8  `mapstructure:"READ_TIMEOUT"`
	WriteTimeout          uint8  `mapstructure:"WRITE_TIMEOUT"`
	IdleTimeout           uint8  `mapstructure:"IDLE_TIMEOUT"`
	MenuServicePort       string `mapstructure:"cbe-super-app-budget_PORT"`
	ServerPublicKey       string `mapstructure:"SERVER_PUBLIC_KEY"`
	JWTSecretKey          string `mapstructure:"JWT_SECRET_KEY"`
	ServerTimeout         int    `mapstructure:"SERVER_TIMEOUT"`
	RedisConfig           RedisConfig
	Key                   string `mapstructure:"KEY"`
	IV                    string `mapstructure:"IV"`
}

func LoadVault(logger logger.Logger) (*VaultConfig, error) {

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

	var (
		VaultAddress = viper.GetString("VAULT_ADDR")
		VaultToken   = viper.GetString("VAULT_TOKEN")
		VaultPath    = viper.GetString("VAULT_PATH")
	)

	client, err := api.NewClient(&api.Config{
		Address: VaultAddress,
	})

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

	vault, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Warn().Msg("Unexpected Vault secret format")
		return nil, fmt.Errorf("unexpected secret format at path: %s", VaultPath)
	}

	log.Info().Msg("Vault Secrets Successfully Retrieved!")

	var config VaultConfig
	if err := mapstructure.Decode(vault, &config); err != nil {
		log.Warn().Msg("Failed to Decode Vault Secret!")
		return nil, err
	}

	log.Info().Msg("Vault Secrets Decode Successfully!")
	return &config, nil
}
