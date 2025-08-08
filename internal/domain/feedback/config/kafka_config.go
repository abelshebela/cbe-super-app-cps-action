package config

import (
	"fmt"
	"os"
	"strconv"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
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

func LoadKafkaConfig(cfg *config.VaultConfig) *KafkaConfig {
	fmt.Printf("DEBUG: Starting LoadKafkaConfig\n")
	fmt.Printf("DEBUG: Using VaultConfig from initiator\n")

	return loadKafkaConfigFromEnv(cfg )

}

func loadKafkaConfigFromEnv(cfg *config.VaultConfig) *KafkaConfig {
	requiredAcks, _ := strconv.Atoi(getEnv("KAFKA_REQUIRED_ACKS", "1"))
	retryMax, _ := strconv.Atoi(getEnv("KAFKA_RETRY_MAX", "3"))
	sessionTimeout, _ := strconv.Atoi(getEnv("KAFKA_SESSION_TIMEOUT", "30000"))
	heartbeatInterval, _ := strconv.Atoi(getEnv("KAFKA_HEARTBEAT_INTERVAL", "3000"))

	return &KafkaConfig{
        Brokers:           cfg.KafkaBrokers,
        FeedbackTopic:     "feedback-events",
        ConsumerGroup:     "cps-action-service",
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
