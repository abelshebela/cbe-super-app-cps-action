package lib

import (
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage/kafka"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type NotificationStore struct {
	logger   utils.Logger
	producer *kafka.NotificationProducer
	cfg      *config.VaultConfig
}

func InitNotificationStore(logger utils.Logger, cfg *config.VaultConfig, producer *kafka.NotificationProducer) *NotificationStore {
	return &NotificationStore{
		logger:   logger,
		cfg:      cfg,
		producer: producer,
	}
}

func (n *NotificationStore) PublishMessage(ctx context.Context, message interface{}) error {

	switch message.(type) {
	case types.SMSKafkaMessage:
		return n.PublishSMSMessage(ctx, message.(types.SMSKafkaMessage))
	case types.EmailKafkaMessage:
		return n.PublishEmailMessage(ctx, message.(types.EmailKafkaMessage))
	default:
		return errors.New("invalid message type")
	}
}

func (n *NotificationStore) PublishSMSMessage(ctx context.Context, smsMsg types.SMSKafkaMessage) error {

	// Validate required fields
	if smsMsg.Recipient == "" {
		n.logger.Errorf("Recipient is required")
		return errors.New("recipient is required")
	}

	if smsMsg.MessageBody == "" {
		n.logger.Errorf("Message body is required")
		return errors.New("message body is required")
	}

	// Publish message to Kafka
	if err := n.producer.PublishMessage(ctx, smsMsg, "sms", n.cfg.KafkaSMSTopic, "SMS"); err != nil {
		return err
	}
	return nil
}

func (n *NotificationStore) PublishEmailMessage(ctx context.Context, emailMsg types.EmailKafkaMessage) error {

	// Validate required fields
	if len(emailMsg.Recipients) == 0 {
		n.logger.Errorf("Recipients are required")
		return errors.New("recipients are required")
	}

	if emailMsg.Subject == "" {
		n.logger.Errorf("Subject is required")
		return errors.New("subject is required")
	}

	if emailMsg.Type == "" {
		n.logger.Errorf("Type is required")

		return errors.New("type is required")
	}

	// Publish message to Kafka
	if err := n.producer.PublishMessage(ctx, emailMsg, "email", n.cfg.KafkaEmailTopic, "Email"); err != nil {
		return errors.New("failed to publish email message to Kafka")
	}
	return nil

}
