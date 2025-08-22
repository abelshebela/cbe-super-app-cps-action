package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// DeadLetterMessage represents a message that failed processing
type DeadLetterMessage struct {
	OriginalTopic   string                  `json:"original_topic"`
	OriginalMessage *sarama.ConsumerMessage `json:"original_message"`
	Error           string                  `json:"error"`
	FailedAt        time.Time               `json:"failed_at"`
	RetryCount      int                     `json:"retry_count"`
}

// SimpleDeadLetterQueue implements DeadLetterQueue interface
type SimpleDeadLetterQueue struct {
	logger utils.Logger
}

func NewSimpleDeadLetterQueue(logger utils.Logger) *SimpleDeadLetterQueue {
	return &SimpleDeadLetterQueue{
		logger: logger,
	}
}

func (dlq *SimpleDeadLetterQueue) SendToDeadLetterQueue(topic string, message *sarama.ConsumerMessage, err error) error {
	deadLetterMsg := &DeadLetterMessage{
		OriginalTopic:   topic,
		OriginalMessage: message,
		Error:           err.Error(),
		FailedAt:        time.Now(),
		RetryCount:      3, // Assuming max retries reached
	}

	// Log the failed message
	dlq.logger.Errorf("Message sent to dead letter queue - Topic: %s, Partition: %d, Offset: %d, Error: %v",
		message.Topic, message.Partition, message.Offset, err)

	if messageBytes, err := json.Marshal(deadLetterMsg); err == nil {
		dlq.logger.Errorf("Dead letter message details: %s", string(messageBytes))
	}

	return nil
}

func (dlq *SimpleDeadLetterQueue) GetDeadLetterMessages(ctx context.Context) ([]*DeadLetterMessage, error) {
	return []*DeadLetterMessage{}, nil
}

func (dlq *SimpleDeadLetterQueue) RetryDeadLetterMessage(ctx context.Context, msg *DeadLetterMessage) error {

	dlq.logger.Infof("Retrying dead letter message from topic: %s", msg.OriginalTopic)
	return fmt.Errorf("retry functionality not implemented")
}