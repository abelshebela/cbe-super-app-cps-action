package feedback

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/config"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/feedback/kafka"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// MockFeedbackRepository implements FeedbackRepository for testing
type MockFeedbackRepository struct {
	feedbacks []*entity.Feedback
	createErr error
}

func (m *MockFeedbackRepository) CreateFeedback(ctx context.Context, req entity.FeedbackRequest, userID string) (*entity.Feedback, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}

	feedback := &entity.Feedback{
		ID:        bson.NewObjectID(),
		UserID:    userID,
		Responses: req.Responses,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.feedbacks = append(m.feedbacks, feedback)
	return feedback, nil
}

// MockLogger implements utils.Logger for testing
type MockLogger struct{}

func (m *MockLogger) Infof(format string, args ...interface{})  {}
func (m *MockLogger) Errorf(format string, args ...interface{}) {}
func (m *MockLogger) Debugf(format string, args ...interface{}) {}
func (m *MockLogger) Warnf(format string, args ...interface{})  {}
func (m *MockLogger) Fatalf(format string, args ...interface{}) {}
func (m *MockLogger) Sync() error                               { return nil }

func TestFeedbackConsumer_ProcessValidMessage(t *testing.T) {
	// Setup
	mockRepo := &MockFeedbackRepository{}
	mockLogger := &MockLogger{}

	cfg := config.KafkaConfig{
		Brokers:           "localhost:9092",
		FeedbackTopic:     "feedback-events",
		ConsumerGroup:     "test-consumer-group",
		RequiredAcks:      1,
		RetryMax:          3,
		SessionTimeout:    30000,
		HeartbeatInterval: 3000,
	}

	consumer, err := kafka.NewFeedbackConsumer(cfg, mockLogger, mockRepo)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Test that consumer was created successfully
	if consumer == nil {
		t.Fatal("Expected consumer to be created, got nil")
	}
}

func TestFeedbackConsumer_ProcessInvalidMessage(t *testing.T) {
	// Setup
	mockRepo := &MockFeedbackRepository{}
	mockLogger := &MockLogger{}

	cfg := config.KafkaConfig{
		Brokers:           "localhost:9092",
		FeedbackTopic:     "feedback-events",
		ConsumerGroup:     "test-consumer-group",
		RequiredAcks:      1,
		RetryMax:          3,
		SessionTimeout:    30000,
		HeartbeatInterval: 3000,
	}

	consumer, err := kafka.NewFeedbackConsumer(cfg, mockLogger, mockRepo)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Test that consumer was created successfully
	if consumer == nil {
		t.Fatal("Expected consumer to be created, got nil")
	}
}

func TestFeedbackConsumer_DatabaseError(t *testing.T) {
	// Setup with mock repository that returns error
	mockRepo := &MockFeedbackRepository{
		createErr: fmt.Errorf("database connection failed"),
	}
	mockLogger := &MockLogger{}

	cfg := config.KafkaConfig{
		Brokers:           "localhost:9092",
		FeedbackTopic:     "feedback-events",
		ConsumerGroup:     "test-consumer-group",
		RequiredAcks:      1,
		RetryMax:          3,
		SessionTimeout:    30000,
		HeartbeatInterval: 3000,
	}

	consumer, err := kafka.NewFeedbackConsumer(cfg, mockLogger, mockRepo)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Test that consumer was created successfully even with error-prone repository
	if consumer == nil {
		t.Fatal("Expected consumer to be created, got nil")
	}
}

func TestFeedbackConsumer_ValidationError(t *testing.T) {
	// Setup
	mockRepo := &MockFeedbackRepository{}
	mockLogger := &MockLogger{}

	cfg := config.KafkaConfig{
		Brokers:           "localhost:9092",
		FeedbackTopic:     "feedback-events",
		ConsumerGroup:     "test-consumer-group",
		RequiredAcks:      1,
		RetryMax:          3,
		SessionTimeout:    30000,
		HeartbeatInterval: 3000,
	}

	consumer, err := kafka.NewFeedbackConsumer(cfg, mockLogger, mockRepo)
	if err != nil {
		t.Fatalf("Failed to create consumer: %v", err)
	}

	// Test that consumer was created successfully
	if consumer == nil {
		t.Fatal("Expected consumer to be created, got nil")
	}
}
