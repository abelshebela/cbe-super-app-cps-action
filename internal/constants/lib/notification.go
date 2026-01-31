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

// PublishInAppBroadcast publishes an in-app broadcast payload to Kafka using
// the shared notification wrapper with type "in_app_broadcast". The payload
// Type should be "inapp".
func (n *NotificationStore) PublishInAppBroadcast(ctx context.Context, inapp types.InAppBroadcastMessage) error {
	// Validate minimal required fields
	if inapp.Title == "" || inapp.Message == "" {
		n.logger.Errorf("in-app broadcast: title and message are required")
		return errors.New("title and message are required")
	}
	if inapp.Type == "" {
		inapp.Type = "inapp"
	}

	// Note: Using KafkaEmailTopic for generic notification fanout unless
	// a dedicated in-app topic is provided in configuration.
	return n.producer.PublishMessage(ctx, inapp, "in_app_broadcast", "inapp_notifications", "InAppBroadcast")
}

// PublishInAppBroadcastWithLog publishes an in-app broadcast with a custom logType
// (e.g., notification.CREATE_PUBLIC_NOTIFICATION.maker|checker) so downstream
// systems can segment maker vs checker events per module/action.
func (n *NotificationStore) PublishInAppBroadcastWithLog(ctx context.Context, inapp types.InAppBroadcastMessage, logType string) error {
	if inapp.Title == "" || inapp.Message == "" {
		n.logger.Errorf("in-app broadcast: title and message are required")
		return errors.New("title and message are required")
	}
	if inapp.Type == "" {
		inapp.Type = "inapp"
	}
	return n.producer.PublishMessage(ctx, inapp, "in_app_broadcast", "inapp_notifications", logType)
}
