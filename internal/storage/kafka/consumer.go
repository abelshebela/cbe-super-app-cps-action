package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/config"
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	"cbe-super-app-cps-action/internal/constants/model"

	"github.com/IBM/sarama"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// FeedbackRepository interface for database operations
type FeedbackRepository interface {
	CreateFeedback(ctx context.Context, req feedback.FeedbackRequest, userID string) (*model.Feedback, error)
}

// DeadLetterQueue interface for handling failed messages
type DeadLetterQueue interface {
	SendToDeadLetterQueue(topic string, message *sarama.ConsumerMessage, error error) error
}

type FeedbackConsumer struct {
	consumerGroup sarama.ConsumerGroup
	logger        utils.Logger
	config        config.KafkaConfig
	feedbackRepo  FeedbackRepository
	deadLetterQ   DeadLetterQueue
	maxRetries    int
}

// NewFeedbackConsumer creates a new Kafka consumer for feedback
func NewFeedbackConsumer(cfg config.KafkaConfig, logger utils.Logger, feedbackRepo FeedbackRepository, deadLetterQ DeadLetterQueue) (*FeedbackConsumer, error) {
	// Configure Sarama consumer group
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // Changed to OffsetNewest for production
	config.Consumer.Group.Session.Timeout = time.Duration(cfg.SessionTimeout) * time.Millisecond
	config.Consumer.Group.Heartbeat.Interval = time.Duration(cfg.HeartbeatInterval) * time.Millisecond
	config.Version = sarama.V2_6_0_0

	// Parse brokers
	brokers := strings.Split(cfg.Brokers, ",")
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}

	consumerGroup, err := sarama.NewConsumerGroup(brokers, cfg.ConsumerGroup, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer group: %w", err)
	}

	return &FeedbackConsumer{
		consumerGroup: consumerGroup,
		logger:        logger,
		config:        cfg,
		feedbackRepo:  feedbackRepo,
		deadLetterQ:   deadLetterQ,
		maxRetries:    3, // Configurable retry count
	}, nil
}

// Start starts consuming messages from Kafka
func (fc *FeedbackConsumer) Start(ctx context.Context) error {
	fc.logger.Infof("Starting Kafka consumer for topic: %s", fc.config.FeedbackTopic)

	topics := []string{fc.config.FeedbackTopic}
	handler := &ConsumerGroupHandler{
		consumer: fc,
	}

	for {
		err := fc.consumerGroup.Consume(ctx, topics, handler)
		if err != nil {
			fc.logger.Errorf("Error from consumer: %v", err)
			return err
		}

		// Check if context was cancelled
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

// Stop stops the Kafka consumer
func (fc *FeedbackConsumer) Stop() error {
	fc.logger.Infof("Stopping Kafka consumer")
	return fc.consumerGroup.Close()
}

// handleFeedbackMessage processes a feedback message from Kafka with retry logic
func (fc *FeedbackConsumer) handleFeedbackMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	fc.logger.Infof("Received message from topic %s, partition %d, offset %d",
		message.Topic, message.Partition, message.Offset)

	// Parse the Kafka message wrapper first
	var kafkaMsg model.KafkaMessage
	if err := json.Unmarshal(message.Value, &kafkaMsg); err != nil {
		fc.logger.Errorf("Failed to unmarshal Kafka message wrapper: %v", err)
		return fmt.Errorf("invalid message wrapper format: %w", err)
	}

	// Validate the message type
	if kafkaMsg.Type != "feedback" {
		fc.logger.Errorf("Unexpected message type: %s, expected: feedback", kafkaMsg.Type)
		return fmt.Errorf("unexpected message type: %s", kafkaMsg.Type)
	}

	// Parse the payload into FeedbackKafkaMessage
	var feedbackMsg model.FeedbackKafkaMessage
	if err := json.Unmarshal(kafkaMsg.Payload, &feedbackMsg); err != nil {
		fc.logger.Errorf("Failed to unmarshal feedback payload: %v", err)
		return fmt.Errorf("invalid feedback payload format: %w", err)
	}

	// Validate the message
	if err := fc.validateFeedbackMessage(&feedbackMsg); err != nil {
		fc.logger.Errorf("Message validation failed: %v", err)
		return fmt.Errorf("message validation failed: %w", err)
	}

	fc.logger.Infof("Processing feedback message - UserID: %s, FeedbackID: %s",
		feedbackMsg.UserID, feedbackMsg.FeedbackID)

	// Create feedback request from Kafka message
	feedbackRequest := feedback.FeedbackRequest{
		Responses: feedbackMsg.Responses,
	}

	// Validate the feedback request
	if err := feedbackRequest.Validate(); err != nil {
		return fmt.Errorf("invalid feedback request: %w", err)
	}

	// Save to database using the feedback repository with retry logic
	var lastErr error
	for retry := 0; retry < fc.maxRetries; retry++ {
		feedback, err := fc.feedbackRepo.CreateFeedback(ctx, feedbackRequest, feedbackMsg.UserID)
		if err != nil {
			lastErr = err
			fc.logger.Errorf("Failed to save feedback to database (attempt %d/%d): %v", retry+1, fc.maxRetries, err)

			// If this is the last retry, break and handle as permanent failure
			if retry == fc.maxRetries-1 {
				break
			}

			// Wait before retry with exponential backoff
			backoff := time.Duration(retry+1) * time.Second
			time.Sleep(backoff)
			continue
		}

		fc.logger.Infof("Successfully saved feedback to database - ID: %s", feedback.ID.Hex())
		return nil
	}

	// If all retries failed, send to dead letter queue
	if fc.deadLetterQ != nil {
		if err := fc.deadLetterQ.SendToDeadLetterQueue(message.Topic, message, lastErr); err != nil {
			fc.logger.Errorf("Failed to send message to dead letter queue: %v", err)
		}
	}

	return fmt.Errorf("failed to process message after %d retries: %w", fc.maxRetries, lastErr)
}

// validateFeedbackMessage validates the feedback message structure
func (fc *FeedbackConsumer) validateFeedbackMessage(msg *model.FeedbackKafkaMessage) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}

	if msg.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if len(msg.Responses) == 0 {
		return fmt.Errorf("responses are required")
	}

	// Validate each response
	for key, response := range msg.Responses {
		if response.Question == "" {
			return fmt.Errorf("question is required for response key: %s", key)
		}
		if response.Answer == nil {
			return fmt.Errorf("answer is required for response key: %s", key)
		}
	}

	return nil
}

// ConsumerGroupHandler implements sarama.ConsumerGroupHandler
type ConsumerGroupHandler struct {
	consumer *FeedbackConsumer
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	h.consumer.logger.Infof("Consumer group session setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.consumer.logger.Infof("Consumer group session cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages()
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		h.consumer.logger.Infof("Processing message: topic=%s partition=%d offset=%d",
			message.Topic, message.Partition, message.Offset)

		// Use session context instead of Background
		if err := h.consumer.handleFeedbackMessage(session.Context(), message); err != nil {
			h.consumer.logger.Errorf("Failed to process message: %v", err)
			// Continue processing other messages but don't mark as processed
			continue
		}

		// Mark message as processed only if successful
		session.MarkMessage(message, "")
	}

	return nil
}
