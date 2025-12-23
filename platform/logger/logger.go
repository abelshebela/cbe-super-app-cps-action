package logger

import (
	"context"
	"time"

	"cbe-super-app-cps-action/internal/constants"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Logger interface {
	GetZapLooger() *zap.Logger
	Named(s string) *logger
	With(fields ...zap.Field) *logger
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Panic(ctx context.Context, msg string, fields ...zap.Field)
	Fatal(ctx context.Context, msg string, fields ...zap.Field)
	extact(ctx context.Context) []zap.Field
}

type logger struct {
	logger  *zap.Logger
	options Options
}

// GetZapLooger implements Logger.
func (l *logger) GetZapLooger() *zap.Logger {
	return l.logger
}

// Named implements Logger.
func (l *logger) Named(s string) *logger {
	l2 := l.logger.Named(s)
	return &logger{
		logger:  l2,
		options: l.options,
	}
}

func (l *logger) With(fields ...zap.Field) *logger {
	l2 := l.logger.With(fields...)
	return &logger{
		logger:  l2,
		options: l.options,
	}
}

func (l *logger) extact(ctx context.Context) []zap.Field {
	var fields []zap.Field
	fields = append(fields, zap.String("time", time.Now().Format(time.RFC3339)))

	if ctx != nil {
		for _, field := range l.options.ExtractFields {
			if v := ctx.Value(field.KeyInContext); v != nil {
				fields = append(fields, field.Func(v))
			}
		}
	}

	return fields
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extact(ctx)...).Debug(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extact(ctx)...).Error(msg, fields...)
}

func (l *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extact(ctx)...).Fatal(msg, fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extact(ctx)...).Info(msg, fields...)
}

func (l *logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extact(ctx)...).Panic(msg, fields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(l.extact(ctx)...).Warn(msg, fields...)
}

type ExtractField struct {
	KeyInContext any
	Func         func(any) zap.Field
}

type Options struct {
	ExtractFields                 []ExtractField
	IngnoreDefacutlEEtactedFields bool
}

func New(l *zap.Logger, options Options) Logger {
	if !options.IngnoreDefacutlEEtactedFields {
		options.ExtractFields = append(options.ExtractFields, []ExtractField{
			{
				KeyInContext: middleware.RequestIDKey,
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-request-id", vString)
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: constants.ContextKey("trace_id"),
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("trace_id", vString)
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: constants.ContextKey("span_id"),
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("span_id", vString)
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: constants.ContextKey("x-user-id"),
				Func: func(v any) zap.Field {
					if vString, ok := v.(string); ok {
						return zap.String("x-user-id", vString)
					}

					return zap.Skip()
				},
			},
			{
				KeyInContext: constants.ContextKey("request-start-time"),
				Func: func(v any) zap.Field {
					if vTime, ok := v.(time.Time); ok {
						return zap.Float64("time-since-request", float64(time.Since(vTime).Milliseconds()))
					}

					return zap.Skip()
				},
			},
		}...)
	}

	return &logger{
		logger:  l,
		options: options,
	}
}
