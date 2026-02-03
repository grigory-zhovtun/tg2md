package writer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grigoriizhovtun/tg2md/internal/sanitizer"
)

// monthNames maps month number to English lowercase name.
var monthNames = []string{
	"",         // 0 - not used
	"january",
	"february",
	"march",
	"april",
	"may",
	"june",
	"july",
	"august",
	"september",
	"october",
	"november",
	"december",
}

// Writer handles file output with monthly splitting.
type Writer struct {
	outputDir     string
	groupName     string
	sanitizedName string
	currentFile   *os.File
	currentWriter *bufio.Writer
	currentMonth  string
	stats         map[string]int
	fileCount     int
}

// New creates a new Writer for the given group.
func New(basePath, groupName string) (*Writer, error) {
	sanitizedName := sanitizer.SanitizeName(groupName)
	outputDir := filepath.Join(basePath, sanitizedName)

	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	return &Writer{
		outputDir:     outputDir,
		groupName:     groupName,
		sanitizedName: sanitizedName,
		stats:         make(map[string]int),
	}, nil
}

// GetOutputDir returns the output directory path.
func (w *Writer) GetOutputDir() string {
	return w.outputDir
}

// WriteMessage writes a message to the appropriate monthly file.
func (w *Writer) WriteMessage(formattedLine string, timestamp time.Time) error {
	monthKey := getMonthKey(timestamp)

	// Switch file if month changed
	if monthKey != w.currentMonth {
		if err := w.switchToMonth(monthKey); err != nil {
			return err
		}
	}

	// Write message with blank line separator
	_, err := w.currentWriter.WriteString(formattedLine + "\n\n")
	if err != nil {
		return fmt.Errorf("write message: %w", err)
	}

	w.stats[monthKey]++
	return nil
}

// GetStats returns the monthly breakdown of written messages.
func (w *Writer) GetStats() map[string]int {
	return w.stats
}

// GetFileCount returns the number of files created.
func (w *Writer) GetFileCount() int {
	return w.fileCount
}

// Close flushes and closes any open files.
func (w *Writer) Close() error {
	if w.currentWriter != nil {
		if err := w.currentWriter.Flush(); err != nil {
			return fmt.Errorf("flush writer: %w", err)
		}
	}
	if w.currentFile != nil {
		if err := w.currentFile.Close(); err != nil {
			return fmt.Errorf("close file: %w", err)
		}
	}
	return nil
}

// switchToMonth closes current file and opens a new one for the given month.
func (w *Writer) switchToMonth(monthKey string) error {
	// Close current file if open
	if w.currentWriter != nil {
		if err := w.currentWriter.Flush(); err != nil {
			return fmt.Errorf("flush writer: %w", err)
		}
	}
	if w.currentFile != nil {
		if err := w.currentFile.Close(); err != nil {
			return fmt.Errorf("close file: %w", err)
		}
	}

	// Create new file
	filename := fmt.Sprintf("%s_%s.md", w.sanitizedName, monthKey)
	filePath := filepath.Join(w.outputDir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("create file %s: %w", filename, err)
	}

	w.currentFile = file
	w.currentWriter = bufio.NewWriter(file)
	w.currentMonth = monthKey
	w.stats[monthKey] = 0
	w.fileCount++

	return nil
}

// getMonthKey generates the month key (e.g., "january_2024").
func getMonthKey(t time.Time) string {
	month := strings.ToLower(monthNames[t.Month()])
	year := t.Year()
	return fmt.Sprintf("%s_%d", month, year)
}
