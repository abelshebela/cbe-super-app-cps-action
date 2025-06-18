package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	kafkaWriter *kafka.Writer
	once        sync.Once
)

type ISMS struct {
	Recipient   string `json:"recipient"`
	MessageBody string `json:"messageBody"`
}

type TransactionDetails struct {
	SenderName      string    `json:"senderName"`
	ReceiverName    string    `json:"receiverName"`
	ReceiverAccount string    `json:"receiverAccount"`
	TransactionType string    `json:"transactionType,omitempty"`
	Date            time.Time `json:"date"`
	Amount          string    `json:"amount"`
	SenderAccount   string    `json:"senderAccount,omitempty"`
}

type IEmail struct {
	Recipients         []string            `json:"recipients,omitempty"`
	Type               string              `json:"type"`
	OtpCode            string              `json:"otpcode,omitempty"`
	CustomerName       string              `json:"customerName,omitempty"`
	MessageBody        string              `json:"messageBody,omitempty"`
	TransactionDetails *TransactionDetails `json:"transactionDetails,omitempty"`
	Link               string              `json:"link,omitempty"`
	Subject            string              `json:"subject,omitempty"`
	Receiver           string              `json:"receiver"`
}

type IInApp struct {
	UserID            bson.ObjectID  `json:"userID"`
	NotificationType  string         `json:"notificationType"`
	NotificationBody  string         `json:"notificationBody"`
	NotificationParts map[string]any `json:"notificationParts,omitempty"`
}

type KafkaMessageType string

const (
	TopicSendSMS   KafkaMessageType = "sendSMS"
	TopicSendEmail KafkaMessageType = "sendEmail"
	TopicSendInApp KafkaMessageType = "sendInApp"
)

func getKafkaWriter() (*kafka.Writer, error) {
	var err error
	once.Do(func() {
		kafkaAddress := os.Getenv("KAFKA_ADDRESS")
		if kafkaAddress == "" {
			err = errors.New("Kafka address is undefined")
			return
		}
		kafkaWriter = &kafka.Writer{
			Addr:         kafka.TCP(kafkaAddress),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Async:        false,
		}
	})
	return kafkaWriter, err
}

func KafkaProducer(ctx context.Context, topic KafkaMessageType, message any) error {
	writer, err := getKafkaWriter()
	if err != nil {
		return err
	}

	value, err := json.Marshal(message)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Topic: string(topic),
		Value: value,
		Time:  time.Now(),
	}

	if err := writer.WriteMessages(ctx, kafkaMsg); err != nil {
		return err
	}

	fmt.Printf("Message sent to topic %q: %v\n", topic, message)
	return nil
}
