package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log *zap.Logger
)

func setLogLevel(cfg *zap.Config, logLevelString string) {
	var logLevel zapcore.Level
	switch logLevelString {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	default:
		logLevel = zap.InfoLevel
	}
	cfg.Level.SetLevel(logLevel)
}

// InitClient initializes log module
func InitClient(logLevel string, logPath string) {
	log, _ = zap.NewProduction()
	log.Debug("Initializing logging")

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	setLogLevel(&cfg, logLevel)

	_, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	if err == nil {
		cfg.OutputPaths = append(cfg.OutputPaths, logPath)
	} else {
		log.Warn("Unable to open log file, logging to stdout", zap.Error(err))
		cfg.OutputPaths = []string{"stdout"}
	}

	log, err = cfg.Build()
	if err != nil {
		panic(err)
	}

	log.Info("Logging Initialized.", zap.Strings("destination", cfg.OutputPaths))
}

// Warn provides a helper method to log error messages using zap
func Warn(msg string, fields ...zapcore.Field) {
	log.Warn(msg, fields...)
}

// Error provides a helper method to log error messages using zap
func Error(msg string, fields ...zapcore.Field) {
	log.Error(msg, fields...)
}

// Fatal provides a helper method to log fatal messages using zap
func Fatal(msg string, fields ...zapcore.Field) {
	log.Fatal(msg, fields...)
}

// Info provides a helper method to log info messages using zap
func Info(msg string, fields ...zapcore.Field) {
	log.Info(msg, fields...)
}
