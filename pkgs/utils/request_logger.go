package utils

import (
	"context"

	"cbe-super-app-cps-action/internal/constants"

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
// Called by middleware after creating a logger with WithStr("user_id", ...).
func CtxWithLogger(ctx context.Context, l utils.Logger) context.Context {
	return context.WithValue(ctx, requestLoggerKey, l)
}
