package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/vault/api"
)

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers           string `json:"brokers"`
	FeedbackTopic     string `json:"feedback_topic"`
	ConsumerGroup     string `json:"consumer_group"`
	RequiredAcks      int    `json:"required_acks"`
	RetryMax          int    `json:"retry_max"`
	SessionTimeout    int    `json:"session_timeout"`
	HeartbeatInterval int    `json:"heartbeat_interval"`
}

// VaultClient wraps the vault client with caching capabilities
type VaultClient struct {
	client     *api.Client
	path       string
	secretData map[string]interface{}
}

// NewVaultClient creates a new vault client
func NewVaultClient() (*VaultClient, error) {
	config := &api.Config{
		Address: getEnv("VAULT_ADDR", ""),
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(getEnv("VAULT_TOKEN", ""))

	valut := &VaultClient{
		client: client,
		path:   getEnv("VAULT_PATH", ""),
	}

	if err := valut.fetchSecrets(); err != nil {
		return nil, err
	}

	return valut, nil
}

// fetchSecrets retrieves secrets from Vault with caching
func (v *VaultClient) fetchSecrets() error {
	secret, err := v.client.Logical().Read(v.path)
	if err != nil {
		return err
	}

	if secret == nil || secret.Data == nil {
		return fmt.Errorf("no secrets found at path: %s", v.path)
	}

	if data, ok := secret.Data["data"].(map[string]interface{}); ok {
		v.secretData = data
		return nil
	}

	return nil
}

// GetSecret retrieves a secret from Vault
func (v *VaultClient) GetSecret(key string) (string, error) {
	if value, ok := v.secretData[key].(string); ok {
		return value, nil
	}
	return "", nil
}

// LoadKafkaConfig loads Kafka configuration from Vault with fallback
func LoadKafkaConfig() *KafkaConfig {
	vaultClient, err := NewVaultClient()
	if err != nil {
		// Fallback to environment variables if Vault is not available
		return loadKafkaConfigFromEnv()
	}

	// Helper function to get config values with Vault fallback
	getConfigValue := func(vaultKey, defaultValue string) string {
		// Try Vault first
		vaultValue, err := vaultClient.GetSecret(vaultKey)
		if err == nil && vaultValue != "" {
			return vaultValue
		}
		// Fallback to environment variable
		return getEnv(vaultKey, defaultValue)
	}

	// Helper for integer values
	getConfigInt := func(vaultKey string, defaultValue int) int {
		// Try Vault first
		if vaultValue, err := vaultClient.GetSecret(vaultKey); err == nil && vaultValue != "" {
			if intValue, err := strconv.Atoi(vaultValue); err == nil {
				return intValue
			}
		}
		// Fallback to environment variable
		if envValue := getEnv(vaultKey, ""); envValue != "" {
			if intValue, err := strconv.Atoi(envValue); err == nil {
				return intValue
			}
		}
		return defaultValue
	}

	return &KafkaConfig{
		Brokers:           getConfigValue("KAFKA_BROKERS", ""),
		FeedbackTopic:     getConfigValue("KAFKA_FEEDBACK_TOPIC", "feedback-events"),
		ConsumerGroup:     getConfigValue("KAFKA_CONSUMER_GROUP", "cps_action_service"),
		RequiredAcks:      getConfigInt("KAFKA_REQUIRED_ACKS", 1),
		RetryMax:          getConfigInt("KAFKA_RETRY_MAX", 3),
		SessionTimeout:    getConfigInt("KAFKA_SESSION_TIMEOUT", 30000),
		HeartbeatInterval: getConfigInt("KAFKA_HEARTBEAT_INTERVAL", 3000),
	}
}

// loadKafkaConfigFromEnv loads Kafka configuration from environment variables (fallback)
func loadKafkaConfigFromEnv() *KafkaConfig {
	requiredAcks, _ := strconv.Atoi(getEnv("KAFKA_REQUIRED_ACKS", "1"))
	retryMax, _ := strconv.Atoi(getEnv("KAFKA_RETRY_MAX", "3"))
	sessionTimeout, _ := strconv.Atoi(getEnv("KAFKA_SESSION_TIMEOUT", "30000"))
	heartbeatInterval, _ := strconv.Atoi(getEnv("KAFKA_HEARTBEAT_INTERVAL", "3000"))

	return &KafkaConfig{
		Brokers:           getEnv("KAFKA_BROKERS", ""),
		FeedbackTopic:     getEnv("KAFKA_FEEDBACK_TOPIC", "feedback-events"),
		ConsumerGroup:     getEnv("KAFKA_CONSUMER_GROUP", "cps-action-service"),
		RequiredAcks:      requiredAcks,
		RetryMax:          retryMax,
		SessionTimeout:    sessionTimeout,
		HeartbeatInterval: heartbeatInterval,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
 