package initiator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	feedbackConfig "cbe-super-app-cps-action/config"
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage/kafka"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// feedbackServiceAdapter adapts service.FeedbackService to kafka.FeedbackRepository interface
type feedbackServiceAdapter struct {
	svc service.FeedbackService
}

func (a *feedbackServiceAdapter) CreateFeedback(ctx context.Context, req feedback.FeedbackRequest, userID string) (*model.Feedback, error) {
	return a.svc.CreateFeedback(ctx, req, userID)
}

func InitFeedbackConsumer(feedbackSvc service.FeedbackService, cfg *config.VaultConfig, logger utils.Logger) error {
	kafkaConfig := feedbackConfig.LoadKafkaConfig(cfg, logger)
	if kafkaConfig.Brokers == "" {
		logger.Infof("Kafka brokers not configured, skipping feedback consumer")
		return nil
	}

	adapter := &feedbackServiceAdapter{svc: feedbackSvc}

	// Create dead letter queue
	deadLetterQueue := kafka.NewSimpleDeadLetterQueue(logger)

	consumer, err := kafka.NewFeedbackConsumer(*kafkaConfig, logger, adapter, deadLetterQueue)
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
