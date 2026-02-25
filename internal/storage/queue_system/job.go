package queue

import (
	"encoding/json"
	"fmt"
)

type Job struct {
	ID                  string          `json:"id"`
	Type                string          `json:"type"`
	Service             string          `json:"service"`
	Method              string          `json:"method"`
	Action              string          `json:"action"`
	Payload             json.RawMessage `json:"payload"`
	Retry               int             `json:"retry"`
	MaxRetry            int             `json:"max_retry"`
	CreatedAt           int64           `json:"created_at"`
	ProcessingStartedAt int64           `json:"processing_started_at,omitempty"`
}

func (j Job) LogPrefix() string {
	return fmt.Sprintf("id=%s service=%s method=%s action=%s", j.ID, j.Service, j.Method, j.Action)
}
