package telemetry

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/event"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	mongoTracerName = "cbe-super-app-member-auth/mongodb"
)

// MongoCommandMonitor creates a MongoDB command monitor for tracing
func MongoCommandMonitor() *event.CommandMonitor {
	tracer := otel.Tracer(mongoTracerName)

	return &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			_, span := tracer.Start(ctx, "mongodb."+evt.CommandName,
				trace.WithSpanKind(trace.SpanKindClient),
				trace.WithAttributes(
					semconv.DBSystemMongoDB,
					semconv.DBName(evt.DatabaseName),
					semconv.DBOperation(evt.CommandName),
					attribute.Int64("db.request_id", evt.RequestID),
					attribute.String("db.connection_id", evt.ConnectionID),
				),
			)

			// Store span in context for later use
			if span != nil {
				// we might want to store this span somewhere accessible in Succeeded/Failed callbacks
				// For now, we'll just start it
			}
		},
		Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
			span := trace.SpanFromContext(ctx)
			if span.IsRecording() {
				span.SetAttributes(
					attribute.Int64("db.duration_ms", int64(evt.Duration)/1000000),
					attribute.Int64("db.request_id", evt.RequestID),
				)
				span.SetStatus(codes.Ok, "")
				span.End()
			}
		},
		Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
			span := trace.SpanFromContext(ctx)
			if span.IsRecording() {
				span.SetAttributes(
					attribute.Int64("db.duration_ms", int64(evt.Duration)/1000000),
					attribute.Int64("db.request_id", evt.RequestID),
					attribute.String("db.error", evt.Failure.Error()),
				)
				span.SetStatus(codes.Error, evt.Failure.Error())
				span.RecordError(nil, trace.WithAttributes(
					attribute.String("error.message", evt.Failure.Error()),
				))
				span.End()
			}
		},
	}
}

// MongoPoolMonitor creates a MongoDB connection pool monitor for tracing
func MongoPoolMonitor() *event.PoolMonitor {
	return &event.PoolMonitor{
		Event: func(evt *event.PoolEvent) {
			// we can add pool monitoring metrics here if needed
			// For now, we'll keep it simple
		},
	}
}
