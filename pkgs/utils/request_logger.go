package utils

import (
	"context"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

// requestLoggerKey is the context key used to store the request-scoped logger.
var requestLoggerKey = constants.ContextKey("request_logger")

// LoggerFromCtx retrieves the request-scoped logger that was injected by
// the auth middleware. If no scoped logger is found in the context, it
// returns the provided fallback logger.
func LoggerFromCtx(ctx context.Context, fallback utils.Logger) utils.Logger {
	if l, ok := ctx.Value(requestLoggerKey).(utils.Logger); ok {
		return l
	}
	return fallback
}

// CtxWithLogger stores a scoped logger in the context.
func CtxWithLogger(ctx context.Context, l utils.Logger) context.Context {
	return context.WithValue(ctx, requestLoggerKey, l)
}

// NewSessionLogger wraps a utils.Logger and prepends [session_id=<id>] to
// every log line. Stored in context via CtxWithLogger so all handlers that
// call LoggerFromCtx automatically carry the session trace.
func NewSessionLogger(inner utils.Logger, sessionID string) utils.Logger {
	return &sessionLogger{inner: inner, prefix: "[session_id=" + sessionID + "] "}
}

type sessionLogger struct {
	inner  utils.Logger
	prefix string
}

func (s *sessionLogger) Infof(template string, args ...interface{}) {
	s.inner.Infof(s.prefix+template, args...)
}
func (s *sessionLogger) Warnf(template string, args ...interface{}) {
	s.inner.Warnf(s.prefix+template, args...)
}
func (s *sessionLogger) Errorf(template string, args ...interface{}) {
	s.inner.Errorf(s.prefix+template, args...)
}
func (s *sessionLogger) Fatalf(template string, args ...interface{}) {
	s.inner.Fatalf(s.prefix+template, args...)
}
func (s *sessionLogger) Debugf(template string, args ...interface{}) {
	s.inner.Debugf(s.prefix+template, args...)
}
func (s *sessionLogger) Sync() error {
	return s.inner.Sync()
}
