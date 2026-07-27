package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cbe-super-app-cps-action/config"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/IBM/sarama"
	cfg "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// UtilityRepository is the narrow interface the consumer calls into the service layer.
type UtilityRepository interface {
	HandleKafkaMessage(ctx context.Context, msg imodel.UtilityKafkaMessage) error
}

// UtilityConsumer consumes from the utility action Kafka topic and delegates
// to UtilityRepository for maker-checker CPS action creation.
type UtilityConsumer struct {
	consumerGroup sarama.ConsumerGroup
	logger        utils.Logger
	kafkaConfig   config.KafkaConfig
	vaultCfg      *cfg.VaultConfig
	utilityRepo   UtilityRepository
	deadLetterQ   DeadLetterQueue
	maxRetries    int
}

func NewUtilityConsumer(
	kafkaConfig config.KafkaConfig,
	vaultConfig *cfg.VaultConfig,
	logger utils.Logger,
	utilityRepo UtilityRepository,
	deadLetterQ DeadLetterQueue,
) (*UtilityConsumer, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	saramaCfg.Consumer.Group.Session.Timeout = time.Duration(kafkaConfig.SessionTimeout) * time.Millisecond
	saramaCfg.Consumer.Group.Heartbeat.Interval = time.Duration(kafkaConfig.HeartbeatInterval) * time.Millisecond
	saramaCfg.Version = sarama.V2_6_0_0

	brokers := strings.Split(kafkaConfig.Brokers, ",")
	for i, b := range brokers {
		brokers[i] = strings.TrimSpace(b)
	}

	group, err := sarama.NewConsumerGroup(brokers, kafkaConfig.ConsumerGroup, saramaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create utility Kafka consumer group: %w", err)
	}

	return &UtilityConsumer{
		consumerGroup: group,
		logger:        logger,
		kafkaConfig:   kafkaConfig,
		vaultCfg:      vaultConfig,
		utilityRepo:   utilityRepo,
		deadLetterQ:   deadLetterQ,
		maxRetries:    3,
	}, nil
}

func (uc *UtilityConsumer) Start(ctx context.Context) error {
	log := local_util.LoggerFromCtx(ctx, uc.logger)
	log.Infof("Starting utility Kafka consumer for topic: %s", uc.kafkaConfig.UtilityActionTopic)

	topics := []string{uc.kafkaConfig.UtilityActionTopic}
	handler := &utilityConsumerGroupHandler{consumer: uc}

	for {
		if err := uc.consumerGroup.Consume(ctx, topics, handler); err != nil {
			log.Errorf("[UtilityConsumer] error from consumer: %v", err)
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

func (uc *UtilityConsumer) Stop() error {
	uc.logger.Infof("Stopping utility Kafka consumer")
	return uc.consumerGroup.Close()
}

// utilityConsumerGroupHandler implements sarama.ConsumerGroupHandler.
type utilityConsumerGroupHandler struct {
	consumer *UtilityConsumer
}

func (h *utilityConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *utilityConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *utilityConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.consumer.handleUtilityMessage(session.Context(), message); err != nil {
			h.consumer.logger.Errorf("[UtilityConsumer] failed to process message offset=%d: %v", message.Offset, err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (uc *UtilityConsumer) handleUtilityMessage(ctx context.Context, message *sarama.ConsumerMessage) error {
	log := local_util.LoggerFromCtx(ctx, uc.logger)
	log.Infof("[UtilityConsumer] received message topic=%s partition=%d offset=%d",
		message.Topic, message.Partition, message.Offset)

	log.Infof("[UtilityConsumer] received message=%+v", message)

	var kafkaMsg imodel.UtilityKafkaMessage
	if err := json.Unmarshal(message.Value, &kafkaMsg); err != nil {
		log.Errorf("[UtilityConsumer] failed to unmarshal message: %v", err)
		return fmt.Errorf("invalid utility message format: %w", err)
	}

	log.Infof("[UtilityConsumer] received kafka body=%+v", kafkaMsg)

	if kafkaMsg.Action == "" {
		log.Errorf("[UtilityConsumer] message missing action field")
		return fmt.Errorf("action field is required")
	}

	var lastErr error
	for retry := 0; retry < uc.maxRetries; retry++ {
		if err := uc.utilityRepo.HandleKafkaMessage(ctx, kafkaMsg); err != nil {
			lastErr = err
			log.Errorf("[UtilityConsumer] HandleKafkaMessage failed (attempt %d/%d): %v", retry+1, uc.maxRetries, err)
			if retry < uc.maxRetries-1 {
				time.Sleep(time.Duration(retry+1) * time.Second)
			}
			continue
		}
		log.Infof("[UtilityConsumer] successfully processed utility message action=%s", kafkaMsg.Action)
		return nil
	}

	if uc.deadLetterQ != nil {
		if dlqErr := uc.deadLetterQ.SendToDeadLetterQueue(message.Topic, message, lastErr); dlqErr != nil {
			log.Errorf("[UtilityConsumer] failed to send to DLQ: %v", dlqErr)
		}
	}

	return fmt.Errorf("failed to process utility message after %d retries: %w", uc.maxRetries, lastErr)
}
