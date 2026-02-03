package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogError_WritesToFile(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "errors.log")

	log, err := New(logFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	log.LogError(123, "Test error message")
	log.LogError(456, "Another error")

	log.Close()

	// Check file content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Message ID 123: Test error message") {
		t.Errorf("Log file should contain first error message")
	}
	if !strings.Contains(contentStr, "Message ID 456: Another error") {
		t.Errorf("Log file should contain second error message")
	}
}

func TestNew_CreatesFile(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "errors.log")

	log, err := New(logFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer log.Close()

	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Errorf("Expected log file to be created")
	}
}

func TestNew_InvalidPath_ReturnsError(t *testing.T) {
	_, err := New("/nonexistent/dir/errors.log")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

func TestLogger_Close(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "errors.log")

	log, err := New(logFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	err = log.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Calling close again should not panic
	err = log.Close()
	// Second close may return error (already closed) but should not panic
}

func TestColorize_NoColorsWhenDisabled(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "errors.log")

	log, err := New(logFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer log.Close()

	// Force disable colors
	log.useColors = false

	result := log.colorize(colorRed, "test")
	if result != "test" {
		t.Errorf("colorize with disabled colors should return plain text, got %q", result)
	}
}

func TestColorize_WithColors(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "errors.log")

	log, err := New(logFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer log.Close()

	// Force enable colors
	log.useColors = true

	result := log.colorize(colorRed, "test")
	expected := colorRed + "test" + colorReset
	if result != expected {
		t.Errorf("colorize with enabled colors = %q, want %q", result, expected)
	}
}
