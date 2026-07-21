package initiator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	feedbackConfig "cbe-super-app-cps-action/config"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage/kafka"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// utilityServiceAdapter adapts service.UtilityService to kafka.UtilityRepository.
type utilityServiceAdapter struct {
	svc service.UtilityService
}

func (a *utilityServiceAdapter) HandleKafkaMessage(ctx context.Context, msg imodel.UtilityKafkaMessage) error {
	return a.svc.HandleKafkaMessage(ctx, msg)
}

func InitUtilityConsumer(utilitySvc service.UtilityService, cfg *config.VaultConfig, logger utils.Logger) error {
	kafkaConfig := feedbackConfig.LoadKafkaConfig(cfg, logger)
	if kafkaConfig.Brokers == "" {
		logger.Infof("Kafka brokers not configured, skipping utility consumer")
		return nil
	}
	if kafkaConfig.UtilityActionTopic == "" {
		logger.Infof("KAFKA_UTILITY_ACTION_TOPIC not configured, skipping utility consumer")
		return nil
	}

	adapter := &utilityServiceAdapter{svc: utilitySvc}
	deadLetterQueue := kafka.NewSimpleDeadLetterQueue(logger)

	consumer, err := kafka.NewUtilityConsumer(*kafkaConfig, cfg, logger, adapter, deadLetterQueue)
	if err != nil {
		logger.Errorf("Failed to create utility Kafka consumer: %v", err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Infof("Starting utility consumer for topic: %s", kafkaConfig.UtilityActionTopic)
		if err := consumer.Start(ctx); err != nil {
			logger.Errorf("Utility consumer error: %v", err)
		}
	}()

	go func() {
		sig := <-sigChan
		logger.Infof("Received signal: %v, shutting down utility consumer...", sig)
		cancel()

		if err := consumer.Stop(); err != nil {
			logger.Errorf("Error stopping utility consumer: %v", err)
		}

		logger.Infof("Utility consumer stopped successfully")
	}()

	return nil
}
