package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/grigoriizhovtun/tg2md/internal/converter"
	"github.com/grigoriizhovtun/tg2md/internal/logger"
	"github.com/grigoriizhovtun/tg2md/internal/parser"
	"github.com/grigoriizhovtun/tg2md/internal/sanitizer"
	"github.com/grigoriizhovtun/tg2md/internal/writer"
)

func main() {
	// Parse arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: tg2md <input.json> [output_path]")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputPath := "."
	if len(os.Args) >= 3 {
		outputPath = os.Args[2]
	}

	// Run conversion
	if err := run(inputFile, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(inputFile, outputPath string) error {
	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputFile)
	}

	// Create parser
	p, err := parser.New(inputFile)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer p.Close()

	// Get chat info
	chatName, _, err := p.GetChatInfo()
	if err != nil {
		return fmt.Errorf("parse chat info: %w", err)
	}

	// Sanitize group name and create output directory
	sanitizedName := sanitizer.SanitizeName(chatName)
	groupDir := filepath.Join(outputPath, sanitizedName)

	if err := os.MkdirAll(groupDir, 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Initialize logger
	log, err := logger.New(filepath.Join(groupDir, "errors.log"))
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer log.Close()

	log.Info("Загрузка: %s", inputFile)
	log.Info("Группа: %s", chatName)

	// Initialize converter and writer
	conv := converter.New()
	w, err := writer.New(outputPath, chatName)
	if err != nil {
		return fmt.Errorf("init writer: %w", err)
	}
	defer w.Close()

	// Process messages
	var totalCount, processedCount, skippedCount int

	for result := range p.StreamMessages() {
		totalCount++

		if result.Error != nil {
			log.LogError(0, result.Error.Error())
			skippedCount++
			continue
		}

		msg := result.Message

		// Convert message
		formatted, timestamp, err := conv.ConvertMessage(msg)
		if err != nil {
			log.LogError(msg.ID, err.Error())
			skippedCount++
			continue
		}

		// Write to file
		if err := w.WriteMessage(formatted, timestamp); err != nil {
			log.LogError(msg.ID, err.Error())
			skippedCount++
			continue
		}

		processedCount++
	}

	// Print stats
	log.Info("Найдено сообщений: %d", totalCount)

	// Print monthly breakdown
	for month, count := range w.GetStats() {
		if count > 0 {
			log.Info("Обработка: %s (%d сообщений)", month, count)
		}
	}

	log.Success("Готово! Создано %d файлов, пропущено %d сообщений",
		w.GetFileCount(), skippedCount)

	return nil
}
