package logger

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Init() {
	cfg := zap.Config{
		Encoding:          "json",
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		DisableCaller:     false,
		DisableStacktrace: true,
		OutputPaths:       []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,    // Capitalize the log level names
			EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC timestamp format
			EncodeDuration: zapcore.SecondsDurationEncoder, // Duration in seconds
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger = zap.Must(cfg.Build())
}

func Sync() {
	logger.Sync()
}

func Info(message string, fields ...zap.Field) {
	logger.Info(message, fields...)
}

func Debug(message string, fields ...zap.Field) {
	logger.Debug(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	logger.Warn(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	logger.Error(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	logger.Fatal(message, fields...)
}

// WithError helper function that logs error and returns it
func WithError(err error, message string) error {
	logger.Error(message, zap.Error(err))
	return err
}

type LogFunc func(message string, fields ...zap.Field)

// WithTraceID helper function that calls LoggerFunc with traceID and spanID
func WithTraceID(ctx context.Context, f LogFunc, message string, fields ...zap.Field) {
	fields = FieldsWithTraceID(ctx, fields...)
	f(message, fields...)
}

// FieldsWithTraceID function that returns fields enriched with traceID and spanID
func FieldsWithTraceID(ctx context.Context, fields ...zap.Field) []zap.Field {
	spanCtx := trace.SpanContextFromContext(ctx)
	fields = append(fields, zap.String("traceID", spanCtx.TraceID().String()))
	fields = append(fields, zap.String("spanID", spanCtx.SpanID().String()))
	return fields
}
