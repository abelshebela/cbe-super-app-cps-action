package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/config"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	"github.com/IBM/sarama"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// FeedbackRepository interface for database operations
type FeedbackRepository interface {
	CreateFeedback(ctx context.Context, req entity.FeedbackRequest, userID string) (*entity.Feedback, error)
}

type FeedbackConsumer struct {
	consumerGroup sarama.ConsumerGroup
	logger        utils.Logger
	config        config.KafkaConfig
	feedbackRepo  FeedbackRepository
}

// NewFeedbackConsumer creates a new Kafka consumer for feedback
func NewFeedbackConsumer(cfg config.KafkaConfig, logger utils.Logger, feedbackRepo FeedbackRepository) (*FeedbackConsumer, error) {
	// Configure Sarama consumer group
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
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

// handleFeedbackMessage processes a feedback message from Kafka
func (fc *FeedbackConsumer) handleFeedbackMessage(message *sarama.ConsumerMessage) error {
	fc.logger.Infof("Received message from topic %s, partition %d, offset %d",
		message.Topic, message.Partition, message.Offset)

	// Parse the Kafka message
	var feedbackMsg entity.FeedbackKafkaMessage
	if err := json.Unmarshal(message.Value, &feedbackMsg); err != nil {
		fc.logger.Errorf("Failed to unmarshal feedback message: %v", err)
		return err
	}

	// Validate the message
	if feedbackMsg.UserID == "" {
		return fmt.Errorf("invalid feedback message: user_id is required")
	}

	if len(feedbackMsg.Responses) == 0 {
		return fmt.Errorf("invalid feedback message: responses are required")
	}

	fc.logger.Infof("Processing feedback message - UserID: %s, FeedbackID: %s",
		feedbackMsg.UserID, feedbackMsg.FeedbackID)

	// Create feedback request from Kafka message
	feedbackRequest := entity.FeedbackRequest{
		Responses: feedbackMsg.Responses,
	}

	// Validate the feedback request
	if err := feedbackRequest.Validate(); err != nil {
		return fmt.Errorf("invalid feedback request: %w", err)
	}

	// Save to database using the feedback repository
	ctx := context.Background()
	feedback, err := fc.feedbackRepo.CreateFeedback(ctx, feedbackRequest, feedbackMsg.UserID)
	if err != nil {
		fc.logger.Errorf("Failed to save feedback to database: %v", err)
		return err
	}

	fc.logger.Infof("Successfully saved feedback to database - ID: %s", feedback.ID.Hex())
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

		if err := h.consumer.handleFeedbackMessage(message); err != nil {
			h.consumer.logger.Errorf("Failed to process message: %v", err)
			// Continue processing other messages
			continue
		}

		// Mark message as processed
		session.MarkMessage(message, "")
	}

	return nil
}
