package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	var cfg = &LogConfig{
		ServerName: "test",
		Appid:      1,
		LogLevel:   "info",
		LogFile:    false,
	}
	l, err := InitLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn")
	l.Error("error")
}
