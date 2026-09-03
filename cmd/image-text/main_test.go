package main

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		arg   string
		level slog.Level
		err   bool
	}{
		{
			name:  "empty defaults to info",
			arg:   "",
			level: slog.LevelInfo,
		},
		{
			name:  "debug",
			arg:   "debug",
			level: slog.LevelDebug,
		},
		{
			name:  "DEBUG uppercase",
			arg:   "DEBUG",
			level: slog.LevelDebug,
		},
		{
			name:  "info",
			arg:   "info",
			level: slog.LevelInfo,
		},
		{
			name:  "warn",
			arg:   "warn",
			level: slog.LevelWarn,
		},
		{
			name:  "error",
			arg:   "error",
			level: slog.LevelError,
		},
		{
			name:  "invalid falls back to info with error",
			arg:   "verbose",
			level: slog.LevelInfo,
			err:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := parseLogLevel(tt.arg)
			assert.Equal(t, tt.err, err != nil)
			assert.Equal(t, tt.level, level)
		})
	}
}
