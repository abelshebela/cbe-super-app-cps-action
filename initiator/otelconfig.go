package initiator

import (
	"os"
	"strconv"
	"strings"

	"cbe-super-app-cps-action/platform/telemetry"
)

// NewOtelConfig constructs a telemetry.Config from environment variables with sane local defaults.
// This keeps OTEL setup isolated from the main application config and is intentionally minimal and explicit.
func NewOtelConfig() telemetry.Config {
	// defaults
	serviceName := "cbe-super-app-cps-action"
	serviceVersion := "1.0.0"
	env := "development"
	endpoint := "localhost:4317"
	enabled := true

	if v := os.Getenv("SERVICE_NAME"); v != "" {
		serviceName = v
	}
	if v := os.Getenv("SERVICE_VERSION"); v != "" {
		serviceVersion = v
	}
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		env = v
	}
	if v := os.Getenv("OTLP_ENDPOINT"); v != "" {
		endpoint = v
	}
	// allow OTLP endpoint to be provided as OTLP_ENDPOINT or OTLP
	if v := os.Getenv("OTLP"); v != "" && endpoint == "localhost:4317" {
		endpoint = v
	}
	endpoint = "192.168.136.1:4317"

	if v := os.Getenv("OTEL_ENABLED"); v != "" {
		lv := strings.ToLower(v)
		enabled = (lv == "1" || lv == "true" || lv == "yes")
	}

	samplingRatio := 1.0
	if v := os.Getenv("OTEL_SAMPLING_RATIO"); v != "" {
		if s, err := strconv.ParseFloat(v, 64); err == nil {
			samplingRatio = s
		}
	}

	return telemetry.Config{
		ServiceName:    serviceName,
		ServiceVersion: serviceVersion,
		Environment:    env,
		OTLPEndpoint:   endpoint,
		Enabled:        enabled,
		SamplingRatio:  samplingRatio,
	}
}
