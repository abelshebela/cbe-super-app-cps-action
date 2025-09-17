package initiator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	feedbackRepo "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/persistence/feedback"
	feedbackConfig "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/config"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/kafka"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func InitFeedbackConsumer(mongoClient *mongo.Client, cfg *config.VaultConfig, logger utils.Logger) error {
	kafkaConfig := feedbackConfig.LoadKafkaConfig(cfg)
	if kafkaConfig.Brokers == "" {
		logger.Infof("Kafka brokers not configured, skipping feedback consumer")
		return nil
	}

	feedbackRepository := feedbackRepo.InitFeedback(mongoClient, cfg.MongoDBDatabase, "feedback", logger)

	// Create dead letter queue
	deadLetterQueue := kafka.NewSimpleDeadLetterQueue(logger)

	consumer, err := kafka.NewFeedbackConsumer(*kafkaConfig, logger, feedbackRepository, deadLetterQueue)
	if err != nil {
		logger.Errorf("Failed to create Kafka consumer: %v", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting feedback consumer for topic: %s", kafkaConfig.FeedbackTopic)
		if err := consumer.Start(ctx); err != nil {
			logger.Errorf("Consumer error: %v", err)
		}
	}()

	go func() {
		sig := <-sigChan
		logger.Infof("Received signal: %v", sig)
		logger.Infof("Shutting down feedback consumer...")
		cancel()

		if err := consumer.Stop(); err != nil {
			logger.Errorf("Error stopping consumer: %v", err)
		}

		logger.Infof("Feedback consumer stopped successfully")
	}()

	return nil
}
