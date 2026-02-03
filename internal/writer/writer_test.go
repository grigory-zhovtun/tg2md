package writer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriter_CreatesSanitizedDirectory(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Рабочий чат")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer w.Close()

	expectedDir := filepath.Join(tempDir, "Рабочий_чат")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %q to exist", expectedDir)
	}
}

func TestWriter_CreatesSanitizedDirectoryWithSpecialChars(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Test: Chat / Group")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer w.Close()

	expectedDir := filepath.Join(tempDir, "Test_Chat_Group")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("Expected directory %q to exist", expectedDir)
	}
}

func TestWriter_CreatesMonthlyFiles(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Test Chat")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Write message for January 2024
	jan := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	err = w.WriteMessage("[2024-01-15 14:30] Иван: Привет!", jan)
	if err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}

	w.Close()

	// Check file was created
	expectedFile := filepath.Join(tempDir, "Test_Chat", "Test_Chat_january_2024.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %q to exist", expectedFile)
	}

	// Check content
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !strings.Contains(string(content), "Привет!") {
		t.Errorf("File content should contain message text")
	}
}

func TestWriter_SwitchesFilesOnMonthChange(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Test Chat")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Write message for January 2024
	jan := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	err = w.WriteMessage("[2024-01-15 14:30] Иван: Сообщение в январе", jan)
	if err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}

	// Write message for February 2024
	feb := time.Date(2024, time.February, 10, 10, 0, 0, 0, time.UTC)
	err = w.WriteMessage("[2024-02-10 10:00] Мария: Сообщение в феврале", feb)
	if err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}

	w.Close()

	// Check both files were created
	janFile := filepath.Join(tempDir, "Test_Chat", "Test_Chat_january_2024.md")
	febFile := filepath.Join(tempDir, "Test_Chat", "Test_Chat_february_2024.md")

	if _, err := os.Stat(janFile); os.IsNotExist(err) {
		t.Errorf("Expected January file %q to exist", janFile)
	}

	if _, err := os.Stat(febFile); os.IsNotExist(err) {
		t.Errorf("Expected February file %q to exist", febFile)
	}

	// Check content of each file
	janContent, _ := os.ReadFile(janFile)
	febContent, _ := os.ReadFile(febFile)

	if !strings.Contains(string(janContent), "январе") {
		t.Errorf("January file should contain January message")
	}
	if strings.Contains(string(janContent), "феврале") {
		t.Errorf("January file should NOT contain February message")
	}

	if !strings.Contains(string(febContent), "феврале") {
		t.Errorf("February file should contain February message")
	}
	if strings.Contains(string(febContent), "январе") {
		t.Errorf("February file should NOT contain January message")
	}
}

func TestWriter_GetStats_ReturnsCorrectCounts(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Test Chat")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Write 3 messages for January
	jan := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	for i := 0; i < 3; i++ {
		err = w.WriteMessage("Message", jan)
		if err != nil {
			t.Fatalf("WriteMessage failed: %v", err)
		}
	}

	// Write 2 messages for February
	feb := time.Date(2024, time.February, 10, 10, 0, 0, 0, time.UTC)
	for i := 0; i < 2; i++ {
		err = w.WriteMessage("Message", feb)
		if err != nil {
			t.Fatalf("WriteMessage failed: %v", err)
		}
	}

	w.Close()

	stats := w.GetStats()

	if stats["january_2024"] != 3 {
		t.Errorf("January count = %d, want 3", stats["january_2024"])
	}
	if stats["february_2024"] != 2 {
		t.Errorf("February count = %d, want 2", stats["february_2024"])
	}
}

func TestWriter_GetFileCount(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Test Chat")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	// Write to 3 different months
	jan := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)
	feb := time.Date(2024, time.February, 15, 0, 0, 0, 0, time.UTC)
	mar := time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)

	w.WriteMessage("Jan message", jan)
	w.WriteMessage("Feb message", feb)
	w.WriteMessage("Mar message", mar)

	w.Close()

	if w.GetFileCount() != 3 {
		t.Errorf("FileCount = %d, want 3", w.GetFileCount())
	}
}

func TestWriter_CyrillicGroupName(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Рабочий чат команды")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	jan := time.Date(2024, time.January, 15, 14, 30, 0, 0, time.UTC)
	err = w.WriteMessage("Тестовое сообщение", jan)
	if err != nil {
		t.Fatalf("WriteMessage failed: %v", err)
	}

	w.Close()

	// Check file was created with Cyrillic name
	expectedFile := filepath.Join(tempDir, "Рабочий_чат_команды", "Рабочий_чат_команды_january_2024.md")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Errorf("Expected file %q to exist", expectedFile)
	}
}

func TestWriter_GetOutputDir(t *testing.T) {
	tempDir := t.TempDir()

	w, err := New(tempDir, "Test Chat")
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer w.Close()

	expectedDir := filepath.Join(tempDir, "Test_Chat")
	if w.GetOutputDir() != expectedDir {
		t.Errorf("GetOutputDir() = %q, want %q", w.GetOutputDir(), expectedDir)
	}
}

func TestGetMonthKey(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "january",
			time:     time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC),
			expected: "january_2024",
		},
		{
			name:     "december",
			time:     time.Date(2023, time.December, 31, 23, 59, 59, 0, time.UTC),
			expected: "december_2023",
		},
		{
			name:     "february",
			time:     time.Date(2025, time.February, 1, 0, 0, 0, 0, time.UTC),
			expected: "february_2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMonthKey(tt.time)
			if result != tt.expected {
				t.Errorf("getMonthKey() = %q, want %q", result, tt.expected)
			}
		})
	}
}
