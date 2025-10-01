package model

import "time"

type APIResponse struct {
	Ok        bool        `json:"ok"`
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}