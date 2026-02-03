package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestTextContent_UnmarshalString(t *testing.T) {
	input := `"Hello, World!"`

	var tc TextContent
	err := json.Unmarshal([]byte(input), &tc)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if tc.Plain != "Hello, World!" {
		t.Errorf("Plain = %q, want %q", tc.Plain, "Hello, World!")
	}
	if len(tc.Entities) != 0 {
		t.Errorf("Entities should be empty, got %d items", len(tc.Entities))
	}
}

func TestTextContent_UnmarshalArray(t *testing.T) {
	input := `[{"type": "plain", "text": "Hello "}, {"type": "bold", "text": "World"}]`

	var tc TextContent
	err := json.Unmarshal([]byte(input), &tc)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if tc.Plain != "" {
		t.Errorf("Plain should be empty, got %q", tc.Plain)
	}
	if len(tc.Entities) != 2 {
		t.Errorf("Entities length = %d, want 2", len(tc.Entities))
	}
	if tc.Entities[0].Type != "plain" {
		t.Errorf("First entity type = %q, want %q", tc.Entities[0].Type, "plain")
	}
	if tc.Entities[1].Type != "bold" {
		t.Errorf("Second entity type = %q, want %q", tc.Entities[1].Type, "bold")
	}
}

func TestTextContent_UnmarshalEmptyString(t *testing.T) {
	input := `""`

	var tc TextContent
	err := json.Unmarshal([]byte(input), &tc)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if tc.Plain != "" {
		t.Errorf("Plain = %q, want empty string", tc.Plain)
	}
}

func TestTextContent_UnmarshalEmptyArray(t *testing.T) {
	input := `[]`

	var tc TextContent
	err := json.Unmarshal([]byte(input), &tc)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(tc.Entities) != 0 {
		t.Errorf("Entities should be empty, got %d items", len(tc.Entities))
	}
}

func TestParser_GetChatInfo(t *testing.T) {
	// Create temp file with test data
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")

	testData := `{
		"name": "Рабочий чат",
		"type": "private_supergroup",
		"id": 1234567890,
		"messages": []
	}`

	if err := os.WriteFile(tempFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p, err := New(tempFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer p.Close()

	name, chatType, err := p.GetChatInfo()
	if err != nil {
		t.Fatalf("GetChatInfo failed: %v", err)
	}

	if name != "Рабочий чат" {
		t.Errorf("name = %q, want %q", name, "Рабочий чат")
	}
	if chatType != "private_supergroup" {
		t.Errorf("chatType = %q, want %q", chatType, "private_supergroup")
	}
}

func TestParser_StreamMessages_Basic(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")

	testData := `{
		"name": "Test Chat",
		"type": "private_group",
		"id": 123,
		"messages": [
			{
				"id": 1,
				"type": "message",
				"date": "2024-01-15T14:30:00",
				"from": "Иван",
				"text": "Привет!"
			},
			{
				"id": 2,
				"type": "message",
				"date": "2024-01-15T14:31:00",
				"from": "Мария",
				"text": "Привет! Как дела?"
			}
		]
	}`

	if err := os.WriteFile(tempFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p, err := New(tempFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer p.Close()

	// Get chat info first
	_, _, err = p.GetChatInfo()
	if err != nil {
		t.Fatalf("GetChatInfo failed: %v", err)
	}

	messages := make([]*Message, 0)
	for result := range p.StreamMessages() {
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
			continue
		}
		messages = append(messages, result.Message)
	}

	if len(messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(messages))
	}

	if messages[0].From != "Иван" {
		t.Errorf("First message from = %q, want %q", messages[0].From, "Иван")
	}
	if messages[0].Text.Plain != "Привет!" {
		t.Errorf("First message text = %q, want %q", messages[0].Text.Plain, "Привет!")
	}

	if messages[1].From != "Мария" {
		t.Errorf("Second message from = %q, want %q", messages[1].From, "Мария")
	}
}

func TestParser_StreamMessages_WithEntities(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")

	testData := `{
		"name": "Test Chat",
		"type": "private_group",
		"id": 123,
		"messages": [
			{
				"id": 1,
				"type": "message",
				"date": "2024-01-15T14:30:00",
				"from": "Анна",
				"text": [
					{"type": "plain", "text": "Это "},
					{"type": "bold", "text": "важный"},
					{"type": "plain", "text": " текст"}
				]
			}
		]
	}`

	if err := os.WriteFile(tempFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p, err := New(tempFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer p.Close()

	_, _, _ = p.GetChatInfo()

	var msg *Message
	for result := range p.StreamMessages() {
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
			continue
		}
		msg = result.Message
	}

	if msg == nil {
		t.Fatal("No message received")
	}

	if len(msg.Text.Entities) != 3 {
		t.Fatalf("Expected 3 entities, got %d", len(msg.Text.Entities))
	}

	if msg.Text.Entities[1].Type != "bold" {
		t.Errorf("Second entity type = %q, want %q", msg.Text.Entities[1].Type, "bold")
	}
	if msg.Text.Entities[1].Text != "важный" {
		t.Errorf("Second entity text = %q, want %q", msg.Text.Entities[1].Text, "важный")
	}
}

func TestParser_StreamMessages_ServiceMessage(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.json")

	testData := `{
		"name": "Test Chat",
		"type": "private_group",
		"id": 123,
		"messages": [
			{
				"id": 1,
				"type": "service",
				"date": "2024-01-15T14:30:00",
				"actor": "Иван",
				"action": "invite_members",
				"text": ""
			}
		]
	}`

	if err := os.WriteFile(tempFile, []byte(testData), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	p, err := New(tempFile)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	defer p.Close()

	_, _, _ = p.GetChatInfo()

	var msg *Message
	for result := range p.StreamMessages() {
		if result.Error != nil {
			t.Errorf("Unexpected error: %v", result.Error)
			continue
		}
		msg = result.Message
	}

	if msg == nil {
		t.Fatal("No message received")
	}

	if msg.Type != "service" {
		t.Errorf("Message type = %q, want %q", msg.Type, "service")
	}
	if msg.Action != "invite_members" {
		t.Errorf("Action = %q, want %q", msg.Action, "invite_members")
	}
	if msg.Actor != "Иван" {
		t.Errorf("Actor = %q, want %q", msg.Actor, "Иван")
	}
}

func TestParser_FileNotFound(t *testing.T) {
	_, err := New("/nonexistent/file.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}
