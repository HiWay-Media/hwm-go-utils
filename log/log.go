package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

// GetLogger returns a console logger writing to stderr at the given level
// (debug, info, warn, error, fatal, panic; default info), with caller info.
// Levels are plain text: ANSI colour codes break log aggregation.
func GetLogger(logLevel string) *zap.SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   "message",
		LevelKey:     "level",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		TimeKey:      "time",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.Lock(os.Stderr), getLevel(logLevel))
	return zap.New(core, zap.AddCaller()).Sugar()
}

func getLevel(level string) zapcore.Level {
	lower := strings.ToLower(level)
	switch lower {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}
