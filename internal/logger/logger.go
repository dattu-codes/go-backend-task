package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger instance used across the application.
var Log *zap.Logger

// InitLogger builds an Uber Zap structured JSON logger.
// It parses the level config and configures clean ISO-8601 timestamps.
func InitLogger(level string) {
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	// Change timestamp encoder to human-readable ISO8601 instead of epoch float
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "timestamp"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// AddCaller includes the filename and line number of the log invocation
	Log = zap.New(core, zap.AddCaller())
}
