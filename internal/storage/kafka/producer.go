package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.uber.org/zap"
)

// NotificationProducer handles publishing notification messages to Kafka
type NotificationProducer struct {
	producer sarama.SyncProducer
	logger   utils.Logger
	config   *config.VaultConfig
}

// NewNotificationProducer creates a new Kafka producer for notifications
func NewNotificationProducer(cfg *config.VaultConfig, producer sarama.SyncProducer, logger utils.Logger) *NotificationProducer {

	return &NotificationProducer{
		producer: producer,
		logger:   logger,
		config:   cfg,
	}
}

// Close closes the Kafka producer
func (np *NotificationProducer) Close(ctx context.Context) {
	if err := np.producer.Close(); err != nil {
		np.logger.Errorf("Failed to close Kafka producer", zap.Error(err))
	} else {
		np.logger.Infof("Kafka producer closed successfully")
	}
}

// publishMessage is a helper function to eliminate code duplication
func (np *NotificationProducer) PublishMessage(ctx context.Context, msg interface{}, msgType, topic, logType string) error {
	// Create notification message
	notificationMsg, err := dto.NewNotificationMessage(
		fmt.Sprintf("%s-%d", msgType, time.Now().UnixNano()),
		msgType,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to create notification message: %w", err)
	}

	// Marshal message to JSON
	messageBytes, err := json.Marshal(notificationMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Sarama producer message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
		Headers: []sarama.RecordHeader{
			{Key: []byte("message_id"), Value: []byte(notificationMsg.ID)},
			{Key: []byte("type"), Value: []byte(msgType)},
			{Key: []byte("timestamp"), Value: []byte(notificationMsg.CreatedAt.Format(time.RFC3339))},
		},
	}

	return np.produceAndWait(ctx, kafkaMsg, notificationMsg.ID, topic, logType)
}

// produceAndWait handles message production and delivery confirmation
func (np *NotificationProducer) produceAndWait(ctx context.Context,
	kafkaMsg *sarama.ProducerMessage,
	messageID, topic, logType string) error {
	// Create a channel for handling the send with timeout
	done := make(chan error, 1)

	go func() {
		partition, offset, err := np.producer.SendMessage(kafkaMsg)
		if err != nil {
			done <- fmt.Errorf("failed to produce message: %w", err)

			return
		}

		np.logger.Infof(
			fmt.Sprintf("%s message published successfully", logType),
			zap.String("message_id", messageID),
			zap.String("topic", topic),
			zap.Int32("partition", partition),
			zap.Int64("offset", offset),
		)
		done <- nil
	}()

	// Wait for either completion, context cancellation, or timeout
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(30 * time.Second):
		return fmt.Errorf("timeout waiting for message delivery")
	}
}
