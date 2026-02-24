package config

import (
	"os"
	"strconv"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// KafkaConfig holds Kafka configuration
type KafkaConfig struct {
	Brokers             string `json:"brokers"`
	FeedbackTopic       string `json:"feedback_topic"`
	SurveyFeedbackTopic string `json:"survey_feedback_topic"`
	ConsumerGroup       string `json:"consumer_group"`
	RequiredAcks        int    `json:"required_acks"`
	RetryMax            int    `json:"retry_max"`
	SessionTimeout      int    `json:"session_timeout"`
	HeartbeatInterval   int    `json:"heartbeat_interval"`
}

func LoadKafkaConfig(cfg *config.VaultConfig, logger utils.Logger) *KafkaConfig {
	logger.Infof("Loading Kafka configuration from VaultConfig")

	return loadKafkaConfigFromEnv(cfg, logger)

}

func loadKafkaConfigFromEnv(cfg *config.VaultConfig, logger utils.Logger) *KafkaConfig {
	requiredAcks, _ := strconv.Atoi(getEnv("KAFKA_REQUIRED_ACKS", "1", logger))
	retryMax, _ := strconv.Atoi(getEnv("KAFKA_RETRY_MAX", "3", logger))
	sessionTimeout, _ := strconv.Atoi(getEnv("KAFKA_SESSION_TIMEOUT", "30000", logger))
	heartbeatInterval, _ := strconv.Atoi(getEnv("KAFKA_HEARTBEAT_INTERVAL", "3000", logger))

	return &KafkaConfig{
		Brokers:             cfg.KafkaBrokers,
		FeedbackTopic:       cfg.KafkaCustomerFeedbackTopic,
		SurveyFeedbackTopic: cfg.KafkaCustomerSurveyTopic,
		ConsumerGroup:       "cps-action-service",
		RequiredAcks:        requiredAcks,
		RetryMax:            retryMax,
		SessionTimeout:      sessionTimeout,
		HeartbeatInterval:   heartbeatInterval,
	}
}

func getEnv(key, defaultValue string, logger utils.Logger) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	logger.Errorf("the %s not found in the vault", key)
	return defaultValue
}
