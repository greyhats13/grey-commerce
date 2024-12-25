// Path: grey-user/pkg/logger/zap_logger.go

package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLogger implements the Logger interface using Zap.
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger initializes and returns a new Zap logger.
func NewZapLogger() (Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "json"
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.MessageKey = "message"
	cfg.EncoderConfig.LevelKey = "severity"
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return &zapLogger{logger: l}, nil
}

// Info logs an info-level message.
func (z *zapLogger) Info(msg string, fields ...Field) {
	z.logger.Info(msg, convertFields(fields)...)
}

// Warn logs a warning-level message.
func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, convertFields(fields)...)
}

// Error logs an error-level message.
func (z *zapLogger) Error(msg string, fields ...Field) {
	z.logger.Error(msg, convertFields(fields)...)
}

// Fatal logs a fatal-level message and exits the application.
func (z *zapLogger) Fatal(msg string, err error, fields ...Field) {
	allFields := append(fields, Field{Key: "error", Value: err.Error()})
	z.logger.Fatal(msg, convertFields(allFields)...)
}

// convertFields converts custom Fields to Zap's Fields.
func convertFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}
