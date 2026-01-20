package logging

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestNewLogger_ValidLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug", "debug"},
		{"info", "info"},
		{"warn", "warn"},
		{"error", "error"},
		{"Debug", "Debug"},
		{"Info", "Info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{Level: tt.level}
			logger, err := NewLogger(cfg)
			if err != nil {
				t.Fatalf("NewLogger() error = %v", err)
			}
			if logger == nil {
				t.Fatal("NewLogger() returned nil logger")
			}
			logger.Sync() // Flush any buffered logs
		})
	}
}

func TestNewLogger_InvalidLevel(t *testing.T) {
	cfg := Config{Level: "invalid-level"}
	logger, err := NewLogger(cfg)

	if err == nil {
		t.Error("NewLogger() should return error for invalid level")
		logger.Sync()
	}
	if logger != nil {
		t.Error("NewLogger() should return nil logger for invalid level")
		logger.Sync()
	}
}

func TestNewLogger_WithFile(t *testing.T) {
	// Create a temporary log file
	tmpFile, err := os.CreateTemp("", "test-log-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cfg := Config{Level: "info", File: tmpFile.Name()}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Fatal("NewLogger() returned nil logger")
	}

	// Log something and verify it doesn't error
	logger.Info("test log message")
	logger.Sync()

	// Verify file exists and has content
	info, err := os.Stat(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat log file: %v", err)
	}
	if info.Size() == 0 {
		t.Error("Log file should have content after logging")
	}
}

func TestNewLogger_EmptyFile(t *testing.T) {
	cfg := Config{Level: "info", File: ""}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Fatal("NewLogger() returned nil logger")
	}
	logger.Sync()
}

func TestNewLogger_DefaultLevel(t *testing.T) {
	cfg := Config{Level: ""}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}
	if logger == nil {
		t.Fatal("NewLogger() returned nil logger")
	}
	logger.Sync()
}

func TestConfig(t *testing.T) {
	cfg := Config{
		Level: "debug",
		File:  "/tmp/test.log",
	}

	if cfg.Level != "debug" {
		t.Errorf("Level = %q, want %q", cfg.Level, "debug")
	}
	if cfg.File != "/tmp/test.log" {
		t.Errorf("File = %q, want %q", cfg.File, "/tmp/test.log")
	}
}

func TestNewLogger_CanLog(t *testing.T) {
	cfg := Config{Level: "info"}
	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	// Test all log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// Test structured logging
	logger.Info("structured log", zap.String("key", "value"), zap.Int("number", 42))

	// Sync may fail in test environments due to stdout redirection
	_ = logger.Sync()
}
