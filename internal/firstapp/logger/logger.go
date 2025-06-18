package logger

import (
	"event-collector/internal/firstapp/config"
	"event-collector/pkg/parse"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
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
func NewZapLogger(cfg *config.Config) (Logger, error) {
	level := parseLevel(cfg.GetEnv(cfg.Logger.Level))
	isDevelopment, _ := parse.ToPrimary[bool](cfg.GetEnv(cfg.Logger.Level))
	disableCaller, _ := parse.ToPrimary[bool](cfg.GetEnv(cfg.Logger.DisableCaller))
	disableStacktrace, _ := parse.ToPrimary[bool](cfg.GetEnv(cfg.Logger.DisableStacktrace))
	encoding := cfg.GetEnv(cfg.Logger.Encoding)
	outputPaths := parse.ToStringArray(cfg.GetEnv(cfg.Logger.OutputPaths))
	errorOutputPaths := parse.ToStringArray(cfg.GetEnv(cfg.Logger.ErrorOutputPaths))
	samplingInitial, _ := parse.ToPrimary[int](cfg.GetEnv(cfg.Logger.Sampling.Initial))
	samplingThereafter, _ := parse.ToPrimary[int](cfg.GetEnv(cfg.Logger.Sampling.Thereafter))

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
		Development:       isDevelopment,
		DisableCaller:     disableCaller,
		DisableStacktrace: disableStacktrace,
		Encoding:          encoding,
		EncoderConfig:     getEncoder(isDevelopment),
		OutputPaths:       outputPaths,
		ErrorOutputPaths:  errorOutputPaths,
		Sampling: &zap.SamplingConfig{
			Initial:    samplingInitial,
			Thereafter: samplingThereafter,
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
	switch level {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "dpanic":
		return zap.DPanicLevel
	case "panic":
		return zap.PanicLevel
	case "fatal":
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
