package initiator

import (
	"cbe-super-app-cps-action/internal/storage/kafka"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func InitKafkaService(cfg *config.VaultConfig, logger utils.Logger) (*kafka.NotificationProducer, *kafka.ClientOrchestrationProducer) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Version = sarama.V2_6_0_0

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	for i, broker := range brokers {
		brokers[i] = strings.TrimSpace(broker)
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Errorf("Failed to create Kafka producer: %v", err)
		return nil, nil
	}

	return kafka.NewNotificationProducer(cfg, producer, logger), kafka.NewClientOrchestrationProducer(cfg, producer, logger)
}
