package logger

import (
	"event-collector/internal/config"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the interface that wraps the basic logging methods.
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}

// zapLogger is a wrapper around zap.Logger that implements Logger interface.
type zapLogger struct {
	logger *zap.Logger
}

// NewZapLogger creates a new zap logger with the given configuration.
func NewZapLogger(cfg *configs.Config) (Logger, error) {
	level := parseLevel(cfg.Logger.Level)
	outputPaths := cfg.Logger.OutputPaths
	errorOutputPaths := cfg.Logger.ErrorOutputPaths

	// Ensure log files exist before Zap tries to use them
	for _, path := range outputPaths {
		if err := ensureLogFile(path); err != nil {
			return nil, err
		}
	}
	for _, path := range errorOutputPaths {
		if err := ensureLogFile(path); err != nil {
			return nil, err
		}
	}

	zapCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Logger.Development,
		DisableCaller:     cfg.Logger.DisableCaller,
		DisableStacktrace: cfg.Logger.DisableStacktrace,
		Encoding:          cfg.Logger.Encoding,
		EncoderConfig:     getEncoder(cfg.Logger.Development),
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
		Sampling: &zap.SamplingConfig{
			Initial:    cfg.Logger.Sampling.Initial,
			Thereafter: cfg.Logger.Sampling.Thereafter,
		},
	}

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	log, err := zapCfg.Build()
	if err != nil {
		return nil, err
	}

	return &zapLogger{logger: log}, nil
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

func (l *zapLogger) With(fields ...zap.Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}

func parseLevel(level string) zapcore.Level {
	switch {
	case strings.EqualFold(level, "debug"):
		return zap.DebugLevel
	case strings.EqualFold(level, "info"):
		return zap.InfoLevel
	case strings.EqualFold(level, "warn"):
		return zap.WarnLevel
	case strings.EqualFold(level, "error"):
		return zap.ErrorLevel
	case strings.EqualFold(level, "dpanic"):
		return zap.DPanicLevel
	case strings.EqualFold(level, "panic"):
		return zap.PanicLevel
	case strings.EqualFold(level, "fatal"):
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func getEncoder(isDevelopment bool) zapcore.EncoderConfig {
	if isDevelopment {
		return zap.NewDevelopmentEncoderConfig()
	}
	return zap.NewProductionEncoderConfig()
}

func ensureLogFile(path string) error {
	if path == "stdout" || path == "stderr" {
		return nil // Skip for stdout/stderr
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists, if not create it
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		file.Close()
	}

	return nil
}
