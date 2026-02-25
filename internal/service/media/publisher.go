package media

import (
	"encoding/json"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type KafkaProducerService interface {
	PublishShortVideoEvent(path, shortVideoID string) error
}

type KafkaProducer struct {
	producer *kafka.Producer
	logger   utils.Logger
	config   *config.VaultConfig
}

func CreateKafkaProducer(logger utils.Logger, cfg *config.VaultConfig) KafkaProducerService {
	confluentProducer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": cfg.KafkaBrokers,
			"security.protocol": "PLAINTEXT",
		})

	if err != nil {
		logger.Fatalf("Failed to create Kafka producer: %v", err)
	}

	go func() {
		for e := range confluentProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					logger.Errorf("[MediaPubSvc][Events] delivery err: %v", ev.TopicPartition.Error)
				}
			}
		}
	}()

	return &KafkaProducer{
		producer: confluentProducer,
		logger:   logger,
		config:   cfg,
	}
}

func (r *KafkaProducer) PublishShortVideoEvent(path, shortVideoID string) error {
	topic := "short-video-events-create"
	shortVideoTopic := "short-video-events-update"

	payload := map[string]string{
		"short_video_id": shortVideoID,
		"path":           path,
		"topic_name":     shortVideoTopic,
		"upload_path":    "news-reel-videos",
	}

	return r.produceMessage(topic, shortVideoID, payload)
}

// produceMessage is a generic helper to marshal and send messages using the local confluent-kafka-go producer
func (r *KafkaProducer) produceMessage(topic string, key string, payload interface{}) error {
	if topic == "" {
		return fmt.Errorf("topic is empty, check configuration")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		r.logger.Errorf("[MediaPubSvc][Produce] marshal err topic %s: %v", topic, err)
		return err
	}

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: int32(kafka.PartitionAny),
		},
		Value: data,
	}

	if key != "" {
		msg.Key = []byte(key)
	}

	if err := r.producer.Produce(msg, nil); err != nil {
		r.logger.Errorf("[MediaPubSvc][Produce] send err topic %s: %v", topic, err)
		return err
	}

	r.logger.Infof("[MediaPubSvc][Produce] sent topic: %s", topic)
	return nil
}
