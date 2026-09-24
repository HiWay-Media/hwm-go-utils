package log

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func TestGetLevel(t *testing.T) {
	cases := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel, "INFO": zapcore.InfoLevel, "Warn": zapcore.WarnLevel,
		"error": zapcore.ErrorLevel, "fatal": zapcore.FatalLevel, "panic": zapcore.PanicLevel, "": zapcore.InfoLevel,
	}
	for in, want := range cases {
		if got := getLevel(in); got != want {
			t.Errorf("getLevel(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestGetLoggerLevel(t *testing.T) {
	l := GetLogger("warn").Desugar()
	if l.Core().Enabled(zapcore.InfoLevel) || !l.Core().Enabled(zapcore.WarnLevel) {
		t.Error("logger level not applied")
	}
}
