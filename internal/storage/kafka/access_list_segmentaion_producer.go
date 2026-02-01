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

// AccessListSegmentationMessage represents the base structure for Access List Segmentation Kafka messages
type AccessListSegmentationMessage struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"eventType"` // "member.sync.cps", "linked_account.sync.cps", "access_control.sync.cps", "kyc.sync.cps"
	Data      json.RawMessage        `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	Headers   map[string]interface{} `json:"headers,omitempty"`
}

// NewAccessListSegmentationMessage creates a new AccessListSegmentationMessage with marshaled payload
func NewAccessListSegmentationMessage(id, eventType string, data interface{}) (*AccessListSegmentationMessage, error) {
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &AccessListSegmentationMessage{
		ID:        id,
		EventType: eventType,
		Data:      payloadBytes,
		CreatedAt: time.Now(),
		Headers:   make(map[string]interface{}),
	}, nil
}

// AccessListSegmentationProducer handles publishing Access List Segmentation messages to Kafka
type AccessListSegmentationProducer struct {
	producer sarama.SyncProducer
	logger   utils.Logger
	config   *config.VaultConfig
}

// NewAccessListSegmentationProducer creates a new Kafka producer for Access List Segmentation
func NewAccessListSegmentationProducer(cfg *config.VaultConfig, producer sarama.SyncProducer, logger utils.Logger) *AccessListSegmentationProducer {
	return &AccessListSegmentationProducer{
		producer: producer,
		logger:   logger,
		config:   cfg,
	}
}

// Close closes the Access List Segmentation Kafka producer
func (asp *AccessListSegmentationProducer) Close(ctx context.Context) {
	if err := asp.producer.Close(); err != nil {
		asp.logger.Errorf("Failed to close Access List Segmentation Kafka producer", zap.Error(err))
	} else {
		asp.logger.Infof("Access List Segmentation Kafka producer closed successfully")
	}
}

// PublishMessage publishes an Access List Segmentation message to Kafka
func (asp *AccessListSegmentationProducer) PublishMessage(ctx context.Context, msg interface{}, msgType, topic, logType string) error {
	// Create Access List Segmentation message
	accessListSegmentationMsg, err := NewAccessListSegmentationMessage(
		fmt.Sprintf("%s-%d", msgType, time.Now().UnixNano()),
		msgType,
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to create Access List Segmentation message: %w", err)
	}

	// Marshal message to JSON
	messageBytes, err := json.Marshal(accessListSegmentationMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal Access List Segmentation message: %w", err)
	}

	if topic == "" {
		// topic = asp.config.AccessListSegmentationTopic
	}

	// Create Sarama producer message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(messageBytes),
		Headers: []sarama.RecordHeader{
			{Key: []byte("message_id"), Value: []byte(accessListSegmentationMsg.ID)},
			{Key: []byte("type"), Value: []byte(msgType)},
			{Key: []byte("timestamp"), Value: []byte(accessListSegmentationMsg.CreatedAt.Format(time.RFC3339))},
		},
	}

	return asp.produceAndWait(ctx, kafkaMsg, accessListSegmentationMsg.ID, topic, logType)
}

// produceAndWait handles Access List Segmentation message production and delivery confirmation
func (asp *AccessListSegmentationProducer) produceAndWait(ctx context.Context,
	kafkaMsg *sarama.ProducerMessage,
	messageID, topic, logType string) error {
	// Create a channel for handling the send with timeout
	done := make(chan error, 1)

	go func() {
		partition, offset, err := asp.producer.SendMessage(kafkaMsg)
		if err != nil {
			done <- fmt.Errorf("failed to produce Access List Segmentation message: %w", err)
			return
		}

		asp.logger.Infof(
			fmt.Sprintf("%s Access List Segmentation message published successfully", logType),
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
		return fmt.Errorf("timeout waiting for Access List Segmentation message delivery")
	}
}
