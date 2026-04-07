package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const defaultMaxDLQSize = 1000

// DeadLetterMessage represents a message that failed processing
type DeadLetterMessage struct {
	OriginalTopic   string                  `json:"original_topic"`
	OriginalMessage *sarama.ConsumerMessage `json:"original_message"`
	Error           string                  `json:"error"`
	FailedAt        time.Time               `json:"failed_at"`
	RetryCount      int                     `json:"retry_count"`
}

// SimpleDeadLetterQueue implements DeadLetterQueue interface with in-memory storage
type SimpleDeadLetterQueue struct {
	logger   utils.Logger
	producer sarama.SyncProducer

	mu       sync.RWMutex
	messages []*DeadLetterMessage
	maxSize  int
}

type DLQOption func(*SimpleDeadLetterQueue)

func WithProducer(p sarama.SyncProducer) DLQOption {
	return func(dlq *SimpleDeadLetterQueue) {
		dlq.producer = p
	}
}

func WithMaxDLQSize(size int) DLQOption {
	return func(dlq *SimpleDeadLetterQueue) {
		if size > 0 {
			dlq.maxSize = size
		}
	}
}

func NewSimpleDeadLetterQueue(logger utils.Logger, opts ...DLQOption) *SimpleDeadLetterQueue {
	dlq := &SimpleDeadLetterQueue{
		logger:  logger,
		maxSize: defaultMaxDLQSize,
	}
	for _, opt := range opts {
		opt(dlq)
	}
	return dlq
}

func (dlq *SimpleDeadLetterQueue) SendToDeadLetterQueue(topic string, message *sarama.ConsumerMessage, err error) error {
	deadLetterMsg := &DeadLetterMessage{
		OriginalTopic:   topic,
		OriginalMessage: message,
		Error:           err.Error(),
		FailedAt:        time.Now(),
		RetryCount:      3,
	}

	dlq.mu.Lock()
	if len(dlq.messages) >= dlq.maxSize {
		// Evict oldest message to make room
		dlq.messages = dlq.messages[1:]
		dlq.logger.Warnf("[dead_letter_queue] buffer full (%d), evicted oldest message", dlq.maxSize)
	}
	dlq.messages = append(dlq.messages, deadLetterMsg)
	dlq.mu.Unlock()

	dlq.logger.Errorf("[dead_letter_queue] message stored - Topic: %s, Partition: %d, Offset: %d, Error: %v",
		message.Topic, message.Partition, message.Offset, err)

	if messageBytes, mErr := json.Marshal(deadLetterMsg); mErr == nil {
		dlq.logger.Errorf("[dead_letter_queue] details: %s", string(messageBytes))
	}

	return nil
}

func (dlq *SimpleDeadLetterQueue) GetDeadLetterMessages(_ context.Context) ([]*DeadLetterMessage, error) {
	dlq.mu.RLock()
	defer dlq.mu.RUnlock()

	result := make([]*DeadLetterMessage, len(dlq.messages))
	copy(result, dlq.messages)
	return result, nil
}

func (dlq *SimpleDeadLetterQueue) RetryDeadLetterMessage(ctx context.Context, msg *DeadLetterMessage) error {
	if msg == nil {
		return fmt.Errorf("message is nil")
	}
	if dlq.producer == nil {
		return fmt.Errorf("no producer configured; cannot retry dead letter messages")
	}
	if msg.OriginalMessage == nil {
		return fmt.Errorf("original message is nil, cannot retry")
	}

	producerMsg := &sarama.ProducerMessage{
		Topic: msg.OriginalTopic,
		Value: sarama.ByteEncoder(msg.OriginalMessage.Value),
	}
	if msg.OriginalMessage.Key != nil {
		producerMsg.Key = sarama.ByteEncoder(msg.OriginalMessage.Key)
	}

	_, _, err := dlq.producer.SendMessage(producerMsg)
	if err != nil {
		dlq.logger.Errorf("[dead_letter_queue] retry failed for topic %s: %v", msg.OriginalTopic, err)
		return fmt.Errorf("retry send failed: %w", err)
	}

	// Remove from stored messages
	dlq.mu.Lock()
	for i, m := range dlq.messages {
		if m == msg {
			dlq.messages = append(dlq.messages[:i], dlq.messages[i+1:]...)
			break
		}
	}
	dlq.mu.Unlock()

	dlq.logger.Infof("[dead_letter_queue] successfully retried message from topic: %s", msg.OriginalTopic)
	return nil
}

// Len returns the current number of stored dead letter messages.
func (dlq *SimpleDeadLetterQueue) Len() int {
	dlq.mu.RLock()
	defer dlq.mu.RUnlock()
	return len(dlq.messages)
}
