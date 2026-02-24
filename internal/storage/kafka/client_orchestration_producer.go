package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.uber.org/zap"
)

// OrchestrationMessage represents the base structure for Kafka messages
type ClientOrchestrationMessage struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"eventType"` // "member.sync.cps", "linked_account.sync.cps", "access_control.sync.cps", "kyc.sync.cps"
	Data      json.RawMessage        `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
}

// NewOrchestrationMessage creates a new OrchestrationMessage with marshaled payload
func NewOrchestrationMessage(id, eventType string, data interface{}) (*ClientOrchestrationMessage, error) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &ClientOrchestrationMessage{
		ID:        id,
		EventType: eventType,
		Data:      payloadBytes,
		CreatedAt: time.Now(),
		Headers:   make(map[string]interface{}),
	}, nil
}

// ClientOrchestration handles publishing notification messages to Kafka
type ClientOrchestrationProducer struct {
	producer sarama.SyncProducer
	logger   utils.Logger
	config   *config.VaultConfig
}

// NewClientOrchestration creates a new Kafka producer for notifications
func NewClientOrchestrationProducer(cfg *config.VaultConfig, producer sarama.SyncProducer, logger utils.Logger) *ClientOrchestrationProducer {

	return &ClientOrchestrationProducer{
		producer: producer,
		logger:   logger,
		config:   cfg,
	}
}

// Close closes the Kafka producer
func (np *ClientOrchestrationProducer) Close(ctx context.Context) {
	if err := np.producer.Close(); err != nil {
		np.logger.Errorf("Failed to close Kafka producer", zap.Error(err))
	} else {
		np.logger.Infof("Kafka producer closed successfully")
	}
}

// publishMessage is a helper function to eliminate code duplication
func (np *ClientOrchestrationProducer) PublishMessage(ctx context.Context, msg interface{}, msgType, topic, logType string) error {
	// Create notification message
	clientOrchestrationMsg, err := NewOrchestrationMessage(
		fmt.Sprintf("%s-%d", msgType, time.Now().UnixNano()),
		msgType,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to create client orchestration message: %w", err)
	}

	// Marshal message to JSON
	messageBytes, err := json.Marshal(clientOrchestrationMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if topic == "" {
		topic = np.config.UserOrchastrator
	}

	// Create Sarama producer message
	kafkaMsg := &sarama.ProducerMessage{
		// Topic: np.config.UserOrchastrator,
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
		Headers: []sarama.RecordHeader{
			{Key: []byte("message_id"), Value: []byte(clientOrchestrationMsg.ID)},
			{Key: []byte("type"), Value: []byte(msgType)},
			{Key: []byte("timestamp"), Value: []byte(clientOrchestrationMsg.CreatedAt.Format(time.RFC3339))},
		},
	}

	return np.produceAndWait(ctx, kafkaMsg, clientOrchestrationMsg.ID, np.config.UserOrchastrator, logType)
}

// produceAndWait handles message production and delivery confirmation
func (np *ClientOrchestrationProducer) produceAndWait(ctx context.Context,
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
