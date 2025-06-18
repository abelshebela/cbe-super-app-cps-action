package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func AxiosSendSms(phoneNumber, messageBody string) {
	ip := os.Getenv("CONFIG_IP")
	port := os.Getenv("CONFIG_PORT_LDAP_NOTIFICATION")
	nodeEnv := os.Getenv("GO_ENV")

	url := fmt.Sprintf("http://%s:%s/v1.0/chatbirrapi/ldapnotif/sms/send", ip, port)

	type requestBody struct {
		Recipient   string `json:"recipient"`
		MessageBody string `json:"messageBody"`
	}

	prefix := ""
	if nodeEnv == "uat" || nodeEnv == "dev" {
		prefix = "UAT: "
	}

	body := requestBody{
		Recipient:   FormatPhoneNumber(phoneNumber),
		MessageBody: prefix + messageBody,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling SMS body: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Error creating SMS request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending SMS: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("SMS sent successfully")
	} else {
		log.Printf("Failed to send SMS, status: %s", resp.Status)
	}
}
