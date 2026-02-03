package logger

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/term"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

// Logger provides colored console output and error file logging.
type Logger struct {
	errorFile *os.File
	errorMu   sync.Mutex
	useColors bool
}

// New creates a new Logger instance.
// errorLogPath specifies where to write error logs.
func New(errorLogPath string) (*Logger, error) {
	file, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("create error log: %w", err)
	}

	return &Logger{
		errorFile: file,
		useColors: supportsColors(),
	}, nil
}

// supportsColors checks if the terminal supports ANSI colors.
func supportsColors() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// colorize wraps text with ANSI color codes if colors are supported.
func (l *Logger) colorize(color, text string) string {
	if l.useColors {
		return color + text + colorReset
	}
	return text
}

// Info prints blue informational message (progress info).
func (l *Logger) Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(l.colorize(colorBlue, "[INFO] "+msg))
}

// Success prints green success message (success, statistics).
func (l *Logger) Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(l.colorize(colorGreen, "[OK] "+msg))
}

// Warning prints yellow warning message.
func (l *Logger) Warning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(l.colorize(colorYellow, "[WARN] "+msg))
}

// Error prints red error message to console.
func (l *Logger) Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(l.colorize(colorRed, "[ERROR] "+msg))
}

// LogError writes error details to the errors.log file.
func (l *Logger) LogError(msgID int64, reason string) {
	l.errorMu.Lock()
	defer l.errorMu.Unlock()

	if l.errorFile != nil {
		fmt.Fprintf(l.errorFile, "Message ID %d: %s\n", msgID, reason)
	}
}

// Close closes the error log file.
func (l *Logger) Close() error {
	if l.errorFile != nil {
		return l.errorFile.Close()
	}
	return nil
}
