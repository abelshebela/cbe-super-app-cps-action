package initiator

import (
	"os"
	"strings"

	configshared "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
)

// AppConfig is the typed application configuration wrapper.
// It contains the shared VaultConfig plus local typed telemetry fields.
type AppConfig struct {
	Shared                 *configshared.VaultConfig
	OTLPEndpoint           string
	OTelEnabled            bool
	OTelDebugTracePresence bool
	GoEnv                  string
}

// NewAppConfig builds an AppConfig from the shared VaultConfig and environment.
// Priority: if the shared config provides the fields, prefer them (if added later).
// Otherwise fall back to environment variables and developer defaults.
func NewAppConfig(sharedCfg *configshared.VaultConfig) *AppConfig {
	ac := &AppConfig{
		Shared: sharedCfg,
	}

	// Defaults
	ac.GoEnv = "development"
	ac.OTLPEndpoint = "localhost:4317"
	ac.OTelEnabled = true
	ac.OTelDebugTracePresence = false

	// Try to read from environment (developer/local) – these will be used until shared config is extended.
	if v := os.Getenv("ENVIRONMENT"); v != "" {
		ac.GoEnv = v
	}
	if v := os.Getenv("OTLP_ENDPOINT"); v != "" {
		ac.OTLPEndpoint = v
	}
	if v := os.Getenv("OTEL_ENABLED"); v != "" {
		lv := strings.ToLower(v)
		ac.OTelEnabled = (lv == "1" || lv == "true" || lv == "yes")
	}
	if v := os.Getenv("OTEL_DEBUG_TRACE_PRESENCE"); v != "" {
		lv := strings.ToLower(v)
		ac.OTelDebugTracePresence = (lv == "1" || lv == "true" || lv == "yes")
	}

	// If the shared VaultConfig gains these fields in the future, prefer them.
	// We're using struct field checks with reflection avoided; instead, attempt to access common-named getters if available.
	// Since we can't rely on implementation, prefer explicit fields via compile-time in future updates.

	return ac
}
