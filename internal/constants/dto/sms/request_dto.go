package smsDto

type SMSRequest struct {
	Recipient   string `json:"recipient"`
	MessageBody string `json:"message_body"`
}
